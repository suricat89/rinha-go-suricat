package interfaces

import "github.com/suricat89/rinha-2024-q1/src/model"

type DatabaseRepository interface {
	CreateTransaction(transaction *model.TransactionModel) error
	GetCustomerData(customerId int) (*model.CustomerModel, error)
	GetCustomerTransactions(customerId int, amount int) (*model.CustomerModel, error)
}

type CacheRepository interface {
	WaitForCustomerLock(customerId int, reqUuid string) error
	UnlockCustomer(customerId int) error
}
