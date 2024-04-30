package storage

import "time"

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadDate time.Time `json:"uploaded_at,omitempty"`
	Sum        float32   `json:"sum,omitempty"`
}

type Withdrawal struct {
	OrderNumber   string    `json:"order"`
	Sum           float32   `json:"sum"`
	ProcessedDate time.Time `json:"processed_at"`
}
