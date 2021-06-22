package ftp

import (
	"flag"
	"time"
)

// Flags
var (
	Addr     = flag.String("ftp-addr", "", "the ftp addr")
	Password = flag.String("ftp-password", "", "the ftp password")
	Timeout  = flag.Duration("ftp-timeout", 0, "the ftp timeout")
	Username = flag.String("ftp-username", "", "the ftp username")
)

// Configuration represents the FTP configuration
type Configuration struct {
	Addr       string        `json:"addr"`
	Password   string        `json:"password"`
	Timeout    time.Duration `toml:"timeout"`
	Username   string        `json:"username"`
	Persistent bool          `json:"persistent"`
	TTL        time.Duration `json:"ttl"`
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() Configuration {
	return Configuration{
		Addr:     *Addr,
		Password: *Password,
		Timeout:  *Timeout,
		Username: *Username,
	}
}
