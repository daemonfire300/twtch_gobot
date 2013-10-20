package main

import (
	"fmt"
	"github.com/daemonfire300/go-ircevent"
	"time"
)

const (
	TWITCH_URL   = "irc.twitch.tv"
	TWITCH_PORT  = "6667"
	BOT_NAME     = "ComboBot"
	TWITCH_OAUTH = "oauth:g6svfy6x31g1wixsijfbsh1siikpgtg"
	//TWITCH_OAUTH = "oauth:mfjaw7a1hx691ldwd7gpmgkv2zw8lit"
)

var reservedUsers = []string{"jtv", "twitch", "combobot"}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type Command struct {
	Name            string
	Args            map[string]string
	Keyword         string
	permissionLevel int
}

type Channel struct {
	Name            string
	Connection      irc.Connection
	OutStream       chan string
	InStream        chan string
	permissionLevel int
}

type Bot struct {
	Channels map[string]*Channel
}

func (command *Command) Call(arg string) string {
	return fmt.Sprintf("/%s %s", command.Keyword, arg)
}

func (channel *Channel) SendMessage(message string) {
	time.Sleep(1500 * time.Millisecond)
	channel.OutStream <- fmt.Sprintf("#%s", channel.Name, message)
}

func (channel *Channel) Connect(OutStream chan string, InStream chan string) {
	channel.OutStream = OutStream
	channel.InStream = InStream

	//userList := map[string]bool{}
	//keyword := "kekse"

	channel.Connection = *irc.IRC(BOT_NAME, BOT_NAME)
	channel.Connection.Password = TWITCH_OAUTH
	err := channel.Connection.Connect(TWITCH_URL + ":" + TWITCH_PORT)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error can not connect to %s", "Twitch IRC"))
	} else {
		channel.Connection.AddCallback("001", func(e *irc.Event) {
			//fmt.Println(fmt.Sprintf("Joing %s", bot.Name))
			//bot.Connection.Join(fmt.Sprintf("#%s", channel.Name))
			time.Sleep(1000 * time.Millisecond)
			channel.Connection.Join(fmt.Sprintf("#%s", channel.Name))
			//channel.Connection.Privmsg(fmt.Sprintf("#%s", channel.Name), "Hello World")
		})
		channel.Connection.AddCallback("PRIVMSG", func(e *irc.Event) {
			/*fmt.Println(e.Nick)
			fmt.Println(bot.Connection.GetNick())*/
			channel.OutStream <- e.Message
			/*if !stringInSlice(e.Nick, reservedUsers) {
				if userList[e.Nick] {
					e.Message = strings.TrimSpace(e.Message)
					if e.Message != "" && e.Message == keyword {
						userList[e.Nick] = true
					}
				}
				channel.OutStream <- e.Message
			} else {
				//OutStream <- "Sent messages too quickly, gnark"
				channel.OutStream <- e.Message
			}*/
		})

		go channel.Connection.Loop()
	}
}

func (bot *Bot) receiveMessage(message string) {
	fmt.Println(message)
}

func (bot *Bot) fanOut(message string) {
	for _, channel := range bot.Channels {
		//go bot.Channels[i].SendMessage(message)
		go channel.Connection.Privmsg(fmt.Sprintf("#%s", channel.Name), message)
	}
}

func (bot *Bot) connectAll() {
	time.Sleep(2000 * time.Millisecond)
	OutStream := make(chan string)
	for _, channel := range bot.Channels {
		InStream := make(chan string)
		go channel.Connect(OutStream, InStream)
	}
	time.Sleep(2000 * time.Millisecond)

	for {
		message := <-OutStream
		fmt.Println(":::::: " + message)
	}
}
