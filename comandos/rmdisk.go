package comandos

import (
	"fmt"
	"os"
	"strings"
)

type ParametrosRmdisk struct {
	path string
}

type Rmdisk struct {
	params ParametrosRmdisk
}

func (r *Rmdisk) Exe(parametros []string) {
	r.params = r.SaveParams(parametros)
	if r.Rmdisk(r.params.path) {
		fmt.Printf("\nrmdisk realizado con exito para la ruta: %s\n\n", r.params.path)
	} else {
		fmt.Printf("\n[ERROR!] no se logro realizar el comando rmdisk para la ruta: %s\n\n", r.params.path)
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
			r.params.path = v
		}
	}
	return r.params
}

func (r *Rmdisk) Rmdisk(path string) bool {
	// Comprobando si existe una ruta valida para la creacion del disco
	if path == "" {
		fmt.Println("no se encontro una ruta", "path = \"\"")
		return false
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println("no se pudo eliminar el archivo", err.Error())
		return false
	}
	return true
}
