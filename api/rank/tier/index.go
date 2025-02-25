package handler

import (
	"Takluz_API/model"
	"Takluz_API/utils"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")
	puuid := os.Getenv("PUUID")
	responseData, err := utils.GetRankDetail(puuid, apiKey)
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
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	tierList := []string{"IRON", "BRONZE", "SILVER", "GOLD", "PLATINUM", "EMERALD", "DIAMOND", "MASTER", "GRANDMASTER"}
	tierEmoteList := []string{"LeagueIron", "LeagueBronze", "LeagueSilver", "LeagueGold", "LeaguePlatinum", "LeagueEmerald", "LeagueDiamond", "LeagueMaster", "LeagueGrandmaster"}
	index := 0
	for i := 0; i < len(tierList); i++ {
		if tierList[i] == playerRank.Tier {
			index = i
		}
	}
	text := tierEmoteList[index] + " " + tierList[index] + ": " + strconv.Itoa(playerRank.Lp) + " lp"
	fmt.Fprint(w, text)
}
