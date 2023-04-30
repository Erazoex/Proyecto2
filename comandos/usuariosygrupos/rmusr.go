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

type ParametrosRmusr struct {
	User string
}

type Rmusr struct {
	params ParametrosRmusr
}

func (r *Rmusr) Exe(parametros []string) {
	r.params = r.SaveParams(parametros)
	if r.Rmusr(r.params.User) {
		consola.AddToConsole(fmt.Sprintf("\nusuario \"%s\" eliminado con exito\n\n", r.params.User))
	} else {
		consola.AddToConsole(fmt.Sprintf("no se logro eliminar el usuario \"%s\"\n\n", r.params.User))
	}
}

func (r *Rmusr) SaveParams(parametros []string) ParametrosRmusr {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "user") {
			v = strings.ReplaceAll(v, "user=", "")
			r.params.User = v
		}
	}
	return r.params
}

func (r *Rmusr) Rmusr(user string) bool {
	if user == "" {
		consola.AddToConsole("no se encontro ningun nombre\n")
		return true
	}
	if logger.Log.IsLoggedIn() && logger.Log.UserIsRoot() {
		if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return r.RmusrPartition(user, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return r.RmusrPartition(user, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(datos.EBR{})), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		}
	}
	return false
}

func (r *Rmusr) RmusrPartition(user string, whereToStart int64, path string) bool {
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
	if !r.ExisteUsuario(ReadFile(&tablaInodo, path, &superbloque), user) {
		consola.AddToConsole(fmt.Sprintf("no existe usuario con ese nombre %s\n", user))
		return false
	}
	contenido := modFile(&tablaInodo, path, &superbloque)
	nuevoContenido := r.DesactivarUsuario(contenido, user)
	// fmt.Println(nuevoContenido)
	if SetFile(&tablaInodo, path, &superbloque, nuevoContenido) {
		consola.AddToConsole(ReadFile(&tablaInodo, path, &superbloque))
		return true
	}
	return false
}

func (r *Rmusr) DesactivarUsuario(contenido string, userName string) string {
	nuevoContenido := ""
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] == "U" {
			if parametros[3] == userName {
				parametros[0] = "0"
			}
			nuevoContenido += parametros[0] + "," + parametros[1] + "," + parametros[2] + "," + parametros[3] + "," + parametros[4] + "\n"
		} else if parametros[1] == "G" {
			nuevoContenido += parametros[0] + "," + parametros[1] + "," + parametros[2] + "\n"
		}
	}
	return nuevoContenido
}

func (r *Rmusr) ExisteUsuario(contenido string, userName string) bool {
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		if strings.Compare(parametros[3], userName) == 0 {
			return true
		}
	}
	return false
}
