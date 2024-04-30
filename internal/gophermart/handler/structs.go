package handler

type WithdrawalsRequest struct {
	Order int64   `json:"order"`
	Sum   float32 `json:"sum"`
}
