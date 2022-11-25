package model

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math"

	"github.com/opencars/seedwork/logger"
)

type Hashable interface {
	GetBrand() string
	GetModel() string
	GetYear() int32
	GetCapacity() int32
}

func Hash(x Hashable) string {
	capacity := float64(x.GetCapacity()) / 100
	roundedCapacity := math.Round(capacity)

	key := fmt.Sprintf("%s-%s-%d-%f", x.GetBrand(), x.GetModel(), x.GetYear(), roundedCapacity)
	logger.Infof("key to be hashed: %s", key)
	sha1 := sha1.Sum([]byte(key))

	return base64.URLEncoding.EncodeToString(sha1[:])
}
