package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteMarketplaceInstalledItem struct {
	InstalledId             valueObject.MarketplaceItemId `json:"installedId"`
	ShouldUninstallServices bool                          `json:"shouldUninstallServices"`
	OperatorAccountId       tkValueObject.AccountId       `json:"-"`
	OperatorIpAddress       tkValueObject.IpAddress       `json:"-"`
}

func NewDeleteMarketplaceInstalledItem(
	installedId valueObject.MarketplaceItemId,
	shouldUninstallServices bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteMarketplaceInstalledItem {
	return DeleteMarketplaceInstalledItem{
		InstalledId:             installedId,
		ShouldUninstallServices: shouldUninstallServices,
		OperatorAccountId:       operatorAccountId,
		OperatorIpAddress:       operatorIpAddress,
	}
}
