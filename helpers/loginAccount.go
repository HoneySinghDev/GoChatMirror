package helpers

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/pkl/pklgen/sessiontype"
	"github.com/HoneySinghDev/tg-fake-group-conversation-userbot/types"
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"golang.org/x/time/rate"
)

var (
	MinLoginWait = 2 * time.Second
	MaxLoginWait = 5 * time.Second
)

func createProxy(params *types.Proxy) (dcs.Resolver, error) {
	sock5, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", params.Host, params.Port), &proxy.Auth{
		User:     params.Username,
		Password: params.Password,
	}, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("error creating proxy: %w", err)
	}
	dc := sock5.(proxy.ContextDialer)
	return dcs.Plain(dcs.PlainOptions{Dial: dc.DialContext}), nil
}

func createSession(params *types.LoginClientParams) (sessionMaker.SessionConstructor, error) {
	switch params.SessionType {
	case sessiontype.Native:
		return sessionMaker.StringSession(params.SessionString), nil
	case sessiontype.Telethon:
		return sessionMaker.TelethonSession(params.SessionString), nil
	case sessiontype.Tdata:
		absolutePath, err := filepath.Abs(params.SessionString)
		if err != nil {
			return nil, fmt.Errorf("error getting absolute path: %w", err)
		}
		log.Println("Reading tdata session from", absolutePath)
		accounts, err := tdesktop.Read(absolutePath, nil)
		if err != nil {
			return nil, fmt.Errorf("error reading tdata session: %w", err)
		}
		if len(accounts) == 0 {
			return nil, errors.New("no accounts found")
		}
		return sessionMaker.TdataSession(accounts[0]), nil
	case sessiontype.Pyrogram:
		return sessionMaker.PyrogramSession(params.SessionString), nil
	default:
		return nil, errors.New("invalid session type")
	}
}

func stopClient(client *gotgproto.Client) func() {
	return func() {
		if client == nil {
			return
		}

		RandomSleep(MinLoginWait, MaxLoginWait)

		ctx := client.CreateContext()

		_, err := client.API().AccountUpdateStatus(ctx, true)
		if err != nil {
			log.Printf("Error while updating status: %s\n", err.Error())
		}

		client.Stop()
	}
}

func LoginClient(params *types.LoginClientParams) (*gotgproto.Client, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), LoginTimeout)
	defer cancel()

	errChan := make(chan error, 1)
	var client *gotgproto.Client

	go func() {
		clientType := gotgproto.ClientTypePhone(params.PhoneNumber)
		var resolver dcs.Resolver
		var err error

		if params.Proxy != nil {
			if resolver, err = createProxy(params.Proxy); err != nil {
				errChan <- err
				return
			}
		}

		var sType sessionMaker.SessionConstructor
		if sType, err = createSession(params); err != nil {
			errChan <- err
			return
		}

		//// Handler of FLOOD_WAIT that will automatically retry request.
		//waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
		//	// Notifying about flood wait.
		//	log.Println("Got FLOOD_WAIT. Will retry after", wait.Duration)
		//})

		client, err = gotgproto.NewClient(
			CONFIG.UserBot.ApiId,
			CONFIG.UserBot.ApiHash,
			clientType,
			&gotgproto.ClientOpts{
				Session:          sType,
				Resolver:         resolver,
				DisableCopyright: true,
				Middlewares: []telegram.Middleware{
					// Setting up FLOOD_WAIT handler to automatically wait and retry request.
					//waiter,
					// Setting up general rate limits to less likely get flood wait errors.
					ratelimit.New(rate.Every(time.Millisecond*100), 5),
				},
			},
		)
		if err != nil {
			errChan <- err
			return
		}

		RandomSleep(MinLoginWait, MaxLoginWait)
		errChan <- nil
	}()

	select {
	case <-ctx.Done():
		return nil, nil, errors.New("operation timed out")
	case err := <-errChan:
		return client, stopClient(client), err
	}
}

func LoginNewAccount(params *types.LoginNewAccountParams) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), LoginTimeout)
	defer cancel()

	errChan := make(chan error, 1)
	done := make(chan bool)
	var sessionString string

	go func() {
		clientType := gotgproto.ClientTypePhone(params.PhoneNumber)

		var resolver dcs.Resolver
		var err error

		if params.Proxy != nil {
			resolver, err = createProxy(params.Proxy)
			if err != nil {
				errChan <- err
				return
			}
		}

		client, err := gotgproto.NewClient(
			CONFIG.UserBot.ApiId,
			CONFIG.UserBot.ApiHash,
			clientType,
			&gotgproto.ClientOpts{
				Session:          sessionMaker.SqlSession(sqlite.Open(fmt.Sprintf("sessions/session-%s", params.PhoneNumber))),
				Resolver:         resolver,
				DisableCopyright: true,
			},
		)
		if err != nil {
			errChan <- err
			done <- true
			return
		}

		session, err := client.ExportStringSession()
		if err != nil {
			errChan <- err
			done <- true
			return
		}

		log.Println("Client logged in successfully")

		sessionString = session

		errChan <- nil
		done <- true
	}()

	select {
	case <-ctx.Done():
		return "", errors.New("operation timed out")
	case <-done:
	case err := <-errChan:
		return sessionString, err
	}

	<-done
	return sessionString, nil
}
