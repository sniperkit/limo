package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	// "time"

	"github.com/sniperkit/limo/model"

	"github.com/fatih/color"
	"github.com/hoop33/entrevista"

	"github.com/segmentio/stats"
	"github.com/segmentio/stats/influxdb"
)

// Service represents a service
type Service interface {
	Login(ctx context.Context) (string, error)
	GetStars(ctx context.Context, starChan chan<- *model.StarResult, token, user string)
	GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token, language string, verbose bool)
	GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int)
}

var (
	// vcs services
	services = make(map[string]Service)
	// stats
	statsTags    []stats.Tag = []stats.Tag{}
	influxClient *influxdb.Client
	influxConfig influxdb.ClientConfig
	statsEngine  *stats.Engine
)

func init() {
	/*
		influxConfig = influxdb.ClientConfig{
			Database:   "limo-httpstats",
			Address:    "127.0.0.1:8086",
			BufferSize: 2 * 1024 * 1024,
			Timeout:    5 * time.Second,
		}
		influxClient = influxdb.NewClientWith(influxConfig)
		influxClient.CreateDB("limo-httpstats")

		// register engine
		stats.Register(influxClient)
		defer stats.Flush()
	*/
}

func registerService(service Service) {
	services[Name(service)] = service
}

// Name returns the name of a service
func Name(service Service) string {
	parts := strings.Split(reflect.TypeOf(service).String(), ".")
	return strings.ToLower(parts[len(parts)-1])
}

// ForName returns the service for a given name, or an error if it doesn't exist
func ForName(name string) (Service, error) {
	if service, ok := services[strings.ToLower(name)]; ok {
		return service, nil
	}
	return &NotFound{}, fmt.Errorf("service '%s' not found", name)
}

func createInterview() *entrevista.Interview {
	interview := entrevista.NewInterview()
	interview.ShowOutput = func(message string) {
		fmt.Print(color.GreenString(message))
	}
	interview.ShowError = func(message string) {
		color.Red(message)
	}
	return interview
}
