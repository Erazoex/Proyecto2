package consola

var content string

func AddToConsole(nuevoContenido string) {
	content += nuevoContenido
}

func GetConsole() string {
	returnable := content
	content = ""
	return returnable
}
