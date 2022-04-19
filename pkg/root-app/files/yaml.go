package files

import (
	"bytes"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// WriteToYamlFile write data to yaml file from  src.
func WriteToYamlFile(path string, src interface{}) error {
	buf := bytes.NewBuffer(nil)
	encoder := yaml.NewEncoder(buf)
	if err := encoder.Encode(src); err != nil {
		return errors.Wrap(err, "writeToYamlFile: encode config filed")
	}
	defer encoder.Close()

	if err := os.WriteFile(path, buf.Bytes(), DefaultFilePerm); err != nil {
		return errors.Wrap(err, "writeToYamlFile: write config to file filed")
	}

	return nil
}

// ReadFromYamlFile read data from yaml file into dst.
func ReadFromYamlFile(path string, dst interface{}) error {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "readFromYamlFile: open config file failed")
	}
	defer f.Close()

	if err = yaml.NewDecoder(f).Decode(dst); err != nil {
		return errors.Wrap(err, "readFromYamlFile: decode config file filed")
	}

	return nil
}
