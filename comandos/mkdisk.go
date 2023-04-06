package comandos

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/erazoex/proyecto2/datos"
)

type ParametrosMkdisk struct {
	size int
	fit  byte
	unit byte
	path string
}

type Mkdisk struct {
	params ParametrosMkdisk
}

func (m Mkdisk) Exe(parametros []string) {
	m.params = m.SaveParams(parametros)
	if m.Mkdisk(m.params.size, m.params.fit, m.params.unit, m.params.path) {
		fmt.Printf("\nmkdisk realizado con exito para la ruta: %s\n\n", m.params.path)
	} else {
		fmt.Printf("\n[ERROR!] no se logro realizar el comando mkdisk para la ruta: %s\n\n", m.params.path)
	}
}

func (m Mkdisk) SaveParams(parametros []string) ParametrosMkdisk {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			v = strings.ReplaceAll(v, "\"", "")
			m.params.path = v
		} else if strings.Contains(v, "size") {
			v = strings.ReplaceAll(v, "size=", "")
			v = strings.ReplaceAll(v, " ", "")
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("hubo un error al convertir a int", err.Error())
			}
			m.params.size = num
		} else if strings.Contains(v, "unit") {
			v = strings.ReplaceAll(v, "unit=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.params.unit = v[0]
		} else if strings.Contains(v, "fit") {
			v = strings.ReplaceAll(v, "fit=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.params.fit = v[0]
		}
	}
	return m.params
}

func (m Mkdisk) Mkdisk(size int, fit byte, unit byte, path string) bool {
	var fileSize = 0
	var master datos.MBR
	// Comprobando si existe una ruta valida para la creacion del disco
	if path == "" {
		fmt.Println("no se encontro una ruta")
		return false
	}
	// comprobando el tamano del disco, debe ser mayor que cero
	if size <= 0 {
		fmt.Println("el tamano del disco debe ser mayor que 0")
		return false
	}
	// tipo de unidad a utilizar, si el parametro esta vacio se utilizaran MegaBytes como default size
	if unit == 'k' || unit == 'K' {
		fileSize = size
	} else if unit == 'm' || unit == 'M' {
		fileSize = size * 1024
	} else if unit == ' ' {
		fileSize = size * 1024
	} else {
		fmt.Println("se debe ingresar una letra que corresponda un tamano valido")
		return false
	}
	// definiendo el tipo de fit que el disco tendra, como default sera First Fit
	if string(fit) == "bf" || string(fit) == "BF" {
		master.Dsk_fit = 'b'
	} else if string(fit) == "ff" || string(fit) == "FF" {
		master.Dsk_fit = 'f'
	} else if string(fit) == "wf" || string(fit) == "WF" {
		master.Dsk_fit = 'w'
	} else if string(fit) == "" {
		master.Dsk_fit = 'f'
	} else {
		fmt.Println("se debe ingresar un tipo de fit valido")
		return false
	}
	// llenando el buffer con '0' para indicar que esta vacio.
	bloque := make([]byte, 1024)
	for i := 0; i < len(bloque); i++ {
		bloque[i] = 0
	}

	iterator := 0
	MkDirectory(path) // creando el directorio para el disco sino existe
	binaryFile, err := os.Create(path)
	if err != nil {
		fmt.Println("error al crear el disco")
		return false
	}
	defer binaryFile.Close()
	for iterator < fileSize {
		_, err := binaryFile.Write(bloque[:])
		if err != nil {
			fmt.Println("error al llenar el disco creado")
		}
		iterator++
	}
	master.Mbr_tamano = int64(fileSize * 1024)
	master.Mbr_dsk_signature = GetRandom()
	// formateando el tiempo
	date := time.Now()
	for i := 0; i < len(master.Mbr_fecha_creacion)-1; i++ {
		master.Mbr_fecha_creacion[i] = date.String()[i]
	}
	FillPartitions(&master)
	WriteMBR(&master, path)
	return true
}

func FillPartitions(master *datos.MBR) {
	for _, v := range master.Mbr_partitions {
		v.Part_status = '0'
		v.Part_fit = '0'
		v.Part_start = -1
		v.Part_size = 0
		v.Part_type = '0'
		copy(v.Part_name[:], "0")
	}
}
