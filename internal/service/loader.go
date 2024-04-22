package service

import (
	"bytes"
	"fmt"
	"github.com/ngirchev/settings-loader/internal/domain"
	"github.com/ngirchev/settings-loader/internal/service/json"
	"github.com/ngirchev/settings-loader/internal/util"
	log "github.com/sirupsen/logrus"
	"sync"
)

type ILoaderService interface {
	LoadComponent(cType string, version string, hash []byte) ([]byte, []byte, error)
}

type LoaderService struct {
	hasher       IHasher
	jsonReader   *json.Reader // todo change to interface
	settingsRepo domain.ISettingsRepo
	rootPath     string
}

func NewLoaderService(hasher IHasher, reader *json.Reader, settingsRepo domain.ISettingsRepo, rootPath string) *LoaderService {
	return &LoaderService{hasher: hasher, jsonReader: reader, settingsRepo: settingsRepo, rootPath: rootPath}
}

// LoadComponent retrieves the binary content and checksum of a component given its type and version.
//
// Parameters:
//   - Type: The type of the component to load.
//   - Version: The version of the component to load.
//
// Returns:
//   - []byte: The binary content of the component.
//   - []byte: The checksum or hash for verification.
//   - error: An error message if the component could not be loaded.
func (l *LoaderService) LoadComponent(cType string, version string, expectedHash []byte) ([]byte, []byte, error) {
	content, settingsDTOs, err := l.jsonReader.ReadSettingsJSONFile(getFilePath(l.rootPath, cType, version))
	if err != nil {
		return nil, nil, err
	}
	newHash := l.hasher.Hash(content)

	isValid := expectedHash == nil || bytes.Equal(expectedHash, newHash)
	if isValid {
		// updates db only if hash valid
		log.Infof("Hash is valid. Db will be updated")
		l.saveToDBPerBatch(settingsDTOs)
		return content, newHash, nil
	} else {
		return nil, newHash, nil
	}
}

// todo add transactions support and error handling.
// The assumption is - we don't need transactions for this case. It's expensive and we are hope on success of saving.
func (l *LoaderService) saveToDBPerBatch(settingsDTOs []json.SettingsDTO) {
	var wg sync.WaitGroup
	batchSize := 10 // todo should be in the properties
	batches := util.SplitArrayIntoBatches(settingsDTOs, batchSize)
	for i, batch := range batches {
		wg.Add(1)
		go util.DoWork(i, func() { l.saveToDb(&batch) }, &wg)
	}
	wg.Wait()
	log.Debug("All data saved")
}

func getFilePath(rootPath string, cType string, version string) string {
	path := fmt.Sprintf("%s/%s/%s.json", rootPath, cType, version)
	log.Debugf("Path created: %s", path)
	return path
}

func (l *LoaderService) saveToDb(dtoBatch *[]json.SettingsDTO) {
	for _, dto := range *dtoBatch {
		settingsEntity := MapToSettingEntity(&dto)
		// todo if we didn't find the entity, we should delete it
		err := l.settingsRepo.CreateOrUpdate(settingsEntity)
		if err != nil {
			log.Errorf("Can't save: %d, %s", dto.Id, err)
		}
	}
}
