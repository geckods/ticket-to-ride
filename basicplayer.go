package main

type BasicPlayer struct {
	trackList []Track //my copy of the board
	trackStatus []int //my copy of the status of each track

	myTrainCards []int //number of each color card I have
	myDestinationTickets []DestinationTicket //list of destination ticket I have
	myNumber int //my player ID
	myTrains int
	constants GameConstants
}

func (b* BasicPlayer) initialize(myNumber int, trackList []Track, constants GameConstants) {
	b.myNumber = myNumber
	b.trackList = trackList
	b.constants = constants
}

func (b* BasicPlayer) informCardPickup(int, GameColor) {
	//	do nothing
}

func (b* BasicPlayer) informTrackLay(int, Track) {
	//	do nothing
}

func (b* BasicPlayer) informDestinationTicketPickup(int) {
	//	do nothing
}

func (b* BasicPlayer) whichTrackCanILay() {
	for i,track := b.trackList {
		if
	}
}

askMove([]int, []int) int //Ask the player what move he wants to do: 0 is pick up cards, 1 is place Tracks, 2 is pick destination ticket

askPickup([]int, []int, int) GameColor   //ask this player, given the gamestate, which card he wants to pick up
giveTrainCard(GameColor)                 //tell this player he has another card of given color
giveDestinationTicket(DestinationTicket) //tell this player has a destination card

askTrackLay([]int, []int) (int, GameColor) //ask this player which track he wants to lay, and with what color

offerDestinationTickets([]DestinationTicket, int) []int //offer a list of destination cards and tell the player to take some of them