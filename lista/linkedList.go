package lista

import (
	"fmt"
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

func (m *MountNode) MountNode(ruta string, digits int, pos int, value *datos.Partition, valueL *datos.EBR) {
	m.Ruta = ruta
	m.Digits = digits
	m.Pos = pos
	m.Value = value
	m.ValueL = valueL
	m.Next = nil
	m.Prev = nil
	m.CreateKey()
}

func (m *MountNode) CreateKey() string {
	directory := strings.Split(m.Ruta, "/")
	lastPart := directory[len(directory)-1]
	fileNameParts := strings.Split(lastPart, ".")
	filename := fileNameParts[0]
	return strconv.Itoa(m.Digits) + strconv.Itoa(m.Pos) + filename
}

type MountList struct {
	First, Last *MountNode
	Tamano      int
}

func (m *MountList) IsEmpty() bool {
	return m.First == nil
}

func (m *MountList) Mount(path string, digit int, part *datos.Partition, partL *datos.EBR) {
	newNode := &MountNode{
		Ruta:   path,
		Key:    "",
		Digits: digit,
		Pos:    m.CountPartitions(path),
		Value:  part,
		ValueL: partL,
	}
	newNode.Key = newNode.CreateKey()
	// fmt.Println(m.IsEmpty())
	if !m.IsEmpty() {
		m.Last.Next = newNode
		newNode.Prev = m.Last
		m.Last = newNode
		m.Tamano++
	} else {
		m.First = newNode
		m.First.Next = nil
		m.Last = newNode
		m.Last.Next = nil
		m.Tamano++
	}
	// deberia crear un m.PrintId o guardar en un singleton la consola
	m.GetId(newNode)
}

func (m *MountList) UnMount(key_ string) *MountNode {
	if !m.IsEmpty() {
		temp := m.First
		counter := 0
		for counter < m.GetSize() {
			if key_ == temp.Key {
				if temp == m.First {
					m.First = m.First.Next
				} else if temp == m.Last {
					m.Last = m.Last.Prev
					m.Last.Next = nil
				} else {
					temp.Prev.Next = temp.Next
					temp.Next.Prev = temp.Prev
				}
				m.Tamano--
				return temp
			}
		}
	}
	return nil
}

func (m *MountList) GetNodeById(key_ string) *MountNode {
	if !m.IsEmpty() {
		temp := m.First
		for temp != nil {
			if key_ == temp.Key {
				return temp
			}
			temp = temp.Next
		}
	}
	return nil
}

func (m *MountList) NodeExist(key_ string) bool {
	if !m.IsEmpty() {
		temp := m.First
		for temp != nil {
			if key_ == temp.Key {
				return true
			}
			temp = temp.Next
		}
	}
	return false
}

func (m *MountList) CountPartitions(path string) int {
	contador := 1
	if !m.IsEmpty() {
		var temp *MountNode
		temp = m.First
		for temp != nil {
			if path == temp.Ruta {
				contador++
			}
			temp = temp.Next
		}
	}
	return contador
}

func (m *MountList) GetSize() int {
	return m.Tamano
}

func (m *MountList) GetId(node *MountNode) {
	fmt.Println("Id:", node.Key)
}

func (m *MountList) PrintList() {
	if !m.IsEmpty() {
		temp := m.First
		for temp != nil {
			fmt.Println("key:", temp.Key)
			fmt.Println("path:", temp.Ruta)
			temp = temp.Next
		}
	}
}

var ListaMount = MountList{
	First:  nil,
	Last:   nil,
	Tamano: 0,
}
