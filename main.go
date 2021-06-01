package main

const NUMCOLORCARDS = 12
const NUMRAINBOWCARDS = 14
const NUMSTARTINGTRAINS = 48
const NUMFACEUPTRAINCARDS = 5
const NUMGAMECOLORS = 9
const NUMINITIALTRAINCARDSDEALT = 4
const NUMINITIALDESTINATIONTICKETSDEALT = 3

type Destination int
type GameColor int

// ['Atlanta', 'Boston', 'Calgary', 'Charleston', 'Chicago', 'Dallas', 'Denver', 'Duluth', 'El Paso', 'Helena', 'Houston', 'Kansas City', 'Las Vegas', 'Little Rock', 'Los Angeles', 'Miami', 'Montreal', 'Nashville', 'New Orleans', 'New York', 'Oklahoma City', 'Omaha', 'Phoenix', 'Pittsburgh', 'Portland', 'Raleigh', 'Saint Louis', 'Salt Lake City', 'San Francisco', 'Santa Fe', 'Sault St. Marie', 'Seattle', 'Toronto', 'Vancouver', 'Washington', 'Winnipeg']
//TODO: build enum/array for destination

const (
	Atlanta Destination = iota
)

type DestinationTicket struct {
	d1, d2 Destination //endpoints
	points int         //score
}

type GameConstants struct {
	NumColorCards, NumRainbowCards, NumStartingTrains, NumFaceUpTrainCards, NumGameColors, NumInitialTrainCardsDealt, NumInitialDestinationTicketsDealt int
}

const (
	Red GameColor = iota
	Orange
	Yellow
	Green
	Blue
	Purple
	Black
	White
	Rainbow
	Other
)

var listOfGameColors = [...]GameColor{Red, Orange, Yellow, Green, Blue, Purple, Black, White, Rainbow}

//TODO: build DestinationTicket array
var listOfDestinationTickets = []DestinationTicket{{Atlanta, Atlanta, 1}}

type Track struct {
	idx    int         //position in array
	d1, d2 Destination // two endpoints
	c      GameColor   //what color
	status int         //status represents whether it is unoccupied (-1) or occupied (which player number has it)
}

type Player interface {
	initialize(myNumber, totPlayers int, constants GameConstants) //Tell the player what his number is and the total number of players, as well as the game settings

	//players are stateful, so we may need to inform them of game events (in case they want to keep track of other players' hands or something)
	informCardPickup(int, GameColor)   //inform this player that a player picked up a card of given color
	informTrackLay(int, Track)         //inform this player that a player placed a track
	informDestinationTicketPickup(int) //inform this player that a player picked up a destination card

	askPickup([]Track, [NUMGAMECOLORS]int, int) GameColor //ask this player, given the gamestate, which card he wants to pick up
	giveTrainCard(GameColor)                              //tell this player he has another card of given color
	giveDestinationTicket(DestinationTicket)              //tell this player has a destination card

	askTrackLay([]Track, [NUMGAMECOLORS]int) (Track, bool)       //ask this player which track he wants to lay, if he wants to lay one
	askDestinationTicketPickup([]Track, [NUMGAMECOLORS]int) bool //ask this player if he wants to pick up a destination card
}

type Engine struct {
	playerList   []Player // a list of Player objects, used to simulate the game
	numPlayers   int      // the number of Players
	activePlayer int      // the current Player whose turn it is

	trainCardHands         [][NUMGAMECOLORS]int  //the engine keeps track of who has what cards: TrainCard[i][j]=the ith player has how many of the j'th color of card
	destinationTicketHands [][]DestinationTicket //which player has what destinationTickets: used for scoring purposes
	numTrains              []int                 //how many trains the i'th player has left: the game will end when this is 0 for any player

	trackList        []Track            //the main game board:an array of Tracks, each track is a single edge
	faceUpTrainCards [NUMGAMECOLORS]int //the cards currently face up on the table, indexed by color
	pileOfTrainCards []GameColor        //the facedown stack of train cards

	pileOfDestinationTickets []DestinationTicket //the facedown stack of destination tickets
}
