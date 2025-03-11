package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

type OldSession struct {
	Index         int    `json:"index"`
	PhoneNumber   string `json:"phoneNumber"`
	SessionString string `json:"sessionString"`
	Proxy         bool   `json:"proxy"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

type Proxy struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SessionAccount struct {
	Username      string `json:"-"`
	PhoneNumber   string `json:"phone_number"`
	SessionString string `json:"session_string"`
	SessionType   string `json:"session_type,omitempty"`
	Banned        bool   `json:"banned"`
	AssignedTo    int64  `json:"-"`
	Proxy         Proxy  `json:"proxy"`
}

func main() {
	input := textinput.New("Enter Your Phone Number To Login")
	input.Placeholder = "old_sessions.json"
	fName, err := input.RunPrompt()
	if err != nil {
		log.Fatal("Error while getting phone number from user:", err)
	}

	file, err := os.Open(fName)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var oldSessions []OldSession
	err = json.Unmarshal(bytes, &oldSessions)
	if err != nil {
		log.Fatal("Error unmarshalling old sessions:", err)
	}

	var newSessions []SessionAccount
	for _, oldSession := range oldSessions {
		newSession := SessionAccount{
			Username:      oldSession.Username,
			PhoneNumber:   oldSession.PhoneNumber,
			SessionString: oldSession.SessionString,
			Banned:        false,   // Since Banned field is not in OldSession
			AssignedTo:    0,       // Since AssignedTo field is not in OldSession
			Proxy:         Proxy{}, // Empty Proxy, as we are ignoring it
		}
		newSessions = append(newSessions, newSession)
	}

	newBytes, err := json.MarshalIndent(newSessions, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling new sessions:", err)
	}

	err = os.WriteFile("new_"+fName, newBytes, 0644)
	if err != nil {
		log.Fatal("Error writing new sessions to file:", err)
	}

	log.Println("Conversion complete. New sessions written to 'new_" + fName + "'")
}
