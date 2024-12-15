package segments

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

const pkgfile = "./package.json"

type packageJSON struct {
	Version string `json:"version"`
}

func getNodeVersion() string {
	out, err := exec.Command("node", "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(out), "\n")
}

func getPackageVersion() string {
	stat, err := os.Stat(pkgfile)
	if err != nil {
		return ""
	}
	if stat.IsDir() {
		return ""
	}
	pkg := packageJSON{""}
	raw, err := os.ReadFile(pkgfile)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(raw, &pkg)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(pkg.Version)
}

func Node(cfg config.Config, align config.Alignment) []Segment {
	nodeVersion := getNodeVersion()
	packageVersion := getPackageVersion()

	segments := []Segment{}

	if nodeVersion != "" {
		segments = append(segments, Segment{
			Name:       "node",
			Content:    cfg.Symbols().NodeIndicator + " " + nodeVersion,
			Foreground: cfg.Theme.NodeVersionFg,
			Background: cfg.Theme.NodeVersionBg,
		})
	}

	if packageVersion != "" {
		segments = append(segments, Segment{
			Name:       "node-segment",
			Content:    packageVersion + " " + cfg.Symbols().NodeIndicator,
			Foreground: cfg.Theme.NodeFg,
			Background: cfg.Theme.NodeBg,
		})
	}

	return segments
}
