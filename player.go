package main

type Player interface {
	initialize(myNumber int, trackList []Track, constants GameConstants) //Tell the player what his number is and the total number of players, as well as the game settings

	//players are stateful, so we may need to inform them of game events (in case they want to keep track of other players' hands or something)
	informCardPickup(int, GameColor)   //inform this player that a player picked up a card of given color
	informTrackLay(int, int)         //inform this player that a player placed a track
	informDestinationTicketPickup(int) //inform this player that a player picked up a destination card

	informStatus([]int, []int) //called to inform the playstate before their turn

	askMove() int //Ask the player what move he wants to do: 0 is pick up cards, 1 is place Tracks, 2 is pick destination ticket
	askPickup(int) GameColor   //ask this player, given the gamestate, which card he wants to pick up
	askTrackLay() (int, GameColor) //ask this player which track he wants to lay, and with what color

	giveTrainCard(GameColor)                 //tell this player he has another card of given color
	giveDestinationTicket(DestinationTicket) //tell this player has a destination card
	offerDestinationTickets([]DestinationTicket, int) []int //offer a list of destination cards and tell the player to take some of them
}
