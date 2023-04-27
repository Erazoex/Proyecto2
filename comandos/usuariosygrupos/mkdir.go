package usuariosygrupos

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/erazoex/proyecto2/bitmap"
	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/lista"
	"github.com/erazoex/proyecto2/logger"
)

type ParametrosMkdir struct {
	Path string
	R    bool
}

type Mkdir struct {
	Params ParametrosMkdir
}

func (m *Mkdir) Exe(parametros []string) {
	m.Params = m.SaveParams(parametros)
	if m.Mkdir(m.Params.Path, m.Params.R) {
		fmt.Printf("\nla carpeta con ruta %s se creo correctamente\n\n", m.Params.Path)
	} else {
		fmt.Printf("\nla carpeta con ruta %s no se pudo crear\n\n", m.Params.Path)
	}
}

func (m *Mkdir) SaveParams(parametros []string) ParametrosMkdir {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			m.Params.Path = v
		} else if v == "r" {
			// v = strings.ReplaceAll(v, "r", "")
			m.Params.R = true
		}
	}
	return m.Params
}

func (m *Mkdir) Mkdir(path string, r bool) bool {
	if path == "" {
		fmt.Println("no se encontro una ruta")
		return false
	}
	path = strings.Replace(path, "/", "", 1)
	if !logger.Log.IsLoggedIn() {
		fmt.Println("no se encuentra un usuario loggeado para crear un archivo")
		return false
	}

	if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
		return createDirectory(logger.Log.GetUserName(), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta, path, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, r)
	} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL != nil {
		return createDirectory(logger.Log.GetUserName(), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta, path, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(datos.EBR{})), r)
	}
	return false
}

func createDirectory(name [10]byte, path, ruta string, whereToStart int64, r bool) bool {
	// superbloque de la particion
	var superbloque datos.SuperBloque
	comandos.Fread(&superbloque, path, whereToStart)
	var tablaInodoRoot datos.TablaInodo
	comandos.Fread(&tablaInodoRoot, path, superbloque.S_inode_start)

	// obteniendo el contenido de la tabla de inodos de Users.txt
	var tablaInodoUsers datos.TablaInodo
	comandos.Fread(&tablaInodoUsers, path, superbloque.S_inode_start+superbloque.S_inode_size)
	contenido := ReadFile(&tablaInodoUsers, path, &superbloque)
	// obteniendo el user id y el group id
	userId := GetUserId(contenido, string(name[:]))
	groupdId := GetGroupId(contenido, string(name[:]))
	if r {
		// Create directories
		FindAndCreateDirectories(&tablaInodoRoot, path, ruta, &superbloque, 0, userId, groupdId)
	}

	num := NewInodeDirectory(&superbloque, path, userId, groupdId)
	FindDirs(num, &tablaInodoRoot, path, ruta, &superbloque, 0)
	comandos.Fwrite(&tablaInodoRoot, path, superbloque.S_inode_start)
	comandos.Fwrite(&superbloque, path, whereToStart)
	PrintTree(&tablaInodoRoot, &superbloque, path)
	return true
}

func NewInodeDirectory(superbloque *datos.SuperBloque, path string, userId, groupId int64) int64 {
	var nuevaTabla datos.TablaInodo
	posicionActual := bitmap.WriteInBitmapBlock(path, superbloque)
	// aqui llenaremos la nueva tabla de inodos
	nuevaTabla.I_uid = userId
	nuevaTabla.I_gid = groupId
	nuevaTabla.I_size = 0
	nuevaTabla.I_type = '0'
	nuevaTabla.I_perm = 664
	// llenando las fechas
	atime := time.Now()
	for i := 0; i < len(nuevaTabla.I_atime)-1; i++ {
		nuevaTabla.I_atime[i] = atime.String()[i]
	}
	ctime := time.Now()
	for i := 0; i < len(nuevaTabla.I_atime)-1; i++ {
		nuevaTabla.I_ctime[i] = ctime.String()[i]
	}
	mtime := time.Now()
	for i := 0; i < len(nuevaTabla.I_atime)-1; i++ {
		nuevaTabla.I_mtime[i] = mtime.String()[i]
	}
	// llenando a todos los bloques no utilizados
	for i := 0; i < len(nuevaTabla.I_block); i++ {
		nuevaTabla.I_block[i] = -1
	}

	// aqui escribiremos el contenido y crearemos el nuevo bloque de carpeta
	posicionNuevoBloqueCarpetas := bitmap.WriteInBitmapBlock(path, superbloque)
	nuevaTabla.I_block[0] = posicionNuevoBloqueCarpetas

	nuevoBloqueCarpetas := datos.BloqueDeCarpetas{}

	// llenando la carpeta
	copy(nuevoBloqueCarpetas.B_content[0].B_name[:], ".")
	nuevoBloqueCarpetas.B_content[0].B_inodo = int32(posicionActual)

	copy(nuevoBloqueCarpetas.B_content[1].B_name[:], "..")
	nuevoBloqueCarpetas.B_content[1].B_inodo = -1

	copy(nuevoBloqueCarpetas.B_content[2].B_name[:], "")
	nuevoBloqueCarpetas.B_content[2].B_inodo = -1

	copy(nuevoBloqueCarpetas.B_content[3].B_name[:], "")
	nuevoBloqueCarpetas.B_content[3].B_inodo = -1
	// escribiendo la nueva tabla de inodos
	comandos.Fwrite(&nuevaTabla, path, superbloque.S_inode_start+posicionActual*superbloque.S_inode_size)
	return posicionActual
}

func GetContent(cont string) string {
	// aqui hay que leer el archivo y ejecutarlo
	file, err := os.Open(cont)
	if err != nil {
		fmt.Printf("Error al intentar abrir el archivo: %s\n", cont)
		return ""
	}

	defer file.Close()

	// Crear un scanner para luego leer linea por linea el archivo de entrada
	scanner := bufio.NewScanner(file)
	content := ""
	// Leyendo linea por linea
	for scanner.Scan() {
		// obteniendo la linea actual
		content += scanner.Text()
	}

	// comprobar que no hubo error al leer el archivo
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo: ", err)
		return ""
	}
	return content
}
