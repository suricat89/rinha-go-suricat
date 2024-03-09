package model

import "time"

type TransactionModel struct {
	Id             int
	Customer       *CustomerModel
	Value          int
	Type           string
	Description    string
	Datetime       time.Time
	CurrentBalance int
	CurrentLimit   int
}
