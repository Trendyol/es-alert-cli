package reader

import (
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/model"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type FileReader struct {
}

func NewFileReader() (*FileReader, error) {
	return &FileReader{}, nil
}

func (f *FileReader) ReadLocalYaml(filename string) ([]model.MonitorConfig, error) {
	// Read YAML file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return nil, err
	}

	// Unmarshal YAML into struct
	var config []model.MonitorConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return nil, err
	}

	return config, nil
}
