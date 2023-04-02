package analizador

import (
	"fmt"
	"strings"
)

type Analyzer struct {
}

func (x Analyzer) exe(input string) {
	commandsAndParams := x.split_input(input)
	var command string
	var params []string
	for i, v := range commandsAndParams {
		if i == 0 {
			command = v
		} else {
			params = append(params, v)
		}
	}
	x.matchParams(command, params)
}

func (x Analyzer) matchParams(command string, params []string) {
	var param string = "hola mundo"
	if command == "execute" {
		i := 0
		for i < len(params) {
			fmt.Println(param)
		}
	}
}

func (x Analyzer) split_input(input string) []string {
	new_input := input
	fmt.Println("hola mundo")
	return strings.Split(new_input, " ")
}
