package dockerapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"

	"github.com/productiveops/dokemon/pkg/server/store"

	"github.com/gabemarshall/pty"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func composeGetStaleStatus(projectName string) string {
	l, err := ComposeContainerList(&DockerComposeContainerList{ProjectName: projectName})
	if err != nil {
		log.Error().Err(err).Msg("Error while calling ComposeContainerList")
		return StaleStatusError
	}

	anyError := false
	anyProcessing := false

	for _, item := range l.Items {
		if item.Stale == StaleStatusYes {
			return StaleStatusYes
		}
		if item.Stale == StaleStatusError {
			anyError = true
		}
		if item.Stale == StaleStatusProcessing {
			anyProcessing = true
		}
	}

	if anyError {
		return StaleStatusError
	}

	if anyProcessing {
		return StaleStatusProcessing
	}

	return StaleStatusNo
}

func ComposeList(req *DockerComposeList) (*DockerComposeListResponse, error) {
	cmd := exec.Command("docker-compose", "ls", "-a", "--format=json")
	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var itemsInternal []ComposeItemInternal
	json.Unmarshal(outb.Bytes(), &itemsInternal)

	sort.Slice(itemsInternal, func(i, j int) bool {
		return itemsInternal[i].Name < itemsInternal[j].Name
	})

	items := make([]ComposeItem, len(itemsInternal))
	for i, itemInternal := range itemsInternal {
		items[i] = ComposeItem{
			Name: itemInternal.Name,
			Status: itemInternal.Status,
			ConfigFiles: itemInternal.ConfigFiles,
			Stale: composeGetStaleStatus(itemInternal.Name),
		}
	}

	return &DockerComposeListResponse{Items: items}, nil
}

// Note: Returns nil if project not found. Does not return error.
func ComposeGet(req *DockerComposeGet) (*ComposeItem, error) {
	l, err := ComposeList(&DockerComposeList{})
	if err != nil {
		return nil, err
	}

	var ret ComposeItem
	idx := slices.IndexFunc(l.Items, func (item ComposeItem) bool { return item.Name == req.ProjectName })
	if idx != -1 {
		ret = l.Items[idx]
	} else {
		return nil, nil
	}

	return &ret, nil
}

func ComposeContainerList(req *DockerComposeContainerList) (*DockerComposeContainerListResponse, error) {
	cmd := exec.Command("docker-compose", "-p", req.ProjectName, "ps", "-a", "--format=json")
	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var items []ComposeContainer
	
	scanner := bufio.NewScanner(&outb)
	for scanner.Scan() {
		item := ComposeContainerInternal{}
		err := json.Unmarshal([]byte(scanner.Text()), &item)
		if err != nil {
			return nil, err
		}
		stale, ok := containerStaleStatus[item.Id]
		if !ok {
			stale = StaleStatusProcessing
		}
		items = append(items, ComposeContainer{
			Id: item.Id,
			Name: item.Name,
			Image: item.Image,
			Service: item.Service,
			Status: item.Status,
			State: item.State,
			Ports: item.Ports,
			Stale: stale,
		})
    }

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return &DockerComposeContainerListResponse{Items: items}, nil
}

var composeLogsSessionId uint = 0
func ComposeLogs(req *DockerComposeLogs, ws *websocket.Conn) error {
	composeLogsSessionId += 1
	sessionId := composeLogsSessionId
	log.Debug().Int("sessionId", int(sessionId)).Msg("Starting compose logs session")

	go discardIncomingMessages(ws)
	clientClosedConnection := false
	connectionClosed := make(chan bool)
	mu := setupPinging(ws, &connectionClosed)

	cmd := exec.Command("docker-compose", "-p", req.ProjectName, "logs", "-f")
	// We use pty so that we get colours in the output. Note that it has bugs on windows and some
	// parts of the lines are not aligned correctly. We only support Linux so this is fine.
	f, err := pty.Start(cmd)
	if err != nil {
		log.Debug().Err(err).Int("sessionId", int(sessionId)).Msg("Error while running compose logs")
		return err
	}

	go func() {
		for range connectionClosed {
			clientClosedConnection = true
			f.Close()
		}
	}()

	b := make([]byte, 1024)
	for {
			n, err := f.Read(b)
			if clientClosedConnection {
				log.Debug().Err(err).Int("sessionId", int(sessionId)).Msg("Compose logs session ended as client closed connection")
				break
			}
			if n == 0 || err != nil {
				log.Debug().Err(err).Int("sessionId", int(sessionId)).Msg("Compose logs read error")
				break
			}
			mu.Lock()
			err = ws.WriteMessage(websocket.BinaryMessage, b[:n])
			if err != nil {
				if err.Error() != "websocket: close sent" {
					log.Debug().Err(err).Int("sessionId", int(sessionId)).Msg("Compose logs write error")
				}
				log.Debug().Err(err).Int("sessionId", int(sessionId)).Msg("Compose logs session ended as client closed connection")
				return err
			}
			mu.Unlock()
	}

	return nil
}

