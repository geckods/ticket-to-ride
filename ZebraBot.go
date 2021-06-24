package main

import "fmt"

//import "flag"

type ZebraBot struct {
	trackList []Track //my copy of the board
	trackStatus []int //my copy of the status of each track
	faceUpCards []int

	myTrainCards []int //number of each color card I have
	myDestinationTickets []DestinationTicket //list of destination ticket I have
	DestinationTicketStatus []int
	myNumber int //my player ID
	myTrains int
	constants GameConstants

	cardsAnimalsWant [][]int
	adjacencyList [][]int

}

func (b* ZebraBot) initialize(myNumber int, trackList []Track, adjList [][]int, constants GameConstants) {
	b.myNumber = myNumber
	b.trackList = trackList
	b.constants = constants
	b.myTrains = constants.NumStartingTrains

	b.myTrainCards = make([]int, b.constants.NumGameColors)
	b.myDestinationTickets=make([]DestinationTicket,0)
	b.DestinationTicketStatus=make([]int,0) //initialise
	b.cardsAnimalsWant= make([][]int, constants.NumPlayers)
	for i:=0;i<constants.NumPlayers;i++{
		b.cardsAnimalsWant[i]=make([]int, constants.NumGameColors)
	}
	b.adjacencyList = adjList
	//b.populateAdjacencyList()
}

func (b* ZebraBot) informStatus(trackStatus []int, faceUpCards []int) {
	b.faceUpCards=faceUpCards
	b.trackStatus=trackStatus

}

func (b* ZebraBot) populateAdjacencyList() {
	b.adjacencyList = make([][]int, b.constants.NumDestinations)
	for i := 0; i < b.constants.NumDestinations; i++ {
		b.adjacencyList[i] = make([]int, 0)
	}

	for i, edge := range b.trackList {
		b.adjacencyList[edge.d1] = append(b.adjacencyList[edge.d1], i)
		b.adjacencyList[edge.d2] = append(b.adjacencyList[edge.d2], i)
	}

}

func (b* ZebraBot) informCardPickup(int, GameColor) {
	//	do nothing
}

func (b* ZebraBot) informTrackLay(int, int) {
	//	do nothing
}

func (b* ZebraBot) informDestinationTicketPickup(int) {
	//	do nothing
}


func (b* ZebraBot) whichTrackCanILay() (int, GameColor) {
	//fmt.Println("Inside whichTrackCAniLay")
	//fmt.Println("I have destinations:", b.myDestinationTickets[0],b.myDestinationTickets[1])
	tracksZebraWants:=make([]int,0)
	for x:=0;x<b.constants.NumGameColors;x++{
		b.cardsAnimalsWant[b.myNumber][x]=0
	}

	for i, dt:= range b.myDestinationTickets{
	if b.DestinationTicketStatus[i]==1 {
		continue
	}else{
		tracksZebraWants=append(tracksZebraWants,b.bfs(dt.d1,dt.d2)...)
	}

	for _,wantedTrackID:=range tracksZebraWants{
		if b.trackStatus[wantedTrackID]!=-1{
			continue}
		wantedTrack:=b.trackList[wantedTrackID]
		//fmt.Println("Zebra Bot wants", wantedTrack.d1, wantedTrack.d2, wantedTrack.c)
		if wantedTrack.length>b.myTrains{
			continue
		}
		if wantedTrack.c==Other {
			for _, allcolor:=range listOfGameColors{
				if allcolor!=Rainbow {
					if wantedTrack.length <= b.myTrainCards[allcolor]+b.myTrainCards[Rainbow]{
						return wantedTrackID, allcolor
					}

				}
			}
		} else {
			if b.myTrainCards[wantedTrack.c] + b.myTrainCards[Rainbow] >= wantedTrack.length {
				return wantedTrackID, wantedTrack.c
			}else{
				b.cardsAnimalsWant[b.myNumber][wantedTrack.c]+=wantedTrack.length
			}
		}

	}
		for i,track := range b.trackList { //rainbow
			if b.trackStatus[i]!=-1 || b.myTrains<track.length {
				continue
			}
			if track.c==Other && track.length>4 {
				for _, allcolor:=range listOfGameColors{
					if allcolor!=Rainbow {
						if track.length <= b.myTrainCards[allcolor]+b.myTrainCards[Rainbow]{
							return i, allcolor
						}
					}
				}
			} else if track.length>4{
				if b.myTrainCards[track.c] + b.myTrainCards[Rainbow] >= track.length {
					return i, track.c
				}else{
					b.cardsAnimalsWant[b.myNumber][track.c]+=track.length
				}
			}
		}

	}
	return -1,Other
}

