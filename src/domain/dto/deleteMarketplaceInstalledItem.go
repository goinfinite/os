package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteMarketplaceInstalledItem struct {
	InstalledId             valueObject.MarketplaceItemId
	ShouldUninstallServices bool
	OperatorAccountId       valueObject.AccountId
	OperatorIpAddress       valueObject.IpAddress
}

func NewDeleteMarketplaceInstalledItem(
	installedId valueObject.MarketplaceItemId,
	shouldUninstallServices bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteMarketplaceInstalledItem {
	return DeleteMarketplaceInstalledItem{
		InstalledId:             installedId,
		ShouldUninstallServices: shouldUninstallServices,
		OperatorAccountId:       operatorAccountId,
		OperatorIpAddress:       operatorIpAddress,
	}
}
