package model

// GameInfo is model for ui
type GameInfo struct {
	Link        string `json:"link"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	ExServiceId string `json:"ex_service_id"`
	Token       string `json:"token"`
}

// GameBOConfig all game
type GameBOConfig struct {
	Query struct {
	} `json:"query"`
	Paging struct {
		From   int    `json:"from"`
		Size   int    `json:"size"`
		Select string `json:"select"`
		Total  int    `json:"total"`
	} `json:"paging"`
	Data []GameInfoBO `json:"data"`
}

// GameInfoBO game info bo
type GameInfoBO struct {
	ServiceName string `json:"serviceName"`
	Category    string `json:"serviceType"`
	ID          string `json:"_id"`
	ServiceID   string `json:"serviceId"`
	ExServiceID string `json:"exServiceId,omitempty"`
	Token       string `json:"token"`
	Link        string `json:"link"`
}
