package utils

import (
	"Takluz_API/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetMiddayTodayUnixThai(setTime time.Time, clock int) (int64, int64, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return 0, 0, fmt.Errorf("error loading timezone: %w", err)
	}

	now := setTime.In(loc)

	year, month, day := now.Date()

	middayToday := time.Date(year, month, day, clock, 0, 0, 0, loc)

	return middayToday.Unix(), (middayToday.Add(18 * time.Hour)).Unix(), nil
}

func GetGameDetail(matchID string, API_KEY string) (model.FullGameDetail, error) {
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
		return model.FullGameDetail{}, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return model.FullGameDetail{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return model.FullGameDetail{}, fmt.Errorf("riot API error: status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.FullGameDetail{}, fmt.Errorf("error reading response body: %w", err)
	}

	var gameDetail model.FullGameDetail
	if err := json.Unmarshal(content, &gameDetail); err != nil {
		return model.FullGameDetail{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return gameDetail, nil

}

func GetMatch(puuid string, startTime int64, endTime int64, gameType string, start int, count int, API_KEY string) ([]string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "sea.api.riotgames.com",
		Path:   "/lol/match/v5/matches/by-puuid",
	}
	u.Path, _ = url.JoinPath(u.Path, puuid, "ids")
	q := u.Query()
	q.Add("startTime", strconv.Itoa(int(startTime)))
	q.Add("endTime", strconv.Itoa(int(endTime)))
	q.Add("type", gameType)
	q.Add("start", strconv.Itoa(start))
	q.Add("count", strconv.Itoa(count))
	q.Add("api_key", API_KEY)
	u.RawQuery = q.Encode()
	client := http.Client{}
	fmt.Println(u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("riot API error: status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var gameString []string
	if err := json.Unmarshal(content, &gameString); err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return gameString, nil
}
func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
