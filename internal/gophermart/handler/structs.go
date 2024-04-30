package handler

type WithdrawalsRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type OrderAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
