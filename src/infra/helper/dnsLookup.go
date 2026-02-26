package infraHelper

import (
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

func DnsLookup(recordName string, recordType *string) ([]string, error) {
	hostname, err := tkValueObject.NewUnixHostname(recordName)
	if err != nil {
		return []string{}, err
	}

	var dnsRecordTypePtr *tkValueObject.DnsRecordType
	if recordType != nil {
		dnsRecordType, err := tkValueObject.NewDnsRecordType(*recordType)
		if err != nil {
			return []string{}, err
		}
		dnsRecordTypePtr = &dnsRecordType
	}

	return tkInfra.NewDnsLookup(hostname, dnsRecordTypePtr).Execute()
}
