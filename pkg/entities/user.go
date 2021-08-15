package entities

import (
	"log"
	"unicode"

	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" validate:"required"`
	LastName string             `json:"lastName" validate:"required"`
	Email    string             `json:"email" validate:"required,email"`
	Password string             `json:"password" validate:"required,password"`
}

type UserData struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewUserData(logger *log.Logger) *UserData {
	return &UserData{
		database.Instance.Collection("users"),
		logger,
	}
}

func (u *User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword)
	return validate.Struct(u)
}

func (u *User) HashPassword() error {
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) DoPasswordMatch(currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(currPassword))
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
