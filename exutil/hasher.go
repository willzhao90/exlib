package exutil

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"time"
)

// more coming
func HashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

//GetTimeMemo function will produce a 6 digits digest of time.now (not garantee unique)
func GetTimeMemo() string {
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := fnv.New32a()
	h.Write([]byte(s))
	s2 := fmt.Sprint(h.Sum32())
	return s2[len(s2)-6:]
}
