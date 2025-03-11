// Code generated from Pkl module `userBotConfig.pkl`. DO NOT EDIT.
package grouptype

import (
	"encoding"
	"fmt"
)

type GroupType string

const (
	Group GroupType = "group"
	File  GroupType = "file"
)

// String returns the string representation of GroupType
func (rcv GroupType) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(GroupType)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for GroupType.
func (rcv *GroupType) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "group":
		*rcv = Group
	case "file":
		*rcv = File
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid GroupType`, str)
	}
	return nil
}
