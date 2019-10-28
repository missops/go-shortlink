package main

import (
	"shortlink/api"
	"shortlink/utils"
)

//Env  ...
type Env struct {
	S api.Storage
}

func getEnv() *Env {
	r := utils.NewRedisCli("10.40.0.200:6379", "", 1)
	return &Env{S: r}
}
