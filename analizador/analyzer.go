package analizador

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/erazoex/proyecto2/comandos"
)

type Analyzer struct {
}

func (a Analyzer) Exe(input string) {
	fmt.Println(input)
	commandsAndParams := a.Split_input(input)
	fmt.Printf("aqui imprime esto: %v\n", commandsAndParams)
	var command string
	var params []string
	for i, v := range commandsAndParams {
		fmt.Println(i)
		fmt.Println(v)
		if i == 0 {
			command = v
		} else {
			params = append(params, v)
		}
	}
	a.MatchParams(command, params)
}

func (a Analyzer) MatchParams(command string, params []string) {
	command = strings.Replace(command, " ", "", 1)
	if command == "execute" {
		for _, v := range params {
			if strings.Contains(v, "path") {
				v = strings.Replace(v, "path=", "", 1)
				v = strings.ReplaceAll(v, "\"", "")
				a.Read(v)
			}
		}
	} else if command == "pause" {
		fmt.Println("")
	} else if command == "mkdisk" {
		m := comandos.Mkdisk{}
		m.Exe(params)
	} else if command == "rmdisk" {
		fmt.Println("")
	} else if command == "fdisk" {
		fmt.Println("")
	} else if command == "mount" {
		fmt.Println("")
	} else if command == "unmount" {
		fmt.Println("")
	} else if command == "mkfs" {
		fmt.Println("")
	} else if command == "login" {
		fmt.Println("")
	} else if command == "logout" {
		fmt.Println("")
	} else if command == "mkgrp" {
		fmt.Println("")
	} else if command == "rmgrp" {
		fmt.Println("")
	} else if command == "mkuser" {
		fmt.Println("")
	} else if command == "rmuser" {
		fmt.Println("")
	} else if command == "rep" {
		fmt.Println("")
	} else if strings.Contains(command, "#") {
		fmt.Println("")
	}
}

func (a Analyzer) Split_input(input string) []string {
	// fmt.Println("haciendo split al input")
	return strings.Split(input, ">")
}

func (a Analyzer) Read(path string) {
	// aqui hay que leer el archivo y ejecutarlo
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error al intentar abrir el archivo: %s\n", path)
		return
	}

	defer file.Close()

	// Crear un scanner para luego leer linea por linea el archivo de entrada
	scanner := bufio.NewScanner(file)

	// Leyendo linea por linea
	for scanner.Scan() {
		// obteniendo la linea actual
		linea := scanner.Text()
		// ejecutar la linea usando a.exe()
		a.Exe(linea)
	}

	// comprobar que no hubo error al leer el archivo
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo: ", err)
		return
	}
}
