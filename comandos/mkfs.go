package comandos

import (
	"fmt"
	"strings"

	"github.com/erazoex/proyecto2/lista"
)

type ParametrosMkfs struct {
	id string
	t  string
}

type Mkfs struct {
	params ParametrosMkfs
}

func (m Mkfs) Exe(parametros []string) {
	m.params = m.SaveParams(parametros)
	if m.Mkfs(m.params.id, m.params.t) {
		fmt.Printf("\nel formateo con EXT2 de la particion con id %s fue exitoso\n\n", m.params.id)
	} else {
		fmt.Println("no se logro formatear la particion con id %s\n", m.params.id)
	}
}

func (m Mkfs) SaveParams(parametros []string) ParametrosMkfs {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		if strings.Contains(v, "id") {
			v = strings.ReplaceAll(v, "id=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.params.id = v
		} else if strings.Contains(v, "type") {
			v = strings.ReplaceAll(v, "type=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.params.t = v
		}
	}
	return m.params
}

func (m Mkfs) Mkfs(id string, t string) bool {
	// comprobando que id no este vacio
	if id == "" {
		fmt.Println("no se encontro el id entre los comandos")
		return false
	}
	// comprobando que type no lleve un valor incorrecto
	if t != "full" && t != "FULL" && t != "" {
		fmt.Println("el valor del comando type no es permitido")
		return false
	}
	if t == "" || t == "full" {
		t = "FULL"
	}
	//creando nuestro nodo auxiliar
	nodo := &lista.MountNode{}
	if lista.ListaMount.NodeExist(id) {
		nodo = lista.ListaMount.GetNodeById(id)
	} else {
		fmt.Printf("el id %s no coincide con ninguna particion montada\n", id)
		return false
	}
	m.Ext2(nodo)
	return true
}

func (m Mkfs) Ext2(nodo *lista.MountNode) {
	whereToStart := 0
	partSize := 0
	if nodo.Value != nil {
		whereToStart = int(nodo.Value.Part_start)
		partSize = int(nodo.Value.Part_size)
	} else if nodo.ValueL != nil {
		whereToStart = int(nodo.ValueL.Part_start)
		partSize = int(nodo.ValueL.Part_size)
	}
}
