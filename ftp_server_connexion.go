package ftp

import (
	"io"

	"github.com/jlaffaye/ftp"
)

type ServerConnexion interface {
	Login(sUsername string, sPwd string) error
	Retr(path string) (*ftp.Response, error)
	FileSize(path string) (int64, error)
	Stor(path string, oReader io.Reader) error
	MakeDir(sSource string) error
	RemoveDir(sSource string) error
	Rename(sSource string, sDestination string) error
	Delete(oath string) error
	Quit() error
	List(sPath string) ([]*ftp.Entry, error)
}
