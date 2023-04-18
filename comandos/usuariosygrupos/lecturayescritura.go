package usuariosygrupos

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/erazoex/proyecto2/bitmap"
	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/datos"
)

func ReadFile(tablaInodo *datos.TablaInodo, path string, superbloque *datos.SuperBloque) string {
	contenido := ""
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo datos.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
		contenido += string(parteArchivo.B_content[:])
	}
	return contenido
}

func AppendFile(path string, superbloque *datos.SuperBloque, tablaInodo *datos.TablaInodo, contenido string) bool {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] != -1 {
			var parteArchivo datos.BloqueDeArchivos
			comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
			if strlen(parteArchivo.B_content) == 63 {
				continue
			} else if strlen(parteArchivo.B_content) < 63 {
				nuevoContenido := string(trimArray(parteArchivo.B_content[:])) + contenido
				// nuevoContenidoArray := createArray([]byte(nuevoContenido))
				if StrlenBytes([]byte(nuevoContenido)) > 63 {
					copy(parteArchivo.B_content[:], nuevoContenido[:63])
					comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
					AppendFile(path, superbloque, tablaInodo, string(nuevoContenido[63:]))
				} else {
					copy(parteArchivo.B_content[:], []byte(nuevoContenido))
					comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
				}
				return true
			}
		} else if tablaInodo.I_block[i] == -1 {
			var nuevoBloque datos.BloqueDeArchivos
			nuevaPosicion := bitmap.WriteInBitmapBlock(path, superbloque)
			nuevoContenido := trimArray([]byte(contenido))
			tablaInodo.I_block[i] = nuevaPosicion
			if StrlenBytes([]byte(contenido)) > 63 {
				copy(nuevoBloque.B_content[:], nuevoContenido[:63])
				comandos.Fwrite(&nuevoBloque, path, superbloque.S_block_start+nuevaPosicion*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
				AppendFile(path, superbloque, tablaInodo, string(nuevoContenido[63:]))
			} else {
				copy(nuevoBloque.B_content[:], nuevoContenido)
				comandos.Fwrite(&nuevoBloque, path, superbloque.S_block_start+nuevaPosicion*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
			}
			return true
		}
	}
	return false
}

func modFile(tablaInodo *datos.TablaInodo, path string, superbloque *datos.SuperBloque) string {
	contenido := ""
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo datos.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
		contenido += string(parteArchivo.B_content[:])
		parteArchivo.B_content = [64]byte{}
		// fmt.Println(parteArchivo.B_content)
		copy(parteArchivo.B_content[:], "modificar")
		comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
	}
	return contenido
}

func SetFile(tablaInodo *datos.TablaInodo, path string, superbloque *datos.SuperBloque, contenido string) bool {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo datos.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
		comparador := [64]byte{}
		copy(comparador[:], []byte("modificar"))
		if bytes.Equal(parteArchivo.B_content[:], comparador[:]) {
			if StrlenBytes([]byte(contenido)) > 63 {
				copy(parteArchivo.B_content[:], contenido[:63])
				comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
				return SetFile(tablaInodo, path, superbloque, string(contenido[63:]))
			} else {
				copy(parteArchivo.B_content[:], []byte(contenido))
				comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
			}
			return true
		}

	}
	return false
}

func FindAndCreateDirectories(tablaInodo *datos.TablaInodo, path, ruta string, superbloque *datos.SuperBloque, posicionActual, userId, groupId int64) {
	var rutaParts []string
	if !strings.Contains(ruta, "/") {
		// si no contiene un "/" quiere decir que ya estamos con el nombre del archivo
		// por lo tanto ya no hay necesidad de crear un directorio
		return
	}
	rutaParts = strings.SplitN(ruta, "/", 2)
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueDeCarpetas datos.BloqueDeCarpetas
		if tablaInodo.I_block[i] == -1 {
			// esto significa que no encontro el directorio que se buscaba
			posicion := bitmap.WriteInBitmapInode(path, superbloque)
			tablaInodo.I_block[i] = posicion
			comandos.Fwrite(&tablaInodo, path, posicionActual*int64(unsafe.Sizeof(datos.TablaInodo{})))
			// despues de haber obtenido la nueva posicion de la nueva tabla de inodos
			// ahora la crearemos y la llenaremos
			var nuevaTablaInodo datos.TablaInodo
			CreateNewDirectory(&nuevaTablaInodo, path, rutaParts[0], superbloque, posicion, posicionActual, userId, groupId)
			FindAndCreateDirectories(&nuevaTablaInodo, path, rutaParts[1], superbloque, posicion, userId, groupId)
		}
		comandos.Fread(&bloqueDeCarpetas, path, tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeCarpetas{})))
		num, compare := CompareDirectories(rutaParts[0], &bloqueDeCarpetas)
		if compare {
			var nuevaTablaInodo datos.TablaInodo
			comandos.Fread(&nuevaTablaInodo, path, num*int64(unsafe.Sizeof(datos.TablaInodo{})))
			FindAndCreateDirectories(&nuevaTablaInodo, path, rutaParts[1], superbloque, num, userId, groupId)
			break
		}
	}
}

