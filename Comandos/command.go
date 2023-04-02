package comandos

import "fmt"

type executer interface {
	exe(parametros []string)
}

func writeMBR(master *MBR) {
	fmt.Println(("hola mundo!"))
}
