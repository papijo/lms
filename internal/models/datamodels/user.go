package datamodels

type CreateUser struct {
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Email        string `json:"email" sql:"unique"`
	MobileNumber string `json:"mobilenumber" sql:"unique"`
	HomeAddress  string `json:"address"`
	State        string `json:"state"`
}
