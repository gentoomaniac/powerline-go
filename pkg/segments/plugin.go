package segments

import (
	"encoding/json"
	"os/exec"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func Plugin(theme config.Theme, plugin string) (segment, bool) {
	output, err := exec.Command("powerline-go-" + plugin).Output()
	if err != nil {
		return nil, false
	}
	segments := segment{}
	err = json.Unmarshal(output, &segments)
	if err != nil {
		// The plugin was found but no valid data was returned. Ignore it
		return []segment{}, true
	}
	return segments, true
}
