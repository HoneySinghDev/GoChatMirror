package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
)

var LoginTimeout = 1 * time.Minute

var Proxies []*types.Proxy

var Accounts []*types.SessionAccount

func LoadAccounts() {
	log.Printf("Loading accounts...%s", CONFIG.UserBot.SessionsFilePath)
	a, err := LoadSessionFile(CONFIG.UserBot.SessionsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Loaded", len(a), "accounts")

	Accounts = a
}

// AddAccount function
func AddAccount(sFile string, newAccount *types.SessionAccount) error {
	// Check if the account already exists
	for _, account := range Accounts {
		if account.PhoneNumber == newAccount.PhoneNumber {
			return fmt.Errorf("account already exists: %s", newAccount.PhoneNumber)
		}
	}

	// Add the new account
	Accounts = append(Accounts, newAccount)

	// Convert the accounts to JSON
	accountData, err := json.Marshal(Accounts)
	if err != nil {
		return fmt.Errorf("error marshalling accounts to JSON: %v", err)
	}

	// Write the new JSON data to the file
	err = os.WriteFile(sFile, accountData, 0644)
	if err != nil {
		return fmt.Errorf("error writing new account data to file: %v", err)
	}

	return nil
}

// RemoveAccount removes an account from a slice of *types.SessionAccount
func RemoveAccount(sFile string, PhoneNumber string) error {
	// Find the index of the account with the specified ID
	index := -1
	for i, account := range Accounts {
		if account.PhoneNumber == PhoneNumber {
			index = i
			break
		}
	}

	if index == -1 {
		// Account isn't found, return the original slice
		return fmt.Errorf("account not found: %s", PhoneNumber)
	}

	// Swap the account to remove with the last account in the slice
	Accounts[index] = Accounts[len(Accounts)-1]

	// Reduce the length of the slice by one
	Accounts = Accounts[:len(Accounts)-1]

	// Convert the accounts to JSON
	accountData, err := json.Marshal(Accounts)
	if err != nil {
		return fmt.Errorf("error marshalling accounts to JSON: %v", err)
	}

	// Write the new JSON data to the file
	err = os.WriteFile(sFile, accountData, 0644)
	if err != nil {
		return fmt.Errorf("error writing updated account data to file: %v", err)
	}

	return nil
}
