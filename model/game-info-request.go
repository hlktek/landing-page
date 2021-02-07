package model

import "time"

// TopWinnerRequest top winner request
type TopWinnerRequest struct {
	UserID      string    `json:"userId"`
	UserType    string    `json:"userType"`
	UserBrand   string    `json:"userBrand"`
	Category    string    `json:"category"`
	SelectValue int       `json:"selectValue"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	OrderBy     string    `json:"orderBy"`
	OrderType   string    `json:"orderType"`
	From        int       `json:"from"`
	Size        int       `json:"size"`
}
type TopWinnerChessRequest struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

// JackpotRequest jackpot request
type JackpotRequest struct {
	Query struct {
		UserID    string    `json:"userId"`
		UserType  string    `json:"userType"`
		StartDate time.Time `json:"startDate"`
		EndDate   time.Time `json:"endDate"`
	} `json:"query"`
	Paging struct {
		From int `json:"from"`
		Size int `json:"size"`
	} `json:"paging"`
}

type FeedBack struct {
	UserID    string `json:"userId"`
	FeedBack  string `json:"feedBack"`
	ServiceID string `json:"serviceId"`
}
