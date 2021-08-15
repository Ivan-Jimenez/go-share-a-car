package user

import (
	"context"
	"fmt"

	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/Ivan-Jimenez/go-share-a-car/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "users"

type Error struct {
	message string
}

func (e *Error) Error() string {
	return e.message
}

type Repository interface {
	NewUser(user *entities.User) (*entities.User, error)
	Login(credentials *entities.LoginCredentials) (*entities.User, string, error)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepo() Repository {
	collection := database.Instance.Collection(collectionName)
	return &repository{
		Collection: collection,
	}
}

func (r *repository) NewUser(user *entities.User) (*entities.User, error) {
	filter := bson.D{{Key: "email", Value: user.Email}}
	if err := r.Collection.FindOne(context.Background(), filter).Err(); err == nil {
		return nil, &Error{
			message: fmt.Sprintf("Username: %s has been taken!", user.Email),
		}
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.ID = primitive.NewObjectID()
	user.HashPassword()

	if _, err := r.Collection.InsertOne(context.Background(), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) Login(credentials *entities.LoginCredentials) (*entities.User, string, error) {
	filter := bson.D{{Key: "email", Value: credentials.Email}}
	result := r.Collection.FindOne(context.Background(), filter)

	if err := result.Err(); err != nil {
		return nil, "", &Error{message: "Invalid credentials"}
	}

	user := &entities.User{}
	result.Decode(user)

	if !user.DoPasswordMatch(credentials.Password) {
		return nil, "", &Error{message: "Invalid credentials"}
	}

	// TODO: Session token
	token := "super secret token"

	return user, token, nil
}
