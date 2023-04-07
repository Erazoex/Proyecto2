package comandos

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/erazoex/proyecto2/datos"
)

type ParametrosFdisk struct {
	size   int
	unit   byte
	path   string
	type_p byte
	fit    byte
	name   [16]byte
}

type Fdisk struct {
	params ParametrosFdisk
}

func (f Fdisk) Exe(parametros []string) {
	f.params = f.SaveParams(parametros)
	if f.Fdisk(f.params.name, f.params.path, f.params.size, f.params.unit, f.params.fit, f.params.type_p) {
		fmt.Printf("\nfdisk realizado con exito para la ruta: %s\n\n", f.params.path)
	} else {
		fmt.Printf("\n[ERROR!] no se logro realizar el comando fdisk para la ruta: %s\n\n", f.params.path)
	}
}

func (f Fdisk) SaveParams(parametros []string) ParametrosFdisk {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			v = strings.ReplaceAll(v, "\"", "")
			f.params.path = v
		} else if strings.Contains(v, "size") {
			v = strings.ReplaceAll(v, "size=", "")
			v = strings.ReplaceAll(v, " ", "")
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("hubo un error al convertir a int", err.Error())
			}
			f.params.size = num
		} else if strings.Contains(v, "unit") {
			v = strings.ReplaceAll(v, "unit=", "")
			v = strings.ReplaceAll(v, " ", "")
			if v == "" {
				f.params.unit = ' '
			} else {
				f.params.unit = v[0]
			}
		} else if strings.Contains(v, "fit") {
			v = strings.ReplaceAll(v, "fit=", "")
			v = strings.ReplaceAll(v, " ", "")
			if v == "" {
				f.params.fit = ' '
			} else {
				f.params.fit = v[0]
			}

		} else if strings.Contains(v, "type") {
			v = strings.ReplaceAll(v, "type=", "")
			v = strings.ReplaceAll(v, " ", "")
			if v == "" {
				f.params.type_p = ' '
			} else {
				f.params.type_p = v[0]
			}
		} else if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			v = v[:16]
			copy(f.params.name[:], v)
		}
	}
	return f.params
}

