package config

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/jkrus/Test_Simple_Multuplexor/pkg/root-app/files"
)

// Load ...
func Load(appRootPath, appName, appUsage, appVersion, configFileName string, cfg interface{}) error {
	file := filepath.Join(files.OsAppRootPath(appRootPath, appName, appUsage, appVersion), configFileName)
	if err := files.ReadFromYamlFile(file, cfg); err != nil {
		return errors.Wrap(err, "Load: read config from file filed")
	}

	return nil
}
