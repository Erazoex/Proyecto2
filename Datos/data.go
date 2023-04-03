package datos

// Contiene informacion sobre el disco
type MBR struct {
	Mbr_tamano         []byte      // Tamano total del disco en bytes
	Mbr_fecha_creacion []byte      // Fecha y hora de creacion del disco
	Mbr_dsk_signature  []byte      // Numero random, que identifica de forma unica a cada disco
	Dsk_fit            []byte      // Tipo de ajuste del disco. Tendra los valores B(Best), F(First) o W(Worst)
	Mbr_partitions     []Partition // Estructura con la informacion de las particiones
}

// Contiene informacion sobre una particion primaria o extendida
type Partition struct {
	Part_status []byte // Indica si la particion esta activa o no
	Part_type   []byte // Indica el tipo de la particion, primaria o extendida. Tendra los valores P o E
	Part_fit    []byte // Tipo de ajuste de la particion. Tendra los valores B(Best), F(First) o W(Worst)
	Part_start  []byte // Indica en que byte del disco inicia la particion
	Part_size   []byte // Contiene el tamano total de la particion en bytes
	Part_name   []byte // Nombre de la particion
}

// Contiene informacion sobre una particion logica
type EBR struct {
	Part_status []byte // Indica si la particion esta activa o no
	Part_fit    []byte // Tipo de ajuste de la particion. Tendra los valores B(Best), F(First) o W(Worst)
	Part_start  []byte // Indica en que byte del disco Inicia la particion
	Part_size   []byte // Contiene el tamano total de la particion en bytes.
	Part_next   []byte // Byte en el que esta el proximo EBR. -1 si no hay siguiente
	Part_name   []byte // Nombre de la particion
}

// Contiene informacion sobre la configuracion del sistema de archivos
type SuperBloque struct {
	S_filesystem_type   []byte // Guarda el numero que identifica el sistema de archivos utilizado
	S_inodes_count      []byte // Guarda el numero total de inodos
	S_blocks_count      []byte // Guarda el numero total de bloques
	S_free_blocks_count []byte // Contiene el numero de bloques libres
	S_free_inodes_count []byte // Contiene el numero de indoos libres
	S_mtime             []byte // Ultima fecha en el que el sistema fue montado
	S_mnt_count         []byte // Indica cuantas veces se ha montado el sistema
	S_magic             []byte // Valor que identifica el sistema de archivos, tendra el valor 0xEF53
	S_inode_size        []byte // Tamano del inodo
	S_block_size        []byte // Tamano del bloque
	S_first_ino         []byte // Primer inodo libre
	S_first_blo         []byte // Primer bloque libre
	S_bm_inode_start    []byte // Guardara el inicio del bitmap de inodos
	S_bm_block_start    []byte // Guardara el inicio del bitmap de bloques
	S_inode_start       []byte // Guardara el inicio de la tabla de inodos
	S_block_start       []byte // Guardara el inicio de la tabla de bloques
}

// Contiene informacion sobre una tabla de inodo
type tablaInodo struct {
	I_uid   []byte // UID del usuario propietario del archivo o carpeta
	I_gid   []byte // GID del grupo al que pertenece el archivo o carpeta
	I_size  []byte // Tamano del archivo en bytes
	I_atime []byte // Ultima fecha en que se leyo el inodo sin modificarlo
	I_ctime []byte // Fecha en la que se creo el inodo
	I_mtime []byte // Ultima fecha en que se modifica el inodo
	I_block []byte // Array en los que los primeros 16 registros son bloques directos
	I_type  []byte // Indica si es archivo o carpeta. Tendra los siguientes valores: 1 = Archivo, 0 = Carpeta
	I_perm  []byte // Guardara los permisos del archivo o carpeta usando la nomenclatura de UGO
}

// Contiene informacion sobre un bloque de carpetas
type BloqueDeCarpetas struct {
	B_content [4]Content // Array con el contenido de la carpeta
}

// Contiene informacion sobre el contenido de un archivo o carpeta
type Content struct {
	B_name  []byte // Nombre de la carpeta o archivo
	B_inodo []byte // Apuntador hacia el inodo asociado al archivo o carpeta
}

// Contiene informacion sobre un bloque de archivos
type BloqueDeArchivos struct {
	B_content []byte // Array con el contenido del archivo con capacidad de solo 64 bytes
}
