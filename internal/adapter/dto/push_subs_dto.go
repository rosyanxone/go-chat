package dto

type PushSubscriptionRequest struct {
	// UserID          uint   `json:"user_id" binding:"required"`
	Endpoint        string `json:"endpoint" binding:"required"`
	ContentEncoding string `json:"content_encoding"`
	Keys            struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// PushPayload is the JSON body delivered to the browser's service worker
// on the "push" event (lookup public/sw.js -> event.data.json()).
type PushAction struct {
	Title  string `json:"title"`
	Action string `json:"action"`
}

type PushData struct {
	URL string `json:"url,omitempty"`
}

type PushPayload struct {
	Actions []PushAction `json:"actions"`
	Title   string       `json:"title"`
	Body    string       `json:"body"`
	Icon    string       `json:"icon,omitempty"`
	Data    PushData     `json:"data"`
}

type NotifyRequest struct {
	Message     string  `json:"message" binding:"required"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Url         *string `json:"url"`
	Code        *string `json:"code"`
}
