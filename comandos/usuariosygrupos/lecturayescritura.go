package usuariosygrupos

import (
	"bytes"
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
			if strlen(parteArchivo.B_content) == 64 {
				continue
			} else if strlen(parteArchivo.B_content) < 64 {
				var nuevoContenido string
				nuevoContenido += string(parteArchivo.B_content[:])
				nuevoContenido += contenido
				if len(nuevoContenido) > 64 {
					copy(parteArchivo.B_content[:], nuevoContenido[:64])
					comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
					AppendFile(path, superbloque, tablaInodo, nuevoContenido[64:])
				} else {
					copy(parteArchivo.B_content[:], []byte(nuevoContenido))
					comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
				}
				return true
			}
		} else if tablaInodo.I_block[i] == -1 {
			var nuevoBloque datos.BloqueDeArchivos
			nuevaPosicion := bitmap.WriteInBitmapBlock(path, superbloque)
			tablaInodo.I_block[i] = nuevaPosicion
			if len(contenido) > 64 {
				copy(nuevoBloque.B_content[:], contenido[:64])
				comandos.Fwrite(&nuevoBloque, path, superbloque.S_block_start+nuevaPosicion*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
				AppendFile(path, superbloque, tablaInodo, contenido[64:])
			} else {
				copy(nuevoBloque.B_content[:], []byte(contenido))
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
		byteArr := parteArchivo.B_content[:]
		if bytes.Equal(byteArr, []byte("modificar")) {
			if len(contenido) > 64 {
				copy(parteArchivo.B_content[:], contenido[:64])
				comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))
				SetFile(tablaInodo, path, superbloque, contenido[64:])
			} else {
				copy(parteArchivo.B_content[:], contenido[:64])
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
