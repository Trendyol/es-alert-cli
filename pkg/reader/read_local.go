package reader

import (
	"io"
	"os"

	"github.com/labstack/gommon/log"

	"github.com/Trendyol/es-alert-cli/pkg/model"
	mapset "github.com/deckarep/golang-set"
	"gopkg.in/yaml.v3"
)

type FileReader struct{}

type FileReaderInterface interface {
	ReadLocalYaml(filename string) (map[string]model.Monitor, mapset.Set, error)
}

func NewFileReader() (*FileReader, error) {
	return &FileReader{}, nil
}

func (f *FileReader) ReadLocalYaml(filename string) (map[string]model.Monitor, mapset.Set, error) {
	// Read YAML file
	file, err := os.Open(filename)
	if err != nil {
		log.Errorf("Error opening file:", err.Error())
		return nil, nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("Error reading file:", err)
		return nil, nil, err
	}

	var monitors []model.Monitor
	err = yaml.Unmarshal(data, &monitors)
	if err != nil {
		log.Errorf("Error unmarshalling YAML:", err)
		return nil, nil, err
	}

	localMonitorSet := mapset.NewSet()
	monitorMap := make(map[string]model.Monitor)
	for _, monitor := range monitors {
		localMonitorSet.Add(monitor.Name)
		monitorMap[monitor.Name] = monitor
	}

	return monitorMap, localMonitorSet, nil
}
