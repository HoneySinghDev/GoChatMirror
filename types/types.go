package types

import (
	"github.com/celestix/gotgproto"
	"github.com/gotd/td/tg"
)

type MediaSupportType string

const (
	MediaSupportTypePhoto MediaSupportType = "photo"
	MediaSupportTypeVideo MediaSupportType = "video"
	MediaSupportTypeAudio MediaSupportType = "audio"
	MediaSupportTypeVoice MediaSupportType = "voice"
	MediaSupportTypeFile  MediaSupportType = "file"
	MediaSupportTypeNone  MediaSupportType = ""
)

func GetMediaSupportType(s string) MediaSupportType {
	switch s {
	case "photo":
		return MediaSupportTypePhoto
	case "video":
		return MediaSupportTypeVideo
	case "audio":
		return MediaSupportTypeAudio
	case "voice":
		return MediaSupportTypeVoice
	case "file":
		return MediaSupportTypeFile
	default:
		return MediaSupportTypeNone
	}
}

type SendMessageParams struct {
	Client       *gotgproto.Client
	UserName     string
	Text         string
	ReplyToMsgID int
	MediaLink    string
	MediaType    MediaSupportType
}

type SendMessagesIDParams struct {
	Client         *gotgproto.Client
	MessageID      int
	SourceChannel  string
	TargetUsername string
}

type SendMessageParamsByAnotherMsg struct {
	Client         *gotgproto.Client
	Message        *tg.Message
	SourceChannel  string
	TargetUsername string
}

type Proxy struct {
	Host     string
	Port     int
	Username string
	Password string
}
