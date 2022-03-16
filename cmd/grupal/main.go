package main

import (
	"github.com/matjam/grupal/api"
	"github.com/matjam/grupal/model"
)

func main() {
	model.Init()
	api.Start()
}
