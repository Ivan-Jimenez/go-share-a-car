package models

import (
	"context"
	"log"
	"unicode"

	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" validate:"required"`
	LastName string `json:"lastName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type UserData struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewUserData(logger *log.Logger) *UserData {
	return &UserData{
		database.Instance.Database.Collection("users"),
		logger,
	}
}

func (data *UserData) SaveUser(ctx context.Context, user *User) (*User, error) {
	user.ID = ""

	res, err := data.collection.InsertOne(ctx, user)
	if err != nil {
		data.logger.Panicf("[ERROR][SaveUser-UserData] %s", err.Error())
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: res.InsertedID}}
	return data.FindUser(ctx, filter)
}

func (data *UserData) FindUser(ctx context.Context, filter interface{}) (*User, error) {
	findUser := data.collection.FindOne(ctx, filter)
	if err := findUser.Err(); err != nil {
		return nil, err
	}

	user := &User{}
	findUser.Decode(user)
	return user, nil
}

func (user *User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword)
	return validate.Struct(user)
}

func (user *User) HashPassword() error {
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (user *User) DoPasswordMatch(currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currPassword))
	return err == nil
}

func validatePassword(fl validator.FieldLevel) bool {
	letters := 0
	eightOrMore, number, upper, especial := false, false, false, false

	password := fl.Field().String()
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			especial = true
		}
		letters++
	}

	eightOrMore = letters >= 8

	return eightOrMore && number && upper && especial
}
