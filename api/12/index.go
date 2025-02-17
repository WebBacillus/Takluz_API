package handler

import (
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

	if apiKey == "" || puuid == "" {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "API_KEY or PUUID not set",
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
	gameString, err := utils.GetMatch(puuid, startTime, endTime, "ranked", 0, 20, apiKey)
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
			Message: "No Match Found",
		}
		utils.SendJSONResponse(w, response, http.StatusOK)
		return
	}
	for i := 0; i < len(gameString); i++ {
		println(gameString[i])
		match, err := utils.GetGameDetail(gameString[i], apiKey)
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
	m := "Win:" + strconv.Itoa(winCount) + " Loss:" + strconv.Itoa(lossCount)
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
