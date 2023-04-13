package usuariosygrupos

import (
	"bytes"
	"unsafe"

	"github.com/erazoex/proyecto2/bitmap"
	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/datos"
)

func ReadInode(tablaInodo *datos.TablaInodo, path string, superbloque *datos.SuperBloque) string {
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
				nuevoContenido := string(parteArchivo.B_content[:]) + contenido
				nuevoContenidoArray := createArray([]byte(nuevoContenido))
				if strlenBytes([]byte(nuevoContenido)) > 63 {
					copy(parteArchivo.B_content[:], nuevoContenidoArray[:63])
					comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
					AppendFile(path, superbloque, tablaInodo, string(nuevoContenidoArray[63:]))
				} else {
					copy(parteArchivo.B_content[:], nuevoContenidoArray)
					comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
				}
				return true
			}
		} else if tablaInodo.I_block[i] == -1 {
			var nuevoBloque datos.BloqueDeArchivos
			nuevaPosicion := bitmap.WriteInBitmapBlock(path, superbloque)
			nuevoContenido := createArray([]byte(contenido))
			tablaInodo.I_block[i] = nuevaPosicion
			if strlenBytes([]byte(contenido)) > 63 {
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

func GetFile(tablaInodo *datos.TablaInodo, path string, superbloque *datos.SuperBloque) string {
	contenido := ""
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo datos.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
		contenido += string(parteArchivo.B_content[:])
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
		if bytes.Equal(parteArchivo.B_content[:], []byte("modificar")) {
			nuevoContenido := createArray([]byte(contenido))
			if strlenBytes([]byte(contenido)) > 64 {
				copy(parteArchivo.B_content[:], nuevoContenido[:64])
				comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
				SetFile(tablaInodo, path, superbloque, string(nuevoContenido[64:]))
			} else {
				copy(parteArchivo.B_content[:], nuevoContenido)
				comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
			}
			return true
		}

	}
	return false
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

func strlenBytes(arr []byte) int {
	count := 0
	for _, c := range arr {
		if c != 0 {
			count++
		}
	}
	return count
}

func createArray(arr []byte) []byte {
	var result []byte
	for _, v := range arr {
		if v != 0 {
			result = append(result, v)
		}
	}
	return result
}
