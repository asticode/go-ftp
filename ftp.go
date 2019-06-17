package ftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/molotovtv/go-logger"
	astiio "github.com/molotovtv/go-astitools/io"
)

// FTP represents an FTP
type FTP struct {
	Addr     string
	Password string
	Timeout  time.Duration
	Username string
}

// New creates a new FTP connection based on a configuration
func New(c Configuration) *FTP {
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
	log.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		log.Debugf("[End] %s in %s", l, time.Since(now))
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

// DownloadReader returns the reader built from the download of a file
func (f *FTP) DownloadReader(src string) (conn *ftp.ServerConn, r io.ReadCloser, err error) {
	// Connect
	if conn, err = f.Connect(); err != nil {
		return
	}

	// Download file
	if r, err = conn.Retr(src); err != nil {
		return
	}
	return
}

// Download downloads a file from the remote server
func (f *FTP) Download(ctx context.Context, src, dst string) (err error) {
	// Log
	l := fmt.Sprintf("FTP download from %s to %s", src, dst)
	log.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		log.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Check context error
	if err = ctx.Err(); err != nil {
		return
	}

	// Connect
	var conn *ftp.ServerConn
	if conn, err = f.Connect(); err != nil {
		return
	}
	defer conn.Quit()

	// Check context error
	if err = ctx.Err(); err != nil {
		return
	}

	// Download file
	var r io.ReadCloser
	log.Debugf("Downloading %s", src)
	if r, err = conn.Retr(src); err != nil {
		return
	}
	defer r.Close()

	// Check context error
	if err = ctx.Err(); err != nil {
		return
	}

	// Create the destination file
	var dstFile *os.File
	log.Debugf("Creating %s", dst)
	if dstFile, err = os.Create(dst); err != nil {
		return
	}
	defer dstFile.Close()

	// Check context error
	if err = ctx.Err(); err != nil {
		return
	}

	// Copy to dst
	var n int64
	log.Debugf("Copying downloaded content to %s", dst)
	n, err = astiio.Copy(ctx, r, dstFile)
	log.Debugf("Copied %dkb", n/1024)
	return
}

// Remove removes a file
func (f *FTP) Remove(src string) (err error) {
	// Log
	l := fmt.Sprintf("FTP Remove of %s", src)
	log.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		log.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Connect
	var conn *ftp.ServerConn
	if conn, err = f.Connect(); err != nil {
		return
	}
	defer conn.Quit()

	// Remove
	log.Debugf("Removing %s", src)
	if err = conn.Delete(src); err != nil {
		return
	}
	return
}

// Upload uploads a source path content to a destination
func (f *FTP) Upload(ctx context.Context, src, dst string) (err error) {
	// Log
	l := fmt.Sprintf("FTP Upload to %s", dst)
	log.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		log.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Check context error
	if err = ctx.Err(); err != nil {
		return
	}

	// Connect
	var conn *ftp.ServerConn
	if conn, err = f.Connect(); err != nil {
		return
	}
	defer conn.Quit()

	// Check context error
	if err = ctx.Err(); err != nil {
		return
	}

	// Open file
	var srcFile *os.File
	log.Debugf("Opening %s", src)
	if srcFile, err = os.Open(src); err != nil {
		return
	}
	defer srcFile.Close()

	// Check context error
	if err = ctx.Err(); err != nil {
		return
	}

	// Store
	log.Debugf("Uploading %s to %s", src, dst)
	if err = conn.Stor(dst, astiio.NewReader(ctx, srcFile)); err != nil {
		return
	}
	return
}

// Remove removes a file
func (f *FTP) FileSize(src string) (s int64, err error) {
	// Log
	l := fmt.Sprintf("FTP file size of %s", src)
	log.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		log.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Connect
	var conn *ftp.ServerConn
	if conn, err = f.Connect(); err != nil {
		return
	}
	defer conn.Quit()

	// File size
	return conn.FileSize(src)
}
