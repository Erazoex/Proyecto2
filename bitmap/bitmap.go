package bitmap

import (
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"

	"github.com/erazoex/proyecto2/datos"
)

func WriteInBitmapInode(path string, superbloque *datos.SuperBloque) int64 {
	file := OpenNewFile(path)
	valor := byte('1')
	position := superbloque.S_first_ino
	if position == -1 {
		fmt.Println("no se encontro posicion vacia, bitmap inode")
		return -1
	}
	file.Seek(superbloque.S_bm_inode_start+position*int64(unsafe.Sizeof(valor)), 0)
	FwriteByte(file, &valor)
	superbloque.S_first_ino = SearchFirstFreeBitmapInodePos(file, superbloque)
	superbloque.S_free_inodes_count--
	return position
}

func WriteInBitmapBlock(path string, superbloque *datos.SuperBloque) int64 {
	file := OpenNewFile(path)
	valor := byte('1')
	position := superbloque.S_first_blo
	if position == -1 {
		fmt.Println("no se encontro posicion vacia, bitmap block")
		return -1
	}
	file.Seek(superbloque.S_bm_block_start+position*int64(unsafe.Sizeof(valor)), 0)
	FwriteByte(file, &valor)
	superbloque.S_first_blo = SearchFirstFreeBitmapBlockPos(file, superbloque)
	superbloque.S_free_blocks_count--
	return position
}

// funciones para borrar en los bitmap de bloques

func DeleteBitmapInode(path string, superbloque *datos.SuperBloque, posicion int64) {
	file := OpenNewFile(path)
	valor := byte('0')
	file.Seek(superbloque.S_bm_inode_start+(posicion*int64(unsafe.Sizeof(valor))), 0)
	FwriteByte(file, &valor)
	superbloque.S_first_ino = SearchFirstFreeBitmapInodePos(file, superbloque)
	superbloque.S_free_inodes_count++
}

func DeleteBitmapBlock(path string, superbloque *datos.SuperBloque, posicion int64) {
	file := OpenNewFile(path)
	valor := byte('0')
	file.Seek(superbloque.S_bm_block_start+(posicion*int64(unsafe.Sizeof(valor))), 0)
	FwriteByte(file, &valor)
	superbloque.S_first_blo = SearchFirstFreeBitmapBlockPos(file, superbloque)
	superbloque.S_free_blocks_count++
}

// buscar primer bit libre en los bitmaps

func SearchFirstFreeBitmapInodePos(file *os.File, superbloque *datos.SuperBloque) int64 {
	contar := 0
	flag := true
	for flag && contar < int(superbloque.S_inodes_count) {
		i := byte('0')
		FreadByte(file, &i)
		if i == '0' {
			return int64(contar)
		}
		contar++
	}
	return -1
}

func SearchFirstFreeBitmapBlockPos(file *os.File, superbloque *datos.SuperBloque) int64 {
	contar := 0
	flag := true
	for flag && contar < int(superbloque.S_blocks_count) {
		i := byte('0')
		FreadByte(file, &i)
		if i == '0' {
			return int64(contar)
		}
		contar++
	}
	return -1
}

// leer un byte en archivo

func FreadByte(file *os.File, temp *byte) {
	err := binary.Read(file, binary.LittleEndian, temp)
	if err != nil {
		fmt.Println("no se pudo leer,", err.Error())
	}
}

// escribir un byte en archivo

func FwriteByte(file *os.File, temp *byte) {
	err := binary.Write(file, binary.LittleEndian, temp)
	if err != nil {
		fmt.Println("no se pudo escribir,", err.Error())
	}
}

// abrir el archivo
func OpenNewFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("no se pudo abrir el archivo para Bitmap", err.Error())
		return nil
	}
	return file
}
