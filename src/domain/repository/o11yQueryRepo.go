package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
)

type O11yQueryRepo interface {
	ReadOverview() (entity.O11yOverview, error)
}
