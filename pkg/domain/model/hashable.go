package model

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

type Hashable interface {
	GetBrand() string
	GetModel() string
	GetYear() int32
	GetCapacity() int32
}

func Hash(x Hashable) string {
	key := fmt.Sprintf("%s-%s-%d-%d", x.GetBrand(), x.GetModel(), x.GetYear(), x.GetCapacity())
	sha1 := sha1.Sum([]byte(key))

	return base64.URLEncoding.EncodeToString(sha1[:])
}
