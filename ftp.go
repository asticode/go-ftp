package ftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	astilog "github.com/molotovtv/go-astilog"
	astiio "github.com/molotovtv/go-astitools/io"
	"github.com/molotovtv/go-logger"
)

// FTP represents an FTP
type FTP struct {
	Addr     string
	Password string
	Timeout  time.Duration
	Username string
	dialer   Dialer
}

// New creates a new FTP connection based on a configuration
func New(c Configuration, dialer Dialer) *FTP {
	return &FTP{
		Addr:     c.Addr,
		Password: c.Password,
		Timeout:  c.Timeout,
		Username: c.Username,
		dialer:   dialer,
	}
}

// Connect connects to the FTP and logs in
func (f *FTP) Connect() (conn ServerConnexion, err error) {
	// Log
	l := fmt.Sprintf("FTP connect to %s with timeout %s", f.Addr, f.Timeout)
	log.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		log.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Dial
	if f.Timeout > 0 {
		conn, err = f.dialer.DialTimeout(f.Addr, f.Timeout)
	} else {
		conn, err = f.dialer.Dial(f.Addr)
	}
	if err != nil {
		return conn, err
	}

	// Login
	if err = conn.Login(f.Username, f.Password); err != nil {
		conn.Quit()
	}
	// fmt.Print(conn)
	// os.Exit(0)

	return conn, err
}

// DownloadReader returns the reader built from the download of a file
func (f *FTP) DownloadReader(src string) (conn ServerConnexion, r io.ReadCloser, err error) {
	// Connect
	if conn, err = f.Connect(); err != nil {
		return conn, nil, err
	}

	// Download file
	if r, err = conn.Retr(src); err != nil {
		return conn, nil, err
	}
	return conn, r, nil
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
	var conn ServerConnexion
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
	var conn ServerConnexion
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
	var conn ServerConnexion
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
	var conn ServerConnexion
	if conn, err = f.Connect(); err != nil {
		return
	}
	defer conn.Quit()

	// File size
	return conn.FileSize(src)
}

// var FTPConnect = func(f *FTP) (conn *ftp.ServerConn, err error) {
// 	return nil, f.Connect
// }

func (f *FTP) List(sFolder string, aExtensionsAllowed []string, sPattern string) []*ftp.Entry {

	// Log
	l := fmt.Sprintf("FTP file list of %s", sFolder)
	astilog.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		astilog.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	var aFiles, aFilesRaw []*ftp.Entry

	// Connect
	var conn ServerConnexion
	conn, err := f.Connect()
	if err != nil {
		return aFilesRaw
	}
	defer conn.Quit()
	aFilesRaw, _ = conn.List(sFolder)

	aExtensions := make(map[string]string)
	for _, sExtension := range aExtensionsAllowed {
		sExtension = strings.ToLower(sExtension)
		aExtensions[sExtension] = sExtension
	}
	bExtension := len(aExtensions) > 0

	bPattern := len(sPattern) > 0

	for _, oFile := range aFilesRaw {
		if oFile.Type != ftp.EntryTypeFile {
			// aFiles = append(aFiles[:k], aFiles[k+1:]...)
			continue
		}

		sExtension := f.GetExtensionFile(oFile)
		if bExtension {
			if _, err := aExtensions[sExtension]; !err {
				// aFiles = append(aFiles[:k], aFiles[k+1:]...)
				continue
			}
		}

		if bPattern {
			sFileName := f.GetFileNameWithoutExtension(oFile.Name)
			if bMatch, _ := regexp.MatchString(sPattern, sFileName); !bMatch {
				continue
			}
		}

		aFiles = append(aFiles, oFile)
	}

	return aFiles
	// return nil
	// File size
	// return conn.ListFileSize(src)
}

func (f *FTP) GetFileNameWithoutExtension(sFileName string) string {
	aFileName := strings.Split(sFileName, ".")
	return strings.Join(aFileName[:len(aFileName)-1], ".")
}

func (f *FTP) GetExtensionFile(oFile *ftp.Entry) string {
	aFileName := strings.Split(oFile.Name, ".")
	sExtension := aFileName[len(aFileName)-1]
	return strings.ToLower(sExtension)
}

func (f *FTP) Exists(sFilePath string) (b bool, err error) {
	// Log
	l := fmt.Sprintf("FTP file exists of %s", sFilePath)
	astilog.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		astilog.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// Connect
	var conn ServerConnexion
	if conn, err = f.Connect(); err != nil {
		return false, err
	}
	defer conn.Quit()

	aFilePath := strings.Split(sFilePath, "/")

	sFileName := aFilePath[len(aFilePath)-1]
	aFilePath = aFilePath[:len(aFilePath)-1]
	sFolder := strings.Join(aFilePath, "/")

	var aExtensions = []string{}

	aFiles := f.List(sFolder, aExtensions, sFileName)

	// File size
	return len(aFiles) > 0, nil
}
