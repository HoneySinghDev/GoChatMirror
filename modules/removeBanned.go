package modules

import (
	"log"
	"math/rand"
	"time"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/helpers"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
)

func RemoveBannedAccount(_ *helpers.UserBotConfig) error {

	for i, acc := range helpers.Accounts {

		if i%5 == 0 {
			//	Sleeping for 20-30 seconds to avoid flood
			rTime := rand.Int()*10 + 20
			log.Printf("Sleeping for %d seconds to avoid flood", rTime)
			time.Sleep(time.Duration(rTime) * time.Second)
		}

		p := &types.LoginClientParams{
			PhoneNumber:   acc.PhoneNumber,
			SessionString: acc.SessionString,
			SessionType:   acc.SessionType,
			Proxy:         &acc.Proxy,
		}

		_, stop, err := helpers.LoginClient(p)
		if err != nil {
			log.Printf("Account %s-%s Login Failed: %s", acc.Username, acc.PhoneNumber, err.Error())
			err := helpers.RemoveAccount(helpers.CONFIG.UserBot.SessionsFilePath, acc.PhoneNumber)
			if err != nil {
				log.Printf("Remove Account %s-%s Failed: %s", acc.Username, acc.PhoneNumber, err.Error())
			}
			log.Printf("Account %s-%s Has Been Removed From List", acc.Username, acc.PhoneNumber)
			return err
		}

		time.Sleep(2 * time.Second)

		stop()

		log.Println("Account Is Working ", acc.PhoneNumber)
		log.Printf("Removing Account %s-%s From List", acc.Username, acc.PhoneNumber)
		time.Sleep(2 * time.Second)
	}

	log.Printf("All Accounts Have Been Checked And Banned Accounts Have Been Removed")

	return nil
}
