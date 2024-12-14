package segments

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
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
	raw, err := ioutil.ReadFile(pkgfile)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(raw, &pkg)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(pkg.Version)
}

func Node(theme config.Theme) []segment {
	nodeVersion := getNodeVersion()
	packageVersion := getPackageVersion()

	segments := segment{}

	if nodeVersion != "" {
		segments = append(segments, segment{
			Name:       "node",
			Content:    p.symbols.NodeIndicator + " " + nodeVersion,
			Foreground: theme.NodeVersionFg,
			Background: theme.NodeVersionBg,
		})
	}

	if packageVersion != "" {
		segments = append(segments, segment{
			Name:       "node-segment",
			Content:    packageVersion + " " + p.symbols.NodeIndicator,
			Foreground: theme.NodeFg,
			Background: theme.NodeBg,
		})
	}

	return segments
}
