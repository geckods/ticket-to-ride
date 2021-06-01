package main

import "math/rand"

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
	length int         // What is the length of the road
}

//TODO: build Track array
var listOfTracks = []Track{{0, Atlanta, Atlanta, Red, -1, 3}}

type Player interface {
	initialize(myNumber, totPlayers int, constants GameConstants) //Tell the player what his number is and the total number of players, as well as the game settings

	//players are stateful, so we may need to inform them of game events (in case they want to keep track of other players' hands or something)
	informCardPickup(int, GameColor)   //inform this player that a player picked up a card of given color
	informTrackLay(int, Track)         //inform this player that a player placed a track
	informDestinationTicketPickup(int) //inform this player that a player picked up a destination card

	askPickup([]Track, [NUMGAMECOLORS]int, int) GameColor //ask this player, given the gamestate, which card he wants to pick up
	giveTrainCard(GameColor)                              //tell this player he has another card of given color
	giveDestinationTicket(DestinationTicket)              //tell this player has a destination card

	askTrackLay([]Track, [NUMGAMECOLORS]int) (int, GameColor, bool) //ask this player which track he wants to lay, and with what color, if he wants to lay one
	askDestinationTicketPickup([]Track, [NUMGAMECOLORS]int) bool    //ask this player if he wants to pick up a destination card
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

	pileOfTrainCards         []GameColor         //the facedown stack of train cards
	pileOfDestinationTickets []DestinationTicket //the facedown stack of destination tickets
}

func (e *Engine) initializePileOfTrainCards(toExclude [NUMGAMECOLORS]int) {
	e.pileOfTrainCards = nil //clear the pile

	//put all the cards into the pile
	for _, c := range listOfGameColors {
		if c == Rainbow {
			for i := 0; i < NUMRAINBOWCARDS-toExclude[c]; i++ {
				e.pileOfTrainCards = append(e.pileOfTrainCards, c)
			}
		} else {
			for i := 0; i < NUMCOLORCARDS-toExclude[c]; i++ {
				e.pileOfTrainCards = append(e.pileOfTrainCards, c)
			}
		}
	}

	//shuffle the deck
	rand.Shuffle(len(e.pileOfTrainCards), func(i, j int) {
		e.pileOfTrainCards[i], e.pileOfTrainCards[j] = e.pileOfTrainCards[j], e.pileOfTrainCards[i]
	})
}

func (e *Engine) drawTopTrainCard() GameColor {
	if len(e.pileOfTrainCards) == 0 {
		//let's figure out what we need to exclude
		var toExclude [NUMGAMECOLORS]int
		for j := 0; j < NUMGAMECOLORS; j++ {
			//first, exclude cards that are face up on the table
			toExclude[j] += e.faceUpTrainCards[j]
			for i, _ := range e.playerList {
				//exclude cards that are in players' hands
				toExclude[j] += e.trainCardHands[i][j]
			}
		}

		e.initializePileOfTrainCards(toExclude)
		if len(e.pileOfTrainCards) == 0 {
			//	if it's still zero, then all cards are in players' hands, we cannot draw any more
			panic("Not Enough Cards in The Deck")
		}
	}

	index := len(e.pileOfTrainCards) - 1              // Get the index of the top most element.
	element := (e.pileOfTrainCards)[index]            // Index into the slice and obtain the element.
	e.pileOfTrainCards = (e.pileOfTrainCards)[:index] // Remove it from the stack by slicing it off.

	return element
}

func (e *Engine) giveCardToPlayer(p int, c GameColor, toHideColorWhenInforming bool) {
	//update the engine's copy
	e.trainCardHands[p][c]++
	//give the player his card
	e.playerList[p].giveTrainCard(c)

	//tell everybody else what happened
	for _, pl := range e.playerList {
		if toHideColorWhenInforming {
			pl.informCardPickup(p, Other)
		} else {
			pl.informCardPickup(p, c)
		}
	}
}

func (e *Engine) giveDestinationTicketToPlayer(p int, ticket DestinationTicket) {
	//update the engine's copy
	e.destinationTicketHands[p] = append(e.destinationTicketHands[p], ticket)
	//give the player his card
	e.playerList[p].giveDestinationTicket(ticket)
	//tell everybody else what happened
	for _, pl := range e.playerList {
		pl.informDestinationTicketPickup(p)
	}
}

func (e *Engine) initializeDestinationTicketPile() {
	//assign
	e.pileOfDestinationTickets = listOfDestinationTickets

	//	shuffle
	rand.Shuffle(len(e.pileOfDestinationTickets), func(i, j int) {
		e.pileOfDestinationTickets[i], e.pileOfDestinationTickets[j] = e.pileOfDestinationTickets[j], e.pileOfDestinationTickets[i]
	})

}

func (e *Engine) drawTopDestinationTicket() (DestinationTicket, bool) {
	if len(e.pileOfDestinationTickets) == 0 {
		return DestinationTicket{}, false
	}

	index := len(e.pileOfDestinationTickets) - 1                      // Get the index of the top most element.
	element := (e.pileOfDestinationTickets)[index]                    // Index into the slice and obtain the element.
	e.pileOfDestinationTickets = (e.pileOfDestinationTickets)[:index] // Remove it from the stack by slicing it off.

	return element, true
}

