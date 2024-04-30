package storage

import "time"

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadDate time.Time `json:"uploaded_at,omitempty"`
	Sum        float64   `json:"sum,omitempty"`
}

type Withdrawal struct {
	OrderNumber   string    `json:"order"`
	Sum           float64   `json:"sum"`
	ProcessedDate time.Time `json:"processed_at"`
}
