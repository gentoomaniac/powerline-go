package segments

import (
	"encoding/json"
	"os/exec"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Plugin(theme config.Theme, plugin string) ([]Segment, bool) {
	output, err := exec.Command("powerline-go-" + plugin).Output()
	if err != nil {
		return nil, false
	}
	segments := []Segment{}
	err = json.Unmarshal(output, &segments)
	if err != nil {
		// The plugin was found but no valid data was returned. Ignore it
		return []Segment{}, true
	}
	return segments, true
}
