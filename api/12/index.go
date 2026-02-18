package handler

import (
	"Takluz_API/model"
	"Takluz_API/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")
	puuid := os.Getenv("PUUID")
	apiHost := os.Getenv("RIOT_HOST")

	responseData, err := utils.GetRankDetail(puuid, apiKey, apiHost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to perform request: %v", err), http.StatusInternalServerError)
		return
	}
	var playerRank model.Rank
	for i := 0; i < len(responseData); i++ {
		if responseData[i].Puuid == puuid && responseData[i].LeagueID == "cf9aeb2e-475c-4ff1-80c8-6d26b9904207" {
			playerRank = responseData[i]
			break
		}
	}
	tierList := []string{"IRON", "BRONZE", "SILVER", "GOLD", "PLATINUM", "EMERALD", "DIAMOND", "MASTER", "GRANDMASTER"}
	tierEmoteList := []string{"LeagueIron", "LeagueBronze", "LeagueSilver", "LeagueGold", "LeaguePlatinum", "LeagueEmerald", "LeagueDiamond", "LeagueMaster", "LeagueGrandmaster"}
	index := 0
	for i := 0; i < len(tierList); i++ {
		if tierList[i] == playerRank.Tier {
			index = i
		}
	}
	text := tierEmoteList[index] + " " + tierList[index] + ": " + strconv.Itoa(playerRank.Lp) + " lp"

	if apiKey == "" || puuid == "" || apiHost == "" {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "API_KEY, PUUID, or RIOT_HOST not set",
		}
		utils.SendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	winCount := 0
	lossCount := 0
	startTime, endTime, err := utils.GetMiddayTodayUnixThai(time.Now().Add(-6*time.Hour), 12)
	if err != nil {
		response := struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("Error getting time: %v", err),
		}
		utils.SendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	gameString, err := utils.GetMatch(puuid, startTime, endTime, "ranked", 0, 20, apiKey, apiHost)
	if err != nil {
		response := struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("Error fetching match data: %v", err),
		}
		utils.SendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	if len(gameString) == 0 {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "Win: 0 | Loss: 0 || " + text,
		}
		utils.SendJSONResponse(w, response, http.StatusOK)
		return
	}
	for i := 0; i < len(gameString); i++ {
		println(gameString[i])
		match, err := utils.GetGameDetail(gameString[i], apiKey, apiHost)
		if err != nil {
			response := struct {
				Message string `json:"message"`
			}{
				Message: fmt.Sprintf("Error fetching match detail: %v", err),
			}
			utils.SendJSONResponse(w, response, http.StatusInternalServerError)
			return
		}
		playerList := match.Info.Participants
		for j := 0; j < 10; j++ {
			if playerList[j].Puuid == puuid {
				if playerList[j].Win {
					winCount += 1
				} else {
					lossCount += 1
				}
			}
		}
	}

	m := "Win: " + strconv.Itoa(winCount) + " | Loss: " + strconv.Itoa(lossCount) + " || " + text
	response := struct {
		Wins    int    `json:"wins"`
		Losses  int    `json:"losses"`
		Message string `json:"message"`
	}{
		Wins:    winCount,
		Losses:  lossCount,
		Message: m,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
