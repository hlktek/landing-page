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
