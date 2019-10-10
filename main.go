package main

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	SchedularService "github.com/nayanmakasare/SchedularService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/robfig/cron.v3"
	"time"
)

const (
	defaultHost = "mongodb://nayan:tlwn722n@cluster0-shard-00-00-8aov2.mongodb.net:27017,cluster0-shard-00-01-8aov2.mongodb.net:27017,cluster0-shard-00-02-8aov2.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin&retryWrites=true&w=majority"
	//defaultHost = "mongodb://192.168.1.9:27017"
	//defaultHost = "mongodb://192.168.1.143:27017"
)

func main() {
	service := grpc.NewService(
		micro.Name("SchedularService"),
		micro.Address(":50055"),
		micro.Version("1.0"),
	)
	service.Init()
	uri := defaultHost
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Debug(err)
	}
	redisClient := GetRedisClient()

	c := cron.New()

	// Register Broker
	publisher := micro.NewPublisher("applySchedule", service.Client())

	h := SchedularServiceHandler{
		ScheduleCollection: mongoClient.Database("cloudwalker").Collection("schedule"),
		RedisConnection:    redisClient,
		EventPublisher:     publisher}

	_, err = c.AddFunc("6-9,9-12,12-15,15-18,18-21,21-24,24-6 * * *", func() {
		log.Info("Cron Triggered " + string(time.Now().Hour()))
		scheduleCollection := mongoClient.Database("cloudwalker").Collection("schedule")
		cur, err := scheduleCollection.Find(context.TODO(), bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(context.TODO()) {
			var cloudwalkerSchedule SchedularService.CloudwalkerScheduler
			err = cur.Decode(&cloudwalkerSchedule)
			if err != nil {
				log.Fatal(err)
			}
			err = publisher.Publish(context.TODO(), &cloudwalkerSchedule)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	//staring cron job
	c.Start()

	//Register Handler
	err = SchedularService.RegisterSchedularServiceHandler(service.Server(), &h)
	if err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
	log.Info("Stopping cron")
	defer c.Stop()
}

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to redis %v", err)
	}
	return client
}
