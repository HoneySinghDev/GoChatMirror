# GoTGChatMirror Golang Telegram Tool

## About
This tool allows you to create fake conversations in Telegram groups using real group or text file messages. It creates the illusion of an active group by using specified accounts to mimic conversation. This can be useful for developing, testing, or any other activity that requires bulk, random messages in a group setting.

## Features

- Create fake conversations from real group messages or text file messages.
- Add multiple accounts to simulate a conversation and make a group look active.
- Automatically copy and send new messages from other groups.
- Load up a text file with messages to simulate a conversation.
- Uses a config.yaml for easy setup and control.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

This project is written in Golang. Ensure you have Go installed on your machine.

### Installation

1. Clone this repository.
```
git clone https://github.com/jaskaur18/golang-telegram-tool.git
```

2. Navigate to the project directory.
```
cd golang-telegram-tool
```

3. Install the required dependencies.
```
go get
```

### Configuration

Update the `config.yaml` file with the details of the accounts that will be used to generate messages.

Example config.yaml:

```yaml
apiID: 123456
apiHash: "your_api_hash"
mainPhoneNumber: "your_phone_number"
mainSessionString: "your_session_string"
```

Please note that the `apiID,apiHash,mainPhoneNumber,mainSessionString` is a placeholder here. Replace it with your actual session string.

## Running the Tool

After you've completed the configuration, you can now run the tool with the command:

```bash
go run main.go
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgements

- The creators of Go programming language for providing such a powerful tool.
- The creators of the Telegram API for their robust and flexible platform.

## Disclaimer

This tool is for development and testing purposes only. Any misuse of this tool will not be the responsibility of the developers.