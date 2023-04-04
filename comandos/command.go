package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/erazoex/proyecto2/datos"
)

// En este paquete se encuentran las funciones en comun
// que se utilizan entre la mayoria de comandos.
func WriteMBR(master *datos.MBR, path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Posicionandonos en el principio del archivo
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Escribiendo el MBR
	var masterBuffer bytes.Buffer
	binary.Write(&masterBuffer, binary.BigEndian, master)
	_, err = file.Write(masterBuffer.Bytes())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()
}

func GetMBR(path string) datos.MBR {
	var mbr datos.MBR
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("no se pudo abrir el archivo para obtener el MBR", err.Error())
		return mbr
	}

	defer file.Close()

	// leyendo el mbr del archivo
	file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Println("no se pudo obtener la informacion del archivo para obtener el MBR", err.Error())
		return mbr
	}
	return mbr
}

func MkDir(fullPath string) {
	directory := path.Dir(fullPath)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0777)
		if err != nil {
			fmt.Println("no se pudo crear el directorio", err.Error())
		}
	}
}

func GetRandom() int64 {
	rand.Seed(time.Now().UnixNano())
	n := 150
	randomNumber := rand.Intn(n)
	return int64(randomNumber)
}