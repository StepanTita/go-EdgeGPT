package chat_bot

type ParsedFrame struct {
	Text          string
	AdaptiveCards string
	Links         []ResponseLink

	Wrap bool
	Skip bool

	SuggestedResponses []string
}

type ResponseLink struct {
	ID    string
	URL   string
	Title string
}
