package chat_hub

type Message struct {
	Author      string `json:"author"`
	InputMethod string `json:"inputMethod"`
	Text        string `json:"text"`
	MessageType string `json:"messageType"`
}

type Participant struct {
	Id string `json:"id"`
}

type Argument struct {
	Source                string      `json:"source"`
	OptionsSets           []string    `json:"optionsSets"`
	AllowedMessageTypes   []string    `json:"allowedMessageTypes"`
	SliceIds              []string    `json:"sliceIds"`
	TraceId               string      `json:"traceId"`
	IsStartOfSession      bool        `json:"isStartOfSession"`
	Message               Message     `json:"message"`
	ConversationSignature string      `json:"conversationSignature"`
	Participant           Participant `json:"participant"`
	ConversationId        string      `json:"conversationId"`
}

type State struct {
	Arguments    []Argument `json:"arguments"`
	InvocationId string     `json:"invocationId"`
	Target       string     `json:"target"`
	Type         int        `json:"type"`
}
