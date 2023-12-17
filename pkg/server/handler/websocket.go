package handler

import (
	"dokemon/pkg/crypto"
	"dokemon/pkg/crypto/ske"
	"dokemon/pkg/messages"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) HandleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	messageType, messageString := messages.ReceiveRaw(ws)

	switch messageType {
	case "ConnectMessage":
		handleConnect(h, ws, messageString)
	 case "TaskSessionMessage":
	 	handleTaskSession(h, ws, messageString)
	default:
		ws.Close()
		return errors.New("invalid message")
	}

	return nil
}

func handleConnect(h *Handler, ws *websocket.Conn, messageString string) {
	m, _ := messages.Parse[messages.ConnectMessage](messageString)

	token, ok := validateConnection(h, ws, m.ConnectionToken)
	if !ok {
		return
	}

	if m.AgentVersion != "" {
		h.nodeStore.UpdateAgentVersion(token.NodeId, m.AgentVersion)
	}

	messages.Send[messages.ConnectResponseMessage](ws, messages.ConnectResponseMessage{Success: true})
	if m.AgentVersion != "" {
		log.Info().Uint("nodeId", token.NodeId).Str("agentVersion", m.AgentVersion).Msg("Agent connected")
	} else {
		log.Info().Uint("nodeId", token.NodeId).Msg("Agent connected")
	}

	err := h.nodeStore.UpdateLastPing(token.NodeId, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("Error while updating initial ping time")
	}

	connectionClosed := make(chan bool)

	go func ()  {
		for {
			_, err := messages.Receive[messages.Ping](ws)
			if err != nil {
				log.Info().Uint("nodeId", token.NodeId).Msg("Agent disconnected")
				err := h.nodeStore.UpdateLastPing(token.NodeId, time.Now().Add(-2 * time.Minute))
				if err != nil {
					log.Error().Err(err).Msg("Error while updating ping time after agent disconnected")
				}
				connectionClosed <- true
				return
			} else {
				err := h.nodeStore.UpdateLastPing(token.NodeId, time.Now())
				if err != nil {
					log.Error().Err(err).Msg("Error while updating ping time")
				}
			}
		}			
	}()

	nodeId := token.NodeId
	if messages.TaskQueue[nodeId] == nil {
		log.Debug().Uint("nodeId", nodeId).Msg("Building TaskQueue")
		messages.TaskQueue[nodeId] = make(chan messages.TaskQueuedMessage)
	}

	outer:
	for  {
		select {
		case <-connectionClosed:
			break outer
		case queuedTask := <- messages.TaskQueue[nodeId]:
			log.Debug().Str("taskId", queuedTask.TaskId).Msg("Sending task")
			err := messages.Send[messages.TaskQueuedMessage](ws, queuedTask)
			if err != nil {
				log.Warn().Str("taskId", queuedTask.TaskId).Msg("Error while sending task to agent")
				return
			}
		}
	}
}

func handleTaskSession(h *Handler, wsAgent *websocket.Conn, messageString string) {
	defer wsAgent.Close()
	
	tsm, _ := messages.Parse[messages.TaskSessionMessage](messageString)

	_, ok := validateConnection(h, wsAgent, tsm.ConnectionToken)
	if !ok {
		return
	}

	if !messages.TaskExists(tsm.TaskId) {
		messages.Send[messages.TaskSessionResponseMessage](wsAgent, messages.TaskSessionResponseMessage{Success: false, Message: "Task does not exist"})
		return
	}

	messages.Send[messages.TaskSessionResponseMessage](wsAgent, messages.TaskSessionResponseMessage{Success: true})

	if !tsm.Stream {
		for {
			messageType, messageString := messages.ReceiveRaw(wsAgent)

			switch messageType {
			case "TaskLogMessage":
				m, _ := messages.Parse[messages.TaskLogMessage](messageString)
				log.Info().Str("level", m.Level).Str("message", m.Text).Msg("Agent Log")
			case "TaskStatusMessage":
				m, _ := messages.Parse[messages.TaskStatusMessage](messageString)
				if strings.HasPrefix(m.Status, "Completed") {
					messages.TaskResponses[tsm.TaskId] <- messageString
					return	
				}
			default:
				log.Error().Str("messageType", messageType).Msg("Unexpected message type received from agent in Task Session")
				return
			}
		}
	} else {
		wsBrowser := messages.TaskSockets[tsm.TaskId]

		var taskStatusMessage messages.TaskStatusMessage

		browseClosedSession := false
		go func ()  {
			// Read from browser and send to agent
			for {
				mt, dat, err := wsBrowser.ReadMessage()
				if err != nil {
					if err == io.EOF {
						message := "EOF"
						taskStatusMessage = messages.TaskStatusMessage{Status: "CompletedWithSuccess", Result: &message}
						browseClosedSession = true
						wsAgent.Close()
						break
					}
					log.Debug().Err(err).Msg("Error while reading streaming message from browser")
					message := err.Error()
					taskStatusMessage = messages.TaskStatusMessage{Status: "CompletedWithFailure", Result: &message}
					browseClosedSession = true
					wsAgent.Close()
					break
				}

				err = wsAgent.WriteMessage(mt, dat)
				if err != nil {
					log.Debug().Err(err).Msg("Error while sending streaming message to agent")
					message := err.Error()
					taskStatusMessage = messages.TaskStatusMessage{Status: "CompletedWithFailure", Result: &message}
					wsAgent.Close()
					break
				}
			}
		}()

		// Read from agent and send to browser
		outer:
		for {
			mt, dat, err := wsAgent.ReadMessage()
			if err != nil {
				if browseClosedSession  {
					log.Debug().Err(err).Msg("Browser closed connection")
				}
				if err == io.EOF {
					message := "EOF"
					taskStatusMessage = messages.TaskStatusMessage{Status: "CompletedWithSuccess", Result: &message}
					break outer
				}
				
				log.Debug().Err(err).Msg("Error while reading streaming message from agent")
				message := err.Error()
				taskStatusMessage = messages.TaskStatusMessage{Status: "CompletedWithFailure", Result: &message}
				break outer
			}

			err = wsBrowser.WriteMessage(mt, dat)
			if err != nil {
				log.Debug().Err(err).Msg("Error while sending streaming message to browser")
				message := err.Error()
				taskStatusMessage = messages.TaskStatusMessage{Status: "CompletedWithFailure", Result: &message}
				break outer
			}
		}

		messages.TaskResponses[tsm.TaskId] <- string(messages.Serialize[messages.TaskStatusMessage](taskStatusMessage))
	}
}

func validateConnection(h *Handler, ws *websocket.Conn, encryptedConnectionToken string) (*Token, bool) {
	connectionTokenJson, err := ske.Decrypt(encryptedConnectionToken)
	if err != nil {
		ws.Close()
		return nil, false
	}

	var connectionToken Token
	json.Unmarshal([]byte(connectionTokenJson), &connectionToken)

	t, err := h.nodeStore.GetById(connectionToken.NodeId)
	if err != nil {
		log.Error().Err(err).Msg("Error while retrieving Node")
		messages.Send[messages.ConnectResponseMessage](ws, messages.ConnectResponseMessage{Success: false, Message: "Internal server error"})
		ws.Close()
		return nil, false
	}

	hash := crypto.HashString(encryptedConnectionToken)
	if hash != *t.TokenHash {
		messages.Send[messages.ConnectResponseMessage](ws, messages.ConnectResponseMessage{Success: false, Message: "Invalid token"})
		ws.Close()
		return nil, false
	}

	return &connectionToken, true
}