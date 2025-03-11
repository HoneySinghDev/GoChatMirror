package modules

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/helpers"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
	"github.com/pterm/pterm"
)

func AddAccountCMD() {

	options := []string{
		"Telethon Session",
		"Native Session",
		"TData Session",
	}

	result, err := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		Show()

	if err != nil {
		log.Fatal(err)
	}

	switch result {
	case "Telethon Session":
		TelethonSession()
	case "Native Session":
		NativeSession()
	case "TData Session":
		TDataSession()
	}
}

func TelethonSession() {
	ph := Input("Enter your Phone Number: ")
	sString := Input("Enter your Telethon Session String: ")

	px, err := helpers.GetRandomProxy()
	if err != nil {
		log.Fatal(err)
	}

	acc := types.SessionAccount{
		PhoneNumber:   ph,
		SessionString: sString,
		SessionType:   "telethon",
		Banned:        false,
		Proxy:         *px,
	}

	err = helpers.AddAccount(helpers.CONFIG.UserBot.SessionsFilePath, &acc)
	if err != nil {
		log.Fatalf("Error while adding account to session file: %s", err)
	}

	pterm.DefaultCenter.Print(pterm.NewStyle(pterm.BgGreen, pterm.FgCyan).Sprint("Account Has Been Added Successfully"))
}

func NativeSession() {
	ph := Input("Enter your Phone Number: ")
	sString := Input("Enter your Native Session String: ")

	px, err := helpers.GetRandomProxy()

	if err != nil {
		log.Fatal(err)
	}

	acc := types.SessionAccount{
		PhoneNumber:   ph,
		SessionString: sString,
		SessionType:   "native",
		Banned:        false,
		Proxy:         *px,
	}

	err = helpers.AddAccount(helpers.CONFIG.UserBot.SessionsFilePath, &acc)
	if err != nil {
		log.Fatalf("Error while adding account to session file: %s", err)
	}

	pterm.DefaultCenter.Print(pterm.NewStyle(pterm.BgGreen, pterm.FgCyan).Sprint("Account Has Been Added Successfully"))
}

func TDataSession() {
	ph := Input("Enter your Phone Number: ")
	sPath := Input("Enter your TData Session Path: ")

	px, err := helpers.GetRandomProxy()

	if err != nil {
		log.Fatal(err)
	}

	acc := types.SessionAccount{
		PhoneNumber:   ph,
		SessionString: sPath,
		SessionType:   "tdata",
		Banned:        false,
		Proxy:         *px,
	}

	_, stop, err := helpers.LoginClient(&types.LoginClientParams{
		PhoneNumber:   ph,
		SessionString: sPath,
		SessionType:   "tdata",
		Proxy:         px,
	})
	if err != nil {
		log.Fatalf("Error while logging in account: %s", err)
	}

	defer stop()

	err = helpers.AddAccount(helpers.CONFIG.UserBot.SessionsFilePath, &acc)
	if err != nil {
		log.Fatalf("Error while adding account to session file: %s", err)
	}

	pterm.DefaultCenter.Print(pterm.NewStyle(pterm.BgGreen, pterm.FgCyan).Sprint("Account Has Been Added Successfully"))
}

// AccountExists checks if an account already exists in the session file
func AccountExists(filePath, phoneNumber string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	var accounts []types.SessionAccount
	if err := json.NewDecoder(file).Decode(&accounts); err != nil {
		return false, err
	}

	for _, acc := range accounts {
		if acc.PhoneNumber == phoneNumber {
			return true, nil
		}
	}
	return false, nil
}

func Exit() {
	pterm.DefaultCenter.Print(pterm.NewStyle(pterm.FgRed).Sprint("Exiting..."))
	pterm.DefaultCenter.Print(pterm.NewStyle(pterm.BgGreen, pterm.FgCyan).Sprint("Bye! Bye! :( "))
	os.Exit(0)
}

func Input(s string) string {
	result, _ := pterm.DefaultInteractiveTextInput.
		WithMultiLine(false).
		WithDefaultText(s).
		WithOnInterruptFunc(Exit).
		WithTextStyle(pterm.NewStyle(pterm.FgLightCyan)).
		Show()

	result = strings.TrimSpace(result)
	result = strings.ReplaceAll(result, " ", "")

	return result
}
