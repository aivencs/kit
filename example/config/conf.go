package main

import (
	"fmt"

	"github.com/aivencs/kit/pkg/config"
)

func main() {
	opts := config.ConfigOptions{Application: "serviceWork", Env: "product"}
	config.InitConfig("consul", opts)
	fmt.Println(config.GetString("application"))
}
