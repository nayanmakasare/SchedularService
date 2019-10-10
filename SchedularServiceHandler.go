package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro"
	SchedularService "github.com/nayanmakasare/SchedularService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"time"
)

type SchedularServiceHandler struct {
	ScheduleCollection *mongo.Collection
	RedisConnection    *redis.Client
	EventPublisher     micro.Publisher
}

func (h *SchedularServiceHandler) GetSchedule(ctx context.Context, req *SchedularService.RequestSchedule, res *SchedularService.UserScheduleResponse) error {
	for _, i := range h.RedisConnection.SMembers(MakeRedisKey(req.Vendor + ":" + req.Brand + ":" + getTimeZone())).Val() {
		var page SchedularService.LauncherPage
		//cvte:shinko:timezone:pageIndex:pageName
		pageParts := strings.Split(i, ":")
		page.PageName = pageParts[4]
		resultInt, _ := strconv.ParseInt(pageParts[3], 10, 32)
		page.PageIndex = int32(resultInt)
		for _, j := range h.RedisConnection.SMembers(MakeRedisKey(i + ":carousel")).Val() {
			var pageCarosuel SchedularService.Carousel
			err := proto.Unmarshal([]byte(j), &pageCarosuel)
			if err != nil {
				return err
			}
			page.Carousel = append(page.Carousel, &pageCarosuel)
		}
		for _, k := range h.RedisConnection.SMembers(MakeRedisKey(i + ":rows")).Val() {
			var pageRow SchedularService.LauncherRows
			//cvte:shinko:timezone:pageIndex:pageName:rows:rowIndex:rowName
			rowParts := strings.Split(k, ":")
			pageRow.RowName = rowParts[7]
			resultInt, _ = strconv.ParseInt(rowParts[3], 10, 32)
			pageRow.RowIndex = int32(resultInt)
			pageRow.RowId = k
			page.LauncherRows = append(page.LauncherRows, &pageRow)
		}
		res.LauncherPage = append(res.LauncherPage, &page)
	}
	return nil
}

func GetBytes(key interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(key)
	return buf.Bytes()
}

func (h *SchedularServiceHandler) MakeSchedule(ctx context.Context, req *SchedularService.CloudwalkerScheduler, res *SchedularService.ServiceResponse) error {
	filter := bson.D{{"vendor", req.Vendor}, {"brand", req.Brand}}
	singleResult := h.ScheduleCollection.FindOne(ctx, filter)
	if singleResult.Err() != nil {
		_, err := h.ScheduleCollection.InsertOne(ctx, req)
		if err != nil {
			return err
		} else {
			res.IsSuccessfull = true
			return nil
		}
	} else {
		_, err := h.ScheduleCollection.ReplaceOne(ctx, filter, req)
		if err != nil {
			return err
		} else {
			res.IsSuccessfull = true
			return nil
		}
	}
}

func getTimeZone() string {
	currentHour := time.Now().Hour()
	var timeZone string
	if currentHour >= 6 && currentHour < 9 {
		timeZone = "6-9"
	} else if currentHour >= 9 && currentHour < 12 {
		timeZone = "9-12"
	} else if currentHour >= 12 && currentHour < 15 {
		timeZone = "12-15"
	} else if currentHour >= 15 && currentHour < 18 {
		timeZone = "15-18"
	} else if currentHour >= 18 && currentHour < 21 {
		timeZone = "18-21"
	} else if currentHour >= 21 && currentHour < 24 {
		timeZone = "21-24"
	} else if currentHour >= 24 {
		timeZone = "1-6"
	}
	return timeZone
}

func MakeRedisKey(badKey string) string {
	return strings.ToLower(strings.Replace(badKey, " ", "_", -1))
}
