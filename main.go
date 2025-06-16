package main

import (
	"github.com/aditya-gupta-dev/oko/config"
	"github.com/aditya-gupta-dev/oko/ui"
)

func main() {
	config.VerifyConfigurationExistence()
	ui.StartApplication()
}
