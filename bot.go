package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/daemonfire300/go-ircevent"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
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

func (user *User) isMod() bool {
	if user.PermissionLevel > 1 {
		return true
	} else {
		return false
	}
}

func (user *User) isAdmin() bool {
	if user.PermissionLevel > 2 {
		return true
	} else {
		return false
	}
}

func (command *Command) Call(arg string) string {
	return fmt.Sprintf("/%s %s", command.Keyword, arg)
}

func NewHook(callback func(*Event), typ string, prio int) *Hook {
	return &Hook{
		Callback: callback,
		Type:     typ,
		Priority: prio,
	}
}

func NewEvent(channel *Channel, information string, e *irc.Event) *Event {
	return &Event{
		Channel:     channel,
		Information: information,
		Message:     e.Message,
		User:        channel.Users[e.Nick],
	}
}

func NewChannel(id int64, name string, db *sql.DB) *Channel {
	channel := &Channel{
		Id:              id,
		Name:            name,
		Activated:       false,
		Users:           make(map[string]User),
		BannedWordList:  make(map[string]bool),
		RepetitionCache: make(map[string]map[string]int),
		Hooks:           make([]*Hook, 10),
		PollCache:       0,
		Database:        db,
	}
	pollCallback := func(e *Event) {
		if strings.HasPrefix(strings.TrimSpace(strings.ToLower(e.Message)), "!poll") && e.User.isAdmin() {
			msg := strings.Split(strings.TrimSpace(strings.ToLower(e.Message)), " ")
			if len(msg) > 1 {
				name := msg[1]
				options := msg[2:]
				fmt.Println("Starting poll on  ", channel.Name, options)
				channel.startPoll(600, name, options)
			}
		}
		if strings.HasPrefix(strings.TrimSpace(strings.ToLower(e.Message)), "!endpoll") && e.User.isAdmin() {
			channel.stopPoll()
		}
	}

	prefTextCallback := func(e *Event) {
		if strings.HasPrefix(strings.TrimSpace(strings.ToLower(e.Message)), "!text") && e.User.isMod() {
			msg := strings.Split(strings.TrimSpace(strings.ToLower(e.Message)), " ")
			if len(msg) > 1 {
				alias := msg[1]
				if len(alias) > 0 {
					row := channel.Database.QueryRow("SELECT id, text FROM predefined_text WHERE alias = $1 AND channel_id = $2", alias, channel.Id)
					var text string
					var id int64
					row.Scan(&id, &text)
					fmt.Println(fmt.Sprintf("%s: %d", text, id))
				}
			}
		}
	}

	subscribedTextCallback := func(e *Event) {
		expr := `^\w+\s+just subscribed!$`
		re, err := regexp.Compile(expr)
		if err != nil {
			log.Fatal(err)
		}
		if re.MatchString(strings.TrimSpace(strings.ToLower(e.Message))) { // && !e.User.isAdmin()
			log.Println("onNewSubscriberHook", "implement your logic here", e.Message, e.Information)
		}
	}

	channel.addHook(NewHook(pollCallback, "ManagePollOnMessage", 10))
	channel.addHook(NewHook(prefTextCallback, "DisplayTextOnMessage", 5))
	channel.addHook(NewHook(subscribedTextCallback, "PerformOnNewSubscriber", 10))
	return channel
}

