package models

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             *uuid.UUID       `bson:"_id"`
	TelegramUserID string           `bson:"telegram_user_id"`
	ChatID         int64            `bson:"chat_id"`
	Coordinates    *CityCoordinates `bson:"_coordinates"`
	Updated        *Modification    `bson:"_updated"`
	Created        *Modification    `bson:"_created"`
}

type CityCoordinates struct {
	Latitude  float64 `bson:"_latitude"`
	Longitude float64 `bson:"_longitude"`
}

type Modification struct {
	At *primitive.DateTime
	By string
}

func Create(chatId int64, latitude, longitude float64) *User {
	id := uuid.New()
	now := primitive.NewDateTimeFromTime(time.Now())
	by := strconv.Itoa(int(chatId))
	modification := &Modification{At: &now, By: by}
	return &User{ID: &id, TelegramUserID: id.String(), ChatID: chatId, Coordinates: &CityCoordinates{Latitude: latitude, Longitude: longitude}, Created: modification}
}

func (person *User) Update() bool {
	now := primitive.NewDateTimeFromTime(time.Now())
	by := strconv.Itoa(int(person.ChatID))
	modification := &Modification{At: &now, By: by}
	person.Updated = modification
	return true
}
