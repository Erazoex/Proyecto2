package usuariosygrupos

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/consola"
	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/lista"
	"github.com/erazoex/proyecto2/logger"
)

type ParametrosRmgrp struct {
	Name string
}

type Rmgrp struct {
	params ParametrosRmgrp
}

func (r *Rmgrp) Exe(parametros []string) {
	r.params = r.SaveParams(parametros)
	if r.Rmgrp(r.params.Name) {
		consola.AddToConsole(fmt.Sprintf("\ngrupo \"%s\" eliminado con exito\n\n", r.params.Name))
	} else {
		consola.AddToConsole(fmt.Sprintf("no se logro eliminar el grupo \"%s\"\n\n", r.params.Name))
	}
}

func (r *Rmgrp) SaveParams(parametros []string) ParametrosRmgrp {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			r.params.Name = v
		}
	}
	return r.params
}

func (r *Rmgrp) Rmgrp(name string) bool {
	if name == "" {
		consola.AddToConsole("no se encontro ningun nombre\n")
		return true
	}
	if logger.Log.IsLoggedIn() && logger.Log.UserIsRoot() {
		if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return r.RmgrpPartition(name, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return r.RmgrpPartition(name, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(datos.EBR{})), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		}
	}
	return false
}

func (r *Rmgrp) RmgrpPartition(name string, whereToStart int64, path string) bool {
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
	if !r.ExisteGrupo(ReadFile(&tablaInodo, path, &superbloque), name) {
		consola.AddToConsole(fmt.Sprintf("no existe grupo con ese nombre %s\n", name))
		return false
	}
	contenido := modFile(&tablaInodo, path, &superbloque)
	nuevoContenido := r.DesactivarGrupo(contenido, name)
	// fmt.Println(nuevoContenido)
	if SetFile(&tablaInodo, path, &superbloque, nuevoContenido) {
		consola.AddToConsole(ReadFile(&tablaInodo, path, &superbloque))
		return true
	}
	return false
}

func (r *Rmgrp) DesactivarGrupo(contenido string, groupName string) string {
	nuevoContenido := ""
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] == "G" {
			if parametros[2] == groupName {
				parametros[0] = "0"
			}
			nuevoContenido += parametros[0] + "," + parametros[1] + "," + parametros[2] + "\n"
		} else if parametros[1] == "U" {
			nuevoContenido += parametros[0] + "," + parametros[1] + "," + parametros[2] + "," + parametros[3] + "," + parametros[4] + "\n"
		}
	}
	return nuevoContenido
}

func (r *Rmgrp) ExisteGrupo(contenido string, groupName string) bool {
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
