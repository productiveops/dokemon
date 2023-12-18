package agent

import (
	"github.com/productiveops/dokemon/pkg/messages"

	"github.com/gorilla/websocket"
)

func completedWithFailure(ws *websocket.Conn, errorMessage string) error {
	err := messages.Send[messages.TaskStatusMessage](ws, messages.TaskStatusMessage{Status: "CompletedWithFailure", Result: &errorMessage})
	return err
}

func completedWithSuccess(ws *websocket.Conn, result *string) error {
	err := messages.Send[messages.TaskStatusMessage](ws, messages.TaskStatusMessage{Status: "CompletedWithSuccess", Result: result})
	return err
}