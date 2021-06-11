package main

import (
	"math"
	"math/rand"
	"sort"
)

//	what we're going to do is compute a score for each track
//	bfs from each destination ticket endpoint, so compute for each destination ticket,track => how far the track is from endpoints
//	the tracks which have minimum of these is useful for connecting destination tickets
//  we score these proportional to the value of the destination ticket as well
//  we also score tracks based on their length (higher length is more points, so useful)
//  and also based on the ease of getting the track (number of cards of that color in hand, and some function of whether that color is on the board)
//  we also add a boost to the score of a track if it was the same track we selected in the last
//  our decision making is to basically select a road based on a probability distribution of these scores for each track, and then play the optimal move for that track
//  i.e. if you can place it, place it, else if you can't, pick up the right card for it

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64
const longerPathMultiplier = 1.0
const destinationTicketMultiplier = 1.0
const trackBonusMultiplier = 1.0
const valueOfCardOnTable = 0.5
const difficultyOfGettingMultiplier = 0.9
const constantForRepeat = 1


const (
	MaxInt  = 1<<(UintSize-1) - 1 // 1<<31 - 1 or 1<<63 - 1
	MaxUint = 1<<UintSize - 1     // 1<<32 - 1 or 1<<64 - 1
)

type AardvarkPlayer struct {
	trackList []Track //my copy of the board
	trackStatus []int //my copy of the status of each track
	faceUpCards []int

	myTrainCards []int //number of each color card I have
	myDestinationTickets []DestinationTicket //list of destination ticket I have
	myNumber int //my player ID
	myTrains int
	constants GameConstants

	adjacencyList [][]int
	trackScores []float32

	lastChosentrack int
}

func (a *AardvarkPlayer) populateAdjacencyList() {
	a.adjacencyList = make([][]int, a.constants.NumDestinations)
	for i := 0; i < a.constants.NumDestinations; i++ {
		a.adjacencyList[i] = make([]int, 0)
	}

	for i, edge := range a.trackList {
		a.adjacencyList[edge.d1] = append(a.adjacencyList[edge.d1], i)
		a.adjacencyList[edge.d2] = append(a.adjacencyList[edge.d2], i)
	}

}


func (a* AardvarkPlayer) initialize(myNumber int, trackList []Track, constants GameConstants) {
	a.myNumber = myNumber
	a.trackList = trackList
	a.constants = constants
	a.myTrains = constants.NumStartingTrains

	a.myTrainCards = make([]int, a.constants.NumGameColors)
	a.myDestinationTickets=make([]DestinationTicket,0)

	a.lastChosentrack = -1

	a.populateAdjacencyList()
}

func (a *AardvarkPlayer) getOtherDestination(d Destination, t Track) Destination {
	if d == t.d1 {
		return t.d2
	} else if d == t.d2 {
		return t.d1
	} else {
		panic("This isn't the right edge")
	}
}

func(a *AardvarkPlayer) getEdgeDistancesFromTarget(d Destination, otherTarget Destination) ([]int, bool) {
	//TODO: using n^2 djikstra, switch to n log n later
	seen := make([]bool, a.constants.NumDestinations)
	dist := make([]int, a.constants.NumDestinations)
	for i:=0;i<a.constants.NumDestinations;i++ {
		dist[i]=MaxInt
	}
	dist[d]=0

	for numIter:=0;numIter<a.constants.NumDestinations;numIter++ {
		cheapestUnseen := -1
		cheapestVal := MaxInt
		for i,val := range dist {
			if val<cheapestVal {
				cheapestVal=val
				cheapestUnseen=i
			}
		}

		seen[cheapestUnseen]=true

		for _,edge := range a.adjacencyList[cheapestUnseen] {
			if a.trackStatus[edge] != -1 && a.trackStatus[edge] != a.myNumber {
				continue
			}
			otherDest := a.getOtherDestination(Destination(cheapestUnseen), a.trackList[edge])
			if a.trackStatus[edge]==-1 {
				dist[otherDest] = min(dist[otherDest], cheapestVal+1)
			} else if a.trackStatus[edge] == a.myNumber {
				dist[otherDest] = min(dist[otherDest], cheapestVal)
			}
		}
	}

	edgeDistances := make([]int,a.constants.NumTracks)

	if dist[otherTarget] == 0 {
		//	we're already done with this destination ticket, return 0
		return edgeDistances, false
	}

	for i,edge := range a.trackList {
		edgeDistances[i]=min(dist[edge.d1], dist[edge.d2])
	}

	return edgeDistances, true

	////	create an empty queue
	//seen := make([]bool, a.constants.NumDestinations)
	//queue := list.New()
	//queue.PushBack(d)
	//
	////	bfs over my edges
	//for queue.Len() > 0 {
	//	currDestInterface := queue.Remove(queue.Front())
	//	currDest := currDestInterface.(Destination)
	//
	//	for _,edge := range a.adjacencyList[currDest] {
	//		if a.trackStatus[edge] != a.myNumber {
	//			continue
	//		}
	//		otherDest := a.getOtherDestination(currDest, a.trackList[edge])
	//		if !seen[otherDest] {
	//			seen[otherDest] = true
	//			queue.PushBack(otherDest)
	//		}
	//	}
	//}
	//
	////mark all vertex Distances as infinity
	//vertexDistances := make([]int,a.constants.NumDestinations)
	//for i:=0;i<a.constants.NumDestinations;i++ {
	//	vertexDistances[i]=MaxInt
	//}
	//
	////mark all edge Distances as infinity
	//edgeDistances := make([]int,a.constants.NumTracks)
	//for i:=0;i<a.constants.NumTracks;i++ {
	//	edgeDistances[i]=MaxInt
	//}
	//
	////	queue is now empty, and we have marked all locations accessible from my location with true in seen
	//for i,hasSeen := range seen {
	//	if hasSeen {
	//		queue.PushBack(i)
	//		vertexDistances[i]=0
	//	}
	//}
	//
	//
	////	now, BFS on all empty edges
	//for queue.Len() > 0 {
	//	currDestInterface := queue.Remove(queue.Front())
	//	currDest := currDestInterface.(Destination)
	//
	//	for _,edge := range a.adjacencyList[currDest] {
	//		if a.trackStatus[edge] != -1 && a.trackStatus[edge] != a.myNumber {
	//			continue
	//		}
	//
	//		edgeDistances[edge]=min(edgeDistances[edge], vertexDistances[currDest]+1)
	//		otherDest := a.getOtherDestination(currDest, a.trackList[edge])
	//		vertexDistances[otherDest]=min(vertexDistances[otherDest], vertexDistances[currDest]+1)
	//		if !seen[otherDest] {
	//			seen[otherDest] = true
	//			queue.PushBack(otherDest)
	//		}
	//	}
	//}
	//
	//return edgeDistances
}

