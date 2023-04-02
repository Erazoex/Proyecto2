package main

// Contiene informacion sobre el disco
type MBR struct {
	mbr_tamano         []byte      // Tamano total del disco en bytes
	mbr_fecha_creacion []byte      // Fecha y hora de creacion del disco
	mbr_dsk_signature  []byte      // Numero random, que identifica de forma unica a cada disco
	dsk_fit            []byte      // Tipo de ajuste del disco. Tendra los valores B(Best), F(First) o W(Worst)
	mbr_partitions     []Partition // Estructura con la informacion de las particiones
}

// Contiene informacion sobre una particion primaria o extendida
type Partition struct {
	part_status []byte // Indica si la particion esta activa o no
	part_type   []byte // Indica el tipo de la particion, primaria o extendida. Tendra los valores P o E
	part_fit    []byte // Tipo de ajuste de la particion. Tendra los valores B(Best), F(First) o W(Worst)
	part_start  []byte // Indica en que byte del disco inicia la particion
	part_size   []byte // Contiene el tamano total de la particion en bytes
	part_name   []byte // Nombre de la particion
}

// Contiene informacion sobre una particion logica
type EBR struct {
	part_status []byte // Indica si la particion esta activa o no
	part_fit    []byte // Tipo de ajuste de la particion. Tendra los valores B(Best), F(First) o W(Worst)
	part_start  []byte // Indica en que byte del disco Inicia la particion
	part_size   []byte // Contiene el tamano total de la particion en bytes.
	part_next   []byte // Byte en el que esta el proximo EBR. -1 si no hay siguiente
	part_name   []byte // Nombre de la particion
}

// Contiene informacion sobre la configuracion del sistema de archivos
type SuperBloque struct {
	s_filesystem_type   []byte // Guarda el numero que identifica el sistema de archivos utilizado
	s_inodes_count      []byte // Guarda el numero total de inodos
	s_blocks_count      []byte // Guarda el numero total de bloques
	s_free_blocks_count []byte // Contiene el numero de bloques libres
	s_free_inodes_count []byte // Contiene el numero de indoos libres
	s_mtime             []byte // Ultima fecha en el que el sistema fue montado
	s_mnt_count         []byte // Indica cuantas veces se ha montado el sistema
	s_magic             []byte // Valor que identifica el sistema de archivos, tendra el valor 0xEF53
	s_inode_size        []byte // Tamano del inodo
	s_block_size        []byte // Tamano del bloque
	s_first_ino         []byte // Primer inodo libre
	s_first_blo         []byte // Primer bloque libre
	s_bm_inode_start    []byte // Guardara el inicio del bitmap de inodos
	s_bm_block_start    []byte // Guardara el inicio del bitmap de bloques
	s_inode_start       []byte // Guardara el inicio de la tabla de inodos
	s_block_start       []byte // Guardara el inicio de la tabla de bloques
}

// Contiene informacion sobre una tabla de inodo
type tablaInodo struct {
	i_uid   []byte // UID del usuario propietario del archivo o carpeta
	i_gid   []byte // GID del grupo al que pertenece el archivo o carpeta
	i_size  []byte // Tamano del archivo en bytes
	i_atime []byte // Ultima fecha en que se leyo el inodo sin modificarlo
	i_ctime []byte // Fecha en la que se creo el inodo
	i_mtime []byte // Ultima fecha en que se modifica el inodo
	i_block []byte // Array en los que los primeros 16 registros son bloques directos
	i_type  []byte // Indica si es archivo o carpeta. Tendra los siguientes valores: 1 = Archivo, 0 = Carpeta
	i_perm  []byte // Guardara los permisos del archivo o carpeta usando la nomenclatura de UGO
}

// Contiene informacion sobre un bloque de carpetas
type BloqueDeCarpetas struct {
	b_content [4]content // Array con el contenido de la carpeta
}

// Contiene informacion sobre el contenido de un archivo o carpeta
type content struct {
	b_name  []byte // Nombre de la carpeta o archivo
	b_inodo []byte // Apuntador hacia el inodo asociado al archivo o carpeta
}

// Contiene informacion sobre un bloque de archivos
type BloqueDeArchivos struct {
	b_content []byte // Array con el contenido del archivo con capacidad de solo 64 bytes
}
