package valueObject

import (
	"errors"
	"log/slog"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const systemResourceIdentifierRegex string = `^sri://(?P<accountId>[\d]{1,64}):(?P<resourceType>[a-zA-Z0-9][\w\-]{1,256})\/(?P<resourceId>[a-zA-Z-0-9\*][\w\.\-]{0,512})$`

type SystemResourceIdentifier string

func NewSystemResourceIdentifier(
	value interface{},
) (sri SystemResourceIdentifier, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sri, errors.New("SystemResourceIdentifierMustBeString")
	}

	re := regexp.MustCompile(systemResourceIdentifierRegex)
	if !re.MatchString(stringValue) {
		return sri, errors.New("InvalidSystemResourceIdentifier")
	}

	return SystemResourceIdentifier(stringValue), nil
}

// Note: this panic is solely to warn about the misuse of the VO, specifically for developer
// auditing, and has nothing to do with user input. This is not a standard and should not be
// followed for the development of other VOs.
func NewSystemResourceIdentifierIgnoreError(value interface{}) SystemResourceIdentifier {
	sri, err := NewSystemResourceIdentifier(value)
	if err != nil {
		panicMessage := "UnexpectedSystemResourceIdentifierCreationError"
		slog.Debug(panicMessage, slog.Any("value", value), slog.String("err", err.Error()))
		panic(panicMessage)
	}

	return sri
}

func NewAccountSri(accountId AccountId) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://0:account/" + accountId.String(),
	)
}

func NewSecureAccessPublicKeySri(
	accountId AccountId,
	secureAccessPublicKeyId SecureAccessPublicKeyId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":secureAccessPublicKey/" +
			secureAccessPublicKeyId.String(),
	)
}

func NewCronSri(
	accountId AccountId,
	cronId CronId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":cron/" + cronId.String(),
	)
}

func NewDatabaseSri(
	accountId AccountId,
	databaseName DatabaseName,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":database/" + databaseName.String(),
	)
}

func NewDatabaseUserSri(
	accountId AccountId,
	databaseUsername DatabaseUsername,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":databaseUser/" + databaseUsername.String(),
	)
}

func NewMarketplaceCatalogItemSri(
	accountId AccountId,
	marketplaceCatalogItemId *MarketplaceItemId,
	marketplaceCatalogItemSlug *MarketplaceItemSlug,
) SystemResourceIdentifier {
	if marketplaceCatalogItemId == nil && marketplaceCatalogItemSlug == nil {
		slog.Debug("MarketplaceCatalogItemSriMustHaveIdOrSlug")
		panic("MarketplaceCatalogItemSriMustHaveIdOrSlug")
	}

	marketplaceCatalogItemSri := "sri://" + accountId.String() + ":marketplaceCatalogItem/"
	if marketplaceCatalogItemId != nil {
		return NewSystemResourceIdentifierIgnoreError(
			marketplaceCatalogItemSri + marketplaceCatalogItemId.String(),
		)
	}

	return NewSystemResourceIdentifierIgnoreError(
		marketplaceCatalogItemSri + marketplaceCatalogItemSlug.String(),
	)
}

func NewMarketplaceInstalledItemSri(
	accountId AccountId,
	marketplaceInstalledItemId MarketplaceItemId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":marketplaceInstalledItem/" +
			marketplaceInstalledItemId.String(),
	)
}

func NewPhpRuntimeSri(
	accountId AccountId,
	virtualHostHostname Fqdn,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":phpRuntime/" + virtualHostHostname.String(),
	)
}

func NewInstallableServiceSri(
	accountId AccountId,
	serviceName ServiceName,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":installableService/" + serviceName.String(),
	)
}

func NewCustomServiceSri(
	accountId AccountId,
	serviceName ServiceName,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":customService/" + serviceName.String(),
	)
}

func NewInstalledServiceSri(
	accountId AccountId,
	serviceName ServiceName,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":installedService/" + serviceName.String(),
	)
}

func NewSslSri(
	accountId AccountId,
	sslPairId SslPairId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":ssl/" + sslPairId.String(),
	)
}

func NewVirtualHostSri(
	accountId AccountId,
	vhostHostname Fqdn,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":virtualHost/" + vhostHostname.String(),
	)
}

func NewMappingSri(
	accountId AccountId,
	mappingId MappingId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":mapping/" + mappingId.String(),
	)
}

func NewMappingSecurityRuleSri(
	accountId AccountId,
	mappingSecurityRuleId MappingSecurityRuleId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":mappingSecurityRule/" + mappingSecurityRuleId.String(),
	)
}

func (vo SystemResourceIdentifier) String() string {
	return string(vo)
}