func (a *AardvarkPlayer) getDTscore(dt DestinationTicket) []float32{
	edgeDistances1,ok := a.getEdgeDistancesFromTarget(dt.d1, dt.d2)
	if !ok {
		return make([]float32, a.constants.NumTracks)
	}
	edgeDistances2,ok := a.getEdgeDistancesFromTarget(dt.d2, dt.d1)
	if !ok {
		return make([]float32, a.constants.NumTracks)
	}
	edgeDistanceSum := make([]int, a.constants.NumTracks)
	sumItems := make(map[int][]int)
	for i:=0;i<a.constants.NumTracks;i++ {
		if edgeDistances1[i] == MaxInt || edgeDistances2[i] == MaxInt {
			edgeDistanceSum[i]=MaxInt
		} else {
			edgeDistanceSum[i]=edgeDistances1[i]+edgeDistances2[i]
		}
		sumItems[edgeDistanceSum[i]]=append(sumItems[edgeDistanceSum[i]],i)
	}

	uniqueValues := make([]int, 0)
	for i := range sumItems {
		uniqueValues = append(uniqueValues, i)
	}

	sort.Ints(uniqueValues)

	ans := make([]float32, a.constants.NumTracks)

	initialMultiplier := float32(1.0)

	for _,val := range uniqueValues {
		for _,edge := range sumItems[val] {
			ans[edge]=initialMultiplier/float32(val+1)
		}
		initialMultiplier *= longerPathMultiplier
	}

	return ans
}

func (a* AardvarkPlayer) difficultyOfGettingTrack(trid int) float32{
	ans := float32(a.trackList[trid].length)
	if a.trackList[trid].c == Other {
		maxVal := float32(0.0)
		temp := float32(0.0)
		for _,c := range listOfGameColors {
			temp=0
			temp += float32(a.myTrainCards[c])
			temp += float32(a.myTrainCards[Rainbow])
			temp += valueOfCardOnTable*float32(a.faceUpCards[c])
			temp += valueOfCardOnTable*float32(a.faceUpCards[Rainbow])
			maxVal = float32(math.Max(float64(maxVal), float64(temp)))
		}
		ans-=maxVal
	} else {
		ans -= float32(a.myTrainCards[a.trackList[trid].c])
		ans -= float32(a.myTrainCards[Rainbow])
		ans -= valueOfCardOnTable*float32(a.faceUpCards[a.trackList[trid].c])
		ans -= valueOfCardOnTable*float32(a.faceUpCards[Rainbow])
	}
	ans = float32(math.Max(float64(ans), 0.0))
	return ans
}

