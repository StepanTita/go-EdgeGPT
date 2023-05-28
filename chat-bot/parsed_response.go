package chat_bot

type ParsedFrame struct {
	Text          string
	AdaptiveCards string
	Links         []string

	Wrap bool
	Skip bool

	SuggestedResponses []string
}
