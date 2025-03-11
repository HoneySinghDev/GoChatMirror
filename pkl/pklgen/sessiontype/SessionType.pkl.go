// Code generated from Pkl module `userBotConfig.pkl`. DO NOT EDIT.
package sessiontype

import (
	"encoding"
	"fmt"
)

type SessionType string

const (
	Native   SessionType = "native"
	String   SessionType = "string"
	Tdata    SessionType = "tdata"
	Telethon SessionType = "telethon"
	Pyrogram SessionType = "pyrogram"
)

// String returns the string representation of SessionType
func (rcv SessionType) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(SessionType)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for SessionType.
func (rcv *SessionType) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "native":
		*rcv = Native
	case "string":
		*rcv = String
	case "tdata":
		*rcv = Tdata
	case "telethon":
		*rcv = Telethon
	case "pyrogram":
		*rcv = Pyrogram
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid SessionType`, str)
	}
	return nil
}
