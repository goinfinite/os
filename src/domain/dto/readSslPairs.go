package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadSslPairsRequest struct {
	Pagination          Pagination                `json:"pagination"`
	SslPairId           *valueObject.SslPairId    `json:"sslPairId"`
	VirtualHostHostname *valueObject.Fqdn         `json:"virtualHostHostname"`
	AltNames            []valueObject.SslHostname `json:"altNames"`
	IssuedBeforeAt      *valueObject.UnixTime     `json:"createdBeforeAt"`
	IssuedAfterAt       *valueObject.UnixTime     `json:"createdAfterAt"`
	ExpiresBeforeAt     *valueObject.UnixTime     `json:"expiresBeforeAt"`
	ExpiresAfterAt      *valueObject.UnixTime     `json:"expiresAfterAt"`
}

type ReadSslPairsResponse struct {
	Pagination Pagination       `json:"pagination"`
	SslPairs   []entity.SslPair `json:"sslPairs"`
}
