package lista

import (
	"strconv"
	"strings"

	"github.com/erazoex/proyecto2/datos"
)

type MountNode struct {
	Key, Ruta   string
	Digits, Pos int
	Value       *datos.Partition
	ValueL      *datos.EBR
	Next, Prev  *MountNode
}

func (m MountNode) MountNode(ruta string, digits int, pos int, value *datos.Partition, valueL *datos.EBR) {
	m.Ruta = ruta
	m.Digits = digits
	m.Pos = pos
	m.Value = value
	m.ValueL = valueL
	m.Next = nil
	m.Prev = nil
	m.CreateKey()
}

func (m MountNode) CreateKey() {
	directory := strings.Split(m.Ruta, "/")
	lastPart := directory[len(directory)-1]
	fileNameParts := strings.Split(lastPart, ".")
	filename := fileNameParts[0]
	m.Key = strconv.Itoa(m.Digits) + strconv.Itoa(m.Pos) + filename
}

type MountList struct {
	first, last *MountNode
	tamano      int
}

func (m MountList) IsEmpty() bool {
	return m.first == nil
}

func (m MountList) Mount(path string, digit int, part *datos.Partition, partL *datos.EBR) {
	newNode := &MountNode{
		Ruta:   path,
		Key:    "",
		Digits: digit,
		Pos:    m.CountPartitions(path),
		Value:  part,
		ValueL: partL,
	}
	if !m.IsEmpty() {
		m.last.Next = newNode
		newNode.Prev = m.last
		m.last = newNode
		m.tamano++
	} else {
		m.first = newNode
		m.last = newNode
		m.tamano++
	}
	// deberia crear un m.PrintId o guardar en un singleton la consola
}

func (m MountList) UnMount(key_ string) *MountNode {
	if !m.IsEmpty() {
		temp := m.first
		counter := 0
		for counter < m.GetSize() {
			if key_ == temp.Key {
				if temp == m.first {
					m.first = m.first.Next
				} else if temp == m.last {
					m.last = m.last.Prev
					m.last.Next = nil
				} else {
					temp.Prev.Next = temp.Next
					temp.Next.Prev = temp.Prev
				}
				m.tamano--
				return temp
			}
		}
	}
	return nil
}

func (m MountList) GetNodeById(key_ string) *MountNode {
	if !m.IsEmpty() {
		temp := m.first
		for temp != nil {
			if key_ == temp.Key {
				return temp
			}
			temp = temp.Next
		}
	}
	return nil
}

func (m MountList) NodeExist(key_ string) bool {
	if !m.IsEmpty() {
		temp := m.first
		for temp != nil {
			if key_ == temp.Key {
				return true
			}
			temp = temp.Next
		}
	}
	return false
}

func (m MountList) CountPartitions(path string) int {
	contador := 1
	if !m.IsEmpty() {
		var temp *MountNode
		temp = m.first
		for temp != nil {
			if path == temp.Ruta {
				contador++
			}
			temp = temp.Next
		}
	}
	return contador
}

func (m MountList) GetSize() int {
	return m.tamano
}

var ListaMount = MountList{
	first:  nil,
	last:   nil,
	tamano: 0,
}
