package main

import (
	"fmt"
	"github.com/daemonfire300/go-ircevent"
)

const (
	TWITCH_URL   = "irc.twitch.tv"
	TWITCH_PORT  = "6667"
	BOT_NAME     = "ComboBot"
	TWITCH_OAUTH = "oauth:g6svfy6x31g1wixsijfbsh1siikpgtg"
)

var reserved_users = []string{"jtv", "twitch", "combobot"}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type Channel struct {
	name       string
	connection irc.Connection
	out_stream chan string
}

type Bot struct {
	channels  []*Channel
	in_stream chan string
}

func (channel *Channel) connect(out_stream chan string) {
	channel.out_stream = out_stream
	channel.connection = *irc.IRC(BOT_NAME, BOT_NAME)
	channel.connection.Password = TWITCH_OAUTH
	err := channel.connection.Connect(TWITCH_URL + ":" + TWITCH_PORT)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error can not connect to %s", channel.name))
	} else {
		channel.connection.AddCallback("001", func(e *irc.Event) {
			//fmt.Println(fmt.Sprintf("Joing %s", channel.name))
			channel.connection.Join(fmt.Sprintf("#%s", channel.name))
			channel.connection.Privmsg(fmt.Sprintf("#%s", channel.name), "Hello World")
		})
		channel.connection.AddCallback("PRIVMSG", func(e *irc.Event) {
			/*fmt.Println(e.Nick)
			fmt.Println(channel.connection.GetNick())*/
			if !stringInSlice(e.Nick, reserved_users) {
				channel.connection.Privmsg(fmt.Sprintf("#%s", channel.name), fmt.Sprintf("echo %s", e.Message))
				//fmt.Println("Echo")
				out_stream <- e.Message
			} else {
				out_stream <- "No echo plx"
			}
		})

		channel.connection.Loop()
	}
}

func (bot *Bot) receive_message(message string) {
	fmt.Println(message)
}

func (bot *Bot) connect_all() {
	bot.in_stream = make(chan string)

	for i := range bot.channels {
		go bot.channels[i].connect(bot.in_stream)
		fmt.Println("______________________________________________________________")
		fmt.Println(i)
		fmt.Println("______________________________________________________________")
	}

	for {
		message := <-bot.in_stream
		bot.receive_message(message)
	}
}
