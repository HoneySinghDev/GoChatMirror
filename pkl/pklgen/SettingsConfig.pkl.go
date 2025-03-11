// Code generated from Pkl module `userBotConfig.pkl`. DO NOT EDIT.
package pklgen

import (
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/pkl/pklgen/grouptype"
	"github.com/apple/pkl-go/pkl"
)

type SettingsConfig struct {
	SourceType grouptype.GroupType `pkl:"sourceType"`

	SourceFile *string `pkl:"sourceFile"`

	ProxyFilePath string `pkl:"proxyFilePath"`

	MinintervalTime *pkl.Duration `pkl:"MinintervalTime"`

	MaxintervalTime *pkl.Duration `pkl:"MaxintervalTime"`

	TypingAction bool `pkl:"TypingAction"`

	ReloadInterval *pkl.Duration `pkl:"ReloadInterval"`
}
