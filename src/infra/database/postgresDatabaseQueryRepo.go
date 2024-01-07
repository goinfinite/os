package databaseInfra

import (
	"github.com/speedianet/os/src/domain/entity"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type PostgresDatabaseQueryRepo struct {
}

func PostgresqlCmd(cmd string) (string, error) {
	return infraHelper.RunCmd(
		"psql",
		"-tAc",
		cmd,
	)
}

func (repo PostgresDatabaseQueryRepo) Get() ([]entity.Database, error) {
	return []entity.Database{}, nil
}
