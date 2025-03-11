package main

import (
	"fmt"
	"log"
	"os"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/modules"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/helpers"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	sl selector.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg {
	case common.DONE:
		return m, tea.Quit
	}

	_, cmd := m.sl.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.sl.View()
}

type TypeMessage struct {
	Type          string
	ZHDescription string
	ENDescription string
}

type TypeOption struct {
	Name        string
	Description string
	Usage       string
	Func        func(config *helpers.UserBotConfig) error
}

func main() {
	helpers.CONFIG = helpers.DefaultServiceConfigFromEnv()

	err := os.MkdirAll("./sessions", 0755)
	if err != nil {
		log.Fatalf("Error creating sessions directory: %s", err.Error())
	}

	helpers.LoadAccounts()

	err = helpers.LoadProxies()
	if err != nil {
		log.Fatalf("Error loading proxies: %s", err.Error())
	}

	m := &model{
		sl: selector.Model{
			Data: []interface{}{
				TypeOption{
					Name:        "Login",
					Description: "Login telegram account",
					Usage:       "Login telegram account with phone number and otp code",
					Func:        modules.LoginAccountCMD,
				},
				TypeOption{
					Name:        "Add",
					Description: "Add telegram account",
					Usage:       "Add telegram account with phone number and session string to sessions file",
					Func:        AddAccount,
				}, TypeOption{
					Name:        "Start",
					Description: "Start Listening",
					Usage:       "Start listening to the group for messages",
					Func:        StartListening,
				}, TypeOption{
					Name:        "Remove",
					Description: "Remove banned accounts",
					Usage:       "Program will check and remove banned accounts",
					Func:        modules.RemoveBannedAccount,
				}, TypeOption{
					Name:        "Exit",
					Description: "Exit this program",
					Usage:       "Program will exit",
					Func:        Exit,
				},
			},
			PerPage: 5,
			// Use the arrow keys to navigate: ↓ ↑ → ←
			// Select Commit Type:
			HeaderFunc: selector.DefaultHeaderFuncWithAppend("What you wanna to do?"),
			// [1] feat (Introducing new features)
			SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				t := obj.(TypeOption)
				return common.FontColor(fmt.Sprintf("[%d] %s (%s)", gdIndex+1, t.Name, t.Description), selector.ColorSelected)
			},
			// 2. fix (Bug fix)
			UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				t := obj.(TypeOption)
				return common.FontColor(fmt.Sprintf(" %d. %s (%s)", gdIndex+1, t.Name, t.Description), selector.ColorUnSelected)
			},

			FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				t := m.Selected().(TypeOption)
				footerTpl := `
	Module: %s
	Usage: %s`
				return common.FontColor(fmt.Sprintf(footerTpl, t.Name, t.Usage), selector.ColorFooter)
			},
			FinishedFunc: func(s interface{}) string {
				return common.FontColor("Current selected: ", selector.ColorFinished) + s.(TypeOption).Name + "\n"
			},
		},
	}

	p := tea.NewProgram(m)
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
	if !m.sl.Canceled() {
		//	log.Printf("selected index => %d\n", m.sl.Index())
		//		log.Printf("selected vaule => %s\n", m.sl.Selected())
		//	Run the selected function
		err := m.sl.Selected().(TypeOption).Func(helpers.CONFIG)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("user canceled...")
	}
}

func AddAccount(_ *helpers.UserBotConfig) error {
	modules.AddAccountCMD()
	return nil
}

func StartListening(c *helpers.UserBotConfig) error {
	if c.Settings.SourceType == "group" {
		log.Println("Starting listening to group")
		modules.StartListener()
	} else {
		log.Println("Starting listening to csv")
		modules.SourceFileListener()
	}

	return nil
}

func Exit(_ *helpers.UserBotConfig) error {
	log.Println("Exiting...")
	os.Exit(0)
	return nil
}
