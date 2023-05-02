package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/erazoex/proyecto2/analizador"
)

func main() {
	var analizador analizador.Analyzer
	running := true
	fmt.Println("Universidad de San Carlos De Guatemala")
	fmt.Println("Facultad de Ingenieria")
	fmt.Println("Escuela de ciencias y sistemas")
	fmt.Println("Seccion A-")
	fmt.Println("Brian Josue Erazo Sagastume")
	fmt.Println("201807253")
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
			// fmt.Println(consola.GetConsole())
		}
	}
	// srv := server.New("8080")
	// err := srv.ListenAndServe()
	// if err != nil {
	// 	panic(err)
	// }
}

// execute >path=./entrada.eaa
