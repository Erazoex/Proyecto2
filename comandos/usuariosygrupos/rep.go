package usuariosygrupos

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/erazoex/proyecto2/comandos"
	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/lista"
)

type ParametrosRep struct {
	Name string
	Path string
	Id   string
	Ruta string
}

type Rep struct {
	Params ParametrosRep
}

func (r *Rep) Exe(parametros []string) {
	r.Params = r.SaveParams(parametros)
	if r.Rep(r.Params.Name, r.Params.Path, r.Params.Id, r.Params.Ruta) {
		fmt.Printf("\nse creo el reporte de tipo %s para la ruta %s correctamente\n\n", r.Params.Name, r.Params.Path)
	} else {
		fmt.Printf("\nno se pudo crear el reporte de tipo %s para la ruta %s\n\n", r.Params.Name, r.Params.Path)
	}
}

func (r *Rep) SaveParams(parametros []string) ParametrosRep {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			r.Params.Name = v
		} else if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			r.Params.Path = v
		} else if strings.Contains(v, "id") {
			v = strings.ReplaceAll(v, "id=", "")
			r.Params.Id = v
		} else if strings.Contains(v, "ruta") {
			v = strings.ReplaceAll(v, "ruta=", "")
			r.Params.Ruta = v
		}
	}
	return r.Params
}

func (r *Rep) Rep(name, path, id, ruta string) bool {
	tiposDeReportes := []string{
		"disk",
		"tree",
		"file",
		"sb",
	}
	esValidoElReporte := false
	for _, reporte := range tiposDeReportes {
		if name == reporte {
			esValidoElReporte = true
		}
	}
	if !esValidoElReporte || name == "" {
		fmt.Println("el tipo de reporte no es valido")
		return false
	}
	if path == "" {
		fmt.Println("el path no puede ser vacio")
		return false
	}
	if !lista.ListaMount.NodeExist(id) {
		fmt.Printf("el id: %s, pertenece a una de las particiones montadas\n", id)
		return false
	}
	if name == "file" && ruta == "" {
		fmt.Println("la ruta a buscar no puede estar vacia")
		return false
	}
	path = strings.Split(path, ".")[0]
	if name == "disk" {
		fmt.Println("disk")
		r.ReporteDisk(path, id)
	} else if name == "tree" {
		fmt.Println("tree")
		r.ReporteTree(path, id)
	} else if name == "file" {
		fmt.Println("file")
		r.ReporteFile(path, id, ruta)
	} else if name == "sb" {
		fmt.Println("sb")
		r.ReporteSuperBloque(path, id)
	}
	return true
}

func (r *Rep) ReporteDisk(path, id string) {
	node := lista.ListaMount.GetNodeById(id)
	master := comandos.GetMBR(node.Ruta)
	tamano_master := master.Mbr_tamano
	contenidoLogicas := ""
	numeroDeLogicas := 0
	contenido := "digraph {\n"
	contenido += "\tnode [shape=plaintext]\n"
	contenido += "\ttable [label=<\n"
	contenido += "\t\t<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
	contenido += "\t\t<TR>\n"
	contenido += "\t\t\t<TD bgcolor=\"yellow\" ROWSPAN=\"2\"><BR/>MBR<BR/></TD>\n"
	existeExtendida := false
	for _, part := range master.Mbr_partitions {
		if part.Part_status == '5' {
			porcentaje := (float64(part.Part_size) / float64(tamano_master)) * 100
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"2\" COLSPAN=\"1\"><BR/>Libre<BR/>" + strconv.Itoa(int(porcentaje)) + "%</TD>\n"
		} else if part.Part_status == '0' {
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"2\" COLSPAN=\"1\"><BR/>Libre<BR/></TD>\n"
		} else if part.Part_type == 'e' || part.Part_type == 'E' {
			existeExtendida = true
			numeroDeLogicas = r.ContarParticiones(node.Ruta, part.Part_start)
			contenidoLogicas = r.RecorrerParticionesDISK(node.Ruta, part.Part_start, tamano_master)
			if numeroDeLogicas == 0 {
				contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"2\" COLSPAN=\"1\">Extendida</TD>\n"
			} else {
				contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"" + strconv.Itoa(2*numeroDeLogicas) + "\">Extendida</TD>\n"
			}
		} else {
			porcentaje := (float64(part.Part_size) / float64(tamano_master)) * 100
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"2\" COLSPAN=\"1\"><BR/>Primaria<BR/>" + strconv.Itoa(int(porcentaje)) + "%</TD>\n"
		}
	}
	contenido += "\t\t</TR>\n"
	if existeExtendida {
		contenido += "\t\t<TR>\n"
		contenido += contenidoLogicas
		contenido += "\t\t</TR>\n"
	}
	contenido += "\t\t</TABLE>\n"
	contenido += "\t>]\n"
	contenido += "}\n"
	directory := path + ".dot"
	// hay que crear los directorios el archivo nuevo
	comandos.MkDirectory(directory)
	comandos.Fopen(directory, contenido)
	// falta mandar el comando para convertirlo en pdf
}

