package chat_hub

import "time"

type ResponseMessage struct {
	Type   int     `json:"type"`
	Target *string `json:"target"`
	Item   struct {
		Result struct {
			Error   interface{} `json:"error"`
			Value   *string     `json:"value"`
			Message *string     `json:"message"`
		}
	}
	Arguments []struct {
		Messages []struct {
			Text          *string    `json:"text"`
			MessageType   *string    `json:"messageType"`
			Author        *string    `json:"author"`
			CreatedAt     *time.Time `json:"createdAt"`
			Timestamp     *time.Time `json:"timestamp"`
			MessageId     *string    `json:"messageId"`
			Offense       *string    `json:"offense"`
			AdaptiveCards []struct {
				Type    *string `json:"type"`
				Version *string `json:"version"`
				Body    []struct {
					Type    *string `json:"type"`
					Text    *string `json:"text"`
					Wrap    *bool   `json:"wrap"`
					Size    *string `json:"size,omitempty"`
					Inlines []struct {
						Type     *string `json:"type"`
						IsSubtle *bool   `json:"isSubtle"`
						Italic   *bool   `json:"italic"`
						Text     *string `json:"text"`
					} `json:"inlines"`
				} `json:"body"`
			} `json:"adaptiveCards"`
			SourceAttributions []struct {
				ProviderDisplayName *string `json:"providerDisplayName"`
				SeeMoreUrl          *string `json:"seeMoreUrl"`
				ImageFaviconUrl     *string `json:"imageFaviconUrl,omitempty"`
				SearchQuery         *string `json:"searchQuery"`
				ImageLink           *string `json:"imageLink,omitempty"`
				ImageWidth          *string `json:"imageWidth,omitempty"`
				ImageHeight         *string `json:"imageHeight,omitempty"`
				ImageFavicon        *string `json:"imageFavicon,omitempty"`
			} `json:"sourceAttributions"`
			Feedback struct {
				Tag       interface{} `json:"tag"`
				UpdatedOn interface{} `json:"updatedOn"`
				Type      *string     `json:"type"`
			} `json:"feedback"`
			ContentOrigin *string     `json:"contentOrigin"`
			Privacy       interface{} `json:"privacy"`
		} `json:"messages"`
		RequestId *string `json:"requestId"`
	} `json:"arguments"`
}
