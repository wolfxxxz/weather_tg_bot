package repositories

import (
	"context"
	"errors"
	"weather_bot/internal/apperrors"
	"weather_bot/internal/config"
	"weather_bot/internal/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const userCollection = "users"

type UserRepoInterface interface {
	SaveUserIfNotExist(ctx context.Context, user *models.User) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateModification(ctx context.Context, user *models.User) error
}

type UserRepo struct {
	log        *logrus.Logger
	collection *mongo.Collection
}

func NewUserRepo(config *config.Config, log *logrus.Logger, mongoDB *mongo.Database) *UserRepo {
	return &UserRepo{log: log, collection: mongoDB.Collection(userCollection)}
}

func (ur *UserRepo) SaveUserIfNotExist(ctx context.Context, user *models.User) error {
	if user == nil {
		return apperrors.MongoSaveUserFailedError.AppendMessage("insert Data If Not Exist user == nil")
	}

	filter := bson.M{"chat_id": user.ChatID}
	res := ur.collection.FindOne(ctx, filter)

	if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		ur.log.Errorf("Cannot find user by chat_id. Err: %+v", res.Err())
		return apperrors.MongoGetFailedError.AppendMessage(res.Err())
	}

	if res.Err() == nil {
		ur.log.Infof("The user already exists in the database with ChatID: %v", user.ChatID)
		return nil
	}

	_, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return apperrors.MongoSaveUserFailedError.AppendMessage(err)
	}

	ur.log.Info("The user has been added, successfully.")
	return nil
}

func (ur *UserRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := ur.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, apperrors.MongoGetFailedError.AppendMessage(err)
	}

	defer cursor.Close(ctx)

	return decodeUsers(ctx, cursor)
}

func (ur *UserRepo) UpdateModification(ctx context.Context, user *models.User) error {
	filter := bson.M{"chat_id": user.ChatID}
	update := bson.M{
		"$set": bson.M{
			"updated": user.Updated,
		},
	}

	res, err := ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return apperrors.MongoUpdateModFailedError.AppendMessage(err)
	}

	if res.ModifiedCount == 0 {
		return apperrors.MongoUpdateModFailedError.AppendMessage("No documents were updated")
	}

	return nil
}

func decodeUsers(ctx context.Context, cursor *mongo.Cursor) ([]*models.User, error) {
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, apperrors.MongoGetFailedError.AppendMessage(err)
		}

		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, apperrors.MongoGetFailedError.AppendMessage(err)
	}
	return users, nil
}
