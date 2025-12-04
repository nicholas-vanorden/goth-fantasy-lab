package handlers

import (
	"goth-ffb-players/internal/models"
	"goth-ffb-players/internal/services/auth"
	api "goth-ffb-players/internal/services/yahooapi"
	"goth-ffb-players/web/components"
	"log"
	"net/http"
	"sort"

	"golang.org/x/oauth2"
)

func OAuthLogin(authCache *auth.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUrl := authCache.Config().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		log.Println("Redirecting to: ", authUrl)
		http.Redirect(w, r, authUrl, http.StatusFound)
	}
}

func OAuthCallback(authCache *auth.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			log.Println("No auth code received in callback")
			http.Redirect(w, r, "/oauth/login", http.StatusFound)
		}

		log.Println("Received auth code: ", code)

		token, err := authCache.ExchangeCode(r.Context(), code)
		if err != nil {
			log.Printf("Token exchange failed: %v", err)
			http.Redirect(w, r, "/oauth/login", http.StatusFound)
		}

		log.Println("Access Token:", token.AccessToken[:50]+"...")
		log.Println("Refresh Token:", token.RefreshToken)
		log.Println("Expires:", token.Expiry)

		http.Redirect(w, r, "/players", http.StatusFound)
	}
}

func Players(authCache *auth.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := authCache.Token(r.Context())
		if err != nil || token == nil || !token.Valid() {
			log.Println("Authentication required.")
			http.Redirect(w, r, "/oauth/login", http.StatusFound)
		}

		client, err := authCache.Client(r.Context())
		if err != nil {
			log.Printf("Failed to create authenticated client: %v", err)
			http.Redirect(w, r, "/oauth/login", http.StatusFound)
		}

		statDefs, err := api.LoadStatDefinitions(client)
		if err != nil {
			http.Error(w, "Failed to fetch stat definitions: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Loaded %d stat definitions:", len(statDefs))
		keys := make([]int, 0, len(statDefs))
		for id := range statDefs {
			keys = append(keys, id)
		}
		sort.Ints(keys)
		for _, key := range keys {
			log.Printf("- ID %d: %s (%s)", key, statDefs[key].DisplayName, statDefs[key].Name)
		}

		players, err := api.FetchPlayers(client)
		if err != nil {
			http.Error(w, "Failed to fetch players: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Players:")
		sort.Slice(players, func(i, j int) bool {
			return *players[i].SeasonPoints > *players[j].SeasonPoints
		})
		for _, p := range players {
			playerXML := []byte("<player>" + p.StatsXML.XML + "</player>")
			stats, _ := api.ParsePlayerStats(playerXML, statDefs)
			log.Printf("- %s (%s), Team: %s, Position: %s, SeasonPoints: %.2f",
				p.Name.Full, p.PlayerKey, p.TeamAbr, p.Position, *p.SeasonPoints)
			for statID, value := range stats {
				log.Printf("    - %s: %.2f", statID, value)
			}
		}

		component := components.Players(players)
		component.Render(r.Context(), w)
	}
}

func SearchPlayers(authCache *auth.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//query := strings.ToLower(r.URL.Query().Get("search"))
		var results []models.Player

		// TODO: cache priviously fetched players instead of calling API again
		// for _, player := range models.Players {
		// 	if strings.Contains(strings.ToLower(player.FirstName), query) ||
		// 		strings.Contains(strings.ToLower(player.LastName), query) {
		// 		results = append(results, player)
		// 	}
		// }

		components.PlayerList(results).Render(r.Context(), w)
	}
}

func PlayerDetail(authCache *auth.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// id := r.PathValue("id")

		// TODO: cache priviously fetched players instead of calling API again
		// for _, player := range models.Players {
		// 	if player.Id == id {
		// 		component := components.PlayerDetail(player)
		// 		component.Render(r.Context(), w)
		// 		return
		// 	}
		// }
		http.NotFound(w, r)
	}
}
