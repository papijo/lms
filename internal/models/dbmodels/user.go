package dbmodels

import "time"

//Used to store information in the database
//Database Table (User)
type User struct {
	ID                  uint64     `gorm:"primary_key;auto_increment" json:"id"`
	FirstName           string     `json:"firstname"`
	LastName            string     `json:"lastname"`
	Email               string     `json:"email" sql:"unique"`
	MobileNumber        string     `json:"mobilenumber" sql:"unique"`
	HomeAddress         string     `json:"address"`
	State               string     `json:"state"`
	Password            string     `json:"password"`
	PasswordChangeToken string     `json:"password_change_token" sql:"password_change_token"`
	VerificationToken   string     `json:"verification_token" sql:"verification_token"`
	TwoFactorAuthToken  string     `json:"two_factor_auth_token"`
	IsVerified          bool       `json:"is_verified"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
	AdminStatus         bool       `json:"admin_status"` // true if the user is an admin, false otherwise
	StaffStatus         bool       `json:"staff_status"`
	StudentStatus       bool       `json:"student_status"`
	IsActive            bool       `json:"is_active"`
}

//Database Table(Bio Data)
type Biodata struct {
	ID                   uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID               uint64    `gorm:"unique;not null"` // Foreign key for User with unique constraint
	DateOfBirth          time.Time `json:"date_of_birth"`
	Gender               string    `json:"gender"`
	Avatar               string    `json:"avatar"`
	NextOfKin            string    `json:"next_of_kin"`
	NextOfKinPhoneNumber string    `json:"next_of_kin_phone"`
	Address              string    `json:"address"`
	City                 string    `json:"city"`
	State                string    `json:"state"`
}
