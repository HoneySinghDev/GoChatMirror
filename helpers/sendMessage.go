package helpers

import (
	"log"
	"time"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

var (
	TypingActionMinWait = 500 * time.Millisecond
	TypingActionMaxWait = 15 * time.Second
)

func SendMessage(params *types.SendMessageParams) (int, error) {

	u := uploader.NewUploader(params.Client.API())
	sender := message.NewSender(params.Client.API()).WithUploader(u)

	context := params.Client.CreateContext()

	target := sender.Resolve(params.UserName)

	if _, err := target.AsInputPeer(context); err != nil {
		return 0, err
	}

	if params.ReplyToMsgID != 0 {
		target.Reply(params.ReplyToMsgID)
	}

	if CONFIG.Settings.TypingAction {
		typingAction := target.TypingAction()

		_ = typingAction.Typing(context)

		dynamicMaxWait := time.Duration(len(params.Text)) * 100 * time.Millisecond
		if dynamicMaxWait > TypingActionMaxWait {
			dynamicMaxWait = TypingActionMaxWait
		}

		RandomSleep(TypingActionMinWait, dynamicMaxWait)

		_ = typingAction.Cancel(context)
	}

	var _msg tg.UpdatesClass
	var err error

	switch params.MediaType {
	case types.MediaSupportTypeNone:
		_msg, err = target.Text(context, params.Text)
	case types.MediaSupportTypePhoto:
		_msg, err = target.PhotoExternal(context, params.MediaLink, styling.Plain(params.Text))
	case types.MediaSupportTypeVideo:
		if file, err2 := u.FromURL(context, params.MediaLink); err2 == nil {
			_msg, err = target.Video(context, file, styling.Plain(params.Text))
		}
	case types.MediaSupportTypeAudio:
		if file, err2 := u.FromURL(context, params.MediaLink); err2 == nil {
			_msg, err = target.Audio(context, file, styling.Plain(params.Text))
		}
	}

	if err != nil {
		log.Println("Error while sending message:", err)
		return 0, err
	}

	if msg, ok := _msg.(*tg.UpdateShortSentMessage); ok {
		return msg.ID, nil
	}

	return 0, nil
}

func SendMessageByID(params *types.SendMessagesIDParams) error {
	ctx := params.Client.CreateContext()
	sender := message.NewSender(params.Client.API())

	messages, err := params.Client.API().MessagesGetMessages(ctx, []tg.InputMessageClass{
		&tg.InputMessageID{
			ID: params.MessageID,
		},
	})

	if err != nil {
		return err
	}

	msges := messages.(*tg.MessagesMessages)

	context := params.Client.CreateContext()

	target := sender.Resolve(params.TargetUsername)

	_, err = target.AsInputPeer(context)

	if err != nil {
		return err
	}

	for _, msg := range msges.Messages {
		m := msg.(*tg.Message)
		m.GetMedia()
		_, err = target.ForwardMessages(nil, msges.Messages[0]).Send(context)
	}

	if err != nil {
		return err
	}

	return nil
}

func SendMessageByAnotherMsg(p *types.SendMessageParamsByAnotherMsg) error {

	//entities, ok := p.Message.GetEntities()
	//if !ok {
	//	log.Println("Error while getting entities from reply message")
	//	return nil
	//}

	sender := message.NewSender(p.Client.API())

	context := p.Client.CreateContext()

	target := sender.Resolve(p.TargetUsername)

	_, err := target.AsInputPeer(context)

	if err != nil {
		return err
	}

	_, err = target.ForwardMessages(nil, p.Message).Send(context)
	if err != nil {
		return err
	}

	return nil
}
