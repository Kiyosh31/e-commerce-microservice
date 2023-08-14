package store

import (
	"context"

	"github.com/Kiyosh31/e-commerce-microservice/inventory/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductStore struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewProductStore(client *mongo.Client, database string, collection string) *ProductStore {
	return &ProductStore{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (store *ProductStore) getProductCollection() *mongo.Collection {
	return store.client.Database(store.database).Collection(store.collection)
}

func (store *ProductStore) Create(ctx context.Context, product types.Product) (*mongo.InsertOneResult, error) {
	col := store.getProductCollection()

	res, err := col.InsertOne(ctx, product)
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}

	return res, nil
}

func (store *ProductStore) GetOne(ctx context.Context, id primitive.ObjectID) (types.Product, error) {
	col := store.getProductCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	var res types.Product
	err := col.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return types.Product{}, err
	}

	return res, nil
}

func (store *ProductStore) Update(ctx context.Context, productToUpdate types.Product) (*mongo.UpdateResult, error) {
	col := store.getProductCollection()
	filter := bson.D{{Key: "_id", Value: productToUpdate.ID}}
	update := bson.D{{Key: "$set", Value: productToUpdate}}

	res, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		return &mongo.UpdateResult{}, err
	}

	return res, nil
}

func (store *ProductStore) Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	col := store.getProductCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	res, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return &mongo.DeleteResult{}, err
	}

	return res, nil
}
