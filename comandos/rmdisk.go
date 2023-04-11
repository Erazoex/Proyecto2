package comandos

import (
	"fmt"
	"os"
	"strings"
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
		fmt.Printf("\nrmdisk realizado con exito para la ruta: %s\n\n", r.Params.Path)
	} else {
		fmt.Printf("\n[ERROR!] no se logro realizar el comando rmdisk para la ruta: %s\n\n", r.Params.Path)
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
