package main

import (
	"context"
	"fmt"

	"github.com/aivencs/kit/pkg/filter"
)

func main() {
	ctx := context.WithValue(context.Background(), "trace", "ctx-filter-001")
	opt := filter.FilterOption{
		Host:     "localhost:6379",
		Auth:     true,
		Username: "",
		Password: "password",
		Database: "",
		Table:    "",
		DB:       1,
		Key:      "seeds",
	}
	filter.InitFilter("redis", opt)
	payload := "19619c9e08f0ed4cc147e211efa8c3fb"
	res, err := filter.Add(ctx, payload)
	fmt.Println(res, err) // output: false nil
	ex, err := filter.Exist(ctx, payload)
	fmt.Println(ex, err) // output: true nil
}