func (f Fdisk) Fdisk(name [16]byte, path string, size int, unit byte, fit byte, t byte) bool {
	if path == "" {
		fmt.Println("no se encontro una ruta")
		return false
	}
	master := GetMBR(path)
	newPartition := datos.Partition{}
	fileSize := 0
	// tipo de unidad a utilizar, si el parametro esta vacio se utilizaran Kilobytes como default size.
	if unit == 'b' || unit == 'B' {
		fileSize = size
	} else if unit == 'k' || unit == 'K' {
		fileSize = size * 1024
	} else if unit == 'm' || unit == 'M' {
		fileSize = size * 1024 * 1024
	} else if unit == ' ' {
		fileSize = size * 1024
	} else {
		fmt.Println("debe ingresar una letra de tamano correcta")
		return false
	}
	// se debe comprobar que no exista ninguna particion con el mismo nombre
	if ExisteParticion(&master, name) {
		fmt.Printf("ya existe una particion con nombre: \"%v\"\n", name)
		return false
	}
	// comprobando el tamano de la particion, este debe ser mayor que cero
	if size <= 0 {
		fmt.Println("el tamano de la particion tiene que ser mayor a 0")
		return false
	}
	// definiendo el tipo de fit que la particion tendra, como default se utilizara Worst Fit
	if fit == 'w' || fit == 'W' {
		newPartition.Part_fit = 'w'
	} else if fit == 'b' || fit == 'B' {
		newPartition.Part_fit = 'b'
	} else if fit == 'f' || fit == 'F' {
		newPartition.Part_fit = 'f'
	} else if fit == ' ' {
		newPartition.Part_fit = 'w'
	} else {
		fmt.Println("se debe ingresar un tipo de fit valido")
		return false
	}

	// verificando que el tamano de la particion a crear sea menor
	// o igual que el tamano que queda en el disco.
	totalSize := unsafe.Sizeof(datos.MBR{})
	for _, v := range master.Mbr_partitions {
		if v.Part_status == '1' {
			totalSize += uintptr(v.Part_start)
		}
	}
	if t != 'l' && t != 'L' {
		if fileSize > int(master.Mbr_tamano)-int(totalSize) {
			fmt.Println("el tamano de la particion es mas grande que el disco")
			return false
		}
	}

	// indicando el tipo de particion
	if t == ' ' {
		t = 'p'
	} else if t != 'p' && t != 'e' && t != 'l' && t != 'P' && t != 'E' && t != 'L' {
		fmt.Printf("el tipo de la particion no es valido: \"%c\"\n", t)
		return false
	}
	newPartition.Part_size = int64(fileSize)
	newPartition.Part_type = t
	newPartition.Part_status = '1'
	copy(newPartition.Part_name[:], name[:])

	// revisando que no exista mas de una particion Extendida y que Exista en caso de que se vaya a crear una particion logica
	existeParticionExtendida := false //esta variable se utiliza para encontrar si existe una particion extendida
	var whereToStart int              // con este valor le vamos a pasar a la particion logica donde comienza la particion extendida
	var partitionSize int             // con este valor le indicamos a la particion logica cuanto espacio ocupa la particion extendida
	var extendedFit byte              // con este valor le indicamos a la particion logica el tipo de ajuste que tiene la particion extendida
	var extendedName [16]byte         // con este valor le indicamos el nombre de la particion extendida a la particion logica
	// aqui le agregamos a las variables anteriores su correspondiente valor
	for _, v := range master.Mbr_partitions {
		if v.Part_type == 'e' || v.Part_type == 'E' {
			copy(extendedName[:], v.Part_name[:])
			existeParticionExtendida = true
			extendedFit = v.Part_fit
			whereToStart = int(v.Part_start)
			partitionSize = int(v.Part_size)
		}
	}

	// comprobamos que exista una particion libre
	existeParticionLibre := false
	if t != 'l' && t != 'L' {
		for _, v := range master.Mbr_partitions {
			if v.Part_status == '0' {
				existeParticionLibre = true
			}
		}
	} else if t == 'l' || t == 'L' {
		existeParticionLibre = true
	}
	// sino se encuentra un espacio libre para particion
	if !existeParticionLibre {
		fmt.Println("no se encuentra ninguna particion libre para crear dentro del disco")
		return false
	}
	// comprobamos que tipo de particion es, luego la creamos
	if t == 'p' || t == 'P' {
		f.CreatePrimaryPartition(&master, newPartition)
	} else if t == 'e' || t == 'E' {
		if existeParticionExtendida {
			fmt.Println("no puede haber mas de una particion extendida")
			return false
		}
		f.CreateExtendedPartition(&master, newPartition, path)
	} else if t == 'l' || t == 'L' {
		if !existeParticionExtendida {
			fmt.Println("no existe una particion extendida para crear una particion logica")
			return false
		}
		particionLogica := datos.EBR{}
		particionLogica.Part_fit = newPartition.Part_fit
		particionLogica.Part_next = -1
		particionLogica.Part_size = newPartition.Part_size
		particionLogica.Part_status = newPartition.Part_status
		copy(particionLogica.Part_name[:], newPartition.Part_name[:])
		// vamos a mandar que tipo de ajuste tiene la particion
		// dentro de este metodo se le indica donde es que comienza la particion logica
		return f.CreateLogicPartition(&particionLogica, path, whereToStart, partitionSize, extendedFit, extendedName)

	}
	WriteMBR(&master, path)
	return true
}

