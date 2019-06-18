package ftp

import (
	"time"

	"github.com/jlaffaye/ftp"
)

type Dialer interface {
	Dial(addr string) (conn ServerConnexion, err error)
	DialTimeout(addr string, timeout time.Duration) (conn ServerConnexion, err error)
}

type defaultDialer struct{}

func (d *defaultDialer) Dial(addr string) (conn ServerConnexion, err error) {
	return ftp.Dial(addr)
}
func (d *defaultDialer) DialTimeout(addr string, timeout time.Duration) (conn ServerConnexion, err error) {
	return ftp.DialTimeout(addr, timeout)
}

// Comment
func NewDefaultDialer() Dialer {
	return &defaultDialer{}
}
