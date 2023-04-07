package lista

import (
	"strconv"
	"strings"

	"github.com/erazoex/proyecto2/datos"
)

type MountNode struct {
	key, ruta   string
	digits, pos int
	value       *datos.Partition
	valueL      *datos.EBR
	next, prev  *MountNode
}

func (m MountNode) MountNode(ruta string, digits int, pos int, value *datos.Partition, valueL *datos.EBR) {
	m.ruta = ruta
	m.digits = digits
	m.pos = pos
	m.value = value
	m.valueL = valueL
	m.next = nil
	m.prev = nil
	m.CreateKey()
}

func (m MountNode) CreateKey() {
	directory := strings.Split(m.ruta, "/")
	lastPart := directory[len(directory)-1]
	fileNameParts := strings.Split(lastPart, ".")
	filename := fileNameParts[0]
	m.key = strconv.Itoa(m.digits) + strconv.Itoa(m.pos) + filename
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
		ruta:   path,
		key:    "",
		digits: digit,
		pos:    m.CountPartitions(path),
		value:  part,
		valueL: partL,
	}
	if !m.IsEmpty() {
		m.last.next = newNode
		newNode.prev = m.last
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
			if key_ == temp.key {
				if temp == m.first {
					m.first = m.first.next
				} else if temp == m.last {
					m.last = m.last.prev
					m.last.next = nil
				} else {
					temp.prev.next = temp.next
					temp.next.prev = temp.prev
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
			if key_ == temp.key {
				return temp
			}
			temp = temp.next
		}
	}
	return nil
}

func (m MountList) NodeExist(key_ string) bool {
	if !m.IsEmpty() {
		temp := m.first
		for temp != nil {
			if key_ == temp.key {
				return true
			}
			temp = temp.next
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
			if path == temp.ruta {
				contador++
			}
			temp = temp.next
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
