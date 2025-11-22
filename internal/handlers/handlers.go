package handlers

import (
	"goth-ffb-players/internal/models"
	"goth-ffb-players/web/components"
	"net/http"
	"strings"
)

func PlayerList(w http.ResponseWriter, r *http.Request) {
	component := components.Players(models.Players)
	component.Render(r.Context(), w)
}

func PlayerDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for _, player := range models.Players {
		if player.Id == id {
			component := components.PlayerDetail(player)
			component.Render(r.Context(), w)
			return
		}
	}
	http.NotFound(w, r)
}

func SearchPlayers(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("search"))
	var results []models.Player

	for _, player := range models.Players {
		if strings.Contains(strings.ToLower(player.FirstName), query) ||
			strings.Contains(strings.ToLower(player.LastName), query) {
			results = append(results, player)
		}
	}

	components.PlayerList(results).Render(r.Context(), w)
}
