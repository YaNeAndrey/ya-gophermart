package storage

import "time"

type Balance struct {
	Current   float32
	Withdrawn float32
}

type Order struct {
	Number     int
	Status     int
	Accrual    float32
	UploadDate time.Time
}

type Withdrawal struct {
	OrderNumber   int
	Sum           float32
	ProcessedDate time.Time
}
