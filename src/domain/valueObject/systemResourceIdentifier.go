package valueObject

import (
	"log/slog"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

var NewSystemResourceIdentifier = tkValueObject.NewSystemResourceIdentifier

func NewSystemResourceIdentifierIgnoreError(
	value interface{},
) tkValueObject.SystemResourceIdentifier {
	return tkValueObject.NewSystemResourceIdentifierMustCreate(value)
}

func NewSecureAccessPublicKeySri(
	accountId tkValueObject.AccountId,
	secureAccessPublicKeyId SecureAccessPublicKeyId,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":secureAccessPublicKey/" +
			secureAccessPublicKeyId.String(),
	)
}

func NewCronSri(
	accountId tkValueObject.AccountId,
	cronId CronId,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":cron/" + cronId.String(),
	)
}

func NewDatabaseSri(
	accountId tkValueObject.AccountId,
	databaseName DatabaseName,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":database/" + databaseName.String(),
	)
}

func NewDatabaseUserSri(
	accountId tkValueObject.AccountId,
	databaseUsername DatabaseUsername,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":databaseUser/" + databaseUsername.String(),
	)
}

func NewMarketplaceCatalogItemSri(
	accountId tkValueObject.AccountId,
	marketplaceCatalogItemId *MarketplaceItemId,
	marketplaceCatalogItemSlug *MarketplaceItemSlug,
) tkValueObject.SystemResourceIdentifier {
	if marketplaceCatalogItemId == nil && marketplaceCatalogItemSlug == nil {
		slog.Debug("MarketplaceCatalogItemSriMustHaveIdOrSlug")
		panic("MarketplaceCatalogItemSriMustHaveIdOrSlug")
	}

	itemSri := "sri://" + accountId.String() +
		":marketplaceCatalogItem/"
	if marketplaceCatalogItemId != nil {
		return NewSystemResourceIdentifierIgnoreError(
			itemSri + marketplaceCatalogItemId.String(),
		)
	}

	return NewSystemResourceIdentifierIgnoreError(
		itemSri + marketplaceCatalogItemSlug.String(),
	)
}

func NewMarketplaceInstalledItemSri(
	accountId tkValueObject.AccountId,
	marketplaceInstalledItemId MarketplaceItemId,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":marketplaceInstalledItem/" +
			marketplaceInstalledItemId.String(),
	)
}

func NewPhpRuntimeSri(
	accountId tkValueObject.AccountId,
	virtualHostHostname tkValueObject.Fqdn,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":phpRuntime/" + virtualHostHostname.String(),
	)
}

func NewInstallableServiceSri(
	accountId tkValueObject.AccountId,
	serviceName ServiceName,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":installableService/" + serviceName.String(),
	)
}

func NewCustomServiceSri(
	accountId tkValueObject.AccountId,
	serviceName ServiceName,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":customService/" + serviceName.String(),
	)
}

func NewInstalledServiceSri(
	accountId tkValueObject.AccountId,
	serviceName ServiceName,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":installedService/" + serviceName.String(),
	)
}

func NewSslSri(
	accountId tkValueObject.AccountId,
	sslPairId SslPairId,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":ssl/" + sslPairId.String(),
	)
}

func NewVirtualHostSri(
	accountId tkValueObject.AccountId,
	vhostHostname tkValueObject.Fqdn,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":virtualHost/" + vhostHostname.String(),
	)
}

func NewMappingSri(
	accountId tkValueObject.AccountId,
	mappingId MappingId,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":mapping/" + mappingId.String(),
	)
}

func NewMappingSecurityRuleSri(
	accountId tkValueObject.AccountId,
	mappingSecurityRuleId MappingSecurityRuleId,
) tkValueObject.SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":mappingSecurityRule/" + mappingSecurityRuleId.String(),
	)
}
