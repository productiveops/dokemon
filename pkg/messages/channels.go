package messages

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var (
	TaskQueue   = make(map[uint]chan TaskQueuedMessage) 	// NodeId is the map key
	TaskResponses  = make(map[string]chan string)         	// Task GUID is the map key
	TaskSockets  = make(map[string]*websocket.Conn)    		// Task GUID is the map key
)

func queueTask[T interface{}](nodeId uint, message T, ws *websocket.Conn) string {
	taskId := uuid.NewString()
	m := string(Serialize(message))

	if TaskQueue[nodeId] == nil {
		TaskQueue[nodeId] = make(chan TaskQueuedMessage)
	}
	if TaskResponses[taskId] == nil {
		TaskResponses[taskId] = make(chan string)
	}
	if ws != nil {
		TaskSockets[taskId] = ws
	}
	TaskQueue[nodeId] <- TaskQueuedMessage{TaskId: taskId, TaskDefinition: m}

	return taskId
}

func ProcessTaskWithResponse[T interface{}, R interface{}](nodeId uint, message T, timeout time.Duration) (*R, error) {
	taskId := queueTask(nodeId, message, nil)
	log.Debug().Str("taskId", taskId).Msg("Task queued")

	defer func ()  {
		delete(TaskResponses, taskId)
	}()

	afterCh := time.After(timeout)
	for {
		select {
		case m := <-TaskResponses[taskId]:
			taskStatusMessage, err := Parse[TaskStatusMessage](m)
			if err != nil {
				panic(err)
			}
	
			if strings.HasPrefix(taskStatusMessage.Status, "CompletedWithSuccess") {
				if taskStatusMessage.Result == nil {
					return nil, nil
				}
				
				res, err := Parse[R](*taskStatusMessage.Result)
				if err != nil {
					return nil, err
				}
				return res, err
			} else {
				return nil, errors.New(*taskStatusMessage.Result)
			}
			
		case <-afterCh:
			return nil, errors.New("timeout")
		}
	}
}

func ProcessTask[T interface{}](nodeId uint, message T, timeout time.Duration) error {
	taskId := queueTask(nodeId, message, nil)
	log.Debug().Str("taskId", taskId).Msg("Task queued")

	defer func ()  {
		delete(TaskResponses, taskId)
	}()

	afterCh := time.After(timeout)
	for {
		select {
		case m := <-TaskResponses[taskId]:
			taskStatusMessage, err := Parse[TaskStatusMessage](m)
			if err != nil {
				panic(err)
			}
	
			if strings.HasPrefix(taskStatusMessage.Status, "CompletedWithSuccess") {
				return nil
			} else {
				return errors.New(*taskStatusMessage.Result)
			}
			
		case <-afterCh:
			return errors.New("timeout")
		}
	}
}

func ProcessStreamTask[T interface{}](nodeId uint, message T, ws *websocket.Conn) error {
	taskId := queueTask(nodeId, message, ws)
	log.Debug().Str("taskId", taskId).Msg("Task queued")

	defer func ()  {
		delete(TaskResponses, taskId)
	}()

	for m := range TaskResponses[taskId] {
		taskStatusMessage, err := Parse[TaskStatusMessage](m)
		if err != nil {
			panic(err)
		}

		log.Debug().Str("taskId", taskId).Int("nodeId", int(nodeId)).Msg("Compose logs session ended as client closed connection")
		
		if strings.HasPrefix(taskStatusMessage.Status, "CompletedWithSuccess") {
			return nil
		} else {
			return errors.New(*taskStatusMessage.Result)
		}
	}

	return nil
}

func TaskExists(taskId string) bool {
	_, ok := TaskResponses[taskId]
	return ok
}