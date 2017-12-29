package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	// "log"

	"github.com/sniperkit/limo/model"

	// cache
	"github.com/sniperkit/httpcache"
	// "github.com/gregjones/httpcache"
	"github.com/sniperkit/httpcache/badgercache" // BadgerKV default implementation

	// stats
	"github.com/segmentio/stats/httpstats"

	// vcs
	"github.com/google/go-github/github"
	"github.com/hoop33/entrevista"
	"golang.org/x/oauth2"
	// debug
	// "github.com/k0kubun/pp"
)

// Github represents the Github service
type Github struct {
	hcache httpcache.Cache
}

// Login logs in to Github
func (g *Github) Login(ctx context.Context) (string, error) {
	interview := createInterview()
	interview.Questions = []entrevista.Question{
		{
			Key:      "token",
			Text:     "Enter your GitHub API token",
			Required: true,
			Hidden:   true,
		},
	}

	answers, err := interview.Run()
	if err != nil {
		return "", err
	}
	return answers["token"].(string), nil
}

// GetStars returns the stars for the specified user (empty string for authenticated user)
func (g *Github) GetStars(ctx context.Context, starChan chan<- *model.StarResult, token string, user string) {
	client := g.getClient(token)

	// The first response will give us the correct value for the last page
	currentPage := 1
	lastPage := 1

	for currentPage <= lastPage {
		repos, response, err := client.Activity.ListStarred(ctx, user, &github.ActivityListStarredOptions{
			Sort:      "updated",
			Direction: "desc", // desc
			ListOptions: github.ListOptions{
				Page:    currentPage,
				PerPage: 100,
			},
		})

		// If we got an error, put it on the channel
		if err != nil {
			starChan <- &model.StarResult{
				Error: err,
				Star:  nil,
			}
		} else {
			// Set last page only if we didn't get an error
			lastPage = response.LastPage

			// Create a Star for each repository and put it on the channel
			for _, repo := range repos {
				star, err := model.NewStarFromGithub(repo.StarredAt, *repo.Repository)
				starChan <- &model.StarResult{
					Error: err,
					Star:  star,
				}
			}
		}
		// Go to the next page
		currentPage++
	}
	close(starChan)
}

// GetEvents returns the events for the authenticated user
func (g *Github) GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int) {
	client := g.getClient(token)

	currentPage := page
	lastPage := page + count - 1

	for currentPage <= lastPage {
		events, _, err := client.Activity.ListEventsReceivedByUser(ctx, user, false, &github.ListOptions{
			Page:    currentPage,
			PerPage: 100,
		})

		if err != nil {
			eventChan <- &model.EventResult{
				Error: err,
				Event: nil,
			}
		} else {
			for _, event := range events {
				eventChan <- &model.EventResult{
					Error: nil,
					Event: model.NewEventFromGithub(event),
				}
			}
		}
		currentPage++
	}
	close(eventChan)
}

// GetTrending returns the trending repositories
func (g *Github) GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token string, language string, verbose bool) {
	client := g.getClient(token)

	// TODO perhaps allow them to specify multiple pages?
	// Might be overkill -- first page probably plenty

	// TODO Make this more configurable. Sort by stars, forks, default.
	// Search by number of stars, pushed, created, or whatever.
	// Lots of possibilities.

	q := g.getDateSearchString()

	if language != "" {
		q = fmt.Sprintf("language:%s %s", language, q)
	}

	if verbose {
		fmt.Println("q =", q)
	}

	currentPage := 1

	result, _, err := client.Search.Repositories(ctx, q, &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
		ListOptions: github.ListOptions{
			Page:    currentPage,
			PerPage: 100,
		},
	})

	// If we got an error, put it on the channel
	if err != nil {
		trendingChan <- &model.StarResult{
			Error: err,
			Star:  nil,
		}
	} else {
		// Create a Star for each repository and put it on the channel
		for _, repo := range result.Repositories {
			star, err := model.NewStarFromGithub(nil, repo)
			trendingChan <- &model.StarResult{
				Error: err,
				Star:  star,
			}
		}
	}

	close(trendingChan)
}

func (g *Github) getDateSearchString() string {
	// TODO make this configurable
	// Default should be in configuration file
	// and should be able to override from command line
	// TODO should be able to specify whether "created" or "pushed"
	date := time.Now().Add(-7 * (24 * time.Hour))
	return fmt.Sprintf("created:>%s", date.Format("2006-01-02"))
}

func ensureDir(path string) {
	d, err := os.Open(path)
	if err != nil {
		os.MkdirAll(path, os.FileMode(0755))
	}
	d.Close()
}

func (g *Github) getClient(token string) *github.Client {

	cacheStoragePrefixPath := filepath.Join("data", "cache.badger")
	ensureDir(cacheStoragePrefixPath)
	var err error
	if g.hcache == nil {
		g.hcache, err = badgercache.New(
			&badgercache.Config{
				ValueDir:    "api.github.com.v3.gzip", //gzip",
				StoragePath: cacheStoragePrefixPath,
				SyncWrites:  true,
				Compress:    true,
			})
		if err != nil {
			fmt.Println("err: ", err.Error())
			return nil
		}
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	var hc http.Client

	t := httpcache.NewTransport(g.hcache)
	t.MarkCachedResponses = true

	hc.Transport = httpstats.NewTransport(t)
	timeout := time.Duration(10 * time.Second)

	return github.NewClient(&http.Client{
		Transport: &oauth2.Transport{
			Base:   hc.Transport,
			Source: ts,
		},
		Timeout: timeout,
	})
}

func init() {
	registerService(&Github{})
}
