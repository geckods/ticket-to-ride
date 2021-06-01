package main

//Source for rules: https://ticket-to-ride.fandom.com/wiki/Ticket_to_Ride
//Source for map: https://images-eu.ssl-images-amazon.com/images/I/B19d%2BVcYwWS.png

import (
	"math/rand"
	"reflect"
)

const NUMCOLORCARDS = 12
const NUMRAINBOWCARDS = 14
const NUMSTARTINGTRAINS = 48
const NUMFACEUPTRAINCARDS = 5
const NUMGAMECOLORS = 9
const NUMINITIALTRAINCARDSDEALT = 4
const NUMINITIALDESTINATIONTICKETSOFFERED = 3
const NUMINITIALDESTINATIONTICKETSPICKED = 2
const NUMDESTINATIONTICKETSOFFERED = 3
const NUMDESTINATIONTICKETSPICKED = 1

func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

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
	NumColorCards, NumRainbowCards, NumStartingTrains, NumFaceUpTrainCards, NumGameColors, NumInitialTrainCardsDealt, NumInitialDestinationTicketsOffered, NumInitialDestinationTicketsPicked, NumDestinationTicketsOffered, NumDestinationTicketsPicked, NumPlayers int
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
	length int         // What is the length of the road
}

//TODO: build Track array
var listOfTracks = []Track{{0, Atlanta, Atlanta, Red, 3}}

type Player interface {
	initialize(myNumber int, trackList []Track, constants GameConstants) //Tell the player what his number is and the total number of players, as well as the game settings

	//players are stateful, so we may need to inform them of game events (in case they want to keep track of other players' hands or something)
	informCardPickup(int, GameColor)   //inform this player that a player picked up a card of given color
	informTrackLay(int, Track)         //inform this player that a player placed a track
	informDestinationTicketPickup(int) //inform this player that a player picked up a destination card

	askPickup([]int, []int, int) GameColor   //ask this player, given the gamestate, which card he wants to pick up
	giveTrainCard(GameColor)                 //tell this player he has another card of given color
	giveDestinationTicket(DestinationTicket) //tell this player has a destination card

	askTrackLay([]int, []int) (int, GameColor, bool) //ask this player which track he wants to lay, and with what color, if he wants to lay one
	askDestinationTicketPickup([]int, []int) bool    //ask this player if he wants to pick up a destination card

	offerDestinationTickets([]DestinationTicket, int) []int //offer a list of destination cards and tell the player to take some of them
}

type Engine struct {
	playerList   []Player // a list of Player objects, used to simulate the game
	activePlayer int      // the current Player whose turn it is

	trainCardHands         [][]int               //the engine keeps track of who has what cards: TrainCard[i][j]=the ith player has how many of the j'th color of card
	destinationTicketHands [][]DestinationTicket //which player has what destinationTickets: used for scoring purposes
	numTrains              []int                 //how many trains the i'th player has left: the game will end when this is 0 for any player

	trackList        []Track //the main game board:an array of Tracks, each track is a single edge
	trackStatus      []int   //stores the current status of each track: which player owns it, or -1 for unoccupied
	faceUpTrainCards []int   //the cards currently face up on the table, indexed by color

	pileOfTrainCards         []GameColor         //the facedown stack of train cards
	pileOfDestinationTickets []DestinationTicket //the facedown stack of destination tickets

	gameConstants GameConstants
}

