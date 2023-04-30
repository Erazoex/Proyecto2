package usuariosygrupos

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/consola"
	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/lista"
	"github.com/erazoex/proyecto2/logger"
)

type ParametrosMkgrp struct {
	Name string
}

type Mkgrp struct {
	params ParametrosMkgrp
}

func (m *Mkgrp) Exe(parametros []string) {
	m.params = m.SaveParams(parametros)
	if m.Mkgrp(m.params.Name) {
		consola.AddToConsole(fmt.Sprintf("\ngrupo \"%s\" creado con exito\n\n", m.params.Name))
	} else {
		consola.AddToConsole(fmt.Sprintf("no se logro crear el grupo \"%s\"\n\n", m.params.Name))
	}
}

func (m *Mkgrp) SaveParams(parametros []string) ParametrosMkgrp {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			m.params.Name = v
		}
	}
	return m.params
}

func (m *Mkgrp) Mkgrp(name string) bool {
	if name == "" {
		consola.AddToConsole("no se encontro ningun nombre\n")
		return true
	}
	if logger.Log.IsLoggedIn() && logger.Log.UserIsRoot() {
		if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return m.MkgrpPartition(name, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return m.MkgrpPartition(name, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(datos.EBR{})), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		}
	}
	return false
}

func (m *Mkgrp) MkgrpPartition(name string, whereToStart int64, path string) bool {
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
		consola.AddToConsole(fmt.Sprintf("ya existe grupo con ese nombre %s\n", name))
		return false
	}
	numero := m.ContarGrupos(ReadFile(&tablaInodo, path, &superbloque))
	grupo := m.AgregarGrupo(numero, name)
	if AppendFile(path, &superbloque, &tablaInodo, grupo) {
		comandos.Fwrite(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(datos.TablaInodo{})))
		consola.AddToConsole(ReadFile(&tablaInodo, path, &superbloque))
		comandos.Fwrite(&superbloque, path, whereToStart)
		return true
	}
	return false
}

func (m *Mkgrp) AgregarGrupo(groupNumber int, groupName string) string {
	return strconv.Itoa(groupNumber) + ",G," + groupName + "\n"
}

func (m *Mkgrp) ContarGrupos(contenido string) int {
	contador := 1
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
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
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] != "G" {
			continue
		}
		if parametros[2] == groupName {
			return true
		}
	}
	return false
}
