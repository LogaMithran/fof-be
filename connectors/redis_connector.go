package connectors

import (
	"context"
	"fmt"
	"friends-of-friends-be/services"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

var Client *redis.Client

func ConnectToRedis() (bool, error) {
	opt, parseErr := redis.ParseURL(os.Getenv("REDIS_LOCAL"))
	if parseErr != nil {
		log.Fatalf("Error in parsing the url")
	}
	Client = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()
	if cmd := Client.Ping(ctx); cmd != nil {
		fmt.Println("Pong response", cmd.Val())
		return true, nil
	}

	return false, fmt.Errorf("Error in connecting to redis")
}

func GeolocationAdd(key string, location services.Location, name string) {
	println("NAME", name, key, location.Longitude, location.Latitude)
	if res := Client.GeoAdd(context.Background(), key, &redis.GeoLocation{
		Longitude: location.Longitude,
		Latitude:  location.Latitude,
		Name:      name,
	}); res != nil {
		fmt.Println("Geo add response", res.Val())
	}
}

func GeoSearch(key string, lat float64, lng float64, radius float64) []redis.GeoLocation {
	response := Client.GeoRadius(context.Background(), key, lng, lat, &redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        "km",
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
	})

	locations, err := response.Result()
	if err != nil {
		log.Println("Error in getting the location")
		return nil
	}

	return locations
}

func Set(key string, value interface{}, expiry time.Duration) {
	response := Client.Set(context.Background(), key, value, expiry)

	println(response)
}

func Get(key string) string {
	println(key)
	response := Client.Get(context.Background(), key)

	val, err := response.Result()
	if val == "" {
		println("Error in getting the result", err.Error(), val)
		return ""
	}

	return val
}

func Subscribe(channel string) {
	Client.Subscribe(context.Background(), channel)

}

func Publish(channel string, message interface{}) {
	Client.Publish(context.Background(), "", message)
}
