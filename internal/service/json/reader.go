package json

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type Reader struct{}

func NewJsonReader() *Reader {
	return &Reader{}
}

func (r Reader) ReadSettingsJSONFile(filepath string) ([]byte, []SettingsDTO, error) {
	data, err := readJSONFile(filepath)
	if err != nil {
		return nil, nil, err
	}

	var settingsDTOs []SettingsDTO
	if err := json.Unmarshal(data, &settingsDTOs); err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		for _, settingsDTO := range settingsDTOs {
			log.Debugf("Found: %d %s %s %d", settingsDTO.Id, settingsDTO.Name, settingsDTO.Color, settingsDTO.Lives)
		}
	}

	return data, settingsDTOs, err
}

func readJSONFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return data, nil
}