func (r *Rep) ContarParticiones(ruta string, whereToStart int64) int {
	contador := 0
	var temp datos.EBR
	comandos.Fread(&temp, ruta, whereToStart)
	flag := true
	for flag {
		if temp.Part_size == 0 {
			flag = false
		} else if temp.Part_next != -1 {
			contador++
		} else if temp.Part_next == -1 {
			contador++
			flag = false
		}
		if temp.Part_next != -1 {
			comandos.Fread(&temp, ruta, temp.Part_next)
		}
	}
	return contador
}

func (r *Rep) RecorrerParticionesDISK(ruta string, whereToStart, tamano_master int64) string {
	contenido := ""
	var temp datos.EBR
	comandos.Fread(&temp, ruta, whereToStart)
	flag := true
	for flag {
		if temp.Part_size == 0 {
			flag = false
		} else if temp.Part_next != -1 && temp.Part_status != '5' {
			porcentaje := (float64(temp.Part_size) / float64(tamano_master)) * 100
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"1\"><BR/>EBR<BR/></TD>\n"
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"1\"><BR/>Logica<BR/>" + strconv.Itoa(int(porcentaje)) + "%</TD>\n"
		} else if temp.Part_next != -1 && temp.Part_status == '5' {
			porcentaje := (float64(temp.Part_size) / float64(tamano_master)) * 100
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"1\"><BR/>EBR<BR/></TD>\n"
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"1\"><BR/>Libre<BR/>" + strconv.Itoa(int(porcentaje)) + "%</TD>\n"
		} else if temp.Part_next == -1 && temp.Part_status == '5' {
			porcentaje := (float64(temp.Part_size) / float64(tamano_master)) * 100
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"1\"><BR/>EBR<BR/></TD>\n"
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"1\"><BR/>Libre<BR/>" + strconv.Itoa(int(porcentaje)) + "%</TD>\n"
			flag = false
		} else if temp.Part_next == -1 {
			porcentaje := (float64(temp.Part_size) / float64(tamano_master)) * 100
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"1\"><BR/>EBR<BR/></TD>\n"
			contenido += "\t\t\t<TD bgcolor=\"green\" ROWSPAN=\"1\" COLSPAN=\"1\"><BR/>Logica<BR/>" + strconv.Itoa(int(porcentaje)) + "%</TD>\n"
			flag = false
		}
		if temp.Part_next != -1 {
			comandos.Fread(&temp, ruta, temp.Part_next)
		}
	}
	return contenido
}

func (r *Rep) ReporteTree(path, id string) {

}

func (r *Rep) ReporteFile(path, id, ruta string) {
	node := lista.ListaMount.GetNodeById(id)
	var whereToStart int64
	if node.Value != nil {
		whereToStart = node.Value.Part_start
	} else if node.ValueL != nil {
		whereToStart = node.ValueL.Part_start + int64(unsafe.Sizeof(datos.EBR{}))
	}
	var superbloque datos.SuperBloque
	comandos.Fread(&superbloque, node.Ruta, whereToStart)
	archivo := ""
	// leeremos la primera tabla de inodos
	var tablaRoot datos.TablaInodo
	comandos.Fread(&tablaRoot, node.Ruta, superbloque.S_inode_start)
	ruta = strings.Replace(ruta, "/", "", 1)
	archivo += r.RecorrerArchivo(&tablaRoot, path, ruta, &superbloque)
	// ahora iniciaremos el archivo graphviz
	contenido := "digraph {\n"
	contenido += "\tnode [shape=plaintext]\n"
	contenido += "\tarchivo [label=\"" + archivo + "\"];\n"
	contenido += "}\n"

	directory := path + ".dot"
	// hay que crear los directorios el archivo nuevo
	comandos.MkDirectory(directory)
	comandos.Fopen(directory, contenido)
	// falta mandar el comando para convertirlo en pdf
}

