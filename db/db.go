package main

import (
	redis "github.com/pavelrosputko/go-redis"
)

type Db struct {
	redis.Client
}

var db = &Db{}

