package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"settings-loader/internal/util"
)

type DBConf struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func DBOpen(conf DBConf) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s port=%s password=%s sslmode=%s",
		conf.Username, conf.DBName, conf.Host, conf.Port, conf.Password, conf.SSLMode))
	// fail-fast if db is not available
	util.HandleError("open sql connection error", err)
	err = db.Ping()
	util.HandleError("ping database", err)
	return db, nil
}

func DBClose(db *sqlx.DB) {
	err := db.Close()
	util.HandleError("close db error", err)
}
