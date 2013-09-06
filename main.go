package main

func main() {
	df := &Channel{
		name: "dreadyfire",
	}

	cb := &Channel{
		name: "combobot",
	}

	bot := &Bot{
		channels: []*Channel{df, cb},
	}
	bot.connect_all()
	//connect()
}
