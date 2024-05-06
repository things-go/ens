package command

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/things-go/ens/driver"
)

func LoadDriver(URL string) (driver.Driver, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	d, err := driver.LoadDriver(u.Scheme)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func joinFilename(dir, filename, suffix string) string {
	suffix = strings.TrimSpace(suffix)
	if suffix != "" && !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	return filepath.Join(dir, filename) + suffix
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it
// and its upper level paths.
func WriteFile(filename string, data []byte) error {
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0655)
}
