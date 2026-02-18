package handler

import (
	"Takluz_API/utils"
	"fmt"
	"time"
)

func GetDailyWinLoss(apiKey, puuid, apiHost string) (int, int, error) {
	winCount := 0
	lossCount := 0
	startTime, endTime, err := utils.GetMiddayTodayUnixThai(time.Now().Add(-6*time.Hour), 12)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting time: %v", err)
	}
	gameString, err := utils.GetMatch(puuid, startTime, endTime, "ranked", 0, 20, apiKey, apiHost)
	if err != nil {
		return 0, 0, fmt.Errorf("error fetching match data: %v", err)
	}
	if len(gameString) == 0 {
		return 0, 0, nil
	}
	for i := 0; i < len(gameString); i++ {
		match, err := utils.GetGameDetail(gameString[i], apiKey, apiHost)
		if err != nil {
			return 0, 0, fmt.Errorf("error fetching match detail: %v", err)
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
	return winCount, lossCount, nil
}
