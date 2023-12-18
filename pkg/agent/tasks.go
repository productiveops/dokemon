package agent

import (
	"slices"
	"strings"
	"time"

	"github.com/productiveops/dokemon/pkg/messages"

	"github.com/rs/zerolog/log"
)

func startTaskSession(tqm messages.TaskQueuedMessage) {
	messageType := strings.Split(tqm.TaskDefinition, " ")[0]
	taskDefinition := tqm.TaskDefinition

	log.Info().Str("messageType", messageType).Msg("Task received")
	log.Debug().Str("taskId", tqm.TaskId).Msg("Starting task session")
	c := open()
	defer c.Close()

	stream := false
	steamMessageTypes := []string{"DockerContainerLogs", "DockerContainerTerminal",
									"DockerComposePull", "DockerComposeUp", "DockerComposeDown", "DockerComposeLogs"}
	if slices.Contains(steamMessageTypes, messageType) {	
		stream = true
	}

	messages.Send(c, messages.TaskSessionMessage{ConnectionToken:  token, TaskId: tqm.TaskId, Stream: stream})
	taskSessionResponseMessage, err := messages.Receive[messages.TaskSessionResponseMessage](c)
	if err != nil {
		log.Error().Str("taskId", tqm.TaskId).Msg("Error while receiving task session response")
		return
	}

	if !taskSessionResponseMessage.Success {
		log.Error().Str("taskId", tqm.TaskId).Str("serverError", taskSessionResponseMessage.Message).Msg("Task session estabilshment failed")
		return
	}

	// We have successfully established the session. Process the task now.
	switch messageType {
	case "DockerContainerList":
		handleDockerContainerList(c, taskDefinition)
	case "DockerContainerLogs":
		handleDockerContainerLogs(c, taskDefinition)
	case "DockerContainerTerminal":
		handleDockerContainerTerminal(c, taskDefinition)
	case "DockerContainerStart":
		handleDockerContainerStart(c, taskDefinition)
	case "DockerContainerStop":
		handleDockerContainerStop(c, taskDefinition)
	case "DockerContainerRestart":
		handleDockerContainerRestart(c, taskDefinition)
	case "DockerContainerRemove":
		handleDockerContainerRemove(c, taskDefinition)
	case "DockerImageList":
		handleDockerImageList(c, taskDefinition)
	case "DockerImageRemove":
		handleDockerImageRemove(c, taskDefinition)
	case "DockerImagesPrune":
		handleDockerImagesPrune(c, taskDefinition)
	case "DockerVolumeList":
		handleDockerVolumeList(c, taskDefinition)
	case "DockerVolumeRemove":
		handleDockerVolumeRemove(c, taskDefinition)
	case "DockerVolumesPrune":
		handleDockerVolumesPrune(c, taskDefinition)
	case "DockerNetworkList":
		handleDockerNetworkList(c, taskDefinition)
	case "DockerNetworkRemove":
		handleDockerNetworkRemove(c, taskDefinition)
	case "DockerNetworksPrune":
		handleDockerNetworksPrune(c, taskDefinition)
	case "DockerComposeList":
		handleDockerComposeList(c, taskDefinition)
	case "DockerComposeGet":
		handleDockerComposeGet(c, taskDefinition)
	case "DockerComposeContainerList":
		handleDockerComposeContainerList(c, taskDefinition)
	case "DockerComposeLogs":
		handleDockerComposeLogs(c, taskDefinition)
	case "DockerComposePull":
		handleDockerComposePull(c, taskDefinition)
	case "DockerComposeUp":
		handleDockerComposeUp(c, taskDefinition)
	case "DockerComposeDown":
		handleDockerComposeDown(c, taskDefinition)
	case "DockerComposeDownNoStreaming":
		handleDockerComposeDownNoStreaming(c, taskDefinition)
	default:
		panic("Invalid message received")
	}

	// Wait until all messages are sent. If we don't sleep here then routine ends and `defer c.Close()` executes closing the
	// socket before message are sent.
	time.Sleep(1 * time.Second)

	log.Debug().Str("taskId", tqm.TaskId).Msg("Task session ended")
}
