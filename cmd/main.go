package main

import (
	"fmt"

	"github.com/amb1s1/dd/pkg/dd"
)

func main() {
	rule := []dd.Rule{}
	player := dd.NewPlayer("player1")
	player2 := dd.NewPlayer("player2")
	player3 := dd.NewPlayer("player3")
	player4 := dd.NewPlayer("player4")
	team, _ := dd.NewTeam([]*dd.Player{&player, &player2})
	team2, _ := dd.NewTeam([]*dd.Player{&player3, &player4})

	game := dd.NewGame(rule, []dd.Team{team, team2})
	game.Table.DealHand(game)
	fmt.Printf("%v", game)
	for _, player := range game.Players {
		fmt.Println(player.Name)
		for _, tile := range player.Hand.TilesOnHand {
			fmt.Println(tile.Left, tile.Right, tile.IsDouble)
		}
	}
}
