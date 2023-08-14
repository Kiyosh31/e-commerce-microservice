package store

import (
	"context"
	"encoding/json"

	"github.com/Kiyosh31/e-commerce-microservice/inventory/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductCommentStore struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewProductCommentStore(client *mongo.Client, database string, collection string) *ProductCommentStore {
	return &ProductCommentStore{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (store *ProductCommentStore) getProductCommentCollection() *mongo.Collection {
	return store.client.Database(store.database).Collection(store.collection)
}

func (store *ProductCommentStore) Create(ctx context.Context, productComment types.ProductComment) (*mongo.InsertOneResult, error) {
	col := store.getProductCommentCollection()

	res, err := col.InsertOne(ctx, productComment)
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}

	return res, nil
}

func (store *ProductCommentStore) GetOne(ctx context.Context, id primitive.ObjectID) (types.ProductComment, error) {
	col := store.getProductCommentCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	var res types.ProductComment
	err := col.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return types.ProductComment{}, err
	}

	return res, nil
}

func (store *ProductCommentStore) GetAll(ctx context.Context, id primitive.ObjectID) ([]types.ProductComment, error) {
	col := store.getProductCommentCollection()
	filter := bson.D{{Key: "productId", Value: id}}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return []types.ProductComment{}, err
	}

	var productComments []types.ProductComment
	if err = cursor.All(context.TODO(), &productComments); err != nil {
		return []types.ProductComment{}, err
	}

	for _, comment := range productComments {
		cursor.Decode(&comment)
		_, err := json.MarshalIndent(comment, "", "    ")
		if err != nil {
			return []types.ProductComment{}, err
		}
	}

	return productComments, nil
}
