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

	df := NewChannel("dreadyfire")

	cb := NewChannel("combobot")

	/*kal := &Channel{
		Name:  "kalbuir_defiancecentral",
		Users: make(map[string]User),
	}*/
	// "kalbuir_defiancecentral": kal
	nugi := NewChannel("nugiyen")
	riotgamesbrazil := NewChannel("riotgamesbrazil")
	esltv_cs := NewChannel("esltv_cs")
	sacriel := NewChannel("sacriel")
	esltv_lol := NewChannel("esltv_lol")
	mojang := NewChannel("mojang")
	chu := NewChannel("chu8")


	nugi.BannedWordList["poe"] = true
	df.BannedWordList["autoreifen"] = true
	bot := &Bot{
		Channels: map[string]*Channel{"dreadyfire": df, 
		"combobot": cb /*"kalbuir_defiancecentral": kal*/, 
		"nugiyen": nugi, 
		"riotgamesbrazil": riotgamesbrazil, 
		"esltv_cs": esltv_cs, 
		"esltv_lol": esltv_lol,
		"mojang": mojang,
		"chu": chu,
		"sacriel": sacriel,
		"nesl_lol":  NewChannel("nesl_lol"),
		"WCS America":  NewChannel("WCS America"),
		"Trick2g":  NewChannel("Trick2g"),
		"nl_Kripp":  NewChannel("nl_Kripp"),
		"Gassymexican":  NewChannel("Gassymexican"),
		"GoodGuyGarry": NewChannel(" GoodGuyGarry"),},
	}

	bot.connectAll()
	//connect()
}
