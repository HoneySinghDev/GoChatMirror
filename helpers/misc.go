package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/pkl/pklgen"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/ext"
	"github.com/gocarina/gocsv"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"golang.org/x/net/proxy"
)

func LoadProxies() error {
	// Check if proxy file exists
	proxyFile := CONFIG.Settings.ProxyFilePath
	if _, err := os.Stat(proxyFile); err != nil {
		errMsg := fmt.Sprintf("Proxies file not found: %s", proxyFile)
		return errors.New(errMsg)
	}
	// Read proxy file
	proxyContent, err := os.ReadFile(proxyFile)
	if err != nil {
		return errors.New("error reading proxies file")
	}

	// Parse proxies file
	lines := strings.Split(string(proxyContent), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		proxyData := strings.Split(line, ",")
		if len(proxyData) != 4 {
			continue
		}
		port, err := strconv.Atoi(proxyData[1])
		if err != nil {
			continue
		}
		p := &types.Proxy{
			Host:     proxyData[0],
			Port:     port,
			Username: proxyData[2],
			Password: proxyData[3],
		}
		Proxies = append(Proxies, p)
	}

	if len(Proxies) == 0 {
		return errors.New("no proxies found")
	}

	log.Println("Loaded Proxies - ", len(Proxies))

	return nil
}

// GetRandomProxy Get Random Proxy
func GetRandomProxy() (*types.Proxy, error) {
	if len(Proxies) == 0 {
		err := LoadProxies()
		if err != nil {
			return nil, err
		}
		if len(Proxies) == 0 {
			return nil, errors.New("no proxies loaded")
		}
	}

	for i := 0; i < 5; i++ { // maximum of 5 attempts
		p := Proxies[rand.Int()%len(Proxies)]
		if err := CheckSock5Proxy(p); err == nil {
			return p, nil
		}
	}

	return nil, errors.New("all proxies are not working")
}

func CheckSock5Proxy(p *types.Proxy) error {
	// Create a dialer with SOCKS5 proxy
	dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port), &proxy.Auth{
		User:     p.Username,
		Password: p.Password,
	}, proxy.Direct)
	if err != nil {
		return fmt.Errorf("failed to create SOCKS5 proxy: %v", err)
	}

	// Dial a TCP connection using the proxy
	conn, err := dialer.Dial("tcp", "google.com:80")
	if err != nil {
		return fmt.Errorf("failed to connect to proxy: %v", err)
	}
	defer conn.Close()

	return nil
}

func LoadSessionFile(sFile string) ([]*types.SessionAccount, error) {
	var accounts []*types.SessionAccount

	// Check if the session file exists
	if _, err := os.Stat(sFile); err != nil {
		log.Printf("No session file found: %s Creating Empty", sFile)
		return nil, nil
	}

	// Read session file
	sessionContent, err := os.ReadFile(sFile)
	if err != nil {
		return nil, errors.New("error reading session file")
	}

	// Parse session file (json)
	err = json.Unmarshal(sessionContent, &accounts)
	if err != nil {
		return nil, fmt.Errorf("error parsing session file: %v", err)
	}

	return accounts, nil
}

type JoinGroupParams struct {
	Client  *gotgproto.Client
	Ctx     *ext.Context
	Account *types.SessionAccount
	Group   *pklgen.GroupConfig
}

// Cache structure to hold joined group information for each client
var groupJoinCache = make(map[string]map[int64]bool)
var cacheMutex sync.Mutex

func isGroupJoined(phoneNumber string, groupID int64) bool {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if groups, ok := groupJoinCache[phoneNumber]; ok {
		return groups[groupID]
	}
	return false
}

func setGroupJoined(phoneNumber string, groupID int64) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if _, ok := groupJoinCache[phoneNumber]; !ok {
		groupJoinCache[phoneNumber] = make(map[int64]bool)
	}
	groupJoinCache[phoneNumber][groupID] = true
}

