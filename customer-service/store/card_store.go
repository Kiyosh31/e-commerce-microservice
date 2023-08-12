package store

import (
	"encoding/json"

	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type CardStore struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewCardStore(client *mongo.Client, database string, collection string) *CardStore {
	return &CardStore{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (store *CardStore) getCardCollection() *mongo.Collection {
	return store.client.Database(store.database).Collection(store.collection)
}

func (store *CardStore) Create(ctx context.Context, card types.Card) (*mongo.InsertOneResult, error) {
	col := store.getCardCollection()

	res, err := col.InsertOne(ctx, card)
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}

	return res, nil
}

func (store *CardStore) GetOne(ctx context.Context, id primitive.ObjectID) (types.Card, error) {
	col := store.getCardCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	var res types.Card
	err := col.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return types.Card{}, err
	}

	return res, nil
}

func (store *CardStore) GetAll(ctx context.Context, id primitive.ObjectID) ([]types.Card, error) {
	col := store.getCardCollection()
	filter := bson.D{{Key: "userId", Value: id}}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return []types.Card{}, err
	}

	var cards []types.Card
	if err = cursor.All(context.TODO(), &cards); err != nil {
		return []types.Card{}, err
	}

	for _, card := range cards {
		cursor.Decode(&card)
		_, err := json.MarshalIndent(card, "", "    ")
		if err != nil {
			return []types.Card{}, err
		}
	}

	return cards, nil
}

func (store *CardStore) Update(ctx context.Context, cardToUpdate types.Card) (*mongo.UpdateResult, error) {
	col := store.getCardCollection()
	filter := bson.D{{Key: "_id", Value: cardToUpdate.ID}}
	update := bson.D{{Key: "$set", Value: cardToUpdate}}

	res, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &mongo.UpdateResult{}, err
	}

	return res, nil
}

func (store *CardStore) Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	col := store.getCardCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	res, err := col.DeleteOne(context.TODO(), filter)
	if err != nil {
		return &mongo.DeleteResult{}, err
	}

	return res, nil
}
