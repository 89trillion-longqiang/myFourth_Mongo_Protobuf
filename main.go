package main

import (
	"giftCode/mongo"
	"giftCode/route"
	"giftCode/setUpRedis"
	"github.com/go-redis/redis"
)



func main() {
	var rdb *redis.Client
	setUpRedis.InitClient(rdb)
	mongo.SetUpMongo()
	r:=route.SetUpRount()
	r.Run(":8080")

}




