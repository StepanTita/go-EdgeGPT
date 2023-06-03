package chat_hub

import "encoding/json"

const (
	delimeter = "\x1e"
)

func appendIdentifier(message json.RawMessage) json.RawMessage {
	return json.RawMessage(string(message) + delimeter)
}
