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
		return 0, errors.New("InvalidMktItemId")
	}

	mii := MktplaceItemId(mktplaceItemUid)
	if !mii.isValid() {
		return 0, errors.New("InvalidMktItemId")
	}

	return mii, nil
}

func NewMktplaceItemIdPanic(value interface{}) MktplaceItemId {
	mii, err := NewMktplaceItemId(value)
	if err != nil {
		panic(err)
	}

	return mii
}

func (mii MktplaceItemId) isValid() bool {
	return mii >= 1 && mii <= 999999999999
}

func (mii MktplaceItemId) Get() int64 {
	return int64(mii)
}

func (mii MktplaceItemId) String() string {
	return strconv.FormatInt(int64(mii), 10)
}

func (miiPtr *MktplaceItemId) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	mii, err := NewMktplaceItemId(valueStr)
	if err != nil {
		return err
	}

	*miiPtr = mii
	return nil
}
