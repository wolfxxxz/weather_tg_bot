package database

import (
	"context"
	"fmt"
	"time"
	"weather_bot/internal/apperrors"
	"weather_bot/internal/config"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeOut = 5

func InitClient(ctx context.Context, conf *config.Config, log *logrus.Logger /*, db MongoDB*/) (*mongo.Database, error) {
	var clientOpt *options.ClientOptions
	mongoDBURL := fmt.Sprintf("mongodb://%s:%s@%s:%s", conf.UserName, conf.Password, conf.MongoHost, conf.MongoPort)
	credential := options.Credential{
		Username: conf.UserName,
		Password: conf.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeOut))
	defer cancel()

	clientOpt = options.Client().ApplyURI(mongoDBURL).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		return nil, apperrors.MongoInitFailedError.AppendMessage(fmt.Sprintf("error to connect mongoDB [%v]", err))
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, apperrors.MongoInitFailedError.AppendMessage(fmt.Sprintf("error to ping mongoDB [%v]", err))
	}

	log.Info("Ping success")

	mongoDB := client.Database(conf.DBName)
	return mongoDB, nil
}