func NewBot() *Bot {
	return &Bot{
		Channels: make(map[string]*Channel),
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

func (channel *Channel) addHook(hook *Hook) {
	channel.Hooks = append(channel.Hooks, hook)
}

func (channel *Channel) flushRepetitionCache() {
	channel.RepetitionCache = make(map[string]map[string]int)
}

func (channel *Channel) flushRepetitionCacheSpecific(username string) {
	delete(channel.RepetitionCache, username)
}

func (channel *Channel) detectRepetition(e *Event) bool {
	/*if len(e.Message) < 10{
		return false
	}*/
	_, ok := channel.Users[e.User.Name]
	if ok {
		h := md5.New()
		io.WriteString(h, e.Message)
		key := fmt.Sprintf("%x", h.Sum(nil)) + e.User.Name
		channel.RepetitionCache[e.User.Name][key]++
		v := channel.RepetitionCache[e.User.Name][key]
		size_u := len(channel.RepetitionCache[e.User.Name])
		size_all := len(channel.RepetitionCache)
		fmt.Println(fmt.Sprintf("UserCache: %d \t \t ChannelCache: %d", size_u, size_all))
		if v > 3 {
			fmt.Println("SPAM Pattern detected.... (simple)")
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (channel *Channel) startPoll(duration int, name string, options []string) {
	callback := func(e *Event) {
		for _, option := range options {
			if strings.HasPrefix(strings.ToLower(e.Message), option) {
				channel.PollCache++
				fmt.Println("PollCount +1 ", option)
			}
		}
	}
	typ := "PollOnMessage"
	prio := 10
	channel.addHook(NewHook(callback, typ, prio))
}

func (channel *Channel) stopPoll() {
	for i, hook := range channel.Hooks {
		if hook != nil {
			if hook.Type == "PollOnMessage" {
				log.Println(fmt.Sprintf("Closing Poll: Result <%d>", channel.PollCache))
				channel.Hooks = append(channel.Hooks[:i], channel.Hooks[i+1:]...)
				channel.PollCache = 0
			}
		}
	}
}

/*func (channel *Channel) detectPattern(e *Event, pattern string) {
}*/

func (channel *Channel) RcvMessage(e *irc.Event) {
	ev := NewEvent(channel, e.Raw, e)
	for _, hook := range channel.Hooks {
		if hook != nil {
			hook.Callback(ev)
		}
	}

	repetition := channel.detectRepetition(ev)
	if repetition {
		fmt.Println("Ban/Timeout User: " + e.Nick)
	}
	if channel.containsBlacklisted(e.Message) {
		fmt.Println("This message contains blacklisted words")
	}

	channel.OutStream <- Message{Channel: channel, Text: e.Message + fmt.Sprintf(" %d", ev.User.PermissionLevel)}
}

func (channel *Channel) SndMessage(message string) {
	channel.Connection.Privmsg(channel.Self(), message)
}

func (channel *Channel) RemoveUser(user string) {
	_, ok := channel.Users[user]
	if ok {
		fmt.Println("Removing User " + user)
		fmt.Println("Clearing Cache " + user)
		//channel.flushRepetitionCacheSpecific(user)
		delete(channel.Users, user)
	}
}

func (channel *Channel) AddUser(user string) {
	_, ok := channel.Users[user]
	if ok == false {
		fmt.Println("Adding User " + user)
		channel.RepetitionCache[user] = make(map[string]int)
		permissionLevel := 0
		if user == channel.Name {
			permissionLevel = 3
		}
		channel.Users[user] = User{
			Name:            user,
			PermissionLevel: permissionLevel,
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
	channel.Connection = irc.IRC(BOT_NAME, BOT_NAME)
	channel.Connection.Password = TWITCH_OAUTH
	err := channel.Connection.Connect(TWITCH_URL + ":" + TWITCH_PORT)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error can not connect to %s", "Twitch IRC"))
	} else {
		channel.Connection.AddCallback("001", func(e *irc.Event) {
			time.Sleep(1000 * time.Millisecond)
			channel.Connection.Join(channel.Self())
		})
		channel.Connection.AddCallback("PRIVMSG", func(e *irc.Event) {
			channel.RcvMessage(e)
		})

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

func (channel *Channel) ConnectBot(OutStream chan Message, InStream chan string, connection *irc.Connection) {
	channel.OutStream = OutStream
	channel.InStream = InStream
	channel.Connection = connection
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

func (bot *Bot) deactivateChannel(channelName string) {
	channel, ok := bot.Channels[channelName]
	if ok {
		fmt.Println("deactivating channel...")
		channel.Activated = false
		bot.leaveChannel(channelName)
	}
}

func (bot *Bot) leaveChannel(channelName string) {
	channel, ok := bot.Channels[channelName]
	if ok {
		fmt.Println("PARTing channel...")
		channel.Connection.Part(channel.Self())
	}
}

func (bot *Bot) activateChannel(channelName string) {
	channel, ok := bot.Channels[channelName]
	if ok {
		fmt.Println("activating channel...")
		channel.Activated = true
		bot.joinChannel(channelName)
	}
}

func (bot *Bot) joinChannel(channelName string) {
	channel, ok := bot.Channels[channelName]
	if ok {
		if channel.Activated {
			fmt.Println("JOINing channel...")
			channel.Connection.Join(channel.Self())
		}
	}
}

func (bot *Bot) addChannel(channel *Channel) {
	_, ok := bot.Channels[channel.Name]
	if ok {
		log.Printf("Channel %s already exists", channel.Name)
	} else {
		log.Printf("Adding channel %s too bot", channel.Name)
		bot.Channels[channel.Name] = channel
		InStream := make(chan string)
		channel.ConnectBot(bot.OutStream, InStream, bot.Connection)
	}
}

func (bot *Bot) ConnectDatabase() {
	db, err := sql.Open("postgres", "user=postgres dbname=gobot password=abc sslmode=disable")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("\n\n__________________________________\n\n\nConnected to Database\n\n\n__________________________________\n\n")
	}
	bot.Database = db
}

func (bot *Bot) LoadChannels() {
	if bot.Database != nil {
		err := bot.Database.Ping()
		if err != nil {
			log.Fatal(err)
		}
		rows, err := bot.Database.Query("SELECT * FROM channel")
		if err != nil {
			log.Fatal(err)
		}

		var id int64
		var name string
		var enabled bool
		var cnt int

		log.Print("Loading channels")
		for rows.Next() {
			rows.Scan(&id, &name, &enabled)
			bot.Channels[name] = NewChannel(id, name, bot.Database)
			bot.Channels[name].Activated = enabled
			if enabled {
				cnt++
			}
			log.Print(name)
		}
		log.Printf("\nLoaded %d channels", cnt)
	}
}

func (bot *Bot) httpHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// read form value
		channelName := r.FormValue("channel")
		action := r.FormValue("action")
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		if len(channelName) > 0 {
			channelName = strings.ToLower(channelName)
			switch {
			case action == "create":
				var id int64
				existsRow := bot.Database.QueryRow("SELECT id FROM channel WHERE name = $1", channelName).Scan(&id)
				if existsRow == sql.ErrNoRows {
					row := bot.Database.QueryRow("INSERT INTO channel(name) VALUES($1) RETURNING id", channelName)
					id = 0
					row.Scan(&id)
					fmt.Println(id)
					bot.addChannel(NewChannel(id, channelName, bot.Database))
				} else {
					fmt.Println("Channel already exists")
				}
			case action == "list":
				for _, user := range bot.Channels[channelName].Users {
					fmt.Fprintf(w, "* %s\n", user.Name)
				}
			case action == "PART":
				bot.deactivateChannel(channelName)
			case action == "JOIN":
				bot.activateChannel(channelName)
			default:
				fmt.Println("No action specified, doing nothing")
				fmt.Fprintf(w, "No action specified, doing nothing: %s %s", channelName, action)
			}
			//bot.addChannel(NewChannel(1337, channelName, bot.Database)) // broken!
		} else {
			switch {
			case action == "listAll":
				for _, channel := range bot.Channels {
					//fmt.Fprintf(w, "* %s  active: %t", channel.Name, channel.Activated)
					fmt.Fprintln(w, fmt.Sprintf("* %s  active: %t", channel.Name, channel.Activated))
				}
			}
		}
	}
}

func (bot *Bot) ArgsToChannel(args []string) *Channel {
	if len(args) > 0 {
		rawName := args[len(args)-1]
		name := strings.TrimSpace(strings.Replace(rawName, "#", "", -1))

		channel, ok := bot.Channels[name]
		if ok {
			return channel
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (bot *Bot) Callback_001(e *irc.Event) {
	time.Sleep(1000 * time.Millisecond)
	for _, channel := range bot.Channels {
		bot.joinChannel(channel.Name)
	}
}

func (bot *Bot) Callback_PRIVMSG(e *irc.Event) {
	channel := bot.ArgsToChannel(e.Arguments)
	if channel != nil {
		channel.RcvMessage(e)
	}
}

func (bot *Bot) Callback_JOIN(e *irc.Event) {
	channel := bot.ArgsToChannel(e.Arguments)
	//fmt.Println("triggered join")
	//fmt.Println(e.Arguments)
	if channel != nil {
		channel.AddUser(e.Nick)
	}
}

func (bot *Bot) Callback_PART(e *irc.Event) {
	channel := bot.ArgsToChannel(e.Arguments)
	if channel != nil {
		channel.RemoveUser(e.Nick)
	}
}

func (bot *Bot) Callback_353(e *irc.Event) {
	channel := bot.ArgsToChannel(e.Arguments)
	//fmt.Println("triggered history add")
	//fmt.Println(e.Arguments)
	if channel != nil {
		channel.HistoryAddUser(strings.Split(e.Message, " "))
	}
}

func (bot *Bot) SetupDefaultCallbacks() {
	bot.Connection.AddCallback("001", bot.Callback_001)
	bot.Connection.AddCallback("PRIVMSG", bot.Callback_PRIVMSG)
	bot.Connection.AddCallback("JOIN", bot.Callback_JOIN)
	bot.Connection.AddCallback("PART", bot.Callback_PART)
	bot.Connection.AddCallback("353", bot.Callback_353)
}

func (bot *Bot) ConnectIRC() {
	bot.Connection = irc.IRC(BOT_NAME, BOT_NAME)
	bot.Connection.Password = TWITCH_OAUTH
	err := bot.Connection.Connect(TWITCH_URL + ":" + TWITCH_PORT)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error can not connect to %s", "Twitch IRC"))
	} else {
		bot.SetupDefaultCallbacks()
		go bot.Connection.Loop()
	}
}

func (bot *Bot) ConnectAll() {
	bot.ConnectDatabase()
	bot.LoadChannels()
	bot.ConnectIRC()
	time.Sleep(1000 * time.Millisecond)
	OutStream := make(chan Message)
	bot.OutStream = OutStream
	for _, channel := range bot.Channels {
		InStream := make(chan string)
		//go channel.Connect(OutStream, InStream)
		channel.ConnectBot(OutStream, InStream, bot.Connection)
	}
	time.Sleep(1000 * time.Millisecond)

	log.Println("Starting WebAPI Server")

	http.HandleFunc("/", bot.httpHandler())
	go http.ListenAndServe(":8181", nil)

	log.Println("Started WebAPI Server")

	time.Sleep(1000 * time.Millisecond)

	for {
		message := <-OutStream
		fmt.Println(time.Now().String() + message.Channel.Name + " : " + message.Text)
	}
}