func (e *Engine) initializePileOfTrainCards(toExclude []int) {
	e.pileOfTrainCards = nil //clear the pile

	//put all the cards into the pile
	for _, c := range listOfGameColors {
		if c == Rainbow {
			for i := 0; i < e.gameConstants.NumRainbowCards-toExclude[c]; i++ {
				e.pileOfTrainCards = append(e.pileOfTrainCards, c)
			}
		} else {
			for i := 0; i < e.gameConstants.NumColorCards-toExclude[c]; i++ {
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

		toExclude := make([]int, e.gameConstants.NumGameColors)
		for j := 0; j < e.gameConstants.NumGameColors; j++ {
			//first, exclude cards that are face up on the table
			toExclude[j] += e.faceUpTrainCards[j]
			for i := range e.playerList {
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

	e.gameConstants = GameConstants{
		NumColorCards:                       NUMCOLORCARDS,
		NumRainbowCards:                     NUMRAINBOWCARDS,
		NumStartingTrains:                   NUMSTARTINGTRAINS,
		NumFaceUpTrainCards:                 NUMFACEUPTRAINCARDS,
		NumGameColors:                       NUMGAMECOLORS,
		NumInitialTrainCardsDealt:           NUMINITIALTRAINCARDSDEALT,
		NumInitialDestinationTicketsOffered: NUMINITIALDESTINATIONTICKETSOFFERED,
		NumInitialDestinationTicketsPicked:  NUMINITIALDESTINATIONTICKETSPICKED,
		NumDestinationTicketsOffered:        NUMDESTINATIONTICKETSOFFERED,
		NumDestinationTicketsPicked:         NUMDESTINATIONTICKETSPICKED,
		NumPlayers:                          len(playerList),
	}

	e.trackList = listOfTracks
	e.trackStatus = make([]int, len(e.trackList))
	for i := range e.trackStatus {
		e.trackStatus[i] = -1
	}

	for i, p := range e.playerList {
		p.initialize(i, e.trackList, e.gameConstants)
		//	initialize each player
	}

	//	set up the pile of traincards (face up cards are initially all 0)
	e.initializePileOfTrainCards(e.faceUpTrainCards)

	for i := range e.playerList {
		//give each player the initial train cards, don't announce card color
		for j := 0; j < e.gameConstants.NumInitialTrainCardsDealt; j++ {
			e.giveCardToPlayer(i, e.drawTopTrainCard(), true)
		}
	}

	//	set up the pile of destination tickets
	e.initializeDestinationTicketPile()

	//	give each player destination tickets
	for i := range e.playerList {
		//give each player the initial destination tickets
		e.runDestinationTokenCollectionPhase(i, e.gameConstants.NumInitialDestinationTicketsOffered, e.gameConstants.NumInitialDestinationTicketsPicked, false)
	}

	for i := range e.playerList {
		e.numTrains[i] = e.gameConstants.NumStartingTrains
	}
}

func (e *Engine) runCollectionPhase() {
	whichColor := e.playerList[e.activePlayer].askPickup(e.trackStatus, e.faceUpTrainCards, 2)
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

	whichColor = e.playerList[e.activePlayer].askPickup(e.trackStatus, e.faceUpTrainCards, 1)
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
	whichTrack, whichColor, ok := e.playerList[e.activePlayer].askTrackLay(e.trackStatus, e.faceUpTrainCards)
	if !ok {
		return false
	}
	if e.trackStatus[whichTrack] != -1 {
		panic("The player tried to place over an occupied track")
	}
	if whichColor == Rainbow {
		panic("The player is trying to play rainbow: if you want to use only rainbows, select any other color by default, like red")
	}
	if e.trackList[whichTrack].c != whichColor && e.trackList[whichTrack].c != Other {
		panic("The player is trying to play with the wrong color for the track")
	}
	if e.trackList[whichTrack].length > e.numTrains[e.activePlayer] {
		panic("The player does not have enough trains to play this move")
	}
	if e.trackList[whichTrack].length > e.trainCardHands[e.activePlayer][whichColor]+e.trainCardHands[e.activePlayer][Rainbow] {
		panic("The player does not have color cards to play this move")
	}

	//	If we made it this far, I think we're good: do the move

	//	mark the track as occupied
	e.trackStatus[whichTrack] = e.activePlayer

	// remove the cards
	if e.trainCardHands[e.activePlayer][whichColor] >= e.trackList[whichTrack].length {
		e.trainCardHands[e.activePlayer][whichColor] -= e.trackList[whichTrack].length
	} else {
		//	gotta use up them rainbows
		e.trainCardHands[e.activePlayer][Rainbow] -= e.trackList[whichTrack].length - e.trainCardHands[e.activePlayer][whichColor]
		e.trainCardHands[e.activePlayer][whichColor] = 0
	}

	//remove the trains
	e.numTrains[e.activePlayer] -= e.trackList[whichTrack].length

	return e.numTrains[e.activePlayer] == 0
}

func (e *Engine) runDestinationTokenCollectionPhase(playerNumber, numToOffer, numToAccept int, toAsk bool) {

	if toAsk {
		if e.playerList[playerNumber].askDestinationTicketPickup(e.trackStatus, e.faceUpTrainCards) == false {
			return
		}
	}

	//create a slice to offer
	offerSlice := make([]DestinationTicket, 0)

	for j := 0; j < numToOffer; j++ {
		ticket, ok := e.drawTopDestinationTicket()
		if !ok {
			panic("Not Enough Destination Tickets To Deal This Many To Each Player")
		}
		offerSlice = append(offerSlice, ticket)
	}

	//	offer the slice
	acceptedList := e.playerList[playerNumber].offerDestinationTickets(offerSlice, numToAccept)
	if len(acceptedList) < numToAccept {
		panic("The player didn't pick enough destination tickets")
	}

	alreadySeenAccepted := make([]int, 0) //used to prevent skirting the rules by printing duplicates

	for _, accepted := range acceptedList {
		if accepted < 0 || accepted >= len(offerSlice) || itemExists(alreadySeenAccepted, accepted) {
			panic("The player gave invalid indices in selecting destination tickets")
		}
		alreadySeenAccepted = append(alreadySeenAccepted, accepted)
	}

	for i, offered := range offerSlice {
		if itemExists(acceptedList, i) {
			//this is one of the destination cards he wants to pick
			e.giveDestinationTicketToPlayer(playerNumber, offered)
		} else {
			//this is one of the ones he wants to not pick
			e.putDestinationTicketBackInPile(offered)
		}
	}
}

func (e *Engine) putDestinationTicketBackInPile(ticket DestinationTicket) {
	e.pileOfDestinationTickets = append([]DestinationTicket{ticket}, e.pileOfDestinationTickets...)
}

func (e *Engine) runSingleTurn() bool {

	//first, ask the guy whose turn it is to pick up cards
	e.runCollectionPhase()
	//then, ask them to put down some roads
	if e.runTrackLayingPhase() {
		//	The game is done
		return true
	}
	//finally ask them to decide and pick some destination tokens
	e.runDestinationTokenCollectionPhase(e.activePlayer, e.gameConstants.NumDestinationTicketsOffered, e.gameConstants.NumDestinationTicketsPicked, true)
	//next player
	e.activePlayer++
	e.activePlayer %= e.gameConstants.NumPlayers
	return false
}

func main() {
	myEngine := Engine{}
	_ = myEngine
}
