package main

import (
	"crypto/md5"
	"encoding/hex"
)

func CaclTargetId(targetEndPoint string) string {
	md5res := md5.Sum([]byte(targetEndPoint))
	targetId := hex.EncodeToString(md5res[:])
	return targetId
}