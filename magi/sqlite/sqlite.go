package sqlite

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func New(
	v *viper.Viper,
) (
	*sql.DB,
	error,
) {
	toReturn, err := sql.Open(
		"sqlite3",
		v.GetString("file"),
	)

	if err != nil {
		log.WithFields(
			log.Fields{
				"package":  "sqlite",
				"function": "New()",
				"error":    err.Error,
				"file":     v.GetString("file"),
			},
		).Error("Failed to open file")
		return nil, err
	}

	return toReturn, nil
}
