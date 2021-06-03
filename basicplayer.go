package main

import "fmt"

type BasicPlayer struct {
	trackList []Track //my copy of the board
	trackStatus []int //my copy of the status of each track
	faceUpCards []int

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
	b.myTrains = constants.NumStartingTrains

	b.myTrainCards = make([]int, b.constants.NumGameColors)
	b.myDestinationTickets=make([]DestinationTicket,0)
}

func (b* BasicPlayer) informStatus(trackStatus []int, faceUpCards []int) {
 	b.faceUpCards=faceUpCards
 	b.trackStatus=trackStatus

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

func (b* BasicPlayer) whichTrackCanILay() (int, GameColor) {
	for i,track := range b.trackList { //rainbow
		if b.trackStatus[i]!=-1 || b.myTrains<track.length {
			continue
		}
		fmt.Println("Found eligible track")
		if track.c==Other {
			for _, allcolor:=range listOfGameColors{
				if allcolor!=Rainbow {
					if track.length <= b.myTrainCards[allcolor]+b.myTrainCards[Rainbow]{
						return i, allcolor
					}
				}
			}
		} else {
			if b.myTrainCards[track.c] + b.myTrainCards[Rainbow] >= track.length {
				return i, track.c
			}
		}
	}
	fmt.Println("Did not find track to be laid")
	return -1, Other
}
func (b* BasicPlayer) askTrackLay() (int, GameColor){
	trackIndex, trackColor:=b.whichTrackCanILay()
	if trackIndex==-1{
		panic("whichTrackCanILay in error, panic, panic, panic")
	}
	if b.trackList[trackIndex].length>b.myTrainCards[trackColor]{
		b.myTrainCards[Rainbow]-= b.trackList[trackIndex].length-b.myTrainCards[trackColor]
		b.myTrainCards[trackColor]=0
	} else{
		b.myTrainCards[trackColor]-=b.trackList[trackIndex].length
	}
	b.myTrains-=b.trackList[trackIndex].length
	return trackIndex, trackColor


} //ask this player which track he wants to lay, and with what color

func (b* BasicPlayer) askMove() int{
	whichTrack,_ := b.whichTrackCanILay()
	if whichTrack!=-1 {
		return 1
	} else {
		return 0
	}

} //Ask the player what move he wants to do: 0 is pick up cards, 1 is place Tracks, 2 is pick destination ticket

func (b* BasicPlayer) askPickup(howManyLeft int) GameColor {
	return Other
}   //ask this player, given the gamestate, which card he wants to pick up


func (b* BasicPlayer) giveTrainCard(card GameColor) {
	b.myTrainCards[card]++
}                 //tell this player he has another card of given color

func (b* BasicPlayer) giveDestinationTicket(d DestinationTicket) {
	//	basic player doesn't care about destination tickets, so do nothing
} //tell this player has a destination card


func (b* BasicPlayer) offerDestinationTickets(dtlist []DestinationTicket,howmany int) []int {
	//basic player doesn't care about destination tickets, so just pick the first bunch of tickets
	ret := make([]int,0)
	for i:=0;i<howmany;i++ {
		ret=append(ret,i)
	}
	return ret
}//offer a list of destination cards and tell the player to take some of them