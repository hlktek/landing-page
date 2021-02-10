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
	IsNewGame   bool   `isNewGame`
	Active      bool   `active`
}

// TopWinner top winner response model
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

// TopWinnerChess top winner response model
type TopWinnerChess struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Data []struct {
			User struct {
				UserID      string `json:"userId"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				BrandId     string `json:"brandId"`
			} `json:"user"`
			TotalLoseCount int    `json:"totalLoseCount"`
			TotalDrawCount int64  `json:"totalDrawCount"`
			TotalWinCount  int64  `json:"totalWinCount"`
			TotalWinLoss   string `json:"totalWinLoss"`
			WinRate        string `json:"winRate"`
		} `json:"data"`
	} `json:"data"`
	Program   string    `json:"program"`
	Version   string    `json:"version"`
	Datetime  time.Time `json:"datetime"`
	Timestamp int64     `json:"timestamp"`
}

// JackpotHistory jackpot history response model
type JackpotHistory struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Data []struct {
			UserID        string      `json:"userId"`
			ServiceID     string      `json:"serviceId"`
			DisplayName   string      `json:"displayName"`
			JackpotType   string      `json:"jackpotType"`
			BetLevel      interface{} `json:"betLevel"`
			JackpotAmount string      `json:"jackpotAmount"`
			Time          int64       `json:"time"`
			Detail        string      `json:"detail"`
			PlaySessionID string      `json:"playSessionId"`
		} `json:"data"`
		Query struct {
			Paging struct {
				From int `json:"from"`
				Size int `json:"size"`
			} `json:"paging"`
			Query struct {
				UserID    string    `json:"userId"`
				UserType  string    `json:"userType"`
				StartDate time.Time `json:"startDate"`
				EndDate   time.Time `json:"endDate"`
			} `json:"query"`
		} `json:"query"`
	} `json:"data"`
	Program   string    `json:"program"`
	Version   string    `json:"version"`
	Datetime  time.Time `json:"datetime"`
	Timestamp int64     `json:"timestamp"`
}

// Wallet response when add money
type Wallet struct {
	Code int `json:"code"`
	Data struct {
		ResBetCode string `json:"resBetCode"`
		ResWinCode string `json:"resWinCode"`
	} `json:"data"`
	Message string `json:"message"`
}
