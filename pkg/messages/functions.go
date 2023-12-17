package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func Send[T Message](c *websocket.Conn, message T) error {
	err := c.WriteMessage(websocket.TextMessage, Serialize(message))
	if err != nil {
		log.Error().Err(err).Msg("Error while sending message via websocket")
		return err
	}
	return nil
}

func ReceiveRaw(c *websocket.Conn) (messageType string, message string) {
	_, messageOnWire, _ := c.ReadMessage()
	messageOnWireString := string(messageOnWire)
	messageTypeName := strings.Split(messageOnWireString, " ")[0]
	return messageTypeName, messageOnWireString
}

func Receive[T Message](c *websocket.Conn) (m *T, _ error) {
	_, messageOnWire, _ := c.ReadMessage()
	messageOnWireString := string(messageOnWire)

	return Parse[T](messageOnWireString)
}

func Serialize[T interface{}](message T) []byte {
	messageJson, _ := json.Marshal(message)
	typeNamePrefix := strings.Split(fmt.Sprintf("%T ", *new(T)), ".")[1]
	messageOnWire := append([]byte(typeNamePrefix), messageJson...)
	return messageOnWire
}

func Parse[T interface{}](messageOnWireString string) (m *T, _ error) {
	typeName := strings.Split(fmt.Sprintf("%T", *new(T)), ".")[1]
	messageTypeName := strings.Split(messageOnWireString, " ")[0]

	if messageTypeName != typeName {
		return nil, errors.New(fmt.Sprintf("Expected message of type %s but received %s", typeName, messageTypeName))
	}

	messageJson := []byte(messageOnWireString[len(messageTypeName):])
	var receivedMessage *T
	json.Unmarshal(messageJson, &receivedMessage)

	return receivedMessage, nil
}

func GetMessageJson[T interface{}](messageOnWireString string) (m []byte, _ error) {
	typeName := strings.Split(fmt.Sprintf("%T", *new(T)), ".")[1]
	messageTypeName := strings.Split(messageOnWireString, " ")[0]

	if messageTypeName != typeName {
		return nil, errors.New(fmt.Sprintf("Expected message of type %s but received %s", typeName, messageTypeName))
	}

	messageJson := []byte(messageOnWireString[len(messageTypeName):])
	
	return messageJson, nil
}