//bfs should return the empty tracks that are part of the route from dt1 to dt2 with least number of nodes
func (b* ZebraBot) bfs(dt1 Destination, dt2 Destination)[]int{

	//fmt.Println("inside bfs")
	var visited= make([]bool, NUMDESTINATIONS)
	var pred= make([]Destination, NUMDESTINATIONS)
	var dist=make([]int, NUMDESTINATIONS)
	var trackids=make([]int,NUMDESTINATIONS)

	var returnTracks=make([]int,0)

	for i:=0;i<NUMDESTINATIONS;i++{
		visited[i]=false
		pred[i]=-1
		dist[i]=100000
		trackids[i]=-1
	}

	visited[dt1]=true
	dist[dt1]=0

	var q []Destination //a queue of nodes to visit, by Destination
	q=append(q,dt1)
	for  len(q)>0{
		currentDT:=q[0]
		q=q[1:]
		for _,trackid := range b.adjacencyList[currentDT]{
			if b.trackList[trackid].d1!=currentDT && !visited[b.trackList[trackid].d1] && (b.trackStatus[trackid]==b.myNumber || b.trackStatus[trackid]==-1) {
				visited[b.trackList[trackid].d1]=true
				q=append(q,b.trackList[trackid].d1 )
				dist[b.trackList[trackid].d1]= dist[currentDT]+1
				pred[b.trackList[trackid].d1]=currentDT
				trackids[b.trackList[trackid].d1]=trackid

				if b.trackList[trackid].d1==dt2{
					counterMovingBack:=dt2
					for pred[counterMovingBack]!=-1{
						returnTracks=append(returnTracks,trackids[counterMovingBack])
						counterMovingBack=pred[counterMovingBack]
					}
					break
					//DESTINATION FOUND, PROCESS TO RETURN TRACK IDS
				}
			}else if b.trackList[trackid].d2!=currentDT && !visited[b.trackList[trackid].d2] && (b.trackStatus[trackid]==b.myNumber || b.trackStatus[trackid]==-1) {
				visited[b.trackList[trackid].d2]=true
				q=append(q,b.trackList[trackid].d2 )
				dist[b.trackList[trackid].d2]= dist[currentDT]+1
				pred[b.trackList[trackid].d2]=currentDT
				trackids[b.trackList[trackid].d2]=trackid

				if b.trackList[trackid].d2==dt2{
					counterMovingBack:=dt2
					for pred[counterMovingBack]!=-1{
						returnTracks=append(returnTracks,trackids[counterMovingBack])
						counterMovingBack=pred[counterMovingBack]
					}
					break
					//DESTINATION FOUND, PROCESS TO RETURN TRACK IDS
				}
			}
			}

		}


	return returnTracks
}

func (b* ZebraBot) askTrackLay() (int, GameColor){
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



func (b* ZebraBot) askMove() int{
	//
	//fmt.Println("inside askMove")
	whichTrack,_ := b.whichTrackCanILay()
	if whichTrack!=-1 {
		return 1
	} else {
		return 0
	}

} //Ask the player what move he wants to do: 0 is pick up cards, 1 is place Tracks, 2 is pick destination ticket



func (b* ZebraBot) askPickup(howManyLeft int, faceUpCards[]int) GameColor {
	b.faceUpCards = faceUpCards
	mostreq:=0
	mostreqind:=-1
	for i,color:= range b.faceUpCards{
		if color==0 {
			continue
		}else{
			if b.cardsAnimalsWant[b.myNumber][i]-color>mostreq{
				mostreq=b.cardsAnimalsWant[b.myNumber][i]-color
				mostreqind=i
			}
		}
	}
	if mostreqind!=-1{
		b.cardsAnimalsWant[b.myNumber][mostreqind]--
		return listOfGameColors[mostreqind]
	}

	return Other
}   //ask this player, given the gamestate, which card he wants to pick up




func (b* ZebraBot) giveTrainCard(card GameColor) {
	b.myTrainCards[card]++
}                 //tell this player he has another card of given color



func (b* ZebraBot) giveDestinationTicket(d DestinationTicket) {
	fmt.Println("inside destination ticket offering to Player", b.myNumber)
	b.myDestinationTickets=append(b.myDestinationTickets, d)
	b.DestinationTicketStatus=append(b.DestinationTicketStatus,-1)
	fmt.Println("Player", b.myNumber, " got Destination ticket from",destinationNames[d.d1],"TO", destinationNames[d.d2], "worth points ", d.points)
	//fmt.Println("Destination Cards giving")
	//	basic player doesn't care about destination tickets, so do nothing
} //tell this player has a destination card



func (b* ZebraBot) offerDestinationTickets(dtlist []DestinationTicket,howmany int) []int {
	//basic player doesn't care about destination tickets, so just pick the first bunch of tickets
	//min is howmany, max is offered slice size
	ret := make([]int,0)
	mini:=1000
	miniind:=-1
	for i:=0;i<len(dtlist);i++ {
		if dtlist[i].points<mini{
			mini=dtlist[i].points
			miniind=i
		}
	}
	for i:=0;i<len(dtlist);i++ {
		if i!=miniind{
		ret=append(ret,i)}
	}
	return ret
}//offer a list of destination cards and tell the player to take some of them package