func createTempComposeFile(projectName string, definition string, variables map[string]store.VariableValue) (string, string, string, error) {
	dir, err := os.MkdirTemp("", projectName)
	if err != nil {
		log.Error().Err(err).Msg("Error while creating temp directory for compose")
		return "", "", "", err
	}

	composeFilename := filepath.Join(dir, "compose.yaml")
	composeFile, err := os.Create(composeFilename)
	if err != nil {
		log.Error().Err(err).Msg("Error while creating temp compose file")
		return "", "", "", err
	}

	_ , err = composeFile.WriteString(definition)
	if err != nil {
		log.Error().Err(err).Msg("Error while writing to temp compose file")
		return "", "", "", err
	}

	envFilename := filepath.Join(dir, ".env")
	envFile, err := os.Create(envFilename)
	if err != nil {
		log.Error().Err(err).Msg("Error while creating temp compose file")
		return "", "", "", err
	}
	
	envVars := toEnvFormat(variables)
	for _, v := range envVars {
		_ , err = envFile.WriteString(v + "\r\n")
		if err != nil {
			log.Error().Err(err).Msg("Error while writing to temp .env file")
			return "", "", "", err
		}
	}

	return dir, composeFilename, envFilename, nil
}

func toEnvFormat(variables map[string]store.VariableValue) ([]string) {
	var ret = make([]string, len(variables))

	i := 0
	for k, v := range variables {
		ret[i] = fmt.Sprintf("%s=%s", k, *v.Value)
		i++
	}

	return ret
}

func logVars(cmd *exec.Cmd, variables map[string]store.VariableValue, ws *websocket.Conn, print bool) {
	if print {
		ws.WriteMessage(websocket.TextMessage, []byte("*** SETTING BELOW VARIABLES: ***\n\n"))
	}

	keys := make([]string, 0)
	for k, _ := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		val := "[SECRET]"
		if !variables[k].IsSecret {
			val = *variables[k].Value
		}
		if print {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s=%s\n", k, val)))
		}
	}
}

func performComposeAction(action string, projectName string, definition string, variables map[string]store.VariableValue, ws *websocket.Conn, printVars bool) error {
	dir, composefile, envfile, err := createTempComposeFile(projectName, definition, variables)
	log.Debug().Str("composeFileName", composefile).Str("envFileName", envfile).Msg("Created temporary compose file and .env file")
	if err != nil {
		return err
	}
	defer func() { 
		log.Debug().Str("fileName", composefile).Msg("Deleting temporary compose file and .env file")
		os.RemoveAll(dir) 
	}()

	var cmd *exec.Cmd
	switch action {
	case "up":
		cmd = exec.Command("docker-compose", "-p", projectName, "--env-file", envfile, "-f", composefile, action, "-d")
	case "down":
		cmd = exec.Command("docker-compose", "-p", projectName, "--env-file", envfile, action)
	case "pull":
		cmd = exec.Command("docker-compose", "-p", projectName, "--env-file", envfile, "-f", composefile, action)
	default:
		panic(fmt.Errorf("unknown compose action %s", action))
	}
	logVars(cmd, variables, ws, printVars)
	
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\n*** STARTING ACTION: %s ***\n\n", action)))
	f, err := pty.Start(cmd)
	if err != nil {
		log.Error().Err(err).Msg("pty returned error")
		return err
	}

	b := make([]byte, 1024)
	for {
		n, err := f.Read(b)
		if n == 0 {
			break
		}
		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg("Error while reading from pty")
			}
			break
		}
		_ = ws.WriteMessage(websocket.BinaryMessage, b[:n])
		// We ignore websocket write errors. This is because 
		// we don't want to terminate the command execution in between 
		// causing unexpected state
	}

	err = cmd.Wait()
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\n*** COMPLETED ACTION: %s ***\n\n", action)))

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Error executing compose %s", action))
	}

	if action == "up" {
		go ContainerRefreshStaleStatus()
	}

	return nil
}

func ComposeDeploy(req *DockerComposeDeploy, ws *websocket.Conn) error {
	go discardIncomingMessages(ws)

	err := performComposeAction("pull", req.ProjectName, req.Definition, req.Variables, ws, true)
	if err != nil {
		return err
	}
	err = performComposeAction("up", req.ProjectName, req.Definition, req.Variables, ws, false)

	return err
}

func ComposePull(req *DockerComposePull, ws *websocket.Conn) error {
	go discardIncomingMessages(ws)
	err := performComposeAction("pull", req.ProjectName, req.Definition, req.Variables, ws, true)
	return err
}

func ComposeUp(req *DockerComposeUp, ws *websocket.Conn) error {
	go discardIncomingMessages(ws)
	err := performComposeAction("up", req.ProjectName, req.Definition, req.Variables, ws, true)
	return err
}

func ComposeDown(req *DockerComposeDown, ws *websocket.Conn) error {
	go discardIncomingMessages(ws)
	err := performComposeAction("down", req.ProjectName, "", nil, ws, true)
	return err
}

func ComposeDownNoStreaming(req *DockerComposeDownNoStreaming) error {
	cmd := exec.Command("docker-compose", "-p", req.ProjectName, "down")
	err := cmd.Start()
	if err != nil {
		log.Error().Err(err).Msg("Error while starting compose down")
		return err
	}

	err = cmd.Wait()
	if err != nil {
		log.Error().Err(err).Msg("Error executing compose down")
	}

	log.Debug().Msg("compose down session closed")

	return nil
}
