package helpers

import (
	"errors"
	"log"
	"math/rand"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/pkl/pklgen"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
)

func NewMessageHandler(c *ext.Context, u *ext.Update) error {
	if len(Accounts) == 0 {
		return errors.New("no accounts found")
	}

	log.Printf("New message from %s\n", u.EffectiveUser().FirstName)

	var account *types.SessionAccount

	for _, acc := range Accounts {
		if acc.AssignedTo == u.EffectiveUser().ID {
			account = acc
			break
		}
	}

	if account == nil {
		//	Select Random Account
		account = Accounts[rand.Int()%len(Accounts)]
	}

	//var account *types.SessionAccount
	//
	//for {
	//	_acc := helpers.Accounts[rand.Int()%len(helpers.Accounts)]
	//	if _acc.PhoneNumber == helpers.CONFIG.UserBot.MainPhoneNumber {
	//		continue
	//	}
	//	account = _acc
	//	break
	//}

	var proxy *types.Proxy

	if account.Proxy.Host == "" {
		log.Printf("No proxy found for account %s\n", account.PhoneNumber)
		proxy, err := GetRandomProxy()
		if err != nil {
			log.Printf("Error while getting random proxy for account %s: %s\n", account.PhoneNumber, err.Error())
			return err
		}

		err = CheckSock5Proxy(proxy)
		if err != nil {
			return err
		}

		proxy = &types.Proxy{
			Username: proxy.Username,
			Password: proxy.Password,
			Host:     proxy.Host,
			Port:     proxy.Port,
		}
	} else {
		proxy = &account.Proxy
	}

	p := &types.LoginClientParams{
		PhoneNumber:   account.PhoneNumber,
		SessionString: account.SessionString,
		SessionType:   account.SessionType,
		Proxy:         proxy,
	}

	client, stop, err := LoginClient(p)
	if err != nil {
		log.Printf("Error while logging in to account %s: %s\n", account.PhoneNumber, err.Error())
		err := RemoveAccount(CONFIG.UserBot.SessionsFilePath, account.PhoneNumber)
		if err != nil {
			return err
		}
		return err
	}

	//ctx, cancel := context.WithTimeout(context.Background(), LoginTimeout)
	//defer cancel()

	defer stop()

	//d, err := client.API().MessagesGetDialogs(ctx, nil)
	//if err != nil {
	//	log.Printf("Error while getting dialogs for account %s: %s\n", account.PhoneNumber, err.Error())
	//	return err
	//}

	//var dialog *tg.MessagesDialogsSlice
	//
	//switch v := d.(type) {
	//case *tg.MessagesDialogs: // messages.dialogs#15ba6c40
	//case *tg.MessagesDialogsSlice: // messages.dialogsSlice#71e094f3
	//	dialog = v
	//case *tg.MessagesDialogsNotModified: // messages.dialogsNotModified#f0e3e596
	//default:
	//	log.Printf("Error while getting dialogs for account %s: %s\n", account.PhoneNumber, "unknown type")
	//	return nil
	//}
	//
	//if dialog == nil {
	//	log.Printf("Error while getting dialogs for account %s: %s\n", account.PhoneNumber, "unknown type")
	//	return nil
	//}

	var group *pklgen.GroupConfig

	for _, g := range CONFIG.Groups {
		if int64(g.SourceID) == u.EffectiveChat().GetID() {
			group = g
			break
		}
	}

	if group == nil {
		return errors.New("group not found")
	}

	//joined := false

	//for _, dd := range dialog.Chats {
	//	switch v := dd.(type) {
	//	case *tg.ChatEmpty: // chatEmpty#29562865
	//	case *tg.Chat: // chat#41cbf256
	//		for _, g := range CONFIG.Groups {
	//			if g.TargetID == v.ID {
	//				joined = true
	//				break
	//			}
	//		}
	//	case *tg.ChatForbidden: // chatForbidden#6592a1a7
	//	case *tg.Channel: // channel#83259464
	//	case *tg.ChannelForbidden: // channelForbidden#17d493d5
	//	default:
	//		continue
	//	}
	//}

	//sender := message.NewSender(client.API())
	//if !joined {
	//
	//	c := client.CreateContext()
	//
	//	target := sender.Resolve(group.TargetUserName)
	//
	//	inputPeer, err := target.AsInputPeer(c)
	//
	//	if err != nil {
	//		return err
	//	}
	//
	//	ic := inputPeer.(*tg.InputPeerChannel)
	//
	//	//	Join the chat
	//	_, err = client.API().ChannelsJoinChannel(ctx, &tg.InputChannel{
	//		ChannelID:  ic.ChannelID,
	//		AccessHash: ic.AccessHash,
	//	})
	//}

	//params := types.SendMessagesIDParams{
	//	Client:         client,
	//	MessageID:      u.EffectiveMessage.Message.ID,
	//	SourceChannel:  group.SourceUserName,
	//	TargetUsername: group.TargetUserName,
	//}
	//
	////	Send the message
	//err = SendMessageByID(&params)
	//if err != nil {
	//	log.Printf("Error while sending message for account %s: %s\n", account.PhoneNumber, err.Error())
	//	return err
	//}

	params := types.SendMessageParamsByAnotherMsg{
		Client:         client,
		Message:        u.EffectiveMessage.Message,
		SourceChannel:  group.SourceUsername,
		TargetUsername: group.TargetUsername,
	}

	var x *tg.UpdateNewMessage

	x.GetMessage().(*tg.Message).GetReplyTo()

	//	Send the message
	err = SendMessageByAnotherMsg(&params)
	if err != nil {
		log.Printf("Error while sending message for account %s: %s\n", account.PhoneNumber, err.Error())
		return err
	}

	log.Printf("Message sent successfully for account %s\n", account.PhoneNumber)

	return nil
}