func (e *Engine) initializeGame(playerList []Player) {
	e.playerList = playerList
	e.activePlayer = 0
	e.numPlayers = len(playerList)

	gameconstants := GameConstants{
		NumColorCards:                     NUMCOLORCARDS,
		NumRainbowCards:                   NUMRAINBOWCARDS,
		NumStartingTrains:                 NUMSTARTINGTRAINS,
		NumFaceUpTrainCards:               NUMFACEUPTRAINCARDS,
		NumGameColors:                     NUMGAMECOLORS,
		NumInitialTrainCardsDealt:         NUMINITIALTRAINCARDSDEALT,
		NumInitialDestinationTicketsDealt: NUMINITIALDESTINATIONTICKETSDEALT,
	}

	for i, p := range e.playerList {
		p.initialize(i, e.numPlayers, gameconstants)
		//	initialize each player
	}

	//	set up the pile of traincards (face up cards are initially all 0)
	e.initializePileOfTrainCards(e.faceUpTrainCards)

	for i, _ := range e.playerList {
		//give each player the initial train cards, don't announce card color
		for j := 0; j < NUMINITIALTRAINCARDSDEALT; j++ {
			e.giveCardToPlayer(i, e.drawTopTrainCard(), true)
		}
	}

	//	set up the pile of destination tickets
	e.initializeDestinationTicketPile()

	//	give each player destination tickets
	for i, _ := range e.playerList {
		//give each player the initial destination tickets
		for j := 0; j < NUMINITIALDESTINATIONTICKETSDEALT; j++ {
			ticket, ok := e.drawTopDestinationTicket()
			if !ok {
				panic("Not Enough Destination Tickets To Deal This Many To Each Player")
			}
			e.giveDestinationTicketToPlayer(i, ticket)
		}
	}

	for i, _ := range e.playerList {
		e.numTrains[i] = NUMSTARTINGTRAINS
	}

	e.trackList = listOfTracks

}

func (e *Engine) runCollectionPhase() {
	whichColor := e.playerList[e.activePlayer].askPickup(e.trackList, e.faceUpTrainCards, 2)
	if whichColor != Other {
		//	he wants a faceup card
		if e.faceUpTrainCards[whichColor] <= 0 {
			//	he cannot pick that card
			panic("The player picked a missing color")
		}
		e.giveCardToPlayer(e.activePlayer, whichColor, false)
		e.faceUpTrainCards[whichColor]--
		e.faceUpTrainCards[e.drawTopTrainCard()]++

		if whichColor == Rainbow {
			//	picking a rainbow color costs 2, so you're done
			return
		}
	} else {
		//	asking for a random card from the deck
		e.giveCardToPlayer(e.activePlayer, e.drawTopTrainCard(), true)
	}

	whichColor = e.playerList[e.activePlayer].askPickup(e.trackList, e.faceUpTrainCards, 1)
	if whichColor == Rainbow {
		panic("The player picked a rainbow on his second turn")
	}

	if whichColor != Other {
		//	he wants a faceup card
		if e.faceUpTrainCards[whichColor] <= 0 {
			//	he cannot pick that card
			panic("The player picked a missing color")
		}
		e.giveCardToPlayer(e.activePlayer, whichColor, false)
		e.faceUpTrainCards[whichColor]--
		e.faceUpTrainCards[e.drawTopTrainCard()]++
	} else {
		//	asking for a random card from the deck
		e.giveCardToPlayer(e.activePlayer, e.drawTopTrainCard(), true)
	}
}

func (e *Engine) runTrackLayingPhase() bool {
	whichTrack, whichColor, ok := e.playerList[e.activePlayer].askTrackLay(e.trackList, e.faceUpTrainCards)
	if !ok {
		return false
	}
	if e.trackList[whichTrack].status != -1 {
		panic("The player tried to place over an occupied track")
	}
	if e.trackList[whichTrack].c != whichColor && e.trackList[whichTrack].c != Other {
		panic("The player is trying to play with the wrong color for the track")
	}
	if e.trackList[whichTrack].length > e.numTrains[e.activePlayer] {
		panic("The player does not have enough trains to play this move")
	}
	if e.trackList[whichTrack].length > e.trainCardHands[e.activePlayer][whichColor] {
		panic("The player does not have color cards to play this move")
	}

	//	If we made it this far, I think we're good: do the move
	//	mark the track as occupied
	e.trackList[whichTrack].status = e.activePlayer
	// remove the cards
	e.trainCardHands[e.activePlayer][whichColor] -= e.trackList[whichTrack].length
	//remove the trains
	e.numTrains[e.activePlayer] -= e.trackList[whichTrack].length

	return e.numTrains[e.activePlayer] == 0
}

func (e *Engine) runSingleTurn() bool {

	//first, ask the guy whose turn it is to pick up cards
	e.runCollectionPhase()
	if e.runTrackLayingPhase() {
		//	The game is done
		return true
	}

	e.activePlayer++
	e.activePlayer %= e.numPlayers
	return false
}

//TODO: rewrite dev card collection to present a list
//TODO: separate out status (dynamic part of tracks) from trackList (static part of tracks)
//TODO: use slices instead of arrays, so that nothing in engine is static/dependant on constants

func main() {

}
