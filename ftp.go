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
func (f *FTP) Connect() (conn *ftp.ServerConn, err error) {
	// Log
	l := fmt.Sprintf("FTP connect to %s with timeout %s", f.Addr, f.Timeout)
	f.Logger.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		f.Logger.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Dial
	if f.Timeout > 0 {
		conn, err = ftp.DialTimeout(f.Addr, f.Timeout)
	} else {
		conn, err = ftp.Dial(f.Addr)
	}
	if err != nil {
		return
	}

	// Login
	if err = conn.Login(f.Username, f.Password); err != nil {
		conn.Quit()
	}
	return
}

// Download downloads a file from the remote server
func (f *FTP) Download(src, dst string) (err error) {
	// Log
	l := fmt.Sprintf("FTP download from %s to %s", src, dst)
	f.Logger.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		f.Logger.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Connect
	var conn *ftp.ServerConn
	if conn, err = f.Connect(); err != nil {
		return
	}
	defer conn.Quit()

	// Download file
	var r io.ReadCloser
	f.Logger.Debugf("Downloading %s", src)
	if r, err = conn.Retr(src); err != nil {
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
