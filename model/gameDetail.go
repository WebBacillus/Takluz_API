package model

type Participant struct {
	Puuid        string `json:"puuid"`
	ChampionName string `json:"championName"`
	TeamId       int    `json:"teamId"` // Added TeamId for easier win checking
	TeamPosition string `json:"teamPosition"`
	Win          bool   `json:"win"`
}

type Team struct {
	TeamId int  `json:"teamId"`
	Win    bool `json:"win"`
}

type Match struct {
	Participants []Participant `json:"participants"`
	Teams        []Team        `json:"teams"` // Include teams for easier access
}

type Info struct {
	Participants []Participant `json:"participants"`
	Teams        []Team        `json:"teams"`
}

type Rank struct {
	LeagueID string `json:"leagueId"`
	Puuid    string `json:"puuid"`
	Tier     string `json:"tier"`
	Rank     string `json:"rank"`
	Lp       int    `json:"leaguePoints"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
}

type FullGameDetail struct {
	Info Info `json:"info"`
}
