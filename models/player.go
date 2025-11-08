package models

import "fmt"

type Player struct {
	Id        string
	FirstName string
	LastName  string
	Jersey    int
	Position  Position
	Team      string
	Outlook   string
}

func (p Player) String() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}

type Position int

const (
	QB Position = iota
	RB
	WR
	TE
	K
)

func (p Position) String() string {
	switch p {
	case QB:
		return "QB"
	case RB:
		return "RB"
	case WR:
		return "WR"
	case TE:
		return "TE"
	case K:
		return "K"
	default:
		return "Unknown"
	}
}

//temp test data
var Players = []Player{
	{
		Id:        "1",
		FirstName: "Tom",
		LastName:  "Brady",
		Position:  Position(QB),
		Team:      "Tampa Bay Buccaneers",
		Outlook:   "Retired. It's over. Give it up. Just accept it.",
	},
	{
		Id:        "2",
		FirstName: "Derrick",
		LastName:  "Henry",
		Position:  Position(RB),
		Team:      "Baltiomore Ravens",
		Outlook:   "Age is finally catching up to him.",
	},
	{
		Id:        "3",
		FirstName: "Davante",
		LastName:  "Adams",
		Position:  Position(WR),
		Team:      "Los Angeles Rams",
		Outlook:   "Still one of the best in the game.",
	},
}
