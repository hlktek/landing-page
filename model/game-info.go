package model

import "time"

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

// TopWinner top winner
type TopWinner struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Data []struct {
			UserID       string `json:"userId"`
			DisplayName  string `json:"displayName"`
			UserType     string `json:"userType"`
			TotalBet     int    `json:"totalBet"`
			TotalWin     int64  `json:"totalWin"`
			TotalWinLoss int64  `json:"totalWinLoss"`
		} `json:"data"`
	} `json:"data"`
	Program   string    `json:"program"`
	Version   string    `json:"version"`
	Datetime  time.Time `json:"datetime"`
	Timestamp int64     `json:"timestamp"`
}
