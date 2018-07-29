package gitlab

import (
	"context"
	"errors"

	// external
	"github.com/hoop33/entrevista"
	"github.com/k0kubun/pp"
	"github.com/xanzy/go-gitlab"

	// internal
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/model"
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/service"
)

// Gitlab represents the Gitlab service
type Gitlab struct {
}

// Login logs in to Gitlab
func (g *Gitlab) Login(ctx context.Context) (string, error) {
	interview := service.CreateInterview()
	interview.Questions = []entrevista.Question{
		{
			Key:      "token",
			Text:     "Enter your GitLab API token",
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
func (g *Gitlab) GetStars(ctx context.Context, starChan chan<- *model.StarResult, token string, user string) {
	client := g.getClient(token)

	currentPage := 1
	lastPage := 1

	opts := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
		Archived: gitlab.Bool(true),
		OrderBy:  gitlab.String("updated_at"),
		Sort:     gitlab.String("desc"),
		// Search:     string("query"),
		Statistics: gitlab.Bool(true),
		Simple:     gitlab.Bool(true),
		Starred:    gitlab.Bool(true),
		Visibility: gitlab.Visibility(gitlab.PublicVisibility),
	}

	for currentPage <= lastPage {
		opts.ListOptions.Page = currentPage
		projects, response, err := client.Projects.ListProjects(opts)
		// If we got an error, put it on the channel
		if err != nil {
			starChan <- &model.StarResult{
				Error: err,
				Star:  nil,
			}
		} else {
			// Set last page only if we didn't get an error
			lastPage = response.TotalPages
			// Create a Star for each repository and put it on the channel
			for _, project := range projects {
				pp.Println(*project)
				star, err := model.NewStarFromGitlab(*project)
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
func (g *Gitlab) GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int) {
	eventChan <- &model.EventResult{
		Error: errors.New("GitLab not yet supported"),
		Event: nil,
	}
	close(eventChan)
}

// GetTrending returns the trending repositories
func (g *Gitlab) GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token string, language string, verbose bool) {
	close(trendingChan)
}

func (g *Gitlab) getClient(token string) *gitlab.Client {
	return gitlab.NewClient(nil, token)
}

func init() {
	service.RegisterService(&Gitlab{})
}
