package helpers_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/helpers"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/pkl/pklgen/sessiontype"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
	_ "github.com/flashlabs/rootpath"
	"github.com/gotd/td/tg"
)

func TestLoginAccount(t *testing.T) {
	helpers.CONFIG = helpers.DefaultServiceConfigFromEnv()

	p := &types.LoginClientParams{
		PhoneNumber:   "+919917772666",
		SessionString: "eyJWZXJzaW9uIjoxLCJEYXRhIjpudWxsfQo=",
		SessionType:   sessiontype.Telethon,
		Proxy:         nil,
	}

	c, stopClient, err := helpers.LoginClient(p)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	user, err := c.Client.Self(context.Background())
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	t.Logf("Logged in as %+v", user)

	api := c.API()
	context := c.CreateContext()

	resolved, err := api.ContactsResolveUsername(context, "catsgang_bot")
	if err != nil {
		t.Errorf("in ContactsResolveUsername err: %v", err)
		return
	}
	tgUser, ok := resolved.Users[0].AsNotEmpty()
	if !ok {
		t.Errorf("in ContactsResolveUsername err: %v", err)
		return
	}

	log.Printf("Resolved: %#v", resolved)

	resWebView, err := api.MessagesRequestWebView(context, &tg.MessagesRequestWebViewRequest{
		Peer: &tg.InputPeerUser{
			UserID:     tgUser.ID,
			AccessHash: tgUser.AccessHash,
		},
		Bot: &tg.InputUser{
			UserID:     tgUser.ID,
			AccessHash: tgUser.AccessHash,
		},
		Platform:    "android",
		FromBotMenu: false,
		URL:         "https://cats-frontend.tgapps.store",
	})
	if err != nil {
		t.Errorf("in MessagesRequestWebView err: %v", err)
		return
	}

	fmt.Println("URL:", resWebView.GetURL())

	if err != nil {
		t.Errorf("in MessagesRequestWebView err: %v", err)
		return
	}

	stopClient()

	t.Logf("Logged out")
}
