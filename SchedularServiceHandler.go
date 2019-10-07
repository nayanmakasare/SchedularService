package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	SchedularService "github.com/nayanmakasare/SchedularService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type SchedularServiceHandler struct {
	MongoCollection *mongo.Collection
	RedisConnection *redis.Client
	EventPublisher micro.Publisher
}

func(h *SchedularServiceHandler) ApplySchedule(ctx context.Context, req *SchedularService.ApplyRequest, res *SchedularService.ServiceResponse) error {
	result := h.MongoCollection.FindOne(ctx, bson.D{{"vendor", req.Vendor}, {"brand", req.Brand}})
	if result.Err() != nil {
		return result.Err()
	}
	var schedule SchedularService.Schedule
	err := result.Decode(&schedule)
	if err != nil {
		return err
	}
	log.Info("Got schedule for vendor ", req.Vendor, req.Brand)
	err = h.EventPublisher.Publish(ctx, &schedule)
	if err != nil {
		return err
	}
	log.Info("Published successfully ")
	res.IsSuccessfull = true
	return nil
}

func MakeRedisKey(badKey string) string {
	return strings.ToLower(strings.Replace(badKey, " ", "", -1))
}

func(h *SchedularServiceHandler) GetSchedule(ctx context. Context, req *SchedularService.RequestSchedule, res *SchedularService.UserScheduleResponse) error{
	log.Info("Get Schedule ", req.Vendor, req.Brand)
	for _,i := range h.RedisConnection.SMembers(MakeRedisKey(req.Vendor+":"+req.Brand+":schedule")).Val(){
		log.Info("Get Schedule 1 ", i)
		var page SchedularService.LauncherPage
		//cvte:shinko:PageName
		page.PageName = strings.Split(i, ":")[2]
		log.Info("Get Schedule 2 ", page.PageName)
		for _, j := range h.RedisConnection.SMembers(MakeRedisKey(i+":carousel")).Val(){
			log.Info("Get Schedule 3 ", j)
			var pageCarosuel SchedularService.Carousel
			proto.Unmarshal([]byte(j), &pageCarosuel)
			log.Info(pageCarosuel.ImageUrl)
			page.Carousel = append(page.Carousel, &pageCarosuel)
		}
		for _, k := range h.RedisConnection.SMembers(MakeRedisKey(i+":rows")).Val(){
			var pageRow SchedularService.LauncherRows
			// cvte:shinko:PageName:RowName
			pageRow.RowName = strings.Split(k, ":")[3]
			log.Info("rowId ", k, "rowName ", pageRow.RowName)
			pageRow.RowId = k
			page.LauncherRows = append(page.LauncherRows, &pageRow)
		}
		res.LauncherPage = append(res.LauncherPage, &page)
	}
	return nil
}

func GetBytes(key interface{}) ([]byte) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(key)
	return buf.Bytes()
}

func (h *SchedularServiceHandler) MakeSchedule(ctx context.Context, req *SchedularService.Schedule, res *SchedularService.ServiceResponse)error {
	log.Info("Tiggered org")
	_,err := h.MongoCollection.InsertOne(ctx, req)
	if err != nil{
		return err
	}else {
		res.IsSuccessfull = true
		return nil
	}
}


