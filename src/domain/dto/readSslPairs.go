package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadSslPairsRequest struct {
	Pagination          tkDto.Pagination          `json:"pagination"`
	SslPairId           *valueObject.SslPairId    `json:"sslPairId"`
	VirtualHostHostname *tkValueObject.Fqdn       `json:"virtualHostHostname"`
	AltNames            []valueObject.SslHostname `json:"altNames"`
	IssuedBeforeAt      *tkValueObject.UnixTime   `json:"createdBeforeAt"`
	IssuedAfterAt       *tkValueObject.UnixTime   `json:"createdAfterAt"`
	ExpiresBeforeAt     *tkValueObject.UnixTime   `json:"expiresBeforeAt"`
	ExpiresAfterAt      *tkValueObject.UnixTime   `json:"expiresAfterAt"`
}

type ReadSslPairsResponse struct {
	Pagination tkDto.Pagination `json:"pagination"`
	SslPairs   []entity.SslPair `json:"sslPairs"`
}
