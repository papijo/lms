package authcontrollers

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/papijo/lms/internal/models/datamodels"
	"github.com/papijo/lms/internal/models/dbmodels"
	"github.com/papijo/lms/utils/helpers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Create User (Student)
func CreateStudentUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		//Start a Database Transaction

		tx := db.Begin()

		if tx.Error != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to start a database transaction",
			})
		}

		//Defer a function to handle the transaction's success or failure
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback() // Rollback the transaction in case of a panic
			}
		}()

		//Parse Request Body
		var user_registration datamodels.CreateUser

		if err := c.Bind(&user_registration); err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Invalid Request Body",
			})
		}

		//Validate the registration struct
		validate := validator.New()
		if err := validate.Struct(user_registration); err != nil {
			tx.Rollback()
			validationErrors := err.(validator.ValidationErrors)
			var validationMsg []string
			for _, v := range validationErrors {
				validationMsg = append(validationMsg, v.Tag())
			}
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error":  "Validation failed",
				"fields": validationMsg,
			})
		}

		//Check if a user with the same email exists
		var existinguser dbmodels.User

		if err := tx.Where("email = ?", user_registration.Email).First(&existinguser).Error; err == nil {
			//User with the same email already exists
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "User with the same email already exists",
			})
		}

		password, err := helpers.GeneratePassword(8)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to generate password",
			})
		}

		// Hask the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			//Rollback Transaction if Password cannot be hashed
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to hash password",
			})
		}

		//Generate a random 4-digit verification token
		verification_token, err := helpers.GenerateVerificationToken(4)
		if err != nil {
			//Rollback Transaction if Token cannot be generated
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to generate verification token",
			})
		}

		//Hash the verification token
		hashedToken, err := bcrypt.GenerateFromPassword([]byte(verification_token), bcrypt.DefaultCost)
		if err != nil {
			//Rollback Transaction if Token cannot be hashed
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to hash verification token",
			})
		}

		//Create a student user based on the registration data
		user := dbmodels.User{
			FirstName:         user_registration.FirstName,
			LastName:          user_registration.LastName,
			Email:             user_registration.Email,
			MobileNumber:      user_registration.MobileNumber,
			Password:          string(hashedPassword),
			VerificationToken: string(hashedToken),
			IsVerified:        false,
			AdminStatus:       false,
			StaffStatus:       false,
			StudentStatus:     false,
			InternStatus:      false,
			IsActive:          true,
		}

		//Save the user in the database within the transaction
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to create user",
			})
		}

		// Commit the transaction if all operations succeed
		tx.Commit()

		//Send Registration Verification Token SMS (Add this to go-routines so that if the fail, the response would still be successful)

		//Send Registration Verification Token Email

		//Return Success Response
		return c.JSON(http.StatusCreated, echo.Map{
			"message": "Student User Account Created Successfully",
		})

	}
}
