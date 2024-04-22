package service

import (
	"github.com/ngirchev/settings-loader/internal/domain"
	"github.com/ngirchev/settings-loader/internal/service/json"
)

func MapToSettingEntity(dto *json.SettingsDTO) domain.SettingsEntity {
	return domain.SettingsEntity{
		Id:    (*dto).Id,
		Name:  (*dto).Name,
		Color: (*dto).Color,
		Lives: (*dto).Lives,
	}
}
