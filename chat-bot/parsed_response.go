package chat_bot

type ParsedFrame struct {
	Text          string
	AdaptiveCards string
	Wrap          bool
	Skip          bool

	SuggestedResponses []string
}
