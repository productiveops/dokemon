package messages

type Ping struct {
}

type ConnectMessage struct {
	ConnectionToken string `json:"connectionToken"`
	AgentVersion    string `json:"agentVersion"`
}

type ConnectResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TaskQueuedMessage struct {
	TaskId         string `json:"taskId"`
	TaskDefinition string `json:"taskDefinition"`
}

type TaskSessionMessage struct {
	ConnectionToken string `json:"connectionToken"`
	TaskId          string `json:"taskId"`
	Stream          bool   `json:"stream"`
}

type TaskSessionResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TaskLogMessage struct {
	Offset uint   `json:"offset"` // TODO:
	Level  string `json:"level"`
	Text   string `json:"text"`
}

type TaskStatusMessage struct {
	Status string  `json:"status"`
	Result *string `json:"result"`
}

type Message interface {
	Ping |
		ConnectMessage | ConnectResponseMessage |
		TaskQueuedMessage | TaskSessionMessage | TaskSessionResponseMessage | TaskStatusMessage | TaskLogMessage
}
