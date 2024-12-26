package dd

type teamID int32
type playerID int

type Strategy interface {
	PlayTile() (Tile, int)
}

type Rule struct {
	// Define rule-specific fields
}

type Game struct {
	rules         []Rule
	teams         []Team
	players       []*Player
	totalScore    map[teamID]int
	scoresPerHand []map[teamID]int
	table         *Table
}

type Board struct {
	tile    Tile
	left    *Tile
	right   *Tile
	blocked bool
}

type Table struct {
	boards        []*Board
	shuffledTiles []*Tile
	playersHand   map[playerID][]*Hand
	currentPlayer *Player
}

type Hand struct {
	tiles []Tile
}

type Tile struct {
	Left     int
	Right    int
	IsDouble bool
}

type Team struct {
	id      teamID
	players []*Player
	score   int
	won     bool
}

type Player struct {
	id      playerID
	name    string
	hand    *Hand
	partner *Player
}
