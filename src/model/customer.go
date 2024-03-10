package model

type CustomerModel struct {
	Id           int
	Limit        int
	Balance      int
	Transactions []*TransactionModel
}

func (c *CustomerModel) GetMongodb() *CustomerMongodb {
	return &CustomerMongodb{
		Id:      c.Id,
		Limit:   c.Limit,
		Balance: c.Balance,
	}
}

type CustomerMongodb struct {
	Id      int `bson:"id"`
	Limit   int `bson:"limit"`
	Balance int `bson:"balance"`
}

func (c *CustomerMongodb) GetModel() *CustomerModel {
	return &CustomerModel{
		Id:      c.Id,
		Limit:   c.Limit,
		Balance: c.Balance,
	}
}
