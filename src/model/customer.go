package model

type CustomerModel struct {
	Id           int
	Limit        int
	Balance      int
	Transactions []*TransactionModel
}
