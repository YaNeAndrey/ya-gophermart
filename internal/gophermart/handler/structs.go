package handler

type WithdrawalsRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}
