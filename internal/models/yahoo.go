package models

import (
	"encoding/xml"
)

type UsersResponse struct {
	XMLName xml.Name `xml:"fantasy_content"`
	Users   []User   `xml:"users>user"`
}

type User struct {
	Games []Game `xml:"games>game"`
}

type Game struct {
	GameKey string   `xml:"game_key"`
	Leagues []League `xml:"leagues>league"`
}

type League struct {
	LeagueKey string `xml:"league_key"`
}

type PlayersResponse struct {
	XMLName xml.Name `xml:"fantasy_content"`
	League  struct {
		Players struct {
			Players []Player `xml:"player"`
		} `xml:"players"`
	} `xml:"league"`
}

type Player struct {
	StatsXML    PlayerStats `xml:"player_stats"`
	PlayerKey   string      `xml:"player_key"`
	PlayerID    string      `xml:"player_id"`
	Name        Name        `xml:"name"`
	Jersey      string      `xml:"uniform_number"`
	Position    string      `xml:"display_position"`
	TeamAbr     string      `xml:"editorial_team_abbr"`
	Team        string      `xml:"editorial_team_full_name"`
	HeadshotUrl string      `xml:"headshot>url"`
	// Stats/points fields if available:
	SeasonPoints *float64 `xml:"player_points>total,omitempty"`
}

type PlayerStats struct {
	XML string `xml:",innerxml"`
}

type Name struct {
	Full string `xml:"full"`
}

type StatDefinition struct {
	StatID        int
	Name          string
	DisplayName   string
	SortOrder     int
	PositionTypes []string
	IsModifiable  bool
	Category      string
}

type StatCategories struct {
	XMLName xml.Name `xml:"fantasy_content"`
	Game    StatGame `xml:"game"`
}

type StatGame struct {
	StatCategoriesList StatCategoriesList `xml:"stat_categories"`
}

type StatCategoriesList struct {
	Stats []StatXML `xml:"stats>stat"`
}

type StatXML struct {
	ID            int      `xml:"stat_id"`
	Name          string   `xml:"name"`
	DisplayName   string   `xml:"display_name"`
	SortOrder     int      `xml:"sort_order"`
	PositionTypes []string `xml:"position_types>position_type"`
	IsModifiable  int      `xml:"is_modifiable"`
	Category      string   `xml:"category"`
}

type Stat struct {
	StatID string `xml:"stat_id"`
	Value  string `xml:"value"`
}

type Stats struct {
	Stats []Stat `xml:"stats>stat"`
}
