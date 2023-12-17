package dockerapi

import (
	"bufio"
	"bytes"
	"dokemon/pkg/server/store"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"

	"github.com/gabemarshall/pty"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

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
		items[i] = ComposeItem(itemInternal)
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
		items = append(items, ComposeContainer(item))
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

func createTempComposeFile(projectName string, definition string) (string, string, error) {
	dir, err := os.MkdirTemp("", projectName)
	if err != nil {
		log.Error().Err(err).Msg("Error while creating temp directory for compose")
		return "", "", err
	}

	filename := filepath.Join(dir, "compose.yaml")
	composeFile, err := os.Create(filename)
	if err != nil {
		log.Error().Err(err).Msg("Error while creating temp compose file")
		return "", "", err
	}

	_ , err = composeFile.WriteString(definition)
	if err != nil {
		log.Error().Err(err).Msg("Error while writing to temp compose file")
		return "", "", err
	}

	return dir, filename, nil
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

func processVars(cmd *exec.Cmd, variables map[string]store.VariableValue, ws *websocket.Conn) {
	cmd.Env = os.Environ()
	ws.WriteMessage(websocket.TextMessage, []byte("Setting below variables:\n"))

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
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s=%s\n", k, val)))
	}
	for _, v := range toEnvFormat(variables) {
		cmd.Env = append(cmd.Env, v)
	}
}

func ComposePull(req *DockerComposePull, ws *websocket.Conn) error {
	go discardIncomingMessages(ws)

	dir, file, err := createTempComposeFile(req.ProjectName, req.Definition)
	log.Debug().Str("fileName", file).Msg("Created temporary compose file")
	if err != nil {
		return err
	}
	defer func() { 
		log.Debug().Str("fileName", file).Msg("Deleting temporary compose file")
		os.RemoveAll(dir) 
	}()
	
	cmd := exec.Command("docker-compose", "-p", req.ProjectName, "-f", file, "pull")
	processVars(cmd, req.Variables, ws)
	
	ws.WriteMessage(websocket.TextMessage, []byte("\nStarting action: Compose Pull\n"))
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
	if err != nil {
		log.Error().Err(err).Msg("Error executing compose pull")
	}

	log.Debug().Msg("compose pull session closed")
	return nil
}

func ComposeUp(req *DockerComposeUp, ws *websocket.Conn) error {
	go discardIncomingMessages(ws)

	dir, file, err := createTempComposeFile(req.ProjectName, req.Definition)
	log.Debug().Str("fileName", file).Msg("Created temporary compose file")
	if err != nil {
		return err
	}
	defer func() { 
		log.Debug().Str("fileName", file).Msg("Deleting temporary compose file")
		os.RemoveAll(dir) 
	}()
	
	cmd := exec.Command("docker-compose", "-p", req.ProjectName, "-f", file, "up", "-d")
	processVars(cmd, req.Variables, ws)
	
	ws.WriteMessage(websocket.TextMessage, []byte("\nStarting action: Compose Up\n"))
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
	if err != nil {
		log.Error().Err(err).Msg("Error executing compose up")
	}

	log.Debug().Msg("compose up session closed")
	return nil
}

func ComposeDown(req *DockerComposeDown, ws *websocket.Conn) error {
	go discardIncomingMessages(ws)

	cmd := exec.Command("docker-compose", "-p", req.ProjectName, "down")
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
	if err != nil {
		log.Error().Err(err).Msg("Error executing compose down")
	}

	log.Debug().Msg("compose down session closed")
	return nil
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
