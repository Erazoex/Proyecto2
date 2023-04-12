package usuariosygrupos

import (
	"bytes"
	"fmt"
	"strings"
	"unsafe"

	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/functions"
	"github.com/erazoex/proyecto2/lista"
	"github.com/erazoex/proyecto2/logger"
)

type ParametrosLogin struct {
	User [10]byte
	Pwd  [10]byte
	Id   string
}

type Login struct {
	Params ParametrosLogin
}

func (l *Login) Exe(parametros []string) {
	l.Params = l.SaveParams(parametros)
	if l.Login(l.Params.User, l.Params.Pwd, l.Params.Id) {
		fmt.Printf("\nusuario \"%s\" loggeado con exito\n\n", l.Params.User)
	} else {
		fmt.Printf("no se logro loggear el usuario \"%s\"\n\n", l.Params.User)
	}
}

func (l *Login) SaveParams(parametros []string) ParametrosLogin {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "user") {
			v = strings.ReplaceAll(v, "user=", "")
			copy(l.Params.User[:], v[:])
		} else if strings.Contains(v, "pwd") {
			v = strings.ReplaceAll(v, "pwd=", "")
			copy(l.Params.Pwd[:], v[:])
		} else if strings.Contains(v, "id") {
			v = strings.ReplaceAll(v, "id=", "")
			l.Params.Id = v
		}
	}
	return l.Params
}

func (l *Login) Login(User [10]byte, Pwd [10]byte, Id string) bool {
	if bytes.Equal(User[:], []byte("")) {
		fmt.Println("no hay user el cual utilizar")
		return false
	}
	if bytes.Equal(Pwd[:], []byte("")) {
		fmt.Println("el usuario no tiene password")
		return false
	}
	if Id == "" {
		fmt.Println("no hay id para buscar en las particiones montadas")
		return false
	}

	node := lista.ListaMount.GetNodeById(Id)
	if node == nil {
		fmt.Printf("el id %s no coincide con ninguna particion montada\n", Id)
		return false
	}
	if node.Value != nil {
		return l.LoginInPrimaryPartition(node.Ruta, User, Pwd, Id, node.Value)
	} else if node.ValueL != nil {
		return l.LoginInLogicPartition(node.Ruta, User, Pwd, Id, node.ValueL)
	} else {
		// no deberia de entrar aqui nunca
		fmt.Println("no hay particion montada")
	}
	fmt.Printf("no se logro loggear el usuario: %s", User)
	return false
}

func (l *Login) LoginInPrimaryPartition(path string, User [10]byte, Pwd [10]byte, Id string, partition *datos.Partition) bool {
	// leyendo el superbloque
	var superbloque datos.SuperBloque
	comandos.Fread(&superbloque, path, partition.Part_start)

	// tabla de inodos del archivo
	var tablaInodo datos.TablaInodo
	comandos.Fread(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(datos.TablaInodo{})))

	// vamos a recorrer la tabla de inodos del archivo Users.txt
	var contenido string
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo datos.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
		contenido += string(parteArchivo.B_content[:])
	}
	// leeremos el archivo por linea que se encuentre dentro del archivo
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		grupo := parametros[2]
		username := parametros[3]
		password := parametros[4]
		if !functions.Equal(User, username) || !functions.Equal(Pwd, password) {
			continue
		}
		user := &logger.User{
			User: User,
			Pass: Pwd,
			Id:   Id,
		}
		copy(user.Grupo[:], grupo)
		return logger.Log.Login(user)
	}
	fmt.Println("no se encontro el usuario dentro del archivo")
	return false
}

func (l *Login) LoginInLogicPartition(path string, User [10]byte, Pwd [10]byte, Id string, partition *datos.EBR) bool {
	// leyendo el superbloque
	var superbloque datos.SuperBloque
	comandos.Fread(&superbloque, path, partition.Part_start+int64(unsafe.Sizeof(datos.EBR{})))

	// tabla de inodos del archivo
	var tablaInodo datos.TablaInodo
	comandos.Fread(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(datos.TablaInodo{})))

	// vamos a recorrer la tabla de inodos del archivo Users.Txt
	var contenido string
	for i := 0; i < len(tablaInodo.I_block); i++ {
		// fmt.Println(tablaInodo.I_block[i])
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo datos.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
		contenido += string(parteArchivo.B_content[:])
	}
	// leeremos el archivo por linea que se encuentre dentro del archivo
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		grupo := parametros[2]
		username := parametros[3]
		password := parametros[4]
		if !functions.Equal(User, username) || !functions.Equal(Pwd, password) {
			continue
		}
		user := &logger.User{
			User: User,
			Pass: Pwd,
		}
		copy(user.Grupo[:], grupo)
		return logger.Log.Login(user)
	}
	fmt.Println("no se encontro el usuario dentro del archivo")
	return false
}
