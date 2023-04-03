package comandos

import (
	"fmt"
	"os"

	"github.com/erazoex/-MIA-Proyecto2/Datos/data"
)

type executer interface {
	exe(parametros []string)
}

func WriteMBR(master *data.MBR, path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = file.Write(master)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()
}
