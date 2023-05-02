package consola

import "fmt"

var content string

func AddToConsole(nuevoContenido string) {
	fmt.Printf(nuevoContenido)
}

func Nothing(contenido string) {
	contenido = ""
}

func GetConsole() string {
	returnable := content
	content = ""
	return returnable
}
