package main

import (
	"fmt"

	"github.com/erazoex/proyecto2/analizador"
)

func main() {
	var analizador analizador.Analyzer
	running := true
	for running {
		var option string
		fmt.Printf("\n")
		fmt.Printf("%s", "Ingrese un nuevo comando: ")
		fmt.Scanln(&option)
		if option == "exit" {
			running = false
		} else {
			fmt.Printf("imprimiendo la opcion %v\n", option)
			analizador.Exe(option)
		}
	}
}

// mkdisk >size=2 >unit=k >path=/misdiscos/disco3.eaa
