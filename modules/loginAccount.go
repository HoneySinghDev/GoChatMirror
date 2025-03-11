package modules

import (
	"errors"
	"log"
	"regexp"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/helpers"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"

	"github.com/erikgeiser/promptkit/textinput"
)

// Func ValidatePhoneNumber() func(string) error {
func ValidatePhoneNumber(phoneNumber string) error {
	// Create a regular expression to match a valid phone number.
	regex := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	// Check if the phone number matches the regular expression.
	if !regex.MatchString(phoneNumber) || len(phoneNumber) < 11 {
		return errors.New("invalid phone number")
	}

	return nil
}
func LoginAccountCMD(_ *helpers.UserBotConfig) error {
	input := textinput.New("Enter Your Phone Number To Login")
	input.Placeholder = "+19171234567"
	input.Validate = ValidatePhoneNumber
	ph, err := input.RunPrompt()
	if err != nil {
		log.Fatal("Error while getting phone number from user:", err)
		return err
	}

	px, err := helpers.GetRandomProxy()
	if err != nil {
		log.Fatal(err)
	}

	p := &types.LoginNewAccountParams{
		PhoneNumber: ph,
		Proxy:       px,
	}

	sessionString, err := helpers.LoginNewAccount(p)
	if err != nil {
		log.Fatal("Error while logging in:", err)
		return err
	}

	log.Println("Account Has Been Logged In Successfully")

	acc := types.SessionAccount{
		PhoneNumber:   ph,
		SessionString: sessionString,
		SessionType:   "native",
		Banned:        false,
		Proxy:         *px,
	}

	err = helpers.AddAccount(helpers.CONFIG.UserBot.SessionsFilePath, &acc)
	if err != nil {
		log.Fatalf("Error while adding account to session file: %s", err)
		return err
	}

	return nil
}
