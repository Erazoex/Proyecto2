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

type ParametrosMkusr struct {
	User string
	Pwd  string
	Grp  string
}

type Mkusr struct {
	Params ParametrosMkusr
}

func (m *Mkusr) Exe(parametros []string) {
	m.Params = m.SaveParams(parametros)
	if m.Mkusr(m.Params.User, m.Params.Pwd, m.Params.Grp) {
		consola.AddToConsole(fmt.Sprintf("\nusuario \"%s\" creado con exito en el grupo %s\n\n", m.Params.User, m.Params.Grp))
	} else {
		consola.AddToConsole(fmt.Sprintf("no se logro crear el usuario \"%s\"\n\n", m.Params.User))
	}
}

func (m *Mkusr) SaveParams(parametros []string) ParametrosMkusr {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.ReplaceAll(v, "\"", "")
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "user") {
			v = strings.ReplaceAll(v, "user=", "")
			m.Params.User = v
		}
		if strings.Contains(v, "grp") {
			v = strings.ReplaceAll(v, "grp=", "")
			m.Params.Grp = v
		}
		if strings.Contains(v, "pwd") {
			v = strings.ReplaceAll(v, "pwd=", "")
			m.Params.Pwd = v
		}
	}
	m.Params.User = strings.ReplaceAll(m.Params.User, "pwd=", "")
	return m.Params
}

func (m *Mkusr) Mkusr(user string, pwd string, grp string) bool {
	consola.AddToConsole(user + "\n")
	consola.AddToConsole(pwd + "\n")
	consola.AddToConsole(grp + "\n")
	if user == "" {
		consola.AddToConsole("no se encontro ningun nombre\n")
		return true
	}
	if logger.Log.IsLoggedIn() && logger.Log.UserIsRoot() {
		if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return m.MkusrPartition(user, pwd, grp, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return m.MkusrPartition(user, pwd, grp, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(datos.EBR{})), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		}
	}
	return false
}

func (m *Mkusr) MkusrPartition(user string, pwd string, grp string, whereToStart int64, path string) bool {
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
	if m.ExisteUsuario(ReadFile(&tablaInodo, path, &superbloque), user) {
		consola.AddToConsole(fmt.Sprintf("ya existe usuario con ese nombre %s\n", user))
		return false
	}
	if !m.ExisteGrupo(ReadFile(&tablaInodo, path, &superbloque), grp) {
		consola.AddToConsole(fmt.Sprintf("no existe un grupo con el nombre %s\n", grp))
		return false
	}
	numero := m.ContarUsuarios(ReadFile(&tablaInodo, path, &superbloque))
	usuario := m.AgregarUsuario(numero, grp, user, pwd)
	if AppendFile(path, &superbloque, &tablaInodo, usuario) {
		comandos.Fwrite(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(datos.TablaInodo{})))
		consola.AddToConsole(ReadFile(&tablaInodo, path, &superbloque))
		comandos.Fwrite(&superbloque, path, whereToStart)
		return true
	}
	return false
}

func (m *Mkusr) AgregarUsuario(userNumber int, groupName string, userName string, password string) string {
	return strconv.Itoa(userNumber) + ",U," + groupName + "," + userName + "," + password + "\n"
}

func (m *Mkusr) ContarUsuarios(contenido string) int {
	contador := 1
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		contador++
	}
	return contador
}

func (m *Mkusr) ExisteUsuario(contenido string, userName string) bool {
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		if parametros[3] == userName {
			return true
		}
	}
	return false
}
func (m *Mkusr) ExisteGrupo(contenido string, groupName string) bool {
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
