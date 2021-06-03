package main

import (
	"math/rand"
	"os"
	"strconv"
)

//Source for rules: https://www.ultraboardgames.com/ticket-to-ride/game-rules.php

type Engine struct {
	playerList   []Player // a list of Player objects, used to simulate the game
	activePlayer int      // the current Player whose turn it is

	trainCardHands         [][]int               //the engine keeps track of who has what cards: TrainCard[i][j]=the ith player has how many of the j'th color of card
	destinationTicketHands [][]DestinationTicket //which player has what destinationTickets: used for scoring purposes
	numTrains              []int                 //how many trains the i'th player has left: the game will end when this is 0 for any player

	destinationNames []string //the list of destination names
	stringColors     []string //the names of colors
	trackList        []Track  //the main game board:an array of Tracks, each track is a single edge
	trackStatus      []int    //stores the current status of each track: which player owns it, or -1 for unoccupied
	faceUpTrainCards []int    //the cards currently face up on the table, indexed by color

	pileOfTrainCards         []GameColor         //the facedown stack of train cards
	pileOfDestinationTickets []DestinationTicket //the facedown stack of destination tickets

	gameConstants GameConstants

	adjacencyList [][]int
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

func (e *Engine) populateAdjacencyList() {
	e.adjacencyList = make([][]int, e.gameConstants.NumDestinations)
	for i := 0; i < e.gameConstants.NumDestinations; i++ {
		e.adjacencyList[i] = make([]int, 0)
	}

	for i, edge := range e.trackList {
		e.adjacencyList[edge.d1] = append(e.adjacencyList[edge.d1], i)
		e.adjacencyList[edge.d2] = append(e.adjacencyList[edge.d2], i)
	}

}

func (e *Engine) initializeGame(playerList []Player, constants GameConstants) {

	//TODO: some of these things refer to global variables, ideally we don't want that, everything can be a parameter

	e.playerList = playerList
	e.activePlayer = 0

	e.destinationNames=destinationNames
	e.stringColors=stringColors

	e.gameConstants = constants
	e.gameConstants.NumPlayers = len(e.playerList)

	e.trackList = listOfTracks
	e.gameConstants.NumTracks = len(e.trackList)
	e.trackStatus = make([]int, len(e.trackList))
	for i := range e.trackStatus {
		e.trackStatus[i] = -1
	}

	//populate adjacency List
	e.populateAdjacencyList()

	//set numTrains
	e.numTrains = make([]int, len(e.playerList))
	for i := range e.playerList {
		e.numTrains[i] = e.gameConstants.NumStartingTrains
	}

	//ADDING RETURN FOR DEBUGGING PURPOSES
	return

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
		e.runDestinationTokenCollectionPhase(i, e.gameConstants.NumInitialDestinationTicketsOffered, e.gameConstants.NumInitialDestinationTicketsPicked)
	}

}

