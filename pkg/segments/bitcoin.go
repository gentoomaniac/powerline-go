package segments

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

const apiUrl = "https://api.binance.com/api/v1/ticker/price?symbol=BTCUSDT"

type response struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func Bitcoin(cfg config.State, align config.Alignment) []Segment {
	data, err := getBtcPrice(apiUrl)
	if err != nil {
		log.Error().Err(err).Msg("failed fetching BTC price")
		return []Segment{}
	}

	return []Segment{{
		Name:       "btc",
		Content:    "\\$" + strings.Split(data.Price, ".")[0],
		Foreground: cfg.Theme.HomeFg,
		Background: cfg.Theme.HomeBg,
	}}
}

func getBtcPrice(url string) (*response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := &response{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