func (r *Rep) RecorrerArchivo(tablaInodo *datos.TablaInodo, path, ruta string, superbloque *datos.SuperBloque) string {
	contenido := ""
	var rutaParts []string
	if !strings.Contains(ruta, "/") {
		// aqui deberiamos crear el metodo para recolectar el contenido del archivo
		for i := 0; i < len(tablaInodo.I_block); i++ {
			if tablaInodo.I_block[i] == -1 {
				break
			}
			var bloqueArchivos datos.BloqueDeArchivos
			comandos.Fread(&bloqueArchivos, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
			contenido += string(bloqueArchivos.B_content[:])
		}
		return contenido
	}
	rutaParts = strings.SplitN(ruta, "/", 2)
	fmt.Println(rutaParts)
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueCarpeta datos.BloqueDeCarpetas
		fmt.Println("posicion ->", tablaInodo.I_block[i])
		if tablaInodo.I_block[i] == -1 {
			break
		}
		comandos.Fread(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		num, compare := CompareDirectories(rutaParts[0], &bloqueCarpeta)
		if compare {
			var nuevaTablaInodo datos.TablaInodo
			comandos.Fread(&nuevaTablaInodo, path, superbloque.S_inode_start+num*superbloque.S_inode_size)
			contenido += r.RecorrerArchivo(&nuevaTablaInodo, path, rutaParts[1], superbloque)
		}
	}
	return contenido
}

func (r *Rep) ReporteSuperBloque(path, id string) {
	node := lista.ListaMount.GetNodeById(id)
	var whereToStart int64
	if node.Value != nil {
		whereToStart = node.Value.Part_start
	} else if node.ValueL != nil {
		whereToStart = node.ValueL.Part_start + int64(unsafe.Sizeof(datos.EBR{}))
	}
	var superbloque datos.SuperBloque
	comandos.Fread(&superbloque, node.Ruta, whereToStart)
	contenido := ""
	contenido += "\tnode [shape=plaintext]\n"
	contenido += "\ttable [label=<\n"
	contenido += "\t\t<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#1ECB23\" COLSPAN=\"2\"> Reporte de SUPERBLOQUE </TD></TR>\n"
	contenido += "\t\t\t<TR><TD> s_filesystem_type </TD><TD>" + strconv.Itoa(int(superbloque.S_filesystem_type)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_inodes_count </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_inodes_count)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD> s_blocks_count </TD><TD>" + strconv.Itoa(int(superbloque.S_blocks_count)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_free_blocks_count </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_free_blocks_count)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD> s_free_inodes_count </TD><TD>" + strconv.Itoa(int(superbloque.S_free_inodes_count)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_mtime </TD><TD bgcolor=\"#85F388\">" + string(superbloque.S_mtime[:]) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_mnt_count </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_mnt_count)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_magic </TD><TD bgcolor=\"#85F388\">" + strconv.FormatInt(superbloque.S_mnt_count, 16) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_inode_size </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_inode_size)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_block_size </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_block_size)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_first_ino </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_first_blo)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_first_blo </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_first_blo)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_bm_inode_start </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_bm_inode_start)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_bm_block_start </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_bm_block_start)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_inode_start </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_inode_start)) + "</TD></TR>\n"
	contenido += "\t\t\t<TR><TD bgcolor=\"#85F388\"> s_block_start </TD><TD bgcolor=\"#85F388\">" + strconv.Itoa(int(superbloque.S_block_start)) + "</TD></TR>\n"
	contenido += "\t\t</TABLE>\n"
	contenido += "\t>]\n"
	contenido += "}\n"
	directory := path + ".dot"
	// hay que crear los directorios el archivo nuevo
	comandos.MkDirectory(directory)
	comandos.Fopen(directory, contenido)
	// falta mandar el comando para convertirlo en pdf
}
