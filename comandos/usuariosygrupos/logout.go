package usuariosygrupos

import (
	"fmt"

	"github.com/erazoex/proyecto2/logger"
)

type Logout struct {
}

func (l *Logout) Exe(parametros []string) {
	if l.Logout() {
		fmt.Printf("\nse ha cerrado la sesion\n\n")
	} else {
		fmt.Printf("no se logro hacer logout\n\n")
	}
}

func (l *Logout) Logout() bool {
	return logger.Log.Logout()
}
