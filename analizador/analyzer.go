package analizador

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/comandos/usuariosygrupos"
	"github.com/erazoex/proyecto2/consola"
)

type Analyzer struct {
}

func (a *Analyzer) Exe(input string) {
	commandsAndParams := a.Split_input(input)
	var command string
	var params []string
	for i, v := range commandsAndParams {
		if i == 0 {
			command = v
		} else {
			params = append(params, v)
		}
	}
	a.MatchParams(command, params)
}

func (a *Analyzer) MatchParams(command string, params []string) {
	command = strings.ReplaceAll(command, " ", "")
	if command == "execute" {
		for _, v := range params {
			if strings.Contains(v, "path") {
				v = strings.Replace(v, "path=", "", 1)
				v = strings.ReplaceAll(v, "\"", "")
				a.Read(v)
			}
		}
	} else if command == "pause" {
		var option string
		consola.AddToConsole("\n")
		consola.AddToConsole("presione 'ENTER' para continuar: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		option = scanner.Text()
		consola.Nothing(option)
	} else if command == "mkdisk" {
		m := comandos.Mkdisk{}
		m.Exe(params)
	} else if command == "rmdisk" {
		r := comandos.Rmdisk{}
		r.Exe(params)
	} else if command == "fdisk" {
		f := comandos.Fdisk{}
		f.Exe(params)
	} else if command == "mount" {
		m := comandos.Mount{}
		m.Exe(params)
	} else if command == "mkfs" {
		m := comandos.Mkfs{}
		m.Exe(params)
	} else if command == "login" {
		l := usuariosygrupos.Login{}
		l.Exe(params)
	} else if command == "logout" {
		l := usuariosygrupos.Logout{}
		l.Exe(params)
	} else if command == "mkgrp" {
		m := usuariosygrupos.Mkgrp{}
		m.Exe(params)
	} else if command == "rmgrp" {
		r := usuariosygrupos.Rmgrp{}
		r.Exe(params)
	} else if command == "mkusr" {
		m := usuariosygrupos.Mkusr{}
		m.Exe(params)
	} else if command == "rmusr" {
		r := usuariosygrupos.Rmusr{}
		r.Exe(params)
	} else if command == "mkfile" {
		m := usuariosygrupos.Mkfile{}
		m.Exe(params)
	} else if command == "mkdir" {
		m := usuariosygrupos.Mkdir{}
		m.Exe(params)
	} else if command == "rep" {
		r := usuariosygrupos.Rep{}
		r.Exe(params)
	} else if strings.Contains(command, "#") {
		contenido := command
		contenido += "\n\n"
		consola.AddToConsole(contenido)
	}
}

func (a *Analyzer) Split_input(input string) []string {
	// fmt.Println("haciendo split al input")
	return strings.Split(input, ">")
}

func (a *Analyzer) Read(path string) {
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
