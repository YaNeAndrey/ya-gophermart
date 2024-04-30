package handler

type WithdrawalsRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type OrderAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
