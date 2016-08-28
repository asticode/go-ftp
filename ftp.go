package ftp

import (
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/rs/xlog"
)

// FTP represents an FTP
type FTP struct {
	Addr     string
	conn     *ftp.ServerConn
	Logger   xlog.Logger
	Password string
	Timeout  time.Duration
	Username string
}

// NewFromConfig creates a new FTP connection based on a configuration
func NewFromConfig(c Configuration) *FTP {
	return &FTP{
		Addr:     c.Addr,
		Password: c.Password,
		Timeout:  c.Timeout,
		Username: c.Username,
	}
}

// Dial dials in the FTP server
func (f *FTP) Dial() (err error) {
	if f.Timeout > 0 {
		f.conn, err = ftp.DialTimeout(f.Addr, f.Timeout)
	} else {
		f.conn, err = ftp.Dial(f.Addr)
	}
	return
}

// Login signs in to the FTP server
func (f *FTP) Login() error {
	return f.conn.Login(f.Username, f.Password)
}

// Move moves a file from the remote server
func (f *FTP) Move(src, dst string) (err error) {
	// Download file
	var r io.ReadCloser
	if r, err = f.conn.Retr(src); err != nil {
		return
	}
	defer r.Close()

	// Create the destination file
	var dstFile *os.File
	if dstFile, err = os.Create(dst); err != nil {
		return
	}
	defer dstFile.Close()

	// Copy to dst
	if _, err = io.Copy(dstFile, r); err != nil {
		return
	}

	// Delete src
	err = f.conn.Delete(src)
	return
}
