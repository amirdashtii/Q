package models

type User struct {
	DBModel
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	PhoneNumber *string `json:"phone_number"`
	Role        string  `json:"role"`
	Disabled    bool    `json:"disabled"`
}
