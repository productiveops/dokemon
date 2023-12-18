package agent

import (
	"github.com/productiveops/dokemon/pkg/dockerapi"
	"github.com/productiveops/dokemon/pkg/messages"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func handleDockerNetworkList(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerNetworkList](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")		}
		return
	}

	res, err := dockerapi.NetworkList(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerNetworkListResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}

func handleDockerNetworkRemove(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerNetworkRemove](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	err = dockerapi.NetworkRemove(m)
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

func handleDockerNetworksPrune(c *websocket.Conn, messageString string) {
	m, err := messages.Parse[dockerapi.DockerNetworksPrune](messageString)
	if err != nil {
		err := completedWithFailure(c, "Error parsing request message")
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	res, err := dockerapi.NetworksPrune(m)
	if err != nil {
		err := completedWithFailure(c, err.Error())
		if err != nil {
			log.Debug().Err(err).Msg("Error sending message to client")
		}
		return
	}

	resString := string(messages.Serialize[dockerapi.DockerNetworksPruneResponse](*res))
	err = completedWithSuccess(c, &resString)
	if err != nil {
		log.Debug().Err(err).Msg("Error sending message to client")
	}
}