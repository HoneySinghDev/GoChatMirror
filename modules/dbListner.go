package modules

import (
	"log"
	"math/rand"
	"time"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/helpers"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
)

type MessageCSVFormat struct {
	ID        int    `csv:"id"`
	FromUser  int    `csv:"fromUser"`
	Message   string `csv:"message"`
	ReplyTo   int    `csv:"replyTo"`
	MediaLink string `csv:"mediaLink"`
	MediaType string `csv:"mediaType"`
}

func selectAccount(lastUsedAccount map[int]string, fromUser int) *types.SessionAccount {
	if fromUser != -1 {
		if account, exists := lastUsedAccount[fromUser]; exists {
			for _, acc := range helpers.Accounts {
				if acc.PhoneNumber == account {
					return acc
				}
			}
		} else {
			account := helpers.Accounts[rand.Int()%len(helpers.Accounts)]
			lastUsedAccount[fromUser] = account.PhoneNumber
			return account
		}
	}

	usedAccounts := map[string]bool{}
	for _, account := range lastUsedAccount {
		usedAccounts[account] = true
	}

	if len(usedAccounts) == len(helpers.Accounts) {
		return helpers.Accounts[rand.Int()%len(helpers.Accounts)]
	}

	for _, account := range helpers.Accounts {
		if !usedAccounts[account.PhoneNumber] {
			return account
		}
	}

	return helpers.Accounts[rand.Int()%len(helpers.Accounts)]
}

func SourceFileListener() {
	var messages []MessageCSVFormat

	messageIDs := make(map[int]int)
	lastUsedAccount := make(map[int]string)

	if helpers.CONFIG.Settings.SourceFile == nil {
		log.Fatal("Source file not specified in config")
	}

	err := helpers.ReadCSVFile[MessageCSVFormat](*helpers.CONFIG.Settings.SourceFile, &messages)
	if err != nil {
		log.Fatalf("Error while reading source file: %s\n", err.Error())
	}

	log.Println("Loaded", len(messages), "messages")

	for _, msg := range messages {
		account := selectAccount(lastUsedAccount, msg.FromUser)
		lastUsedAccount[msg.ID] = account.PhoneNumber

		log.Printf("Selected account %s\n", account.PhoneNumber)

		var proxy *types.Proxy

		if account.Proxy.Host == "" {
			_proxy, err := helpers.GetRandomProxy()
			if err != nil {
				log.Printf("Error while getting random proxy for account %s: %s\n", account.PhoneNumber, err.Error())
				continue
			}

			err = helpers.CheckSock5Proxy(proxy)
			if err != nil {
				log.Printf("Error while checking proxy for account %s: %s\n", account.PhoneNumber, err.Error())
				continue
			}

			proxy = &types.Proxy{
				Username: _proxy.Username,
				Password: _proxy.Password,
				Host:     _proxy.Host,
				Port:     _proxy.Port,
			}
		} else {
			proxy = &account.Proxy
		}

		p := &types.LoginClientParams{
			PhoneNumber:   account.PhoneNumber,
			SessionString: account.SessionString,
			Proxy:         proxy,
			SessionType:   account.SessionType,
		}

		client, stop, err := helpers.LoginClient(p)
		if err != nil {
			log.Printf("Error while logging in to account %s: %s\n", account.PhoneNumber, err.Error())
			err := helpers.RemoveAccount(helpers.CONFIG.UserBot.SessionsFilePath, account.PhoneNumber)
			if err != nil {
				log.Printf("Error while removing account %s: %s\n", account.PhoneNumber, err.Error())
			}
			continue
		}
		account.Username = client.Self.Username

		var replyToMsgID int
		if msg.ReplyTo != 0 {
			replyToMsgID = messageIDs[msg.ReplyTo]
		}

		err = helpers.JoinGroup(&helpers.JoinGroupParams{
			Client:  client,
			Ctx:     client.CreateContext(),
			Account: account,
			Group:   helpers.CONFIG.Groups[0],
		})
		if err != nil {
			log.Printf("Error while joining group from account %s: %s\n", account.Username, err.Error())
			stop()
			continue
		}

		helpers.RandomSleep(2, 8)

		params := types.SendMessageParams{
			Client:       client,
			UserName:     helpers.CONFIG.Groups[0].TargetUsername,
			Text:         msg.Message,
			ReplyToMsgID: replyToMsgID,
			MediaLink:    msg.MediaLink,
			MediaType:    types.GetMediaSupportType(msg.MediaType),
		}

		msgID, err := helpers.SendMessage(&params)
		if err != nil {
			log.Printf("Error while sending message from account %s: %s\n", account.Username, err.Error())
			stop()
			continue
		}

		msg.ID = msgID
		messageIDs[msg.ID] = msgID

		stop()

		MinRandTime := helpers.CONFIG.Settings.MinintervalTime.GoDuration() + time.Duration(rand.Intn(10))
		MaxRandTime := helpers.CONFIG.Settings.MaxintervalTime.GoDuration() + time.Duration(rand.Intn(10))

		log.Printf("Message sent from account %s, waiting %s to %s\n", account.Username, MinRandTime, MaxRandTime)

		helpers.RandomSleep(MinRandTime, MaxRandTime)
	}

	log.Println("All messages sent")
	log.Println("Exiting")
}
