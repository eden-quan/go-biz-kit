package utils

import (
	"errors"
	"github.com/jinzhu/copier"
	"strconv"
)

var IdConverters = []copier.TypeConverter{{
	SrcType: int64(0),
	DstType: copier.String,
	Fn: func(src interface{}) (dst interface{}, err error) {
		srcValue, ok := src.(int64)
		if !ok {
			return nil, errors.New("src is not int")
		}
		return strconv.FormatInt(srcValue, 10), nil
	},
}, {
	SrcType: int32(0),
	DstType: copier.String,
	Fn: func(src interface{}) (dst interface{}, err error) {
		srcValue, ok := src.(int32)
		if !ok {
			return nil, errors.New("src is not int")
		}
		return strconv.FormatInt(int64(srcValue), 10), nil
	},
}}

func Copy(desc interface{}, src interface{}) error {
	return copier.CopyWithOption(desc, src, copier.Option{
		Converters:    IdConverters,
		IgnoreEmpty:   true,
		CaseSensitive: false,
		DeepCopy:      true,
	})
}
