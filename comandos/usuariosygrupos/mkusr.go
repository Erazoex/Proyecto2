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

type ParametrosMkusr struct {
	User string
	Pwd  string
	Grp  string
}

type Mkusr struct {
	params ParametrosMkusr
}

func (m *Mkusr) Exe(parametros []string) {
	m.params = m.SaveParams(parametros)
	if m.Mkusr(m.params.User, m.params.Pwd, m.params.Grp) {
		fmt.Printf("\nusuario \"%s\" creado con exito en el grupo %s\n\n", m.params.User, m.params.Grp)
	} else {
		fmt.Printf("no se logro crear el usuario \"%s\"\n\n", m.params.User)
	}
}

func (m *Mkusr) SaveParams(parametros []string) ParametrosMkusr {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "user") {
			v = strings.ReplaceAll(v, "user=", "")
			m.params.User = v
		} else if strings.Contains(v, "pwd") {
			v = strings.ReplaceAll(v, "pwd=", "")
			m.params.Pwd = v
		} else if strings.Contains(v, "grp") {
			v = strings.ReplaceAll(v, "grp=", "")
			m.params.Grp = v
		}
	}
	return m.params
}

func (m *Mkusr) Mkusr(user string, pwd string, grp string) bool {
	if user == "" {
		fmt.Println("no se encontro ningun nombre")
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
		fmt.Println("ya existe usuario con ese nombre", user)
		return false
	}
	if !m.ExisteGrupo(ReadFile(&tablaInodo, path, &superbloque), grp) {
		fmt.Printf("no existe un grupo con el nombre %s\n", grp)
		return false
	}
	numero := m.ContarUsuarios(ReadFile(&tablaInodo, path, &superbloque))
	usuario := m.AgregarUsuario(numero, grp, user, pwd)
	if AppendFile(path, &superbloque, &tablaInodo, usuario) {
		comandos.Fwrite(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(datos.TablaInodo{})))
		fmt.Println(ReadFile(&tablaInodo, path, &superbloque))
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
