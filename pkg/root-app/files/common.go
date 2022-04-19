package files

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// IsFileExist reports whether the named file or directory exists.
func IsFileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// OsAppRootPath returns default config path for the app.
func OsAppRootPath(appRootPath, appName, appUsage, appVersion string) string {
	name := appName + "-" + appVersion
	root := appRootPath
	dir, _ := os.UserConfigDir()

	switch runtime.GOOS {
	case "windows", "darwin":
		name = appUsage
		root = strings.ToUpper(root)
	}

	return filepath.Join(dir, root, name)
}

func MakeDirs(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DefaultDirPerm); err != nil {
		return errors.Wrap(err, "writeToYamlFile: create dir filed")
	}

	return nil
}
