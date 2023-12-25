package dockerapi

import (
	"context"
	"encoding/binary"
	"io"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func ContainerList(req *DockerContainerList) (*DockerContainerListResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	dcontainers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: req.All})
	if err != nil {
		return nil, err
	}

	containers := make([]Container, len(dcontainers))
	for i, c := range dcontainers {
		ports := make([]Port, len(c.Ports))
		for j, port := range c.Ports {
			ports[j] = Port(port)
		}

		image := strings.Split(c.Image, "@")[0]
		containers[i] = Container{
			Id:     c.ID,
			Name:   c.Names[0][1:],
			Image: image,
			Status: c.Status,
			State: c.State,
			Ports: ports,
		}
	}

	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Name < containers[j].Name
	  })
	  
	return &DockerContainerListResponse{Items: containers}, nil
}

func ContainerStart(req *DockerContainerStart) (error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.Background(), req.Id, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ContainerStop(req *DockerContainerStop) (error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.ContainerStop(context.Background(), req.Id, container.StopOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ContainerRestart(req *DockerContainerRestart) (error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.ContainerRestart(context.Background(), req.Id, container.StopOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ContainerRemove(req *DockerContainerRemove) (error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(context.Background(), req.Id, types.ContainerRemoveOptions{Force: req.Force})
	if err != nil {
		return err
	}

	return nil
}

func ContainerLogs(req *DockerContainerLogs, ws *websocket.Conn) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	o := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow: true,
		Tail: "40",
	}

	r, err := cli.ContainerLogs(context.Background(), req.Id, o)
	if err != nil {
		return err
	}
	defer r.Close()

	hdr := make([]byte, 8)
	for {
		_, err := r.Read(hdr)
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				panic(err)				
			}
		}

		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = r.Read(dat)
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				panic(err)				
			}
		}
		
		err = ws.WriteMessage(websocket.BinaryMessage, dat)
		if err != nil {
			if err.Error() != "websocket: close sent" {
				log.Debug().Err(err).Msg("Error sending message to client")
			}
			log.Debug().Msg("Container logs session ended as client closed connection")
			return err
		}
	}
}

var terminalSessionId uint = 0
func ContainerTerminal(req *DockerContainerTerminal, wsBrowser *websocket.Conn) (error) {
	terminalSessionId += 1
	sessionId := terminalSessionId
	log.Debug().Uint("sessionId", terminalSessionId).Msg("Starting terminal session")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	execConfig := types.ExecConfig{
		Tty:          true,
		Detach:       false,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"sh"},
	}

	idResponse, err := cli.ContainerExecCreate(context.Background(), req.Id, execConfig)
	if err != nil {
		log.Error().Uint("sessionId", sessionId).Msg("Error while calling ContainerExecCreate")
		return err
	}

	hijackedResponse, err := cli.ContainerExecAttach(context.Background(), idResponse.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		log.Error().Uint("sessionId", sessionId).Msg("Error while calling ContainerExecAttach")
		return err
	}

	browserClosedConnection := false
	connectionClosed := make(chan bool)
	mu := setupPinging(wsBrowser, &connectionClosed)

	go func() {
		for range connectionClosed {
			log.Debug().Uint("sessionId", sessionId).Msg("Closing terminal session as client closed connection")
			hijackedResponse.Conn.Close()
			wsBrowser.Close()
			browserClosedConnection = true
			return 
		}
	}()

	go func() {
		for {
			_, inputBytes, err := wsBrowser.ReadMessage()
			if err != nil {
				if !browserClosedConnection {
					log.Debug().Uint("sessionId", sessionId).Err(err).Msg("Error while reading from client socket")
				}
				return
			}
			hijackedResponse.Conn.Write(inputBytes)
		}
	}()

	for {
        outputBytesBuffer := make([]byte, 1024)
		n, err := hijackedResponse.Conn.Read(outputBytesBuffer)
		if err != nil {
			if !browserClosedConnection {
				log.Debug().Uint("sessionId", sessionId).Msg("Error while reading from terminal. Closing session.")
			}
			return err
		}
		mu.Lock()
		if err := wsBrowser.WriteMessage(websocket.BinaryMessage, outputBytesBuffer[0:n]); err != nil {
			if !browserClosedConnection {
				log.Debug().Uint("sessionId", sessionId).Msg("Error while writing to client socket. Closing session.")
				return err
			}
			return nil
		}
		mu.Unlock()
	}
}