func JoinGroup(p *JoinGroupParams) error {
	// Check if the group is already joined using the cache
	if isGroupJoined(p.Account.PhoneNumber, int64(p.Group.TargetID)) {
		log.Printf("Account %s is already in the group %0.f\n", p.Account.PhoneNumber, p.Group.TargetID)
		return nil
	}

	sender := message.NewSender(p.Client.API())

	joined := false

	d, err := p.Client.API().MessagesGetDialogs(p.Ctx, nil)
	if err != nil {
		if !strings.Contains(err.Error(), "messages.getDialogs#a0f4cb4f as nil") {
			log.Printf("Error while getting dialogs for account %s: %v\n", p.Account.PhoneNumber, err)
			return err
		} else {
			c := p.Client.CreateContext()

			target := sender.Resolve(p.Group.TargetUsername)

			inputPeer, err := target.AsInputPeer(c)
			if err != nil {
				return err
			}

			ic := inputPeer.(*tg.InputPeerChannel)

			// Join the chat
			_, err = p.Client.API().ChannelsJoinChannel(p.Ctx, &tg.InputChannel{
				ChannelID:  ic.ChannelID,
				AccessHash: ic.AccessHash,
			})
			if err != nil {
				return err
			}

			setGroupJoined(p.Account.PhoneNumber, int64(p.Group.TargetID))
			return nil
		}
	}

	var dialog *tg.MessagesDialogsSlice

	switch v := d.(type) {
	case *tg.MessagesDialogs: // messages.dialogs#15ba6c40
	case *tg.MessagesDialogsSlice: // messages.dialogsSlice#71e094f3
		dialog = v
	case *tg.MessagesDialogsNotModified: // messages.dialogsNotModified#f0e3e596
	default:
		log.Printf("Error!!! getting dialogs for account %s: %s\n", p.Account.PhoneNumber, "unknown type")
		return nil
	}

	if dialog == nil {
		log.Printf("Error! getting dialogs for account dialogs null %s: %s\n", p.Account.PhoneNumber, "unknown type")
		return nil
	}

	for _, dd := range dialog.Chats {
		switch v := dd.(type) {
		case *tg.ChatEmpty: // chatEmpty#29562865
		case *tg.Chat: // chat#41cbf256
			if int64(p.Group.TargetID) == v.ID {
				joined = true
				break
			}
		case *tg.ChatForbidden: // chatForbidden#6592a1a7
		case *tg.Channel: // channel#83259464
			if int64(p.Group.TargetID) == v.ID {
				joined = true
				break
			}
		case *tg.ChannelForbidden: // channelForbidden#17d493d5
		default:
			continue
		}
	}

	if !joined {
		c := p.Client.CreateContext()

		target := sender.Resolve(p.Group.TargetUsername)

		inputPeer, err := target.AsInputPeer(c)
		if err != nil {
			return err
		}

		ic := inputPeer.(*tg.InputPeerChannel)

		// Join the chat
		_, err = p.Client.API().ChannelsJoinChannel(p.Ctx, &tg.InputChannel{
			ChannelID:  ic.ChannelID,
			AccessHash: ic.AccessHash,
		})
		if err != nil {
			return err
		}

		setGroupJoined(p.Account.PhoneNumber, int64(p.Group.TargetID))
	}

	return nil
}

// ReadCSVFile reads a CSV file into a slice of structs
func ReadCSVFile[T any](path string, out *[]T) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer file.Close()

	if err := gocsv.UnmarshalFile(file, out); err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	if len(*out) == 0 {
		return fmt.Errorf("source file is empty")
	}

	return nil
}

// SaveCSVFile writes a slice of structs to a CSV file
func SaveCSVFile[T any](path string, data *[]T, overwrite ...bool) error {
	var file *os.File
	var err error

	if len(overwrite) > 0 && overwrite[0] {
		// Overwrite the existing file
		file, err = os.Create(path)
	} else {
		// Append to the existing file
		file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(data, file); err != nil {
		return fmt.Errorf("failed to write to source file: %v", err)
	}

	return nil
}

// RandomSleep randomly wait between x and y second function
func RandomSleep(min, max time.Duration) {
	rand.NewSource(time.Now().UnixNano()) // Properly seed the random number generator

	r := rand.Int63n(int64(max-min)) + int64(min)

	time.Sleep(time.Duration(r))
}
