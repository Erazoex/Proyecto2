package comandos

import (
	"fmt"
	"os"
	"strings"

	"github.com/erazoex/proyecto2/consola"
)

type ParametrosRmdisk struct {
	Path string
}

type Rmdisk struct {
	Params ParametrosRmdisk
}

func (r *Rmdisk) Exe(parametros []string) {
	r.Params = r.SaveParams(parametros)
	if r.Rmdisk(r.Params.Path) {
		consola.AddToConsole(fmt.Sprintf("\nrmdisk realizado con exito para la ruta: %s\n\n", r.Params.Path))
	} else {
		consola.AddToConsole(fmt.Sprintf("\n[ERROR!] no se logro realizar el comando rmdisk para la ruta: %s\n\n", r.Params.Path))
	}
}

func (r *Rmdisk) SaveParams(parametros []string) ParametrosRmdisk {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			v = strings.ReplaceAll(v, "\"", "")
			r.Params.Path = v
		}
	}
	return r.Params
}

func (r *Rmdisk) Rmdisk(path string) bool {
	// Comprobando si existe una ruta valida para la creacion del disco
	if path == "" {
		consola.AddToConsole("no se encontro una ruta\n")
		return false
	}
	err := os.Remove(path)
	if err != nil {
		consola.AddToConsole(fmt.Sprintf("no se pudo eliminar el archivo %s\n", err.Error()))
		return false
	}
	return true
}