func (f Fdisk) CreatePrimaryPartition(master *datos.MBR, newPartition datos.Partition) {
	// Asignacion de que particion es la que se utilizara
	if master.Dsk_fit == 'b' {
		BestFit(master, &newPartition)
	} else if master.Dsk_fit == 'w' {
		WorstFit(master, &newPartition)
	} else if master.Dsk_fit == 'f' {
		FirstFit(master, &newPartition)
	}
}

func (f Fdisk) CreateExtendedPartition(master *datos.MBR, newPartition datos.Partition, path string) {
	// Asignacion de que particion es la que se utilizara
	if master.Dsk_fit == 'b' {
		BestFit(master, &newPartition)
	} else if master.Dsk_fit == 'w' {
		WorstFit(master, &newPartition)
	} else if master.Dsk_fit == 'f' {
		FirstFit(master, &newPartition)
	}
	temp := datos.EBR{}
	temp.Part_status = '0'
	temp.Part_fit = '0'
	temp.Part_start = newPartition.Part_start
	temp.Part_size = 0
	temp.Part_next = -1
	copy(temp.Part_name[:], "")
	WriteEBR(&temp, path, newPartition.Part_start)
}

func BestFit(master *datos.MBR, newPartition *datos.Partition) {
	bestFit := 0
	// para encontrar el mejor fit lo primero que hay que hacer
	// es recorrer la lista de particiones y verificar que exista
	// una particion disponible, si esta particion se encuentra
	// disponible se comprobara que el tamano de esta sea mayor o
	// igual al tamano de la particion que estamos creando. de ser
	// asi, le asignaremos esa posicion, luego seguir iterando para
	// buscar si existe alguna particion con menor cantidad de espacio
	// donde se ajuste la particion que estamos creando.
	encontroParticion := false
	for i, v := range master.Mbr_partitions {
		if v.Part_status == '5' && v.Part_size >= newPartition.Part_size {
			if i != bestFit {
				if v.Part_size < master.Mbr_partitions[bestFit].Part_size {
					encontroParticion = true
					bestFit = i
				}
			}
		}
	}
	if !encontroParticion {
		for i, v := range master.Mbr_partitions {
			if v.Part_start == -1 {
				bestFit = i
				break
			}
		}
	}
	master.Mbr_partitions[bestFit] = *newPartition
	if bestFit == 0 {
		master.Mbr_partitions[bestFit].Part_start = int64(unsafe.Sizeof(datos.MBR{}))

	} else {
		master.Mbr_partitions[bestFit].Part_start = master.Mbr_partitions[bestFit-1].Part_start + master.Mbr_partitions[bestFit-1].Part_size
	}
}

func WorstFit(master *datos.MBR, newPartition *datos.Partition) {
	worstFit := 0
	encontroParticion := false
	for i, v := range master.Mbr_partitions {
		if v.Part_status == '5' && v.Part_size >= newPartition.Part_size {
			if i != worstFit {
				if v.Part_size > master.Mbr_partitions[worstFit].Part_size {
					worstFit = i
					encontroParticion = true
				}
			}
		}
	}
	if !encontroParticion {
		for i, v := range master.Mbr_partitions {
			if v.Part_start == -1 {
				worstFit = i
				break
			}
		}
	}
	master.Mbr_partitions[worstFit] = *newPartition
	if worstFit == 0 {
		master.Mbr_partitions[worstFit].Part_start = int64(unsafe.Sizeof(datos.MBR{}))
	} else {
		master.Mbr_partitions[worstFit].Part_start = master.Mbr_partitions[worstFit-1].Part_start + master.Mbr_partitions[worstFit-1].Part_size
	}
}

func FirstFit(master *datos.MBR, newPartition *datos.Partition) {
	firstFit := 0
	for i, v := range master.Mbr_partitions {
		if v.Part_status == '5' && v.Part_size >= newPartition.Part_size || v.Part_start == -1 {
			firstFit = i
			break
		}
	}
	master.Mbr_partitions[firstFit] = *newPartition
	if firstFit == 0 {
		master.Mbr_partitions[firstFit].Part_start = int64(unsafe.Sizeof(datos.MBR{}))
	} else {
		master.Mbr_partitions[firstFit].Part_start = master.Mbr_partitions[firstFit-1].Part_start + master.Mbr_partitions[firstFit-1].Part_size
	}
}

