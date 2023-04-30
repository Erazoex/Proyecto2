package comandos

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/erazoex/proyecto2/consola"
	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/functions"
	"github.com/erazoex/proyecto2/lista"
)

type ParametrosMount struct {
	Path string
	Name [16]byte
}

type Mount struct {
	Params ParametrosMount
}

func (m *Mount) Exe(parametros []string) {
	m.Params = m.SaveParams(parametros)
	if m.Mount(m.Params.Path, m.Params.Name) {
		consola.AddToConsole(fmt.Sprintf("\nparticion %s montada con exito\n\n", m.Params.Path))
	} else {
		consola.AddToConsole(fmt.Sprintf("no se logro montar la particion %s\n", m.Params.Path))
	}
}

func (m *Mount) SaveParams(parametros []string) ParametrosMount {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			v = strings.ReplaceAll(v, "\"", "")
			m.Params.Path = v
		} else if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			copy(m.Params.Name[:], v[:])
		}
	}
	return m.Params
}

func (m *Mount) Mount(path string, name [16]byte) bool {
	// comprobando que el parametro "path" sea diferente de ""
	if path == "" {
		consola.AddToConsole("no se encontro una ruta\n")
		return false
	}
	// comprobando que el parametro "name" sea diferente de ""
	if bytes.Equal(name[:], []byte("")) {
		consola.AddToConsole("se debe de contar con un nombre para realizar este comando\n")
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
				consola.AddToConsole("la particion ya se encuentra montada\n")
				return false
			}
			if particion.Part_type == 'e' || particion.Part_type == 'E' {
				consola.AddToConsole("no se puede montar una particion extendida\n")
				return false
			}
			particion.Part_type = '2'
			var part *datos.Partition = new(datos.Partition)
			part = &particion
			lista.ListaMount.Mount(path, 53, part, nil)
			partitionMounted = true
			// tener un metodo de MountList que agregue un texto a la consola
			lista.ListaMount.PrintList()
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
				lista.ListaMount.PrintList()
			}
		}
	}
	if !partitionMounted {
		consola.AddToConsole(fmt.Sprintf("no se encontro una particion con el nombre de %s\n", string(functions.TrimArray(name[:]))))
		return false
	}
	WriteMBR(&master, path)
	return true
}

func (m *Mount) MountParticionLogica(path string, whereToStart int, name [16]byte) {
	logicPartitionMounted := false
	temp := ReadEBR(path, int64(whereToStart))
	flag := true
	for flag {
		if bytes.Equal(temp.Part_name[:], name[:]) {
			temp.Part_status = '2'
			var partL *datos.EBR = new(datos.EBR)
			partL = &temp
			lista.ListaMount.Mount(path, 53, nil, partL)
			logicPartitionMounted = true
			flag = false
			break
		} else if temp.Part_next != -1 {
			temp = ReadEBR(path, temp.Part_next)
		} else {
			flag = false
		}
	}
	if !logicPartitionMounted {
		consola.AddToConsole(fmt.Sprintf("no se encontro una particion con el nombre %s\n", string(functions.TrimArray(name[:]))))
		return
	}
}
