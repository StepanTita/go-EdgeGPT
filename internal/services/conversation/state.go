package conversation

type Status struct {
	Value   *string `json:"value"`
	Message *string `json:"message"`
}

type State struct {
	ConversationSignature *string `json:"conversationSignature"`
	ClientID              *string `json:"clientId"`
	ConversationID        *string `json:"conversationId"`
	Result                *Status `json:"result"`
}
