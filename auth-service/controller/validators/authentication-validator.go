package validators

import (
	"errors"

	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/dlclark/regexp2"
)

func RegisterValidation(user *models.User) error {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	const passwordRegexPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`

	emailRe := regexp2.MustCompile(emailRegexPattern, regexp2.None)
	isEmailMatch, err := emailRe.MatchString(user.Email)
	if err != nil {
		return err
	}
	if !isEmailMatch {
		return errors.New("please enter a valid Email")
	}

	if user.PhoneNumber != nil {
		phoneRe := regexp2.MustCompile(`^\d+$`, regexp2.None)
		isPhoneMatch, err := phoneRe.MatchString(*user.PhoneNumber)
		if err != nil {
			return err
		}
		if !isPhoneMatch {
			return errors.New("please enter a valid phone number")
		}
	}

	passwordRe := regexp2.MustCompile(passwordRegexPattern, regexp2.None)
	isPasswordMatch, err := passwordRe.MatchString(user.Password)
	if err != nil {
		return err
	}
	if !isPasswordMatch {
		return errors.New("password must be at least 8 characters long and include at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

func LoginValidation(user *models.User) error {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	const passwordRegexPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`

	emailRe := regexp2.MustCompile(emailRegexPattern, regexp2.None)
	isEmailMatch, err := emailRe.MatchString(user.Email)
	if err != nil {
		return err
	}
	if !isEmailMatch {
		return errors.New("please enter a valid Email")
	}

	passwordRe := regexp2.MustCompile(passwordRegexPattern, regexp2.None)
	isPasswordMatch, err := passwordRe.MatchString(user.Password)
	if err != nil {
		return err
	}
	if !isPasswordMatch {
		return errors.New("password must be at least 8 characters long and include at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}
