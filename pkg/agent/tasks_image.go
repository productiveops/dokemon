package agent

import (
	"dokemon/pkg/dockerapi"
	"dokemon/pkg/messages"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func handleDockerImageList(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerImageList](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.ImageList(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerImageListResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerImageRemove(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerImageRemove](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.ImageRemove(m)
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

func handleDockerImagesPrune(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerImagesPrune](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.ImagesPrune(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerImagesPruneResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}