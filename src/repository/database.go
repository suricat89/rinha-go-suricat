package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suricat89/rinha-2024-q1/src/model"
)

type DatabaseRepository struct {
	DBPool *pgxpool.Pool
}

func NewDatabaseRepository(DBPool *pgxpool.Pool) *DatabaseRepository {
	return &DatabaseRepository{DBPool}
}

func (r *DatabaseRepository) CreateTransaction(transaction *model.TransactionModel) error {
	conn, err := r.DBPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`INSERT INTO "transaction" ("customer_id", "value", "type", "description", "datetime")
     VALUES ($1, $2, $3, $4, $5)`,
		transaction.Customer.Id,
		transaction.Value,
		transaction.Type,
		transaction.Description,
		transaction.Datetime,
	)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE "customer" SET "balance" = $1 WHERE "id" = $2`,
		transaction.Customer.Balance,
		transaction.Customer.Id,
	)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return tx.Commit(context.Background())
}

func (r *DatabaseRepository) GetCustomerData(customerId int) (*model.CustomerModel, error) {
	conn, err := r.DBPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		context.Background(),
		`SELECT "id", "limit", "balance" FROM "customer" WHERE "id" = $1`,
		customerId,
	)
	if err != nil {
		return nil, err
	}

  if !rows.Next() {
		return nil, nil
	}

	result := new(model.CustomerModel)
	err = rows.Scan(
		&result.Id,
		&result.Limit,
		&result.Balance,
	)

	return result, err
}

func (r *DatabaseRepository) GetCustomerTransactions(customerId int, amount int) (*model.CustomerModel, error) {
	conn, err := r.DBPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		context.Background(),
		`SELECT c."id", c."limit", c."balance", t."value", t."type", t."description", t."datetime"
     FROM "customer" c, "transaction" t
     WHERE c.id = t.customer_id
     AND c.id = $1
     ORDER BY t.datetime DESC
     LIMIT $2`,
		customerId,
    amount,
	)
  if err != nil {
    return nil, err
  }

  customer := new(model.CustomerModel)
  customer.Transactions = make([]*model.TransactionModel, 0)
  firstRow := true
  
  for rows.Next() {
    var (
      Id int
      Limit int
      Balance int
      Value int
      Type string
      Description string
      Datetime time.Time
    )

    transaction := new(model.TransactionModel)

    err := rows.Scan(
      &Id,
      &Limit,
      &Balance,
      &Value,
      &Type,
      &Description,
      &Datetime,
    )
    if err != nil {
      return customer, err
    }

    if firstRow {
      customer.Id = Id
      customer.Balance = Balance
      customer.Limit = Limit
      firstRow = false
    }
    transaction.Value = Value
    transaction.Type = Type
    transaction.Description = Description
    transaction.Datetime = Datetime

    customer.Transactions = append(customer.Transactions, transaction)
  }

  return customer, nil
}
