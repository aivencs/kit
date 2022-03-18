package tool

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
)

// 生成消息摘要
func CreateDigest(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 元素是否在数组中
func IsContainString(target string, raw []string) bool {
	sort.Strings(raw)
	index := sort.SearchStrings(raw, target)
	if index < len(raw) && raw[index] == target {
		return true
	}
	return false
}

func IsContainInt(target int, raw []int) bool {
	sort.Ints(raw)
	index := sort.SearchInts(raw, target)
	if index < len(raw) && raw[index] == target {
		return true
	}
	return false
}
