package domain

type ISettingsRepo interface {
	CreateOrUpdate(entity SettingsEntity) error
	GetById(id string) (SettingsEntity, error)
}
