package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/erazoex/proyecto2/analizador"
	"github.com/erazoex/proyecto2/consola"
)

func main() {
	var analizador analizador.Analyzer
	running := true
	for running {
		var option string
		fmt.Printf("\n")
		fmt.Printf("%s", "Ingrese un nuevo comando: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		option = scanner.Text()
		if option == "exit" {
			running = false
		} else {
			analizador.Exe(option)
			fmt.Println(consola.GetConsole())
		}
	}
}

// execute >path=./entrada.eaa
