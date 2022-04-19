package config

import (
	"log"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/jkrus/Test_Simple_Multuplexor/pkg/root-app/files"
)

// Init ...
func Init(appRootPath, appName, appUsage, appVersion, configFileName string, cfg interface{}) error {
	filePath := filepath.Join(files.OsAppRootPath(appRootPath, appName, appUsage, appVersion), configFileName)
	if files.IsFileExist(filePath) {
		log.Println("Read data from config file in path:", filePath)
		if err := files.ReadFromYamlFile(filePath, cfg); err != nil {
			return errors.Wrap(err, "Init: read config file filed")
		}
	} else {
		log.Println("Create default config file in path:", filePath)
	}

	if err := files.MakeDirs(filePath); err != nil {
		return errors.Wrap(err, "Init: can not create dirs")
	}
	if err := files.WriteToYamlFile(filePath, cfg); err != nil {
		return errors.Wrap(err, "Init: create config file filed")
	}

	return nil
}
