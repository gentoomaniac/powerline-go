package segments

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gentoomaniac/powerline-go/pkg/cache"
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
)

const cacheEntry = "github-notifications"

func GithubNotifications(cfg config.State, align config.Alignment) []Segment {
	content := ""

	if time.Since(cache.Age(cacheEntry)) > time.Minute*5 {
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

		data, err := json.Marshal(notifications)
		if err != nil {
			log.Error().Err(err).Msg("failed writinmg github notifications to cache	")
		}
		cache.Save(cacheEntry, data)

		if res.LastPage != res.FirstPage {
			content = fmt.Sprintf("\uea84 >%d", len(notifications))
		} else if len(notifications) == 0 {
			return []Segment{}
		} else {
			content = fmt.Sprintf("\uea84 %d", len(notifications))
		}
	} else {
		data, err := cache.Read(cacheEntry)
		if err != nil {
			log.Error().Err(err).Msg("could not read github notification data from cache")
		}

		var notifications []github.Notification
		err = json.Unmarshal(data, &notifications)
		if err != nil {
			log.Error().Err(err).Msg("failed to load cached github notifications")
		}
		content = fmt.Sprintf("\uea84 %d", len(notifications))
	}

	return []Segment{{
		Name:       "ghnotifications",
		Content:    content,
		Foreground: cfg.Theme.DockerMachineFg,
		Background: cfg.Theme.DockerMachineBg,
	}}
}
