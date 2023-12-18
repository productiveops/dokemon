package agent

import (
	"github.com/productiveops/dokemon/pkg/dockerapi"
	"github.com/productiveops/dokemon/pkg/messages"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func handleDockerContainerList(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerContainerList](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.ContainerList(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerContainerListResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerContainerLogs(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerContainerLogs](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ContainerLogs(m, c)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerContainerTerminal(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerContainerTerminal](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ContainerTerminal(m, c)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerContainerStart(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerContainerStart](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ContainerStart(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = completedWithSuccess(c, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerContainerStop(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerContainerStop](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ContainerStop(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = completedWithSuccess(c, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerContainerRestart(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerContainerRestart](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ContainerRestart(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = completedWithSuccess(c, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerContainerRemove(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerContainerRemove](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ContainerRemove(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = completedWithSuccess(c, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}
