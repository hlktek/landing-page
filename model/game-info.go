package model

// GameInfo is model for ui
type GameInfo struct {
	Link        string `json:"link"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	ExServiceId string `json:"ex_service_id"`
	Token       string `json:"token"`
}
