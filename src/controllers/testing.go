package controllers

import (
	"const/path"
	//"fmt"
	"libraries/common"
	"models/testing"
)

type Testing struct {
}

func (this *Testing) Concurrence() {
	type concurrence struct {
		Start  int
		Offset int
		Max    int
	}
	type item struct {
		Targetdbschema  string
		Groupname       string
		ServerParameter string
		Amount          int
		Concurrence     concurrence
	}
	var config map[string]item
	common.LoadJson(path.CONFIG_PATH+"testing"+path.DS+"schema.json", &config)

	start := config["concurrence"].Concurrence.Start
	offset := config["concurrence"].Concurrence.Offset
	max := config["concurrence"].Concurrence.Max

	c := start
	for c < max {
		testing.Concurrence(config["concurrence"].Groupname, config["concurrence"].ServerParameter, config["concurrence"].Targetdbschema, config["concurrence"].Amount, c)
		c = c + offset - (c % offset)
	}
}
