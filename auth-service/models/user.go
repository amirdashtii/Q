package models

type User struct {
	DBModel
	FirstName   string  `gorm:"size:255" json:"first_name"`
	LastName    string  `gorm:"size:255" json:"last_name"`
	Email       string  `gorm:"unique" json:"email"`
	Password    string  `json:"password"`
	PhoneNumber *string `gorm:"uniqueIndex;null" json:"phone_number"`
	Role        string  `gorm:"default:user" json:"role"`
	Disabled    bool    `json:"disabled"`
}
