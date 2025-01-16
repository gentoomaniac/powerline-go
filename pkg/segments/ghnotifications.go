package segments

import (
	"context"
	"fmt"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
)

func GithubNotifications(cfg config.State, align config.Alignment) []Segment {
	// TODO: caching! keep a state
	ghClient := github.NewClient(nil).WithAuthToken(cfg.GithubToken)
	if ghClient == nil {
		return []Segment{}
	}

	notifications, res, err := ghClient.Activity.ListNotifications(
		context.Background(),
		&github.NotificationListOptions{
			ListOptions: github.ListOptions{},
		})
	if err != nil {
		log.Error().Err(err).Msg("failed gerting github notifications")
	}

	content := ""
	if res.LastPage != res.FirstPage {
		content = fmt.Sprintf("\uea84 >%d", len(notifications))
	} else if len(notifications) == 0 {
		return []Segment{}
	} else {
		content = fmt.Sprintf("\uea84 %d", len(notifications))
	}

	return []Segment{{
		Name:       "ghnotifications",
		Content:    content,
		Foreground: cfg.Theme.DockerMachineFg,
		Background: cfg.Theme.DockerMachineBg,
	}}
}
