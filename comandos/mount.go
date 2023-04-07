package comandos

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/lista"
)

type ParametrosMount struct {
	path string
	name [16]byte
}

type Mount struct {
	params ParametrosMount
}

func (m Mount) Exe(parametros []string) {
	m.params = m.SaveParams(parametros)
	if m.Mount(m.params.path, m.params.name) {
		fmt.Printf("\nparticion %s montada con exito\n\n", m.params.path)
	} else {
		fmt.Printf("no se logro montar la particion %s\n", m.params.path)
	}
}

func (m Mount) SaveParams(parametros []string) ParametrosMount {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			v = strings.ReplaceAll(v, "\"", "")
			m.params.path = v
		} else if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			v = v[:16]
			copy(m.params.name[:], v)
		}
	}
	return m.params
}

func (m Mount) Mount(path string, name [16]byte) bool {
	// comprobando que el parametro "path" sea diferente de ""
	if path == "" {
		fmt.Println("no se encontro una ruta")
		return false
	}
	// comprobando que el parametro "name" sea diferente de ""
	if bytes.Equal(name[:], []byte("")) {
		fmt.Println("se debe de contar con un nombre para realizar este comando")
		return false
	}
	master := GetMBR(path)
	partitionMounted := false
	particionEncontrada := false
	for _, particion := range master.Mbr_partitions {
		// si entro aqui es porque si leyo el MBR del disco
		if bytes.Equal(particion.Part_name[:], name[:]) {
			// comprobaremos que la particion no se haya montado previamente
			particionEncontrada = true
			if particion.Part_status == '2' {
				fmt.Println("la particion ya se encuentra montada")
				return false
			}
			if particion.Part_type == 'e' || particion.Part_type == 'E' {
				fmt.Println("no se puede montar una particion extendida")
				return false
			}
			particion.Part_type = '2'
			var part *datos.Partition = new(datos.Partition)
			lista.ListaMount.Mount(path, 53, part, nil)
			partitionMounted = true
			// tener un metodo de MountList que agregue un texto a la consola
			break

		}
	}
	if !particionEncontrada {
		// buscaremos si existe una particion logica con ese nombre
		for _, particion := range master.Mbr_partitions {
			if particion.Part_type == 'e' || particion.Part_type == 'E' {
				partitionMounted = true
				m.MountParticionLogica(path, int(particion.Part_start), name)
				// tener un metodo de Mount List que agregue un texto a la consola
			}
		}
	}
	if !partitionMounted {
		fmt.Printf("no se encontro una particion con el nombre de %s\n", name)
		return false
	}
	WriteMBR(&master, path)
	return true
}

func (m Mount) MountParticionLogica(path string, whereToStart int, name [16]byte) {
	//TODO: aqui me quede toca seguir manana
}
