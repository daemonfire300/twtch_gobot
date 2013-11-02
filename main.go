package main

import (
	"container/list"
)

var commands = list.New()

func main() {

	commands.PushBack(Command{
		Name:            "Timeout",
		permissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/timeout",
	})

	commands.PushBack(Command{
		Name:            "Ban",
		permissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/ban",
	})

	commands.PushBack(Command{
		Name:            "Unban",
		permissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/unban",
	})

	commands.PushBack(Command{
		Name:            "Slowmode ON",
		permissionLevel: 3,
		Args: map[string]string{
			"seconds": "seconds",
		},
		Keyword: "/slow",
	})

	commands.PushBack(Command{
		Name:            "Slowmode OFF",
		permissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/slowoff",
	})

	commands.PushBack(Command{
		Name:            "Subscribers ONLY",
		permissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/subscribers",
	})

	commands.PushBack(Command{
		Name:            "Subscribers ONLY OFF",
		permissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/subscribersoff",
	})

	commands.PushBack(Command{
		Name:            "Clear",
		permissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/clear",
	})

	commands.PushBack(Command{
		Name:            "Assign Mod",
		permissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/mod",
	})

	commands.PushBack(Command{
		Name:            "Remove Mod",
		permissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/unmod",
	})

	commands.PushBack(Command{
		Name:            "ModList",
		permissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/mods",
	})

	df := &Channel{
		Name:            "dreadyfire",
		Users:           make(map[string]User),
		BannedWordList:  make(map[string]bool),
		RepetitionCache: make(map[string]int),
	}

	cb := &Channel{
		Name:            "combobot",
		Users:           make(map[string]User),
		BannedWordList:  make(map[string]bool),
		RepetitionCache: make(map[string]int),
	}

	/*kal := &Channel{
		Name:  "kalbuir_defiancecentral",
		Users: make(map[string]User),
	}*/
	// "kalbuir_defiancecentral": kal
	nugi := &Channel{
		Name:            "nugiyen",
		Users:           make(map[string]User),
		BannedWordList:  make(map[string]bool),
		RepetitionCache: make(map[string]int),
	}
	riotgamesbrazil := NewChannel("riotgamesbrazil")
	esltv_cs := NewChannel("esltv_cs")
	sacriel := NewChannel("sacriel")

	nugi.BannedWordList["poe"] = true
	df.BannedWordList["autoreifen"] = true
	bot := &Bot{
		Channels: map[string]*Channel{"dreadyfire": df, "combobot": cb /*"kalbuir_defiancecentral": kal*/, "nugiyen": nugi, "riotgamesbrazil": riotgamesbrazil, "esltv_cs": esltv_cs, "sacriel": sacriel},
	}

	bot.connectAll()
	//connect()
}
