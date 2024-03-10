package repository

import (
	"context"

	"github.com/suricat89/rinha-2024-q1/src/constants"
	"github.com/suricat89/rinha-2024-q1/src/interfaces"
	"github.com/suricat89/rinha-2024-q1/src/model"
	"github.com/suricat89/rinha-2024-q1/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) interfaces.DatabaseRepository {
	return &MongoRepository{db}
}

func (r *MongoRepository) CreateTransaction(transaction *model.TransactionModel) error {
	t := utils.NewTraceTime("repository/mongodb", "CreateTransaction")

	t.Start("Collection('transaction').InsertOne")
	trsRes, err := r.db.Collection(constants.MONGODB_COLLECTION_TRANSACTION).InsertOne(
		context.Background(),
		transaction.GetMongodb(),
	)
	t.End()
	if err != nil {
		return err
	}

	t.Start("Collection('customer').UpdateOne")
	_, err = r.db.Collection(constants.MONGODB_COLLECTION_CUSTOMER).UpdateOne(
		context.Background(),
		bson.D{{
			Key:   "id",
			Value: transaction.Customer.Id,
		}},
		bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "balance",
				Value: transaction.Customer.Balance,
			}},
		}},
	)
	t.End()
	if err != nil {
		t.Start("Collection('transaction').DeleteOne")
		r.db.Collection(constants.MONGODB_COLLECTION_TRANSACTION).DeleteOne(
			context.Background(),
			bson.D{{
				Key:   "_id",
				Value: trsRes.InsertedID,
			}},
		)
		t.End()
		return err
	}

	return nil
}

func (r *MongoRepository) GetCustomerData(customerId int) (*model.CustomerModel, error) {
	t := utils.NewTraceTime("repository/mongodb", "GetCustomerData")

	customer := new(model.CustomerMongodb)
	t.Start("Collection('customer').FindOne")
	err := r.db.Collection(constants.MONGODB_COLLECTION_CUSTOMER).FindOne(
		context.Background(),
		bson.D{{
			Key:   "id",
			Value: customerId,
		}},
	).Decode(&customer)
	t.End()
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return customer.GetModel(), nil
}

func (r *MongoRepository) GetCustomerTransactions(customerId int, amount int) (*model.CustomerModel, error) {
	transactions := make([]*model.TransactionMongodb, 0)

	customer, err := r.GetCustomerData(customerId)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, nil
	}

	options := new(options.FindOptions)
	options.SetSort(bson.D{{
		Key:   "datetime",
		Value: -1,
	}})
	options.SetLimit(int64(amount))

	cursor, err := r.db.Collection(constants.MONGODB_COLLECTION_TRANSACTION).Find(
		context.Background(),
		bson.D{{
			Key:   "customerId",
			Value: customerId,
		}},
		options,
	)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &transactions)
	if err != nil {
		return nil, err
	}

	customer.Transactions = make([]*model.TransactionModel, 0)
	for _, transaction := range transactions {
		customer.Transactions = append(
			customer.Transactions,
			transaction.GetModel(),
		)
	}

	return customer, nil
}
