package store

import (
	"encoding/json"

	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type AddressStore struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewAddressStore(client *mongo.Client, database string, collection string) *AddressStore {
	return &AddressStore{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (store *AddressStore) getAddressCollection() *mongo.Collection {
	return store.client.Database(store.database).Collection(store.collection)
}

func (store *AddressStore) Create(ctx context.Context, address types.Address) (*mongo.InsertOneResult, error) {
	col := store.getAddressCollection()

	res, err := col.InsertOne(ctx, address)
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}

	return res, nil
}

func (store *AddressStore) GetOne(ctx context.Context, id primitive.ObjectID) (types.Address, error) {
	col := store.getAddressCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	var res types.Address
	err := col.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return types.Address{}, err
	}

	return res, nil
}

func (store *AddressStore) GetAll(ctx context.Context, id primitive.ObjectID) ([]types.Address, error) {
	col := store.getAddressCollection()
	filter := bson.D{{Key: "userId", Value: id}}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return []types.Address{}, err
	}

	var cards []types.Address
	if err = cursor.All(context.TODO(), &cards); err != nil {
		return []types.Address{}, err
	}

	for _, card := range cards {
		cursor.Decode(&card)
		_, err := json.MarshalIndent(card, "", "    ")
		if err != nil {
			return []types.Address{}, err
		}
	}

	return cards, nil
}

func (store *AddressStore) Update(ctx context.Context, cardToUpdate types.Address) (*mongo.UpdateResult, error) {
	col := store.getAddressCollection()
	filter := bson.D{{Key: "_id", Value: cardToUpdate.ID}}
	update := bson.D{{Key: "$set", Value: cardToUpdate}}

	res, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &mongo.UpdateResult{}, err
	}

	return res, nil
}

func (store *AddressStore) Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	col := store.getAddressCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	res, err := col.DeleteOne(context.TODO(), filter)
	if err != nil {
		return &mongo.DeleteResult{}, err
	}

	return res, nil
}
