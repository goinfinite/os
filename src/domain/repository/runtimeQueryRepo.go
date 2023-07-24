package repository

import "github.com/speedianet/sam/src/domain/valueObject"

type RuntimeQueryRepo interface {
	GetPhpVersions() ([]valueObject.PhpVersion, error)
}
