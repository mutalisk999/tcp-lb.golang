package main

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

func calcTargetId(targetEndPoint string) string {
	md5res := md5.Sum([]byte(targetEndPoint))
	targetId := hex.EncodeToString(md5res[:])
	return targetId
}

func verifyEndPointStr(endPointStr string) bool {
	l := strings.Split(endPointStr, ":")
	if len(l) != 2 {
		return false
	}

	ipAddrStr, portStr := l[0], l[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}
	if port <= 0 || port > 65535 {
		return false
	}

	l = strings.Split(ipAddrStr, ".")
	if len(l) != 4 {
		return false
	}
	for _, str := range l {
		i, err := strconv.Atoi(str)
		if err != nil {
			return false
		}
		if i < 0 || i > 255 {
			return false
		}
	}

	return true
}
