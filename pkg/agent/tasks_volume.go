package agent

import (
	"github.com/productiveops/dokemon/pkg/dockerapi"
	"github.com/productiveops/dokemon/pkg/messages"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func handleDockerVolumeList(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerVolumeList](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.VolumeList(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerVolumeListResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerVolumeRemove(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerVolumeRemove](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.VolumeRemove(m)
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

func handleDockerVolumesPrune(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerVolumesPrune](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.VolumesPrune(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerVolumesPruneResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}