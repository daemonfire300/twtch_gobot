package main

import (
	"container/list"
)

var commands = list.New()

func main() {

	commands.PushBack(Command{
		Name:            "Timeout",
		PermissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/timeout",
	})

	commands.PushBack(Command{
		Name:            "Ban",
		PermissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/ban",
	})

	commands.PushBack(Command{
		Name:            "Unban",
		PermissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/unban",
	})

	commands.PushBack(Command{
		Name:            "Slowmode ON",
		PermissionLevel: 3,
		Args: map[string]string{
			"seconds": "seconds",
		},
		Keyword: "/slow",
	})

	commands.PushBack(Command{
		Name:            "Slowmode OFF",
		PermissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/slowoff",
	})

	commands.PushBack(Command{
		Name:            "Subscribers ONLY",
		PermissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/subscribers",
	})

	commands.PushBack(Command{
		Name:            "Subscribers ONLY OFF",
		PermissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/subscribersoff",
	})

	commands.PushBack(Command{
		Name:            "Clear",
		PermissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/clear",
	})

	commands.PushBack(Command{
		Name:            "Assign Mod",
		PermissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/mod",
	})

	commands.PushBack(Command{
		Name:            "Remove Mod",
		PermissionLevel: 3,
		Args: map[string]string{
			"username": "username",
		},
		Keyword: "/unmod",
	})

	commands.PushBack(Command{
		Name:            "ModList",
		PermissionLevel: 3,
		Args:            map[string]string{},
		Keyword:         "/mods",
	})

	/*df := NewChannel("dreadyfire")

		cb := NewChannel("combobot")*/

	/*kal := &Channel{
		Name:  "kalbuir_defiancecentral",
		Users: make(map[string]User),
	}*/
	// "kalbuir_defiancecentral": kal
	/*nugi := NewChannel("nugiyen")
	riotgamesbrazil := NewChannel("riotgamesbrazil")
	esltv_cs := NewChannel("esltv_cs")
	sacriel := NewChannel("sacriel")
	esltv_lol := NewChannel("esltv_lol")
	mojang := NewChannel("mojang")
	chu := NewChannel("chu8")*/

	bot := NewBot()
	bot.connectAll()
	//connect()
}
