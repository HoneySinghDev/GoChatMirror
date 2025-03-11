package modules

import (
	"log"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/helpers"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
)

func StartListener() {

	var account *types.SessionAccount

	for _, acc := range helpers.Accounts {
		if acc.PhoneNumber == helpers.CONFIG.UserBot.MainPhoneNumber {
			account = acc
			break
		}
	}

	if account == nil {
		log.Fatalf("Main Account %s not found in accounts list", helpers.CONFIG.UserBot.MainPhoneNumber)
	}

	params := &types.LoginClientParams{
		PhoneNumber:   helpers.CONFIG.UserBot.MainPhoneNumber,
		SessionString: account.SessionString,
		SessionType:   account.SessionType,
	}

	log.Printf("Logging in to Main Account : %s\n", params.PhoneNumber)

	mainC, stop, err := helpers.LoginClient(params)
	if err != nil {
		log.Printf("Account %s Login Failed: %s", params.PhoneNumber, err.Error())
		return
	}

	defer stop()

	//sString, err := mainC.ExportStringSession()
	//if err != nil {
	//	log.Printf("Error while exporting session string: %s", err.Error())
	//	return
	//}
	//
	//log.Printf("Session string: \n%s\n", sString)

	log.Println("Logged in to Main Account")
	mainC.Dispatcher.AddHandler(handlers.NewMessage(filters.Message.All, helpers.NewMessageHandler))
	log.Printf("Added Message Handler to Main Account: %s\n", params.PhoneNumber)

	err = mainC.Idle()
	if err != nil {
		log.Fatalf("Error while idling main account: %s", err.Error())
	}
}
