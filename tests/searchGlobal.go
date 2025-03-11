package main

import (
	"context"
	"log"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/helpers"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/pkl/pklgen/sessiontype"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
	"github.com/gotd/td/tg"
)

func main() {
	helpers.CONFIG = helpers.DefaultServiceConfigFromEnv()

	p := types.LoginAccountParams{
		PhoneNumber:   "919917772666",
		Proxy:         nil,
		SessionString: "eyJWZXJzaW9uIjoxLCJEYXRhIjpudWxsfQo=",
		SessionType:   sessiontype.Native,
	}

	client, err := helpers.LoginAccount(&p)
	if err != nil {
		panic(err)
	}

	peers, err := client.API().MessagesSearchGlobal(context.Background(), &tg.MessagesSearchGlobalRequest{
		Q:      "Protein",
		Filter: &tg.InputMessagesFilterEmpty{},
	})
	if err != nil {
		panic(err)
	}

	for _, chats := range peers.(*tg.MessagesMessages).Chats {
		log.Printf("Chat: %s", chats.String())
	}
}
