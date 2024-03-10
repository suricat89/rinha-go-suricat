package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/suricat89/rinha-2024-q1/src/interfaces"
	"github.com/suricat89/rinha-2024-q1/src/model"
	"github.com/suricat89/rinha-2024-q1/src/utils"
)

type CustomerController struct {
	databaseRepository interfaces.DatabaseRepository
	cacheRepository    interfaces.CacheRepository
}

type RequestNewTransaction struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type ResponseNewTransaction struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

type ResponseBalance struct {
	Total    int       `json:"total"`
	Datetime time.Time `json:"data_extrato"`
	Limit    int       `json:"limite"`
}

type ResponseTransactionItem struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	Datetime    time.Time `json:"realizada_em"`
}

type ResponseGetTransactions struct {
	Balance      *ResponseBalance           `json:"saldo"`
	Transactions []*ResponseTransactionItem `json:"ultimas_transacoes"`
}

func NewCustomerController(
	databaseRepository interfaces.DatabaseRepository,
	cacheRepository interfaces.CacheRepository,
) *CustomerController {
	return &CustomerController{
		databaseRepository,
		cacheRepository,
	}
}

func (cc *CustomerController) NewTransaction(c fiber.Ctx) error {
	t := utils.NewTraceTime("controller", "NewTransaction")

	t.Start("Validate request")
	customerId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"message": "Malformed Customer ID",
			"error":   err,
		})
	}

	reqBody := new(RequestNewTransaction)
	c.Bind().Body(reqBody)

	if reqBody.Value <= 0 {
		return c.Status(422).JSON(&fiber.Map{
			"message": "Value must be a positive number",
			"error":   nil,
		})
	}
	if reqBody.Type != "c" && reqBody.Type != "d" {
		return c.Status(422).JSON(&fiber.Map{
			"message": `Type must be "c" or "d"`,
			"error":   nil,
		})
	}
	if reqBody.Description == "" {
		return c.Status(422).JSON(&fiber.Map{
			"message": "Description must be informed",
			"error":   nil,
		})
	}
	if len(reqBody.Description) > 10 {
		return c.Status(422).JSON(&fiber.Map{
			"message": "Description must be 10 characters max",
			"error":   nil,
		})
	}
	t.End()

	transaction := new(model.TransactionModel)
	transaction.Datetime = time.Now()
	transaction.Description = reqBody.Description
	transaction.Type = reqBody.Type
	transaction.Value = reqBody.Value

	reqUuid := uuid.NewString()

	t.Start(fmt.Sprintf("cacheRepository.WaitForCustomerLock [%s]", reqUuid))
	err = cc.cacheRepository.WaitForCustomerLock(customerId, reqUuid)
	t.End()
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"message": "Error waiting for customer lock",
			"error":   err,
		})
	}
	defer func(cause error) {
		t.Start(fmt.Sprintf("cacheRepository.UnlockCustomer [%s]", reqUuid))
		cc.cacheRepository.UnlockCustomer(customerId)
		t.End()
	}(err)

	t.Start("databaseRepository.GetCustomerData")
	customer, err := cc.databaseRepository.GetCustomerData(customerId)
	t.End()
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"message": "Error fetching database",
			"error":   err,
		})
	}
	if customer == nil {
		return c.Status(404).JSON(&fiber.Map{
			"message": "Customer not found",
			"error":   nil,
		})
	}
	transaction.Customer = customer
	transaction.Id = reqUuid

	if transaction.Type == "c" {
		transaction.Customer.Balance += transaction.Value
	} else {
		transaction.Customer.Balance -= transaction.Value

		if (transaction.Customer.Balance * -1) > transaction.Customer.Limit {
			return c.Status(422).JSON(&fiber.Map{
				"message": "This transaction exeeds the Limit",
				"error":   nil,
			})
		}
	}

	t.Start("databaseRepository.CreateTransaction")
	err = cc.databaseRepository.CreateTransaction(transaction)
	t.End()
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"message": "Error trying to persist transaction",
			"error":   err,
		})
	}

	return c.Status(200).JSON(&ResponseNewTransaction{
		Limit:   transaction.Customer.Limit,
		Balance: transaction.Customer.Balance,
	})
}

func (cc *CustomerController) GetTransactions(c fiber.Ctx) error {
	customerId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"message": "Malformed Customer ID",
			"error":   err,
		})
	}

	customer, err := cc.databaseRepository.GetCustomerTransactions(customerId, 10)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"message": "Error fetching customer transactions on database",
			"error":   err,
		})
	}

	if len(customer.Transactions) == 0 {
		customer, err = cc.databaseRepository.GetCustomerData(customerId)
		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"message": "Error fetching customer data on database",
				"error":   err,
			})
		}
		if customer == nil {
			return c.Status(404).JSON(&fiber.Map{
				"message": "Customer not found",
				"error":   nil,
			})
		}
	}

	response := new(ResponseGetTransactions)
	response.Balance = &ResponseBalance{
		Total:    customer.Balance,
		Datetime: time.Now(),
		Limit:    customer.Limit,
	}
	response.Transactions = make([]*ResponseTransactionItem, len(customer.Transactions))

	for i, transaction := range customer.Transactions {
		response.Transactions[i] = &ResponseTransactionItem{
			Value:       transaction.Value,
			Type:        transaction.Type,
			Description: transaction.Description,
			Datetime:    transaction.Datetime,
		}
	}

	return c.Status(200).JSON(response)
}
