package ftp

import (
	"fmt"
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

// Connect connects to the FTP and logs in
func (f *FTP) Connect() (err error) {
	if f.conn == nil {
		// Dial
		if err = f.Dial(); err != nil {
			return
		}

		// Login
		err = f.Login()
	}
	return
}

// Copy copies a file from the remote server
func (f *FTP) Copy(src, dst string) (err error) {
	// Log
	l := fmt.Sprintf("FTP copy from %s to %s", src, dst)
	f.Logger.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		f.Logger.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Connect
	if err = f.Connect(); err != nil {
		return
	}

	// Download file
	var r io.ReadCloser
	f.Logger.Debugf("Downloading %s", src)
	if r, err = f.conn.Retr(src); err != nil {
		return
	}
	defer r.Close()

	// Create the destination file
	var dstFile *os.File
	f.Logger.Debugf("Creating %s", dst)
	if dstFile, err = os.Create(dst); err != nil {
		return
	}
	defer dstFile.Close()

	// Copy to dst
	var n int64
	f.Logger.Debugf("Copying downloaded content to %s", dst)
	n, err = io.Copy(dstFile, r)
	f.Logger.Debugf("Copied %dkb", n/1024)
	return
}

// Dial dials in the FTP server
func (f *FTP) Dial() (err error) {
	// Log
	l := fmt.Sprintf("FTP dial to %s with timeout %s", f.Addr, f.Timeout)
	f.Logger.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		f.Logger.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Dial
	if f.Timeout > 0 {
		f.conn, err = ftp.DialTimeout(f.Addr, f.Timeout)
	} else {
		f.conn, err = ftp.Dial(f.Addr)
	}
	return
}

// Login signs in to the FTP server
func (f *FTP) Login() error {
	// Log
	l := fmt.Sprintf("FTP login to %s", f.Addr)
	f.Logger.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		f.Logger.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Login
	return f.conn.Login(f.Username, f.Password)
}

// Move moves a file from the remote server
func (f *FTP) Move(src, dst string) (err error) {
	// Log
	l := fmt.Sprintf("FTP move from %s to %s", src, dst)
	f.Logger.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		f.Logger.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Connect
	if err = f.Connect(); err != nil {
		return
	}

	// Copy
	if err = f.Copy(src, dst); err != nil {
		return
	}

	// Delete src
	f.Logger.Debugf("Deleting %s", src)
	err = f.conn.Delete(src)
	return
}
