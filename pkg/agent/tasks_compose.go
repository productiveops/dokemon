package agent

import (
	"dokemon/pkg/dockerapi"
	"dokemon/pkg/messages"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func handleDockerComposeList(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerComposeList](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.ComposeList(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerComposeListResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerComposeGet(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerComposeGet](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.ComposeGet(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	if res == nil {
		err = completedWithSuccess(c, nil)
	} else {
		resString := string(messages.Serialize[dockerapi.ComposeItem](*res))
		err = completedWithSuccess(c, &resString)
	}
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerComposeContainerList(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerComposeContainerList](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.ComposeContainerList(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerComposeContainerListResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerComposeLogs(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerComposeLogs](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ComposeLogs(m, c)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerComposePull(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerComposePull](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ComposePull(m, c)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerComposeUp(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerComposeUp](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ComposeUp(m, c)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerComposeDown(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerComposeDown](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ComposeDown(m, c)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerComposeDownNoStreaming(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerComposeDownNoStreaming](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ComposeDownNoStreaming(m)
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