package repository

import (
	"github.com/speedianet/os/src/domain/entity"
)

type O11yQueryRepo interface {
	GetOverview() (entity.O11yOverview, error)
}
