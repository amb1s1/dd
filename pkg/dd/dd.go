package dd

import (
	"time"

	klog "github.com/go-kit/log"
	"golang.org/x/exp/rand"
	"k8s.io/klog"
)

var (
	dbPlayerID = map[uint64]bool{}
	dbTeamID   = map[uint64]bool{}
	tileDigits = []int{0, 1, 2, 3, 4, 5, 6}
)

type teamID uint64
type playerID uint64

type Strategy interface {
	PlayTile() (tile, int)
}

type Rule struct {
	// Define rule-specific fields
}

type Game struct {
	rules         []Rule
	teams         []*team
	players       map[playerID]*player
	totalScore    map[teamID]int
	scoresPerHand []map[teamID]int
	table         *table
	isBeggining   bool
}

type board struct {
	tile    tile
	left    *tile
	right   *tile
	blocked bool
}

type table struct {
	boards        []*board
	shuffledTiles []*tile
	playersID     []playerID
	playersHand   map[playerID][]*hand
	currentPlayer *player
}

type hand struct {
	tiles []tile
}

type tile struct {
	Left     int
	Right    int
	IsDouble bool
}

type team struct {
	id      teamID
	name    string
	players []*player
	score   int
	won     bool
}

type player struct {
	id      playerID
	name    string
	hand    *hand
	partner *player
}

func NewGame() *Game {
	return &Game{
		rules:         make([]Rule, 0),
		players:       make(map[playerID]*player),
		totalScore:    make(map[teamID]int),
		scoresPerHand: make([]map[teamID]int, 0),
		table:         newTable(),
		isBeggining:   true,
	}
}

func NewBoard() *board {
	return &board{
		tile:    tile{},
		left:    &tile{},
		right:   &tile{},
		blocked: false,
	}
}

func newTable() *table {
	t := &table{
		boards:        make([]*board, 0),
		playersHand:   make(map[playerID][]*hand),
		currentPlayer: &player{},
		playersID:     make([]playerID, 0),
	}
	t.shuffleTiles()
	return t
}

func (t *table) shuffleTiles() {
	left := shuffle(tileDigits)
	right := shuffle(left)
	tiles := make([]*tile, 28)
	counter := 0
	for _, l := range left {
		for _, r := range right {
			tiles[counter] = newTile(l, r)
			counter++
		}
		right = right[1:]
	}
	t.shuffledTiles = tiles
}

func shuffle(digits []int) []int {
	rand.Seed(uint64(time.Now().UnixNano()))
	rand.Shuffle(len(digits), func(i, j int) { digits[i], digits[j] = digits[j], digits[i] })
	return digits
}

func newTile(left, right int) *tile {
	return &tile{
		Left:     left,
		Right:    right,
		IsDouble: left == right,
	}
}

func (g *Game) AddTeam(t *team) *Game {
	g.teams = append(g.teams, t)
	for _, player := range t.players {
		g.players[player.id] = player
		g.table.playersID = append(g.table.playersID, player.id)
	}
	return g
}

func (g *Game) AddRule(r *Rule) *Game {
	g.rules = append(g.rules, *r)
	return g
}

func (g *Game) Start() {
	g.table.dealHands(g)
}

func (t *table) dealHands(game *Game) {
	playerIndeces := []int{0, 1, 2, 3}
	for amountPick := 0; amountPick < 7; amountPick++ {
		rand.Seed(uint64(time.Now().UnixNano()))
		shuffle(playerIndeces)
		for _, playerIndex := range playerIndeces {
			player := game.players[t.playersID[playerIndex]]
			// Generate a random index
			shuffleIndex := rand.Intn(len(t.shuffledTiles))
			tile := t.shuffledTiles[shuffleIndex]
			if tile.Left == 6 && tile.Right == 6 && game.isBeggining {
				t.currentPlayer = player
			}
			player.hand.tiles = append(player.hand.tiles, *tile)
		}
		t.shuffledTiles = t.shuffledTiles[1:]
	}
	for _, player := range game.players {
		t.playersHand[player.id] = append(t.playersHand[player.id], player.hand)
	}
}

func NewPlayer(name string) *player {
	return &player{
		name:    name,
		id:      fetchPlayerID(),
		hand:    newHand(),
		partner: &player{},
	}
}

func NewTeam(name string) *team {
	return &team{
		name:  name,
		id:    fetchTeamID(),
		score: 0,
		won:   false,
	}
}

func (t *team) AddPlayer(player *player) {
	t.players = append(t.players, player)
}

func newHand() *hand {
	return &hand{
		tiles: make([]tile, 0),
	}
}

func fetchTeamID() teamID {
	var id teamID
	var b [8]byte

	for {
		_, err := rand.Read(b[:])
		if err != nil {
			panic("unable to generate random bytes")
		}

		id = teamID(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)

		if !dbTeamID[uint64(id)] {
			dbTeamID[uint64(id)] = true
			break
		}
		klog.Infof("Team ID %d already exists\n", id)
	}

	return id
}

func fetchPlayerID() playerID {
	var id playerID
	var b [8]byte

	for {
		_, err := rand.Read(b[:])
		if err != nil {
			panic("unable to generate random bytes")
		}

		id = playerID(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)

		if !dbPlayerID[uint64(id)] {
			dbPlayerID[uint64(id)] = true
			break
		}
		klog.Infof("Player ID %d already exists\n", id)
	}

	return id
}
