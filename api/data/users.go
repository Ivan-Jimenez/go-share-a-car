package data

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" validate:"required"`
	LastName string `json:"lastName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

const emailRegex = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
const passwordRegex = "^(?=.*[\\d\\W])(?=.*[a-z])(?=.*[A-Z]).{8,100}$"

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
	return err != nil
}

func validateEmail(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(emailRegex)
	matches := re.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1
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
