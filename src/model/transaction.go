package model

import "time"

type TransactionModel struct {
	Id          string
	Customer    *CustomerModel
	Value       int
	Type        string
	Description string
	Datetime    time.Time
}

func (t *TransactionModel) GetMongodb() *TransactionMongodb {
	return &TransactionMongodb{
		CustomerId:  t.Customer.Id,
		Value:       t.Value,
		Type:        t.Type,
		Description: t.Description,
		Datetime:    t.Datetime,
	}
}

type TransactionMongodb struct {
	CustomerId  int       `bson:"customerId"`
	Value       int       `bson:"value"`
	Type        string    `bson:"type"`
	Description string    `bson:"description"`
	Datetime    time.Time `bson:"datetime"`
}

func (t *TransactionMongodb) GetModel() *TransactionModel {
	return &TransactionModel{
		Customer: &CustomerModel{
			Id: t.CustomerId,
		},
		Value:       t.Value,
		Type:        t.Type,
		Description: t.Description,
		Datetime:    t.Datetime,
	}
}
