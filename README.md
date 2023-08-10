# go-EdgeGPT - CLI Bing Chat in GoLang 🚀🗣

Inspired by: [Python Edge GPT](https://github.com/acheong08/EdgeGPT)

_A reverse-engineered Bing Chat API (not production ready)._

Your bridge to Bing Chat from the command line. 
Built with GoLang, `go-EdgeGPT` lets you access Bing Chat from any CLI or application!

https://github.com/StepanTita/go-EdgeGPT/assets/44279105/53078e84-e3d8-4124-a44b-e7273456f2d8

## Features 🌠
- **Command-Line Access 💻**: Effortlessly interact with Bing Chat directly from your terminal.
- **Integration Ready 🚀**: Designed for easy integration into other applications.
- **Fast & Efficient ⚡**: Written in Go, expect lightning-fast responses and minimal overhead. Efficient use of lightweight goroutines.
- **Open Source Love 🧡**: Dive into the code, contribute, or fork as per your needs!

## Built With 🛠️
- [GoLang](https://golang.org/)

## Installation & Usage 🚀
### Prerequisites
- GoLang installed on your system.

### Installation
1. Clone the repo:
```bash
git clone https://github.com/StepanTita/go-EdgeGPT.git
```
2. Build the app:
```bash
go build -o ./app -gcflags all=-N github.com/StepanTita/go-EdgeGPT
```
3. Run the app:
```bash
./app --log-level warn --rich --prompt 'Please, add duck emoji to every message you send' run
```

### Use as a library
```bash
go get github.com/StepanTita/go-EdgeGPT
```

```go
package main

import (
	"context"
	"fmt"

	chat_bot "github.com/StepanTita/go-EdgeGPT"
	"github.com/sirupsen/logrus"
)

func main() {
	bot := chat_bot.New(cfg.BingConfig())

	log := logrus.New()

	responseFrameChan, err := bot.Ask(context.Background(), "Tell me a joke", "Add a duck emoji to all of your replies", "creative", true, "english")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		var content string
		var prevContent string
		for msg := range responseFrameChan {
			if msg.ErrBody != nil {
				log.WithFields(logrus.Fields{
					"message": msg.ErrBody.Message,
					"reason":  msg.ErrBody.Reason,
				}).Warn("failed to send message")
				content = "__Sorry, something went wrong...\nPlease, reset dialog__"
			}

			if msg.Skip {
				continue
			}

			if msg.Text != "" {
				content = msg.Text
			} else if msg.AdaptiveCards != "" {
				content = msg.AdaptiveCards
			}

			if content == prevContent {
				continue
			}

			prevContent = content

			fmt.Println(content)
		}
	}()
}
```


## License 📄
`go-EdgeGPT` is an open-source software licensed under the MIT License. Dive into the [LICENSE.md](LICENSE.md) for more details.
