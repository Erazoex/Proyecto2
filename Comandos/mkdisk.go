package comandos

import (
	"fmt"
	"strconv"
	"strings"
)

type ParametrosMkdisk struct {
	size int
	fit  [2]byte
	unit byte
	path string
}

type Mkdisk struct {
	params ParametrosMkdisk
}

func (m Mkdisk) exe() {

}

func (m Mkdisk) saveParams(parametros []string) {
	for _, v := range parametros {
		if strings.Contains(v, ">path") {
			v = strings.Replace(v, ">path=", "", 1)
			if v != "" {
				m.params.path = v
			}
		} else if strings.Contains(v, ">size") {
			v = strings.Replace(v, ">size=", "", 1)
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("hubo un error al convertir a int")
			}
			m.params.size = num
		} else if strings.Contains(v, ">unit") {
			v = strings.Replace(v, ">unit=", "", 1)
			m.params.unit = v[0]
		} else if strings.Contains(v, ">fit") {
			v = strings.Replace(v, ">fit=", "", 1)
			m.params.fit = [2]byte{v[0], v[1]}
		}
	}
}
