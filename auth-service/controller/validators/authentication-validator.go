package validators

import (
	"errors"

	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/dlclark/regexp2"
)

func RegisterValidation(user *models.User) error {
	err := EmailValidation(user.Email)
	if err != nil {
		return err
	}

	err = PasswordValidation(user.Password)
	if err != nil {
		return err
	}

	if user.PhoneNumber != nil {
		err := PhoneNumberValidation(*user.PhoneNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoginValidation(user *models.User) error {
	err := EmailValidation(user.Email)
	if err != nil {
		return err
	}

	err = PasswordValidation(user.Password)
	if err != nil {
		return err
	}

	return nil
}

func UpdateValidation(user *models.User) error {

	if user.Email != "" {
		err := EmailValidation(user.Email)
		if err != nil {
			return err
		}
	}

	if user.Password != "" {

		err := PasswordValidation(user.Password)
		if err != nil {
			return err
		}
	}

	if user.PhoneNumber != nil {
		err := PhoneNumberValidation(*user.PhoneNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

func EmailValidation(email string) error {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRe := regexp2.MustCompile(emailRegexPattern, regexp2.None)
	ok, err := emailRe.MatchString(email)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("please enter a valid Email")
	}
	return nil
}

func PasswordValidation(password string) error {
	const passwordRegexPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`
	passwordRe := regexp2.MustCompile(passwordRegexPattern, regexp2.None)
	ok, err := passwordRe.MatchString(password)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("password must be at least 8 characters long and include at least one uppercase letter, one lowercase letter, and one number")
	}
	return nil
}

func PhoneNumberValidation(phoneNumber string) error {
	phoneRe := regexp2.MustCompile(`^\d+$`, regexp2.None)
	ok, err := phoneRe.MatchString(phoneNumber)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("please enter a valid phone number")
	}
	return nil
}

func TokenValidation(token string) error {
	if token == "" {
		err := errors.New("no valid token")
		return err
	}
	return nil
}