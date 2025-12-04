package yahooapi

import (
	"encoding/xml"
	"fmt"
	m "goth-ffb-players/internal/models"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var (
	cache     map[int]m.StatDefinition
	cacheOnce sync.Once
)

func LoadStatDefinitions(client *http.Client) (map[int]m.StatDefinition, error) {
	var err error

	cacheOnce.Do(func() {
		cache, err = loadStatCategories(client)
	})

	return cache, err
}

func ParsePlayerStats(statsXML []byte, defs map[int]m.StatDefinition) (map[string]float64, error) {
	var s m.Stats
	if err := xml.Unmarshal(statsXML, &s); err != nil {
		return nil, err
	}

	result := map[int]float64{}
	for _, st := range s.Stats {
		id, _ := strconv.Atoi(st.StatID)
		v, _ := strconv.ParseFloat(st.Value, 64)
		result[id] = v
	}

	return translateStats(result, defs), nil
}

func translateStats(raw map[int]float64, defs map[int]m.StatDefinition) map[string]float64 {
	result := map[string]float64{}

	for id, value := range raw {
		def, ok := defs[id]
		if !ok {
			continue // skip unknown IDs
		}

		name := def.DisplayName
		if name == "" {
			name = def.Name
		}

		result[name] = value
	}

	return result
}

func loadStatCategories(client *http.Client) (map[int]m.StatDefinition, error) {
	url := os.Getenv("YAHOO_FANTASY_API_STAT_CATEGORIES_URL")

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching stat_categories: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sc m.StatCategories
	if err := xml.Unmarshal(body, &sc); err != nil {
		return nil, fmt.Errorf("xml unmarshal: %w", err)
	}

	mp := make(map[int]m.StatDefinition)

	for _, s := range sc.Game.StatCategoriesList.Stats {
		mp[s.ID] = m.StatDefinition{
			StatID:        s.ID,
			Name:          s.Name,
			DisplayName:   s.DisplayName,
			SortOrder:     s.SortOrder,
			PositionTypes: s.PositionTypes,
			IsModifiable:  s.IsModifiable == 1,
			Category:      s.Category,
		}
	}

	return mp, nil
}
