package handler

import (
	"Takluz_API/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func GetMiddayTodayUnixThai(setTime time.Time) (int64, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error loading timezone:", err)
		return 0, err
	}

	now := setTime.In(loc)

	year, month, day := now.Date()

	middayToday := time.Date(year, month, day, 12, 0, 0, 0, loc)

	return middayToday.Unix(), nil
}

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

func Handler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")
	puuid := os.Getenv("PUUID")

	if apiKey == "" || puuid == "" {
		// return c.Status(fiber.StatusInternalServerError).SendString("API_KEY or PUUID not set")
	}

	winCount := 0
	lossCount := 0
	getTime, err := GetMiddayTodayUnixThai(time.Now().Add(-12 * time.Hour))
	if err != nil {
		fmt.Println("xdd")
	}
	gameString, err := getMatch(puuid, getTime, "ranked", 0, 20, apiKey)
	if err != nil {
		fmt.Println("Error in getMatch:", err) // Log the error
		// return c.Status(fiber.StatusInternalServerError).SendString("Error fetching match data")
	}
	if len(gameString) == 0 {
		// return c.Status(fiber.StatusNotFound).SendString("No Match Found")
	}
	for i := 0; i < len(gameString); i++ {
		println(gameString[i])
		match, err := getGameDetail(gameString[i], apiKey)
		if err != nil {
			fmt.Println("Error in getMatchDetail:", err)
			// return c.Status(fiber.StatusInternalServerError).SendString("Error fetching match detail")
		}
		playerList := match.Info.Participants
		for j := 0; j < 10; j++ {
			if playerList[j].Puuid == puuid {
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
	// fmt.Fprintf(w, s)
	// Create a response struct.  This is much better than a concatenated string.
	response := struct {
		Wins   int    `json:"wins"`
		Losses int    `json:"losses"`
		String string `json:"string"`
	}{
		Wins:   winCount,
		Losses: lossCount,
		String: s,
	}

	// Marshal the response to JSON.
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json.  This is *essential*.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //  Explicitly set the status code.
	w.Write(jsonResponse)        // Write the JSON response.
}
