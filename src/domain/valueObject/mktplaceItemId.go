package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MktplaceItemId int64

func NewMktplaceItemId(value interface{}) (MktplaceItemId, error) {
	mktplaceItemUid, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidMktplaceItemId")
	}

	mktplaceItemId := MktplaceItemId(mktplaceItemUid)
	if !mktplaceItemId.isValid() {
		return 0, errors.New("InvalidMktplaceItemId")
	}

	return mktplaceItemId, nil
}

func NewMktplaceItemIdPanic(value interface{}) MktplaceItemId {
	mktplaceItemId, err := NewMktplaceItemId(value)
	if err != nil {
		panic(err)
	}

	return mktplaceItemId
}

func (mktplaceItemId MktplaceItemId) isValid() bool {
	return mktplaceItemId >= 1 && mktplaceItemId <= 999999999999
}

func (mktplaceItemId MktplaceItemId) Get() int64 {
	return int64(mktplaceItemId)
}

func (mktplaceItemId MktplaceItemId) String() string {
	return strconv.FormatInt(int64(mktplaceItemId), 10)
}

func (mktplaceItemIdPtr *MktplaceItemId) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	mktplaceItemId, err := NewMktplaceItemId(valueStr)
	if err != nil {
		return err
	}

	*mktplaceItemIdPtr = mktplaceItemId
	return nil
}
