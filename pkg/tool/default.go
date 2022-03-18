package tool

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"time"
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

// 计算时间差值
func CalcTimeDiff(ta, tb string) int64 {
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", ta, time.Local)
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", tb, time.Local)
	return t2.Unix() - t1.Unix()
}
