package helpers

import (
	"context"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/pkl/pklgen"
	"github.com/rs/zerolog/log"
)

var CONFIG *UserBotConfig

type UserBotConfig struct {
	UserBot      *pklgen.UserBotConfig
	Settings     *pklgen.SettingsConfig
	Groups       []*pklgen.GroupConfig
	AllowedMedia *pklgen.AllowedMediaConfig
}

func DefaultServiceConfigFromEnv() *UserBotConfig {
	cfg, err := pklgen.LoadFromPath(context.Background(), "pkl/local/userBotConfig.pkl")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	return &UserBotConfig{
		UserBot:      cfg.UserBot,
		Settings:     cfg.Settings,
		Groups:       cfg.Groups,
		AllowedMedia: cfg.AllowedMedia,
	}
}