func (e *Engine) runCollectionPhase() {
	whichColor := e.playerList[e.activePlayer].askPickup(2)
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

	whichColor = e.playerList[e.activePlayer].askPickup(1)
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
	whichTrack, whichColor := e.playerList[e.activePlayer].askTrackLay()
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

func (e *Engine) runDestinationTokenCollectionPhase(playerNumber, numToOffer, numToAccept int) {

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

	//end condition: since this turn is starting with a player with less than two trains, he must have reached here on the previous move
	if e.numTrains[e.activePlayer] <= 2 {
		return true
	}

	//first, inform the player of the game state
	e.playerList[e.activePlayer].informStatus(e.trackStatus, e.faceUpTrainCards)

	whichMove := e.playerList[e.activePlayer].askMove()
	//first, ask the guy whose turn it is what he wants to d
	if whichMove == 0 {
		// let him pick up cards
		e.runCollectionPhase()
	} else if whichMove == 1 {
		//ask them to put down some tracks
		e.runTrackLayingPhase()
	} else if whichMove == 2 {
		//finally ask them to decide and pick some destination tokens
		e.runDestinationTokenCollectionPhase(e.activePlayer, e.gameConstants.NumDestinationTicketsOffered, e.gameConstants.NumDestinationTicketsPicked)
	} else {
		panic("Invalid move choice")
	}

	//next player
	e.activePlayer++
	e.activePlayer %= e.gameConstants.NumPlayers
	return false
}

func (e *Engine) determinePlayerScore(playerNumber int) int {
	score := 0

	// add all the scores for paths
	for i, status := range e.trackStatus {
		if status == playerNumber {
			score += e.gameConstants.routeLengthScores[e.trackList[i].length]
		}
	}

	//	add or subtract the score for each destination ticket
	for _, ticket := range e.destinationTicketHands[playerNumber] {
		if e.isConnected(ticket.d1, ticket.d2, playerNumber) {
			score += ticket.points
		} else {
			score -= ticket.points
		}
	}

	return score
}

func (e *Engine) dfs(src, dst Destination, playerNumber int, seen []bool) bool {

	if src == dst {
		return true
	}

	seen[src] = true

	var otherDestination Destination

	for _, edgeIndex := range e.adjacencyList[src] {
		if e.trackStatus[edgeIndex] != playerNumber {
			continue
		}

		otherDestination = e.getOtherDestination(src, e.trackList[edgeIndex])

		if seen[otherDestination] {
			continue
		}

		if e.dfs(otherDestination, dst, playerNumber, seen) {
			return true
		}
	}
	return false
}

func (e *Engine) isConnected(d1, d2 Destination, playerNumber int) bool {
	seen := make([]bool, e.gameConstants.NumDestinations)

	return e.dfs(d1, d2, playerNumber, seen)
}

func (e *Engine) getOtherDestination(d Destination, t Track) Destination {
	if d == t.d1 {
		return t.d2
	} else if d == t.d2 {
		return t.d1
	} else {
		panic("This isn't the right edge")
	}
}

func (e *Engine) computeLongestPathRecursive(x Destination, currLen int, adjList [][]int, seenEdge []bool, toUpdate *int) {
	if currLen > (*toUpdate) {
		*toUpdate = currLen
	}

	for _, edgeIndex := range adjList[x] {
		if !seenEdge[edgeIndex] {
			seenEdge[edgeIndex] = true
			e.computeLongestPathRecursive(e.getOtherDestination(x, e.trackList[edgeIndex]), currLen+e.trackList[edgeIndex].length, adjList, seenEdge, toUpdate)
			seenEdge[edgeIndex] = false
		}
	}
}

func (e *Engine) determineLongestPathForPlayer(playerNumber int) int {
	//	There are two ways I thought of doing this and I'm not sure which is better
	//	Way #1: iterate over all subsets of edges, and see if you can get a set of edges which induces a subgraph which is connected and eulerean
	//	The check for eulerean is simple, the check for connectedness is dfs
	//	Way #2: from all starting points, repeatedly try all edges, update ans=max(ans,currPathLen)
	//	 I'm not even sure if way #2 is correct

	//	First, create a playerAdjList, which selects from adjList only those which belong to the player
	playerAdjList := make([][]int, e.gameConstants.NumDestinations)
	for i, adj := range e.adjacencyList {
		for _, edge := range adj {
			if e.trackStatus[edge] == playerNumber {
				playerAdjList[i] = append(playerAdjList[i], edge)
			}
		}
	}

	seenEdge := make([]bool, e.gameConstants.NumTracks)

	ans := 0
	//	now, from each starting point, run the recursive computer
	for i := 0; i < e.gameConstants.NumDestinations; i++ {
		e.computeLongestPathRecursive(Destination(i), 0, playerAdjList, seenEdge, &ans)
	}

	return ans
}

func (e *Engine) getLongestPathPlayers() []int {
	longestPathers := make([]int, 0)
	currLongestPathLength := 0

	for i := range e.playerList {
		sc := e.determineLongestPathForPlayer(i)
		if sc > currLongestPathLength {
			currLongestPathLength = sc
			longestPathers = nil
			longestPathers = append(longestPathers, i)
		} else if sc == currLongestPathLength {
			longestPathers = append(longestPathers, i)
		}
	}
	return longestPathers
}

func (e *Engine) determineWinners() []int {
	winners := make([]int, 0)
	currBestScore := 0

	//figure out which player(s) have longest paths
	longestPathPlayers := e.getLongestPathPlayers()

	for i := range e.playerList {
		sc := e.determinePlayerScore(i)
		if itemExists(longestPathPlayers, i) {
			sc += e.gameConstants.LongestPathScore
		}
		if sc > currBestScore {
			currBestScore = sc
			winners = nil
			winners = append(winners, i)
		} else if sc == currBestScore {
			winners = append(winners, i)
		}
	}
	return winners
}

func (e *Engine) runGame(playerList []Player, constants GameConstants) []int {
	//initialize
	e.initializeGame(playerList, constants)

	gameOver := false
	//run turns until the game is over
	for !gameOver {

		//write the graph to file
		e.writeGraphToFile("graph.txt")

		gameOver = e.runSingleTurn()
	}

	//determine the Winner
	return e.determineWinners()
}

func (e *Engine) writeGraphToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic("failed creating file")
	}
	defer file.Close()
	_, err = file.WriteString(e.getGraphVizString())
	if err != nil {
		panic("failed writing to file")
	}

}

func (e *Engine) getGraphVizString() string {

	//for graphviz testing purposes
	//for i:=0;i<10;i++ {
	//	e.trackStatus[i]=1
	//}
	//for i:=10;i<20;i++ {
	//	e.trackStatus[i]=2
	//}
	graphString := "Graph G {\n"
	graphString += "\toverlap=true\n"
	graphString += "\tmode=KK\n"

	for dest,pos := range mapPositions {
		graphString += "\t"
		graphString += dest
		graphString += " [ pos=\""
		graphString += pos
		graphString += "!\" ];"
		graphString += "\n"
	}

	for i, track := range e.trackList {
		graphString += "\t"
		graphString += e.destinationNames[track.d1]
		graphString += " -- "
		graphString += e.destinationNames[track.d2]
		graphString += " [ len=" + strconv.Itoa(track.length) + ","
		if e.trackStatus[i] == -1 {
			//	the track is empty
			graphString += "style=dotted,color="
			if track.c == Rainbow {
				graphString += "red:green:yellow:blue:orange:purple"
			} else {
				graphString += e.stringColors[track.c]
			}
		} else {
			//there is a player on the track
			graphString += "style=bold,color="
			if e.trackStatus[i] == int(Rainbow) {
				graphString += "red:green:yellow:blue:orange:purple"
			} else {
				graphString += e.stringColors[e.trackStatus[i]]
			}
		}

		_,ok1 := mapPositions[e.destinationNames[track.d1]]
		_,ok2 := mapPositions[e.destinationNames[track.d2]]

		if ok1||ok2 {
			graphString += ", weight=1"
		} else {
			graphString += ", weight=1"
		}
		graphString += "];\n"
	}
	graphString += "}"
	return graphString
}