package main

import (
	"context"
	"fmt"

	"github.com/aivencs/kit/pkg/cache"
)

func main() {
	ctx := context.WithValue(context.Background(), "trace", "ctx-cache-001")
	opt := cache.CacheOption{
		Host:     "localhost:6379",
		Auth:     true,
		Username: "",
		Password: "password",
		Database: "",
		Table:    "",
		DB:       1,
	}
	cache.InitCache("redis", opt)
	payload := "19619c9e08f0ed4cc147e211efa8c3fb"
	r, err := cache.SetEx(ctx, payload, 1, 20)
	fmt.Println(r, err) // output: OK nil
	ov := cache.Overdue(ctx, payload)
	fmt.Println(ov)                             // output: true
	fmt.Println(cache.Set(ctx, payload, "105")) // output: OK nil
	val, err := cache.Get(ctx, payload)
	fmt.Println(string(val.([]uint8)), err) // output: 105 <nil>
}
