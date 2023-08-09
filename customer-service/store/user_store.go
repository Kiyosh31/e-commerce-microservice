package store

import (
	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type UserStore struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewUserStore(client *mongo.Client, database string, collection string) *UserStore {
	return &UserStore{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (store *UserStore) getUserCollection() *mongo.Collection {
	return store.client.Database(store.database).Collection(store.collection)
}

func (store *UserStore) SigningUser(ctx context.Context, email string) (types.User, error) {
	col := store.getUserCollection()
	filter := bson.D{{Key: "email", Value: email}}

	var user types.User
	err := col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (store *UserStore) CreateUser(ctx context.Context, user types.User) (*mongo.InsertOneResult, error) {
	col := store.getUserCollection()

	res, err := col.InsertOne(ctx, user)
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}

	return res, nil
}

func (store *UserStore) GetOneUser(ctx context.Context, id primitive.ObjectID) (types.User, error) {
	col := store.getUserCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	var res types.User
	err := col.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return types.User{}, err
	}

	return res, nil
}

func (store *UserStore) UpdateUser(ctx context.Context, userToUpdate types.User) (*mongo.UpdateResult, error) {
	col := store.getUserCollection()
	filter := bson.D{{Key: "_id", Value: userToUpdate.ID}}
	update := bson.D{{Key: "$set", Value: userToUpdate}}

	res, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &mongo.UpdateResult{}, err
	}

	return res, nil
}

func (store *UserStore) DeleteUser(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	col := store.getUserCollection()
	filter := bson.D{{Key: "_id", Value: id}}

	res, err := col.DeleteOne(context.TODO(), filter)
	if err != nil {
		return &mongo.DeleteResult{}, err
	}

	return res, nil
}