func (f Fdisk) CreateLogicPartition(logicPartition *datos.EBR, path string, whereToStart int, partitionSize int, extendedFit byte, extendedName [16]byte) bool {
	if extendedFit == 'f' {
		return FirstFitLogicPart(logicPartition, path, whereToStart, partitionSize, extendedName)
	} else if extendedFit == 'b' {
		return BestFitLogicPart(logicPartition, path, whereToStart, partitionSize, extendedName)
	} else if extendedFit == 'w' {
		return WorstFitLogicPart(logicPartition, path, whereToStart, partitionSize, extendedName)
	}
	return false
}

func FirstFitLogicPart(logicPartition *datos.EBR, path string, whereToStart int, partitionSize int, extendedName [16]byte) bool {
	temp := datos.EBR{}
	totalSize := 0
	totalSize += int(logicPartition.Part_size)
	temp = ReadEBR(path, int64(whereToStart))
	flag := true
	for flag {
		if temp.Part_size == 0 {
			if partitionSize < int(logicPartition.Part_size) {
				fmt.Println("la particion logica es mas grande que la extendida")
				return false
			}
			logicPartition.Part_start = int64(whereToStart)
			WriteEBR(logicPartition, path, int64(whereToStart))
			flag = false
		} else if temp.Part_status == '5' {
			if temp.Part_size >= logicPartition.Part_size {
				logicPartition.Part_start = temp.Part_start
				logicPartition.Part_next = temp.Part_next
				WriteEBR(logicPartition, path, temp.Part_start)
				flag = false
			}
		} else if temp.Part_next == -1 {
			totalSize += int(temp.Part_size)
			if partitionSize < totalSize {
				fmt.Println("el tamano de todas las particiones logicas unidas son mas grandes que la particion extendida, espacio insuficiente")
				return false
			}
			temp.Part_next = temp.Part_start + temp.Part_size
			logicPartition.Part_start = temp.Part_next
			WriteEBR(&temp, path, temp.Part_start)
			WriteEBR(logicPartition, path, temp.Part_next)
			flag = false
		} else {
			totalSize += int(temp.Part_size)
			temp = ReadEBR(path, temp.Part_next)
		}
	}
	// aqui deberia ir un print a la consola
	return true
}

func BestFitLogicPart(logicPartition *datos.EBR, path string, whereToStart int, partitionSize int, extededName [16]byte) bool {
	var particionesLogicas []datos.EBR
	temp := datos.EBR{}
	totalSize := 0
	totalSize += int(logicPartition.Part_size)
	temp = ReadEBR(path, int64(whereToStart))
	Wrote := false
	flag := true
	for flag {
		if temp.Part_size == 0 {
			if partitionSize < int(logicPartition.Part_size) {
				fmt.Println("la particion logica es mas grande que la extendida")
				return false
			}
			logicPartition.Part_start = int64(whereToStart)
			WriteEBR(logicPartition, path, int64(whereToStart))
			flag = false
			Wrote = true
		} else if temp.Part_status == '5' {
			particionesLogicas = append(particionesLogicas, temp)
		} else if temp.Part_next == -1 {
			flag = false
		} else {
			totalSize += int(temp.Part_size)
			temp = ReadEBR(path, temp.Part_next)
		}
	}
	bestFit := 0
	tempSize := 0
	if len(particionesLogicas) != 0 {
		for i, v := range particionesLogicas {
			if tempSize != 0 {
				bestFit = i
			} else if tempSize > int(v.Part_size) && v.Part_size >= logicPartition.Part_size {
				tempSize = int(v.Part_size)
				bestFit = i
			}
		}
		logicPartition.Part_start = particionesLogicas[bestFit].Part_start
		logicPartition.Part_next = particionesLogicas[bestFit].Part_next
		WriteEBR(logicPartition, path, logicPartition.Part_start)
		Wrote = true
	}
	if !Wrote {
		totalSize = int(logicPartition.Part_size)
		temp = ReadEBR(path, int64(whereToStart))
		flag2 := true
		for flag2 {
			if temp.Part_next == -1 {
				totalSize += int(temp.Part_size)
				if partitionSize < totalSize {
					fmt.Println("el tamano de todas las particiones logicas unidas son mas grandes que la particion extendida, espacio insuficiente")
					return false
				}
				temp.Part_next = temp.Part_start + temp.Part_size
				logicPartition.Part_start = temp.Part_next
				WriteEBR(&temp, path, temp.Part_start)
				WriteEBR(logicPartition, path, temp.Part_next)
				flag2 = false
			} else {
				totalSize += int(temp.Part_size)
				temp = ReadEBR(path, temp.Part_next)
			}
		}
	}
	// aqui deberia ir un print a la consola
	return true
}

