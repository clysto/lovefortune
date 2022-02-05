package plugin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

var (
	URL = "https://devapi.qweather.com/v7/"
	KEY = os.Getenv("QWEATHER_KEY")
)

type WeatherPlugin struct {
	key string
}

func NewWeatherPlugin(key string) *WeatherPlugin {
	return &WeatherPlugin{
		key: key,
	}
}

type WeatherResult struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Now        struct {
		ObsTime   string `json:"obsTime"`
		Temp      string `json:"temp"`
		FeelsLike string `json:"feelsLike"`
		Icon      string `json:"icon"`
		Text      string `json:"text"`
		Wind360   string `json:"wind360"`
		WindDir   string `json:"windDir"`
		WindScale string `json:"windScale"`
		WindSpeed string `json:"windSpeed"`
		Humidity  string `json:"humidity"`
		Precip    string `json:"precip"`
		Pressure  string `json:"pressure"`
		Vis       string `json:"vis"`
		Cloud     string `json:"cloud"`
		Dew       string `json:"dew"`
	} `json:"now"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

func (plugin *WeatherPlugin) Name() string {
	return "Weather Plugin"
}

func (plugin *WeatherPlugin) Funcs() template.FuncMap {
	return template.FuncMap{
		"weather": plugin.Weather,
	}
}

func (plugin *WeatherPlugin) Weather(location string) *WeatherResult {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, URL+"weather/now", nil)
	q := req.URL.Query()
	q.Add("key", KEY)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var result WeatherResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil
	}

	return &result
}
