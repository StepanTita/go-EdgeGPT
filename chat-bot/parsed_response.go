package chat_bot

type ParsedFrame struct {
	Text string
	Wrap bool
	Skip bool

	SuggestedResponses []string
}
