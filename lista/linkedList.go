package lista

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erazoex/proyecto2/consola"
	"github.com/erazoex/proyecto2/datos"
	"github.com/erazoex/proyecto2/functions"
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
	m.PrintList()
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
	consola.AddToConsole(fmt.Sprintf("Id: %s\n", node.Key))
}

func (m *MountList) PrintList() {
	str := ""
	for i := 0; i < 110; i++ {
		str += "-"
	}
	contenido := ""
	contenido += fmt.Sprintf("%s\n", str)
	contenido += fmt.Sprintf("%-15s", "Id")
	contenido += fmt.Sprintf("%-15s", "Name")
	contenido += fmt.Sprintf("%-10s", "Type")
	contenido += fmt.Sprintf("%-10s", "Fit")
	contenido += fmt.Sprintf("%-10s", "Start")
	contenido += fmt.Sprintf("%-10s", "Size")
	contenido += fmt.Sprintf("%-10s", "Status")
	contenido += fmt.Sprintf("%-30s\n", "Ruta")
	if !m.IsEmpty() {
		temp := m.First
		for temp != nil {
			contenido += fmt.Sprintf("%s\n", str)
			contenido += fmt.Sprintf("%-15s", temp.Key)
			if temp.Value != nil {
				contenido += fmt.Sprintf("%-15s", string(functions.TrimArray(temp.Value.Part_name[:])))
				contenido += fmt.Sprintf("%-10s", string(temp.Value.Part_type))
				contenido += fmt.Sprintf("%-10c", temp.Value.Part_fit)
				contenido += fmt.Sprintf("%-10d", temp.Value.Part_start)
				contenido += fmt.Sprintf("%-10d", temp.Value.Part_size)
				contenido += fmt.Sprintf("%-10c", temp.Value.Part_status)

			} else if temp.ValueL != nil {
				contenido += fmt.Sprintf("%-15s", string(functions.TrimArray(temp.ValueL.Part_name[:])))
				contenido += fmt.Sprintf("%-10s", "L")
				contenido += fmt.Sprintf("%-10c", temp.ValueL.Part_fit)
				contenido += fmt.Sprintf("%-10d", temp.ValueL.Part_start)
				contenido += fmt.Sprintf("%-10d", temp.ValueL.Part_size)
				contenido += fmt.Sprintf("%-10c", temp.ValueL.Part_status)
			}
			contenido += fmt.Sprintf("%-30s\n", temp.Ruta)
			temp = temp.Next
		}
	}
	contenido += fmt.Sprintf("%s\n\n", str)
	consola.AddToConsole(contenido)
}

var ListaMount = MountList{
	First:  nil,
	Last:   nil,
	Tamano: 0,
}

// LinkedList de usuarios
type UserID struct {
	uid   string
	gid   string
	uname string
	Next  *UserID
	Prev  *UserID
}

func (u *UserID) GetUID() string {
	return u.uid
}

func (u *UserID) GetGID() string {
	return u.gid
}

func (u *UserID) GetUName() string {
	return u.uname
}

type UserList struct {
	First  *UserID
	Last   *UserID
	Length int
}

func (u *UserList) IsEmpty() bool {
	return u.First == nil
}

func (u *UserList) AddUser(userId, groupId, username string) {
	newNode := &UserID{
		uid:   userId,
		gid:   groupId,
		uname: username,
	}
	if !u.IsEmpty() {
		u.Last.Next = newNode
		newNode.Prev = u.Last
		u.Last = newNode
		u.Length++
	} else {
		u.First = newNode
		u.First.Next = nil
		u.Last = newNode
		u.Last.Next = nil
		u.Length++
	}
}

func (u *UserList) GetUserById(userId_ string) *UserID {
	temp := u.First
	for temp != nil {
		if temp.uid == userId_ {
			return temp
		}
		temp = temp.Next
	}
	return nil
}

func (u *UserList) GetUsersByGroup(groupId_ string) []*UserID {
	var result []*UserID
	temp := u.First
	for temp != nil {
		if temp.gid == groupId_ {
			result = append(result, temp)
		}
	}
	return result
}