func (a* AardvarkPlayer) informStatus(trackStatus []int, faceUpCards []int) {
	a.faceUpCards=faceUpCards
	a.trackStatus=trackStatus

	a.trackScores = make([]float32,a.constants.NumTracks)

	//	first, for each destination ticket, let's compute a score and add that to the current score
	for _,dt := range a.myDestinationTickets {
		thisDtScore := a.getDTscore(dt)
		for i:=0;i<a.constants.NumTracks;i++ {
			a.trackScores[i]+=destinationTicketMultiplier*thisDtScore[i]*float32(dt.points)*2.0
		}
	}

	// add the score for building the road
	for i:=0;i<a.constants.NumTracks;i++ {
		a.trackScores[i]+=trackBonusMultiplier*float32(a.constants.routeLengthScores[a.trackList[i].length])
	}

	//	divide score based on difficulty to get that card
	for i:=0;i<a.constants.NumTracks;i++ {
		a.trackScores[i]*=float32(math.Pow(difficultyOfGettingMultiplier, float64(a.difficultyOfGettingTrack(i))))
	}

	if a.lastChosentrack != -1 {
		a.trackScores[a.lastChosentrack]+=constantForRepeat
	}

	//set to 0 for blocked tracks
	totalScore := float32(0)

	for i:=0;i<a.constants.NumTracks;i++ {
		if a.trackStatus[i]!=-1 {
			a.trackScores[i]=0
		}
		totalScore += a.trackScores[i]
	}

	//	make it so that sum is One
	for i:=0;i<a.constants.NumTracks;i++ {
		a.trackScores[i]/=totalScore
	}

}

func (a* AardvarkPlayer) informCardPickup(int, GameColor) {
	//	do nothing
}

func (a* AardvarkPlayer) informTrackLay(int, int) {
	//	do nothing
}

func (a* AardvarkPlayer) informDestinationTicketPickup(int) {
	//	do nothing
}

func (a* AardvarkPlayer) askTrackLay() (int, GameColor){
	canLay,c := a.canILayThisTrack(a.lastChosentrack)
	if !canLay {
		panic("I THOUGHT I COULD LAY THIS TRACK BUT I CANT")
	}
	return a.lastChosentrack, c
} //ask this player which track he wants to lay, and with what color

func (a* AardvarkPlayer) canILayThisTrack(trid int) (bool, GameColor) {
	bestColor := a.trackList[trid].c
	bestColorVal := -1
	if a.trackList[trid].c==Other {
		for _, allcolor:=range listOfGameColors{
			if allcolor!=Rainbow {
				if a.trackList[trid].length <= a.myTrainCards[allcolor]+a.myTrainCards[Rainbow]{
					return true, allcolor
				} else if a.myTrainCards[allcolor]+a.myTrainCards[Rainbow] > bestColorVal{
					bestColorVal = a.myTrainCards[allcolor]+a.myTrainCards[Rainbow]
					bestColor = allcolor
				}
			}
		}
	} else {
		if a.myTrainCards[a.trackList[trid].c] + a.myTrainCards[Rainbow] >= a.trackList[trid].length {
			return true, a.trackList[trid].c
		}
	}
	return false, bestColor
}

func (a* AardvarkPlayer) askMove() int{
	//	in an askMove, we should have already filled trackScores, so here we just randomly sample from the distribution
	randomNumber := rand.Float32()
	selector := 0
	cumulativeProbability := float32(0)
	for ;selector<a.constants.NumTracks;selector++ {
		cumulativeProbability += a.trackScores[selector]
		if cumulativeProbability>=randomNumber {
			break
		}
	}

	if selector >= a.constants.NumTracks || a.trackStatus[selector]!=-1 {
		panic("SOMETHING BAD HAPPENED")
	}

	a.lastChosentrack = selector

	canLay,_ := a.canILayThisTrack(a.lastChosentrack)

	if canLay {
		return 1
	} else {
		return 0
	}
} //Ask the player what move he wants to do: 0 is pick up cards, 1 is place Tracks, 2 is pick destination ticket

func (a* AardvarkPlayer) askPickup(howManyLeft int) GameColor {
	canLayTrack, c := a.canILayThisTrack(a.lastChosentrack)
	if canLayTrack {
		panic("I Thought I couldn't lay the track but I can")
	}
	if a.faceUpCards[c] > 0 {
		return c
	}
	if a.faceUpCards[Rainbow] > 0 && howManyLeft > 1 {
		return Rainbow
	}
	return Other
}   //ask this player, given the gamestate, which card he wants to pick up


func (a* AardvarkPlayer) giveTrainCard(card GameColor) {
	a.myTrainCards[card]++
} //tell this player he has another card of given color

func (a* AardvarkPlayer) giveDestinationTicket(d DestinationTicket) {
	a.myDestinationTickets = append(a.myDestinationTickets, d)
} //tell this player has a destination card


func (a* AardvarkPlayer) offerDestinationTickets(dtlist []DestinationTicket,howmany int) []int {
	//basic player doesn't care about destination tickets, so just pick the first bunch of tickets
	ret := make([]int,0)
	for i:=0;i<howmany;i++ {
		ret=append(ret,i)
	}
	return ret
}//offer a list of destination cards and tell the player to take some of them