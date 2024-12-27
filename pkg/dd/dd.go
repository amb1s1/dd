package dd

var (
	dbPlayerID = map[uint64]bool{}
	dbTeamID   = map[uint64]bool{}
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
	players       []*player
	totalScore    map[teamID]int
	scoresPerHand []map[teamID]int
	table         *table
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
		players:       make([]*player, 0),
		totalScore:    make(map[teamID]int),
		scoresPerHand: make([]map[teamID]int, 0),
		table:         newTable(),
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
	return &table{
		boards:        make([]*board, 0),
		shuffledTiles: make([]*tile, 0),
		playersHand:   make(map[playerID][]*hand),
		currentPlayer: &player{},
	}
}

func (g *Game) AddTeam(t *team) *Game {
	g.teams = append(g.teams, t)
	g.players = append(g.players, t.players...)
	return g
}

func (g *Game) AddRule(r *Rule) *Game {
	g.rules = append(g.rules, *r)
	return g
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

	id = teamID(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
	for dbTeamID[uint64(id)] {
		id = teamID(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
	}
	dbTeamID[uint64(id)] = true
	return id
}

func fetchPlayerID() playerID {
	var id playerID
	var b [8]byte

	id = playerID(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
	for dbPlayerID[uint64(id)] {
		id = playerID(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
	}
	dbPlayerID[uint64(id)] = true
	return id
}
