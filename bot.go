package main

import (
	"crypto/md5"
	"fmt"
	"github.com/daemonfire300/go-ircevent"
	"io"
	"strings"
	"time"
)

const (
	TWITCH_URL   = "irc.twitch.tv"
	TWITCH_PORT  = "6667"
	BOT_NAME     = "combobot"
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

type User struct {
	Name            string
	permissionLevel int // 0 guest, 1 user, 2 mod, 3 admin, 4 staff
	banned          bool
}

type Command struct {
	Name            string
	Args            map[string]string
	Keyword         string
	permissionLevel int
}

type Message struct {
	Channel *Channel
	Text    string
}

type Event struct {
	Channel     *Channel
	Information string
	Message     string
	User        User
}

type Channel struct {
	Name            string
	Connection      irc.Connection
	OutStream       chan Message
	InStream        chan string
	permissionLevel int
	Users           map[string]User
	BannedWordList  map[string]bool
	Callbacks       map[string]func(*Event)
	RepetitionCache map[string]int
}

type Bot struct {
	Channels map[string]*Channel
}

func (user *User) isMod() bool {
	if user.permissionLevel > 1 {
		return true
	} else {
		return false
	}
}

func (user *User) isAdmin() bool {
	if user.permissionLevel > 2 {
		return true
	} else {
		return false
	}
}

func (command *Command) Call(arg string) string {
	return fmt.Sprintf("/%s %s", command.Keyword, arg)
}

func NewEvent(channel *Channel, information string, e *irc.Event) *Event {
	return &Event{
		Channel:     channel,
		Information: information,
		Message:     e.Message,
		User:        channel.Users[e.Nick],
	}
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:            name,
		Users:           make(map[string]User),
		BannedWordList:  make(map[string]bool),
		RepetitionCache: make(map[string]int),
	}
}

func (channel *Channel) Self() string {
	return fmt.Sprintf("#%s", channel.Name)
}

func (channel *Channel) containsBlacklisted(message string) bool {
	for _, word := range strings.Split(message, " ") {
		_, ok := channel.BannedWordList[strings.ToLower(strings.TrimSpace(word))]
		if ok {
			return true
		}
	}
	return false
}

func (channel *Channel) flushRepetitionCache() {
	channel.RepetitionCache = nil
	channel.RepetitionCache = make(map[string]int)
}

func (channel *Channel) detectRepetion(e *Event) bool {
	h := md5.New()
	io.WriteString(h, e.Message)
	key := fmt.Sprintf("%x", h.Sum(nil)) + e.User.Name
	channel.RepetitionCache[key]++
	v := channel.RepetitionCache[key]
	if v > 3 {
		fmt.Println("SPAM Pattern detected.... (simple)")
		return true
	} else {
		return false
	}
}

func (channel *Channel) RcvMessage(e *irc.Event) {
	ev := NewEvent(channel, "", e)
	repetition := channel.detectRepetion(ev)
	if repetition {
		fmt.Println("Ban/Timeout User: " + e.Nick)
	}
	if channel.containsBlacklisted(e.Message) {
		fmt.Println("This message contains blacklisted words")
	}
	channel.OutStream <- Message{Channel: channel, Text: e.Message}
}

func (channel *Channel) SndMessage(message string) {
	channel.Connection.Privmsg(channel.Self(), message)
}

func (channel *Channel) RemoveUser(user string) {
	_, ok := channel.Users[user]
	if ok {
		fmt.Println("Removing User " + user)
		delete(channel.Users, user)
	}
}

func (channel *Channel) AddUser(user string) {
	_, ok := channel.Users[user]
	if ok == false {
		fmt.Println("Adding User " + user)
		channel.Users[user] = User{
			Name: user,
		}
	}
}

func (channel *Channel) HistoryAddUser(userList []string) {
	if len(userList) > 0 {
		fmt.Println("Adding Users that are already present in the Channel")
		for _, user := range userList {
			channel.AddUser(user)
		}
	}
}

func (channel *Channel) Connect(OutStream chan Message, InStream chan string) {
	channel.OutStream = OutStream
	channel.InStream = InStream
	channel.Connection = *irc.IRC(BOT_NAME, BOT_NAME)
	channel.Connection.Password = TWITCH_OAUTH
	err := channel.Connection.Connect(TWITCH_URL + ":" + TWITCH_PORT)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error can not connect to %s", "Twitch IRC"))
	} else {
		channel.Connection.AddCallback("001", func(e *irc.Event) {
			time.Sleep(1000 * time.Millisecond)
			channel.Connection.Join(channel.Self())
		})
		channel.Connection.AddCallback("PRIVMSG", channel.RcvMessage)

		channel.Connection.AddCallback("JOIN", func(e *irc.Event) {
			channel.AddUser(e.Nick)
		})
		channel.Connection.AddCallback("PART", func(e *irc.Event) {
			channel.RemoveUser(e.Nick)
		})
		channel.Connection.AddCallback("353", func(e *irc.Event) {
			channel.HistoryAddUser(strings.Split(e.Message, " "))
		})

		go channel.Connection.Loop()
	}
}

func (bot *Bot) receiveMessage(message string) {
	fmt.Println(message)
}

func (bot *Bot) fanOut(message string) {
	for _, channel := range bot.Channels {
		go channel.SndMessage(message)
	}
}

func (bot *Bot) writeToChannel(channel string, message string) {
	bot.Channels[channel].SndMessage(message)
}

func (bot *Bot) connectAll() {
	time.Sleep(1000 * time.Millisecond)
	OutStream := make(chan Message)
	for _, channel := range bot.Channels {
		InStream := make(chan string)
		go channel.Connect(OutStream, InStream)
	}
	time.Sleep(1000 * time.Millisecond)

	for {
		message := <-OutStream
		fmt.Println(time.Now().String() + message.Channel.Name + " : " + message.Text)
	}
	for _, channel := range bot.Channels {
		fmt.Println(channel.RepetitionCache)
	}
}