func WorstFitLogicPart(logicPartition *datos.EBR, path string, whereToStart int, partitionSize int, extededName [16]byte) bool {
	var particionesLogicas []datos.EBR
	temp := datos.EBR{}
	totalSize := 0
	totalSize += int(logicPartition.Part_size)
	temp = ReadEBR(path, int64(whereToStart))
	Wrote := false
	flag := true
	for flag {
		if temp.Part_size == 0 {
			if partitionSize < int(logicPartition.Part_size) {
				fmt.Println("la particion logica es mas grande que la extendida")
				return false
			}
			logicPartition.Part_start = int64(whereToStart)
			WriteEBR(logicPartition, path, int64(whereToStart))
			flag = false
			Wrote = true
		} else if temp.Part_status == '5' {
			particionesLogicas = append(particionesLogicas, temp)
		} else if temp.Part_next == -1 {
			flag = false
		} else {
			totalSize += int(temp.Part_size)
			temp = ReadEBR(path, temp.Part_next)
		}
	}
	worstFit := 0
	tempSize := 0
	if len(particionesLogicas) != 0 {
		for i, v := range particionesLogicas {
			if tempSize != 0 {
				worstFit = i
			} else if tempSize < int(v.Part_size) && v.Part_size >= logicPartition.Part_size {
				tempSize = int(v.Part_size)
				worstFit = i
			}
		}
		logicPartition.Part_start = particionesLogicas[worstFit].Part_start
		logicPartition.Part_next = particionesLogicas[worstFit].Part_next
		WriteEBR(logicPartition, path, logicPartition.Part_start)
		Wrote = true
	}
	if !Wrote {
		totalSize = int(logicPartition.Part_size)
		temp = ReadEBR(path, int64(whereToStart))
		flag2 := true
		for flag2 {
			if temp.Part_next == -1 {
				totalSize += int(temp.Part_size)
				if partitionSize < totalSize {
					fmt.Println("el tamano de todas las particiones logicas unidas son mas grandes que la particion extendida, espacio insuficiente")
					return false
				}
				temp.Part_next = temp.Part_start + temp.Part_size
				logicPartition.Part_start = temp.Part_next
				WriteEBR(&temp, path, temp.Part_start)
				WriteEBR(logicPartition, path, temp.Part_next)
				flag2 = false
			} else {
				totalSize += int(temp.Part_size)
				temp = ReadEBR(path, temp.Part_next)
			}
		}
	}
	// aqui deberia ir un print a la consola
	return true
}

func ExisteParticion(master *datos.MBR, name [16]byte) bool {
	for _, v := range master.Mbr_partitions {
		if bytes.Equal(v.Part_name[:], name[:]) {
			return true
		}
	}
	return false
}
