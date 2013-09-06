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
	Channels []*Channel
}

func (command *Command) Call(arg string) string {
	return fmt.Sprintf("/%s %s", command.Keyword, arg)
}

func (channel *Channel) SendMessage(message string) {
	channel.Connection.Privmsg(fmt.Sprintf("#%s", channel.Name), message)
}

func (channel *Channel) Connect(OutStream chan string, InStream chan string) {
	channel.OutStream = OutStream
	channel.InStream = InStream
	channel.Connection = *irc.IRC(BOT_NAME, BOT_NAME)
	channel.Connection.Password = TWITCH_OAUTH
	err := channel.Connection.Connect(TWITCH_URL + ":" + TWITCH_PORT)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error can not connect to %s", channel.Name))
	} else {
		channel.Connection.AddCallback("001", func(e *irc.Event) {
			//fmt.Println(fmt.Sprintf("Joing %s", channel.Name))
			channel.Connection.Join(fmt.Sprintf("#%s", channel.Name))
			time.Sleep(1000 * time.Millisecond)
			channel.Connection.Privmsg(fmt.Sprintf("#%s", channel.Name), "Hello World")
		})
		channel.Connection.AddCallback("PRIVMSG", func(e *irc.Event) {
			/*fmt.Println(e.Nick)
			fmt.Println(channel.Connection.GetNick())*/
			if !stringInSlice(e.Nick, reservedUsers) {
				channel.Connection.Privmsg(fmt.Sprintf("#%s", channel.Name), fmt.Sprintf("echo %s", e.Message))
				//fmt.Println("Echo")
				OutStream <- e.Message
			} else {
				OutStream <- "No echo plx"
			}
		})

		go channel.Connection.Loop()
		for {
			message := <-InStream
			time.Sleep(900 * time.Millisecond)
			channel.SendMessage(message)
		}
	}
}

func (bot *Bot) receiveMessage(message string) {
	fmt.Println(message)
}

func (bot *Bot) fanOut(message string) {
	for i := range bot.Channels {
		//go bot.Channels[i].SendMessage(message)
		bot.Channels[i].InStream <- message
	}
}

func (bot *Bot) connectAll() {
	for i := range bot.Channels {
		inStream := make(chan string)
		outStream := make(chan string)

		go bot.Channels[i].Connect(outStream, inStream)

		fmt.Println("______________________________________________________________")
		fmt.Println(i)
		fmt.Println("______________________________________________________________")
	}
	time.Sleep(2000 * time.Millisecond)
	bot.fanOut("I am from thaa bot's OutStream")
	for {
		for i := range bot.Channels {
			message := <-bot.Channels[i].OutStream
			bot.receiveMessage(message)
		}
	}
}
