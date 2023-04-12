package usuariosygrupos

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/lista"
	"github.com/erazoex/proyecto2/logger"
)

type ParametrosMkgrp struct {
	name string
}

type Mkgrp struct {
	params ParametrosMkgrp
}

func (m *Mkgrp) Exe(parametros []string) {
	m.params = m.SaveParams(parametros)
	if m.Mkgrp(m.params.name) {
		fmt.Printf("\ngrupo \"%s\" creado con exito\n\n", m.params.name)
	} else {
		fmt.Printf("no se logro crear el grupo \"%s\"\n\n", m.params.name)
	}
}

func (m *Mkgrp) SaveParams(parametros []string) ParametrosMkgrp {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "user=", "")
			m.params.name = v
		}
	}
	return m.params
}

func (m *Mkgrp) Mkgrp(name string) bool {
	if name == "" {
		fmt.Println("no se encontro ningun nombre")
		return true
	}
	if logger.Log.IsLoggedIn() && logger.Log.UserIsRoot() {
		fmt.Println("aqui entra")
		fmt.Println(lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil)
		if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			fmt.Println("print")
			return m.MkgrpPartition(name, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			fmt.Println("println")
			return m.MkgrpPartition(name, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(datos.EBR{})), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		}
	}
	return false
}

func (m *Mkgrp) MkgrpPartition(name string, whereToStart int64, path string) bool {
	fmt.Println("hola mundo!")
	// superbloque de la particion
	var superbloque datos.SuperBloque
	comandos.Fread(&superbloque, path, whereToStart)

	// tabla de inodos de archivo Users.txt
	var tablaInodo datos.TablaInodo
	comandos.Fread(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(datos.TablaInodo{})))
	// modificar la fecha en la que se esta modificando el inodo
	mtime := time.Now()
	for i := 0; i < len(tablaInodo.I_mtime); i++ {
		tablaInodo.I_mtime[i] = mtime.String()[i]
	}
	if m.ExisteGrupo(ReadFile(&tablaInodo, path, &superbloque), name) {
		fmt.Println("ya existe grupo con ese nombre", name)
		return false
	}
	numero := m.ContarGrupos(ReadFile(&tablaInodo, path, &superbloque))
	grupo := m.AgregarGrupo(numero, name)
	fmt.Println(ReadFile(&tablaInodo, path, &superbloque))
	return AppendFile(path, &superbloque, &tablaInodo, grupo)
}

func (m *Mkgrp) AgregarGrupo(groupNumber int, groupName string) string {
	return strconv.Itoa(groupNumber) + ",G," + groupName + "\n"
}

func (m *Mkgrp) ContarGrupos(contenido string) int {
	contador := 1
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] != "G" {
			continue
		}
		contador++
	}
	return contador
}

func (m *Mkgrp) ExisteGrupo(contenido string, groupName string) bool {
	fmt.Println("imprimir-->", contenido)
	fmt.Println("imprimir---", groupName)
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		fmt.Println(parametros[1])
		if parametros[1] != "G" {
			continue
		}
		if parametros[2] == groupName {
			return true
		}
	}
	return false
}
