package yahooapi

import (
	"encoding/xml"
	"fmt"
	m "goth-ffb-players/internal/models"
	"io"
	"log"
	"net/http"
	"os"
)

func FetchPlayers(client *http.Client) ([]m.Player, error) {
	leagueKeys, err := fetchUserLeagueKeys(client)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		os.Getenv("YAHOO_FANTASY_API_PLAYERS"),
		leagueKeys[0],
	)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var presp m.PlayersResponse
	if err := xml.Unmarshal(body, &presp); err != nil {
		log.Println("Yahoo API error:", string(body))
		return nil, err
	}

	return presp.League.Players.Players, nil
}

func fetchUserLeagueKeys(client *http.Client) ([]string, error) {
	url := os.Getenv("YAHOO_FANTASY_API_LEAGUES")

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var uresp m.UsersResponse
	if err := xml.Unmarshal(body, &uresp); err != nil {
		log.Println("Yahoo API error:", string(body))
		return nil, err
	}

	var leagueKeys []string
	for _, user := range uresp.Users {
		for _, game := range user.Games {
			for _, league := range game.Leagues {
				leagueKeys = append(leagueKeys, league.LeagueKey)
			}
		}
	}
	return leagueKeys, nil
}
