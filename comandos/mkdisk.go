package comandos

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/erazoex/proyecto2/consola"
	"github.com/erazoex/proyecto2/datos"
)

type ParametrosMkdisk struct {
	Size int
	Fit  byte
	Unit byte
	Path string
}

type Mkdisk struct {
	Params ParametrosMkdisk
}

func (m *Mkdisk) Exe(parametros []string) {
	m.Params = m.SaveParams(parametros)
	if m.Mkdisk(m.Params.Size, m.Params.Fit, m.Params.Unit, m.Params.Path) {
		consola.AddToConsole(fmt.Sprintf("\nmkdisk realizado con exito para la ruta: %s\n\n", m.Params.Path))
	} else {
		consola.AddToConsole(fmt.Sprintf("\n[ERROR!] no se logro realizar el comando mkdisk para la ruta: %s\n\n", m.Params.Path))
	}
}

func (m *Mkdisk) SaveParams(parametros []string) ParametrosMkdisk {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			v = strings.ReplaceAll(v, "\"", "")
			m.Params.Path = v
		} else if strings.Contains(v, "size") {
			v = strings.ReplaceAll(v, "size=", "")
			v = strings.ReplaceAll(v, " ", "")
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("hubo un error al convertir a int", err.Error())
			}
			m.Params.Size = num
		} else if strings.Contains(v, "unit") {
			v = strings.ReplaceAll(v, "unit=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.Params.Unit = v[0]
		} else if strings.Contains(v, "fit") {
			v = strings.ReplaceAll(v, "fit=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.Params.Fit = v[0]
		}
	}
	return m.Params
}

func (m *Mkdisk) Mkdisk(size int, fit byte, unit byte, path string) bool {
	var fileSize = 0
	var master datos.MBR
	// Comprobando si existe una ruta valida para la creacion del disco
	if path == "" {
		consola.AddToConsole("no se encontro una ruta\n")
		return false
	}
	// comprobando el tamano del disco, debe ser mayor que cero
	if size <= 0 {
		consola.AddToConsole("el tamano del disco debe ser mayor que 0\n")
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
		consola.AddToConsole("se debe ingresar una letra que corresponda un tamano valido\n")
		return false
	}
	// definiendo el tipo de fit que el disco tendra, como default sera First Fit
	// fmt.Printf("tipo de la variable fit %T\n", fit)
	// fmt.Println("el fit es:", fit)
	if string(fit) == "bf" || string(fit) == "BF" {
		master.Dsk_fit = 'b'
	} else if string(fit) == "ff" || string(fit) == "FF" {
		master.Dsk_fit = 'f'
	} else if string(fit) == "wf" || string(fit) == "WF" {
		master.Dsk_fit = 'w'
	} else if fit == 0 {
		master.Dsk_fit = 'f'
	} else {
		consola.AddToConsole("se debe ingresar un tipo de fit valido\n")
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
		consola.AddToConsole("error al crear el disco\n")
		return false
	}
	defer binaryFile.Close()
	for iterator < fileSize {
		_, err := binaryFile.Write(bloque[:])
		if err != nil {
			consola.AddToConsole("error al llenar el disco creado\n")
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
	for i := 0; i < len(master.Mbr_partitions); i++ {
		master.Mbr_partitions[i].Part_status = '0'
		master.Mbr_partitions[i].Part_fit = '0'
		master.Mbr_partitions[i].Part_start = 0
		master.Mbr_partitions[i].Part_size = 0
		master.Mbr_partitions[i].Part_type = '0'
		copy(master.Mbr_partitions[i].Part_name[:], "")
	}
}
