package chat_hub

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"

	"github.com/StepanTita/go-EdgeGPT/common"
	"github.com/StepanTita/go-EdgeGPT/common/convert"
	"github.com/StepanTita/go-EdgeGPT/services/conversation"
)

type chatHubRequest struct {
	conversationSignature *string
	clientID              *string
	conversationID        *string
	invocationID          int

	state State
}

func newChatHubRequest(state *conversation.State) chatHubRequest {
	return chatHubRequest{
		conversationSignature: state.ConversationSignature,
		clientID:              state.ClientID,
		conversationID:        state.ConversationID,
		invocationID:          0,
	}
}

func conversationStyleToOptions(s string, options []string) []string {
	switch s {
	case "creative":
		return []string{
			"nlu_direct_response_filter",
			"deepleo",
			"disable_emoji_spoken_text",
			"responsible_ai_policy_235",
			"enablemm",
			"h3imaginative",
			"travelansgnd",
			"dv3sugg",
			"clgalileo",
			"gencontentv3",
			"dv3sugg",
			"responseos",
			"e2ecachewrite",
			"cachewriteext",
			"nodlcpcwrite",
			"travelansgnd",
			"nojbfedge",
		}
	case "balanced":
		return []string{"nlu_direct_response_filter",
			"deepleo",
			"disable_emoji_spoken_text",
			"responsible_ai_policy_235",
			"enablemm",
			"galileo",
			"dv3sugg",
			"responseos",
			"e2ecachewrite",
			"cachewriteext",
			"nodlcpcwrite",
			"travelansgnd",
			"nojbfedge",
		}
	case "precise":
		return []string{"nlu_direct_response_filter",
			"deepleo",
			"disable_emoji_spoken_text",
			"responsible_ai_policy_235",
			"enablemm",
			"galileo",
			"dv3sugg",
			"responseos",
			"e2ecachewrite",
			"cachewriteext",
			"nodlcpcwrite",
			"travelansgnd",
			"h3precise",
			"clgalileo",
			"nojbfedge",
		}
	default:
		return options
	}
}

func (r *chatHubRequest) Update(prompt, conversationStyle string, options []string, searchResults bool) error {
	r.state = State{
		Arguments: []Argument{
			{
				Source:      "cib",
				OptionsSets: conversationStyleToOptions(conversationStyle, options),
				AllowedMessageTypes: []string{
					"Chat",
					"Disengaged",
					"AdsQuery",
					"SemanticSerp",
					"GenerateContentQuery",
					"SearchQuery",
				},
				SliceIds: []string{
					"chk1cf",
					"nopreloadsscf",
					"winlongmsg2tf",
					"perfimpcomb",
					"sugdivdis",
					"sydnoinputt",
					"wpcssopt",
					"wintone2tf",
					"0404sydicnbs0",
					"405suggbs0",
					"scctl",
					"330uaugs0",
					"0329resp",
					"udscahrfon",
					"udstrblm5",
					"404e2ewrt",
					"408nodedups0",
					"403tvlansgnd",
				},
				TraceId:          common.MustGenerateRandomHex(32),
				IsStartOfSession: r.invocationID == 0,
				Message: Message{
					Author:      "user",
					InputMethod: "Keyboard",
					Text:        prompt,
					MessageType: "Chat",
				},
				ConversationSignature: convert.FromPtr(r.conversationSignature),
				Participant: Participant{
					Id: convert.FromPtr(r.clientID),
				},
				ConversationId: convert.FromPtr(r.conversationID),
			},
		},
		InvocationId: strconv.Itoa(r.invocationID),
		Target:       "chat",
		Type:         4,
	}

	if searchResults {
		r.state.Arguments[0].AllowedMessageTypes = append(r.state.Arguments[0].AllowedMessageTypes, []string{
			"InternalSearchQuery",
			"InternalSearchResult",
			"InternalLoaderMessage",
			"RenderCardRequest",
		}...)
	}

	r.invocationID += 1
	return nil
}

func (r *chatHubRequest) EncodeJson() (json.RawMessage, error) {
	b, err := json.Marshal(r.state)
	if err != nil {
		return nil, errors.Wrap(err, "failed to json encode message")
	}
	return b, nil
}
