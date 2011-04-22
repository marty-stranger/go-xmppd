package main

import (
	"g/redis"
)

type Db struct {
	redis.Client
}

var db = &Db{}

