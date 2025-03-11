package types

import (
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/pkl/pklgen/sessiontype"
)

type LoginClientParams struct {
	PhoneNumber   string
	SessionString string
	SessionType   sessiontype.SessionType
	Proxy         *Proxy
}

type LoginNewAccountParams struct {
	PhoneNumber string
	Proxy       *Proxy
}

type SessionAccount struct {
	Username      string                  `json:"-"`
	PhoneNumber   string                  `json:"phone_number"`
	SessionString string                  `json:"session_string"`
	SessionType   sessiontype.SessionType `json:"session_type,omitempty"`
	Banned        bool                    `json:"banned"`
	AssignedTo    int64                   `json:"-"`
	Proxy         Proxy                   `json:"proxy"`
}
