package main

type Destination int
type GameColor int

type DestinationTicket struct {
	d1, d2 Destination //endpoints
	points int         //score
}

type GameConstants struct {
	NumDestinations, NumTracks, NumColorCards, NumRainbowCards, NumStartingTrains, NumFaceUpTrainCards, NumGameColors, NumInitialTrainCardsDealt, NumInitialDestinationTicketsOffered, NumInitialDestinationTicketsPicked, NumDestinationTicketsOffered, NumDestinationTicketsPicked, NumPlayers, LongestPathScore int
	routeLengthScores                                                                                                                                                                                                                                                                                              []int
}

type Track struct {
	idx    int         //position in array
	d1, d2 Destination // two endpoints
	c      GameColor   //what color
	length int         // What is the length of the road
}
