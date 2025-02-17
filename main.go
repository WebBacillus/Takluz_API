package main

import (
	"Takluz_API/model"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func getGameDetail(matchID string, API_KEY string) (model.FullGameDetail, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "sea.api.riotgames.com",
		Path:   "/lol/match/v5/matches",
	}
	u.Path, _ = url.JoinPath(u.Path, matchID)
	q := u.Query()
	q.Add("api_key", API_KEY)
	u.RawQuery = q.Encode()

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return model.FullGameDetail{}, err
	}
	resp, _ := client.Do(req)
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.FullGameDetail{}, err
	}
	var gameDetail model.FullGameDetail
	if err := json.Unmarshal(content, &gameDetail); err != nil {
		return model.FullGameDetail{}, err
	}

	return gameDetail, nil

}

func getMatch(puuid string, time_stamp int64, gameType string, start int, count int, API_KEY string) ([]string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "sea.api.riotgames.com",
		Path:   "/lol/match/v5/matches/by-puuid",
	}
	u.Path, _ = url.JoinPath(u.Path, puuid, "ids")
	q := u.Query()
	q.Add("startTime", strconv.Itoa(int(time_stamp)))
	q.Add("type", gameType)
	q.Add("start", strconv.Itoa(start))
	q.Add("count", strconv.Itoa(count))
	q.Add("api_key", API_KEY)
	u.RawQuery = q.Encode()
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, _ := client.Do(req)
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var gameString []string
	if err := json.Unmarshal(content, &gameString); err != nil {
		return nil, err
	}
	return gameString, nil
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("demo-viper")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	app := fiber.New()
	app.Get("/getWinLOL", func(c *fiber.Ctx) error {
		winCount := 0
		lossCount := 0
		gameString, err := getMatch(viper.GetString("PUUID"), time.Now().Add(-10*time.Hour).Unix(), "ranked", 0, 20, viper.GetString("API_KEY"))
		if err != nil {
			panic(err.Error())
		}
		for i := 0; i < len(gameString); i++ {
			println(gameString[i])
			match, _ := getGameDetail(gameString[i], viper.GetString("API_KEY"))
			playerList := match.Info.Participants
			for j := 0; j < 10; j++ {
				if playerList[j].Puuid == viper.GetString("PUUID") {
					// fmt.Println(playerList[j])
					if playerList[j].Win {
						winCount += 1
					} else {
						lossCount += 1
					}
				}
			}
		}
		s := "Win:" + strconv.Itoa(winCount) + " Loss:" + strconv.Itoa(lossCount)
		return c.SendString(s)
	})
	app.Listen("localhost:4444")
}
