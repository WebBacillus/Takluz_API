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
	apiHost := os.Getenv("RIOT_HOST")
	responseData, err := utils.GetRankDetail(puuid, apiKey, apiHost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to perform request: %v", err), http.StatusInternalServerError)
		return
	}
	var playerRank model.Rank
	for i := 0; i < len(responseData); i++ {
		if responseData[i].QueueType == "RANKED_SOLO_5x5" {
			playerRank = responseData[i]
			break
		}
	}
	text := playerRank.Tier + ": " + strconv.Itoa(playerRank.Lp) + " lp"
	html := fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<head>
		<title>Player Rank</title>
		<style>
			.rank-text {
				font-size: 40px;
				color: white;
				background-color: black;
			}
		</style>
		<script>
			// JavaScript to refresh the page every 10 seconds.
			setTimeout(function() {
				window.location.reload(1); // Force a reload from the server
				console.log("refresh")
			}, 60000); // 60000 milliseconds = 60 seconds
		</script>
	</head>
	<body>
		<p class="rank-text">%s</p>
	</body>
	</html>`, text)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}
