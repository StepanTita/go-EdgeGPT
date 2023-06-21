package chat_bot

type ParsedFrame struct {
	Text          string
	AdaptiveCards string

	Sources   []ResponseLink
	Resources []ResourceLink

	Wrap bool
	Skip bool

	SuggestedResponses []string

	ErrBody *ErrorBody
}

type ErrorBody struct {
	Reason  string
	Message string
}

type ResponseLink struct {
	ID    string
	URL   string
	Title string
}

type ResourceLink struct {
	Type string
	URL  string
}
