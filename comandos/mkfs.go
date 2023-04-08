package comandos

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/lista"
)

type ParametrosMkfs struct {
	id string
	t  string
}

type Mkfs struct {
	params ParametrosMkfs
}

func (m Mkfs) Exe(parametros []string) {
	m.params = m.SaveParams(parametros)
	if m.Mkfs(m.params.id, m.params.t) {
		fmt.Printf("\nel formateo con EXT2 de la particion con id %s fue exitoso\n\n", m.params.id)
	} else {
		fmt.Printf("no se logro formatear la particion con id %s\n", m.params.id)
	}
}

func (m Mkfs) SaveParams(parametros []string) ParametrosMkfs {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		if strings.Contains(v, "id") {
			v = strings.ReplaceAll(v, "id=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.params.id = v
		} else if strings.Contains(v, "type") {
			v = strings.ReplaceAll(v, "type=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.params.t = v
		}
	}
	return m.params
}

func (m Mkfs) Mkfs(id string, t string) bool {
	// comprobando que id no este vacio
	if id == "" {
		fmt.Println("no se encontro el id entre los comandos")
		return false
	}
	// comprobando que type no lleve un valor incorrecto
	if t != "full" && t != "FULL" && t != "" {
		fmt.Println("el valor del comando type no es permitido")
		return false
	}
	if t == "" || t == "full" {
		t = "FULL"
	}
	//creando nuestro nodo auxiliar
	nodo := lista.ListaMount.GetNodeById(id)
	if nodo == nil {
		fmt.Printf("el id %s no coincide con ninguna particion montada\n", id)
		return false
	}
	m.Ext2(nodo)
	return true
}

func (m Mkfs) Ext2(nodo *lista.MountNode) {
	whereToStart := 0
	partSize := 0
	if nodo.Value != nil {
		whereToStart = int(nodo.Value.Part_start)
		partSize = int(nodo.Value.Part_size)
	} else if nodo.ValueL != nil {
		whereToStart = int(nodo.ValueL.Part_start)
		partSize = int(nodo.ValueL.Part_size)
	}
	n := ((partSize - int(unsafe.Sizeof(datos.SuperBloque{}))) / (4 + int(unsafe.Sizeof(datos.TablaInodo{})) + 3*int(unsafe.Sizeof(datos.BloqueDeArchivos{}))))
	if n < 1 {
		fmt.Println("el tamano de la particion es mas pequeno que el sistema de archivos")
		return
	}
	inodesQuantity := int64(n)
	blocksQuantity := int64(3 * inodesQuantity)

	// llenando la estructura del superbloque
	superBlock := datos.SuperBloque{
		S_filesystem_type:   2,
		S_inodes_count:      inodesQuantity,
		S_blocks_count:      blocksQuantity,
		S_free_inodes_count: inodesQuantity,
		S_free_blocks_count: blocksQuantity - 2,
		S_mnt_count:         0,
		S_magic:             0xEF53,
		S_inode_size:        int64(unsafe.Sizeof(datos.TablaInodo{})),
		S_block_size:        int64(unsafe.Sizeof(datos.BloqueDeArchivos{})),
		S_first_ino:         2,
		S_first_blo:         2,
	}
	superBlock.S_bm_inode_start = int64(whereToStart) + int64(unsafe.Sizeof(datos.SuperBloque{}))
	superBlock.S_bm_block_start = superBlock.S_bm_inode_start + inodesQuantity
	superBlock.S_inode_start = superBlock.S_bm_block_start + blocksQuantity
	superBlock.S_block_start = superBlock.S_inode_start + int64(unsafe.Sizeof(datos.TablaInodo{})*uintptr(inodesQuantity))
	date := time.Now()
	for i := 0; i < len(superBlock.S_mtime)-1; i++ {
		superBlock.S_mtime[i] = date.String()[i]
	}

	// escribiendo el superbloque
	Fwrite(&superBlock, nodo.Ruta, int64(whereToStart))

	// buffers para bloques e inodos
	inodos := make([]byte, inodesQuantity)
	bloques := make([]byte, blocksQuantity)

	// llenando los buffers
	for i := 0; i < len(inodos); i++ {
		inodos[i] = '0'
	}
	for i := 0; i < len(bloques); i++ {
		bloques[i] = '0'
	}

	// inodos ocupados
	inodos[0] = '1'
	inodos[1] = '1'
	Fwrite(&inodos, nodo.Ruta, superBlock.S_bm_inode_start)

	// bloques ocupados
	bloques[0] = '1'
	bloques[1] = '1'
	Fwrite(&bloques, nodo.Ruta, superBlock.S_bm_block_start)

	// crear tabla de inodos root
	rootInodeTable := datos.TablaInodo{
		I_uid:  1,
		I_gid:  1,
		I_size: 0,
		I_type: '0',
		I_perm: 664,
	}
	// llenando las fechas
	atime := time.Now()
	for i := 0; i < len(rootInodeTable.I_atime)-1; i++ {
		rootInodeTable.I_atime[i] = atime.String()[i]
	}
	ctime := time.Now()
	for i := 0; i < len(rootInodeTable.I_atime)-1; i++ {
		rootInodeTable.I_ctime[i] = ctime.String()[i]
	}
	mtime := time.Now()
	for i := 0; i < len(rootInodeTable.I_atime)-1; i++ {
		rootInodeTable.I_mtime[i] = mtime.String()[i]
	}
	// llenando a todos los bloques no utilizados
	for i := 0; i < len(rootInodeTable.I_block); i++ {
		rootInodeTable.I_block[i] = -1
	}
	// apuntando al bloque 0 (bloque de carpetas root)
	rootInodeTable.I_block[0] = 0

	// escribiendo la tabla de inodos root
	Fwrite(&rootInodeTable, nodo.Ruta, superBlock.S_inode_start)

	// creando el bloque de carpetas root
	bloqueCarpetasRoot := datos.BloqueDeCarpetas{}

	copy(bloqueCarpetasRoot.B_content[0].B_name[:], ".")
	bloqueCarpetasRoot.B_content[0].B_inodo = 0

	copy(bloqueCarpetasRoot.B_content[1].B_name[:], "..")
	bloqueCarpetasRoot.B_content[1].B_inodo = 0

	copy(bloqueCarpetasRoot.B_content[2].B_name[:], "users.txt")
	bloqueCarpetasRoot.B_content[2].B_inodo = 1

	copy(bloqueCarpetasRoot.B_content[3].B_name[:], "")
	bloqueCarpetasRoot.B_content[3].B_inodo = -1

	// llenando el archivo users.txt
	content := "1,G,root\n1,U,root,root,123\n"

	// crear tabla de inodos de archivo
	fileInodeTable := datos.TablaInodo{
		I_uid:  1,
		I_gid:  1,
		I_size: 0,
		I_type: '1',
		I_perm: 664,
	}
	// llenando las fechas
	atime = time.Now()
	for i := 0; i < len(fileInodeTable.I_atime)-1; i++ {
		fileInodeTable.I_atime[i] = atime.String()[i]
	}
	ctime = time.Now()
	for i := 0; i < len(fileInodeTable.I_atime)-1; i++ {
		fileInodeTable.I_ctime[i] = ctime.String()[i]
	}
	mtime = time.Now()
	for i := 0; i < len(fileInodeTable.I_atime)-1; i++ {
		fileInodeTable.I_mtime[i] = mtime.String()[i]
	}
	// llenando a todos los bloques no utilizados
	for i := 0; i < len(rootInodeTable.I_block); i++ {
		rootInodeTable.I_block[i] = -1
	}
	// apuntando al bloque 1 (primer bloque de archivos creado para users.txt)
	rootInodeTable.I_block[0] = 1

	// crear bloque de archivos y escribiendo el contenido
	bloqueArchivos := datos.BloqueDeArchivos{}
	copy(bloqueArchivos.B_content[:], []byte(content))

	// escribiendo el bloque de carpetas root
	Fwrite(&bloqueCarpetasRoot, nodo.Ruta, superBlock.S_block_start)

	// escribiendo la tabla de inodos del archivo users.txt
	Fwrite(&fileInodeTable, nodo.Ruta, superBlock.S_inode_start+int64(unsafe.Sizeof(datos.TablaInodo{})))

	// escribiendo el bloque 1 del archivo users.txt
	Fwrite(&bloqueArchivos, nodo.Ruta, superBlock.S_block_start+int64(unsafe.Sizeof(datos.BloqueDeArchivos{})))

	if nodo.Value != nil {
		// aqui deberia de ir un metodo para guardar para la consola
		fmt.Println("")
	} else if nodo.ValueL != nil {
		// aqui igual deberia de ir
		fmt.Println("")
	}
	fmt.Println("El formateo fue exitoso")
}
