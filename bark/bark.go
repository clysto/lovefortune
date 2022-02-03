package bark

type BarkPushBody struct {
	Body      string                `json:"body"`
	DeviceKey string                `json:"device_key"`
	Title     string                `json:"title"`
	ExtParams BarkPushBodyExtParams `json:"ext_params"`
}

type BarkPushBodyExtParams struct {
	Badge int    `json:"badge"`
	Icon  string `json:"icon"`
	Group string `json:"group"`
	URL   string `json:"url"`
}
