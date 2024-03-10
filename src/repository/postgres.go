package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suricat89/rinha-2024-q1/src/interfaces"
	"github.com/suricat89/rinha-2024-q1/src/model"
	"github.com/suricat89/rinha-2024-q1/src/utils"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) interfaces.DatabaseRepository {
	return &PostgresRepository{db}
}

func (r *PostgresRepository) CreateTransaction(transaction *model.TransactionModel) error {
	t := utils.NewTraceTime("repository/postgres", "CreateTransaction")

	t.Start("db.Acquire")
	conn, err := r.db.Acquire(context.Background())
	t.End()
	if err != nil {
		return err
	}
	defer func(cause error) {
		t.Start("conn.Release")
		conn.Release()
		t.End()
	}(err)

	t.Start("conn.Begin")
	tx, err := conn.Begin(context.Background())
	t.End()
	if err != nil {
		return err
	}

	t.Start("tx.Exec (INSERT INTO transaction)")
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
	t.End()
	if err != nil {
		t.Start("tx.Rollback")
		tx.Rollback(context.Background())
		t.End()
		return err
	}

	t.Start("tx.Exec (UPDATE customer)")
	_, err = tx.Exec(
		context.Background(),
		`UPDATE "customer" SET "balance" = $1 WHERE "id" = $2`,
		transaction.Customer.Balance,
		transaction.Customer.Id,
	)
	t.End()
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	t.Start("tx.Commit")
	err = tx.Commit(context.Background())
	t.End()

	return err
}

func (r *PostgresRepository) GetCustomerData(customerId int) (*model.CustomerModel, error) {
	t := utils.NewTraceTime("repository/postgres", "GetCustomerData")

	t.Start("db.Acquire")
	conn, err := r.db.Acquire(context.Background())
	t.End()
	if err != nil {
		return nil, err
	}
	defer func(cause error) {
		t.Start("conn.Release")
		conn.Release()
		t.End()
	}(err)

	t.Start("conn.Query")
	rows, err := conn.Query(
		context.Background(),
		`SELECT "id", "limit", "balance" FROM "customer" WHERE "id" = $1`,
		customerId,
	)
	t.End()
	if err != nil {
		return nil, err
	}

	t.Start("rows.Next")
	hasNext := rows.Next()
	t.End()
	if !hasNext {
		return nil, nil
	}

	result := new(model.CustomerModel)
	t.Start("rows.Scan")
	err = rows.Scan(
		&result.Id,
		&result.Limit,
		&result.Balance,
	)
	t.End()

	return result, err
}

func (r *PostgresRepository) GetCustomerTransactions(customerId int, amount int) (*model.CustomerModel, error) {
	conn, err := r.db.Acquire(context.Background())
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
			Id          int
			Limit       int
			Balance     int
			Value       int
			Type        string
			Description string
			Datetime    time.Time
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