func FindDirectories(AgregarTabla int64, tablaInodo *datos.TablaInodo, path, ruta string, superbloque *datos.SuperBloque, posicionActual int64) {
	var rutaParts []string
	if !strings.Contains(ruta, "/") {
		// si no contiene un "/" quiere decir que ya estamos con el nombre del archivo
		// por lo tanto ya no hay necesidad de crear un directorio
		fmt.Println("el archivo", ruta)
		crearArchivoDentroDeTablaInodo(AgregarTabla, tablaInodo, superbloque, path, ruta)
		return
	}
	rutaParts = strings.SplitN(ruta, "/", 2)
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueDeCarpetas datos.BloqueDeCarpetas
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		comandos.Fread(&bloqueDeCarpetas, path, tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeCarpetas{})))
		num, compare := CompareDirectories(rutaParts[0], &bloqueDeCarpetas)
		if compare {
			var nuevaTablaInodo datos.TablaInodo
			comandos.Fread(&nuevaTablaInodo, path, num*int64(unsafe.Sizeof(datos.TablaInodo{})))
			FindDirectories(AgregarTabla, &nuevaTablaInodo, path, rutaParts[1], superbloque, num)
			break
		}
	}
}

func crearArchivoDentroDeTablaInodo(AgregarTabla int64, tablaInodo *datos.TablaInodo, superbloque *datos.SuperBloque, path, nombreArchivo string) {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueCarpeta datos.BloqueDeCarpetas
		if tablaInodo.I_block[i] == -1 {
			posicionBloqueCarpeta := bitmap.WriteInBitmapBlock(path, superbloque)
			tablaInodo.I_block[i] = posicionBloqueCarpeta
			copy(bloqueCarpeta.B_content[0].B_name[:], []byte(nombreArchivo))
			bloqueCarpeta.B_content[0].B_inodo = int32(AgregarTabla)
			comandos.Fwrite(&bloqueCarpeta, path, superbloque.S_block_start+posicionBloqueCarpeta*superbloque.S_block_size)
		} else {
			comandos.Fread(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
			for j := 0; j < len(bloqueCarpeta.B_content); i++ {
				if bloqueCarpeta.B_content[j].B_inodo == -1 {
					bloqueCarpeta.B_content[j].B_inodo = int32(AgregarTabla)
					copy(bloqueCarpeta.B_content[j].B_name[:], []byte(nombreArchivo))
					comandos.Fwrite(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
				}
			}
		}
	}
}

func CreateNewDirectory(nuevaTabla *datos.TablaInodo, path, nombre string, superbloque *datos.SuperBloque, posicionActual, posicionPadre, userId, groupId int64) {
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
	posicionNuevoBloqueCarpetas := bitmap.WriteInBitmapBlock(path, superbloque)
	nuevaTabla.I_block[0] = posicionNuevoBloqueCarpetas

	// escribiendo la nueva tabla de inodos
	comandos.Fwrite(nuevaTabla, path, posicionActual*int64(unsafe.Sizeof(datos.TablaInodo{})))

	nuevoBloqueCarpetas := datos.BloqueDeCarpetas{}

	// llenando la carpeta
	copy(nuevoBloqueCarpetas.B_content[0].B_name[:], ".")
	nuevoBloqueCarpetas.B_content[0].B_inodo = int32(posicionActual)

	copy(nuevoBloqueCarpetas.B_content[1].B_name[:], "..")
	nuevoBloqueCarpetas.B_content[1].B_inodo = int32(posicionPadre)

	copy(nuevoBloqueCarpetas.B_content[2].B_name[:], "")
	nuevoBloqueCarpetas.B_content[2].B_inodo = -1

	copy(nuevoBloqueCarpetas.B_content[3].B_name[:], "")
	nuevoBloqueCarpetas.B_content[3].B_inodo = -1

	// escribiendo el nuevo bloque de carpetas
	comandos.Fwrite(&nuevoBloqueCarpetas, path, posicionNuevoBloqueCarpetas)
}

func SearchFreeSpace(tablaInodo *datos.TablaInodo, path string, superbloque *datos.SuperBloque) datos.BloqueDeCarpetas {
	var bloqueCarpeta datos.BloqueDeCarpetas
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			fmt.Println("no se encontro un espacio libre")
			return bloqueCarpeta
		}
		comandos.Fread(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeCarpetas{})))
		if !SearchFreeSpaceInBlock(&bloqueCarpeta) {
			continue
		}
		// CreateNewFile(&bloqueCarpeta)
	}
	return bloqueCarpeta // esta puesto solo por poner
}

func CreateNewFile(bloqueCarpeta *datos.BloqueDeCarpetas, path string, superbloque *datos.SuperBloque) {
	for i := 0; i < len(bloqueCarpeta.B_content); i++ {
		if bloqueCarpeta.B_content[i].B_inodo != -1 {
			continue
		}
		posicion := bitmap.WriteInBitmapBlock(path, superbloque)
		bloqueCarpeta.B_content[i].B_inodo = int32(posicion)

	}
}

func SearchFreeSpaceInBlock(bloqueCarpeta *datos.BloqueDeCarpetas) bool {
	for i := 0; i < len(bloqueCarpeta.B_content); i++ {
		if bloqueCarpeta.B_content[i].B_inodo != -1 {
			continue
		}
		return true
	}
	return false
}

func CompareDirectories(rutaPart string, bloqueCarpeta *datos.BloqueDeCarpetas) (int64, bool) {
	comparador := [64]byte{}
	copy(comparador[:], []byte(rutaPart))
	for _, content := range bloqueCarpeta.B_content {
		if bytes.Equal(content.B_name[:], comparador[:]) {
			return int64(content.B_inodo), true
		}
	}
	return -1, false
}

func strlen(arr [64]byte) int {
	count := 0
	for _, c := range arr {
		if c != 0 {
			count++
		}
	}
	return count
}

func StrlenBytes(arr []byte) int {
	count := 0
	for _, c := range arr {
		if c != 0 {
			count++
		}
	}
	return count
}

func trimArray(arr []byte) []byte {
	var result []byte
	for _, v := range arr {
		if v != 0 {
			result = append(result, v)
		}
	}
	return result
}
