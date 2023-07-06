package valueObject

type Byte int64

func (b Byte) Get() int64 {
	return int64(b)
}

func (b Byte) ToKiB() int64 {
	return b.Get() / 1024
}

func (b Byte) ToMiB() int64 {
	return b.ToKiB() / 1024
}

func (b Byte) ToGiB() int64 {
	return b.ToMiB() / 1024
}

func (b Byte) ToTiB() int64 {
	return b.ToGiB() / 1024
}
