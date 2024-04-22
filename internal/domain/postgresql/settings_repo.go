package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"settings-loader/internal/domain"
)

type SettingsRepo struct {
	db *sqlx.DB
}

func NewSettingsRepo(props DBConf) *SettingsRepo {
	db, err := DBOpen(props)
	if err != nil {
		panic(err)
	}
	return &SettingsRepo{db: db}
}

func (r *SettingsRepo) Close() {
	DBClose(r.db)
}

func (r *SettingsRepo) CreateOrUpdate(settings domain.SettingsEntity) error {
	sqlStatement := `INSERT INTO settings (id, name, color, lives, created, updated) VALUES ($1, $2, $3, $4, now(), now()) 
						ON CONFLICT (id) DO UPDATE SET name = $2, color = $3, lives = $4, updated = now();`
	_, err := r.db.Exec(sqlStatement, settings.Id, settings.Name, settings.Color, settings.Lives)
	if err != nil {
		return fmt.Errorf("execute create statement: %w", err)
	}

	return nil
}

func (r *SettingsRepo) GetById(id string) (domain.SettingsEntity, error) {
	query := "SELECT * FROM settings WHERE id = $1"

	var settingsEntity domain.SettingsEntity
	err := r.db.Get(&settingsEntity, query, id)
	if err != nil {
		return domain.SettingsEntity{}, fmt.Errorf("execute get statement: %w", err)
	}

	return settingsEntity, nil
}
