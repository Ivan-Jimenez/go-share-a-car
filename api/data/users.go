package data

type User struct {
	ID       string `json:"id,omitempty "bson:"_id,omitempty"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Validate() {

}

// func (user *User) Validate() error {

// }

// func validateEmail(fl validator.FieldValue) {

// }
