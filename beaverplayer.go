package main

import (
	"container/heap"
	"strconv"

	"github.com/wangjia184/sortedset"
	//"fmt"
	"math"
	"math/rand"
	"sort"
)

//	beaverPlayer is gonna be a copy of aardvarkPlayer, but with score parameters passable through a function, to suit GA based optimization
// I spent a lot of time on GA w/ just these params, and it was only ever as good as aadrvarkplayer, not better.
// So instead, the strat now is going to be to sample from the probability distibution X times, and play the majority move.

//related to computation of dt scores
const longerPathMultiplierMin = 0.0
const longerPathMultiplierMax = 1.0

const pathDenominatorPowerMin = 0.0
const pathDenominatorPowerMax = 10.0

// related to computation of difficulty scores
const valueOfCardOnTableMin = 0.0
const valueOfCardOnTableMax = 5.0

const difficultyOfGettingBaseMin = 0.0
const difficultyOfGettingBaseMax = 5.0

//related to merging the three metrics
const destinationTicketMultiplierMin = 0.0
const destinationTicketMultiplierMax = 1.0

const trackBonusMultiplierMin = 0.0
const trackBonusMultiplierMax = 1.0

const difficultyOfGettingMultiplierMin = 0.0
const difficultyOfGettingMultiplierMax = 1.0

// related to later optimization
const constantForRepeatMin = 0.0
const constantForRepeatMax = 100.0

// related to number of samples
const sampleNumberMin = 1.0
const sampleNumberMax = 101.0


type BeaverPlayer struct {
	trackList []Track //my copy of the board
	trackStatus []int //my copy of the status of each track
	faceUpCards []int

	myTrainCards []int //number of each color card I have
	myDestinationTickets []DestinationTicket //list of destination ticket I have
	myNumber int //my player ID
	myTrains int
	constants GameConstants

	adjacencyList [][]int
	trackScores []float64

	lastChosentrack int
	chosenMove int

	//related to computation of dt scores
	longerPathMultiplier float64
	pathDenominatorPower float64

	// related to computation of difficulty scores
	valueOfCardOnTable float64
	difficultyOfGettingBase float64

	//related to merging the three metrics
	destinationTicketMultiplier float64
	trackBonusMultiplier float64
	difficultyOfGettingMultiplier float64

	// related to later optimization
	constantForRepeat float64

	sampleNumber int
}

func (b *BeaverPlayer) populateAdjacencyList() {
	b.adjacencyList = make([][]int, b.constants.NumDestinations)
	for i := 0; i < b.constants.NumDestinations; i++ {
		b.adjacencyList[i] = make([]int, 0)
	}

	for i, edge := range b.trackList {
		b.adjacencyList[edge.d1] = append(b.adjacencyList[edge.d1], i)
		b.adjacencyList[edge.d2] = append(b.adjacencyList[edge.d2], i)
	}

}

func (b *BeaverPlayer) setScoringParameters(inputs [] float64) {
	//related to computation of dt scores
	b.longerPathMultiplier = scaleFloat(inputs[0], longerPathMultiplierMin, longerPathMultiplierMax)
	b.pathDenominatorPower = scaleFloat(inputs[1], pathDenominatorPowerMin, pathDenominatorPowerMax)

	// related to computation of difficulty scores
	b.valueOfCardOnTable = scaleFloat(inputs[2], valueOfCardOnTableMin, valueOfCardOnTableMax)
	b.difficultyOfGettingBase = scaleFloat(inputs[3], difficultyOfGettingBaseMin,difficultyOfGettingBaseMax)

	//related to merging the three metrics
	b.destinationTicketMultiplier = scaleFloat(inputs[4], destinationTicketMultiplierMin, destinationTicketMultiplierMax)
	b.trackBonusMultiplier = scaleFloat(inputs[5], trackBonusMultiplierMin, trackBonusMultiplierMax)
	b.difficultyOfGettingMultiplier = scaleFloat(inputs[6], difficultyOfGettingMultiplierMin, difficultyOfGettingMultiplierMax)

	// related to later optimization
	b.constantForRepeat = scaleFloat(inputs[7], constantForRepeatMin, constantForRepeatMax)

	b.sampleNumber = int(math.Floor(scaleFloat(inputs[8], sampleNumberMin, sampleNumberMax)))
}

func (b * BeaverPlayer) initialize(myNumber int, trackList []Track,adjList [][]int, constants GameConstants) {
	b.myNumber = myNumber
	b.trackList = trackList
	b.constants = constants
	b.myTrains = constants.NumStartingTrains

	b.myTrainCards = make([]int, b.constants.NumGameColors)
	b.myDestinationTickets=make([]DestinationTicket,0)

	b.lastChosentrack = -1

	b.adjacencyList = adjList
}

func (b *BeaverPlayer) getOtherDestination(d Destination, t Track) Destination {
	if d == t.d1 {
		return t.d2
	} else if d == t.d2 {
		return t.d1
	} else {
		panic("This isn't the right edge")
	}
}

func(b *BeaverPlayer) getEdgeDistancesFromTarget(d Destination, otherTarget Destination) ([]int, bool) {
	//TODO: using n^2 djikstra, switch to n log n later
	seen := make([]bool, b.constants.NumDestinations)
	dist := make([]int, b.constants.NumDestinations)
	for i:=0;i< b.constants.NumDestinations;i++ {
		dist[i]=MaxInt
	}
	dist[d]=0

	//zap.S().Debug(dist)

	for numIter:=0;numIter< b.constants.NumDestinations;numIter++ {
		cheapestUnseen := -1
		cheapestVal := MaxInt
		for i,val := range dist {
			if val<=cheapestVal && !seen[i] {
				cheapestVal=val
				cheapestUnseen=i
			}
		}

		//zap.S().Debug(seen)
		//zap.S().Debug(seen[33],destinationNames[33])
		//zap.S().Debug(numIter,cheapestUnseen, cheapestVal, b.constants.NumDestinations)

		seen[cheapestUnseen]=true

		for _,edge := range b.adjacencyList[cheapestUnseen] {
			if b.trackStatus[edge] != -1 && b.trackStatus[edge] != b.myNumber {
				continue
			}
			otherDest := b.getOtherDestination(Destination(cheapestUnseen), b.trackList[edge])
			if b.trackStatus[edge]==-1 {
				dist[otherDest] = min(dist[otherDest], cheapestVal+1)
			} else if b.trackStatus[edge] == b.myNumber {
				dist[otherDest] = min(dist[otherDest], cheapestVal)
			}
		}
	}

	//zap.S().Debug(dist)

	edgeDistances := make([]int, b.constants.NumTracks)

	if dist[otherTarget] == 0 {
		//	we're already done with this destination ticket, return 0
		return edgeDistances, false
	}

	for i,edge := range b.trackList {
		edgeDistances[i]=min(dist[edge.d1], dist[edge.d2])
	}

	return edgeDistances, true

}

func(b *BeaverPlayer) getEdgeDistancesFromTargetFast(d Destination, otherTarget Destination) ([]int, bool) {
	//TODO: using n^2 djikstra, switch to n log n later

	djikstraItemPointers := make([]*DjikstraItem, b.constants.NumDestinations)

	pq := make(PriorityQueue, b.constants.NumDestinations)
	for i:=0;i< b.constants.NumDestinations;i++ {
		djikstraItemPointers[i]=new(DjikstraItem)
		djikstraItemPointers[i].priority=MaxInt
		djikstraItemPointers[i].vertexIndex=i
		djikstraItemPointers[i].index=i
		pq[i]=djikstraItemPointers[i]
	}

	heap.Init(&pq)

	pq.update(djikstraItemPointers[d],int(d),0)
	//zap.S().Debug(dist)

	for numIter:=0;numIter< b.constants.NumDestinations;numIter++ {

		currVertex := heap.Pop(&pq).(*DjikstraItem)
		cheapestUnseen := currVertex.vertexIndex
		cheapestVal := currVertex.priority

		//zap.S().Debug(seen)
		//zap.S().Debug(seen[33],destinationNames[33])
		//zap.S().Debug(numIter,cheapestUnseen, cheapestVal, b.constants.NumDestinations)

		for _,edge := range b.adjacencyList[cheapestUnseen] {
			if b.trackStatus[edge] != -1 && b.trackStatus[edge] != b.myNumber {
				continue
			}
			otherDest := b.getOtherDestination(Destination(cheapestUnseen), b.trackList[edge])
			if b.trackStatus[edge]==-1 {
				if cheapestVal+1 < djikstraItemPointers[otherDest].priority {
					pq.update(djikstraItemPointers[otherDest], int(otherDest), cheapestVal+1)
				}
			} else if b.trackStatus[edge] == b.myNumber {
				if cheapestVal < djikstraItemPointers[otherDest].priority {
					pq.update(djikstraItemPointers[otherDest], int(otherDest), cheapestVal)
				}
			}
		}
	}

	//zap.S().Debug(dist)

	edgeDistances := make([]int, b.constants.NumTracks)

	if djikstraItemPointers[otherTarget].priority == 0 {
		//	we're already done with this destination ticket, return 0
		return edgeDistances, false
	}

	for i,edge := range b.trackList {
		edgeDistances[i]=min(djikstraItemPointers[edge.d1].priority, djikstraItemPointers[edge.d2].priority)
	}

	return edgeDistances, true
}


func(b *BeaverPlayer) getEdgeDistancesFromTargetFast2(d Destination, otherTarget Destination) ([]int, bool) {
	//TODO: using n^2 djikstra, switch to n log n later
	dist := make([]int, b.constants.NumDestinations)
	for i:=0;i< b.constants.NumDestinations;i++ {
		dist[i]=MaxInt
	}
	dist[d]=0

	set := sortedset.New()

	strings := make([]string, b.constants.NumDestinations)


	for i:=0;i< b.constants.NumDestinations;i++ {
		strings[i] = strconv.Itoa(i)
		set.AddOrUpdate(strings[i],MaxInt, i)
	}

	set.AddOrUpdate(strings[int(d)],0, int(d))

	for numIter:=0;numIter< b.constants.NumDestinations;numIter++ {

		node:=set.PopMin()
		cheapestUnseen := node.Value.(int)
		cheapestVal := int(node.Score())

		//zap.S().Debug(seen)
		//zap.S().Debug(seen[33],destinationNames[33])
		//zap.S().Debug(numIter,cheapestUnseen, cheapestVal, b.constants.NumDestinations)

		for _,edge := range b.adjacencyList[cheapestUnseen] {
			if b.trackStatus[edge] != -1 && b.trackStatus[edge] != b.myNumber {
				continue
			}
			otherDest := b.getOtherDestination(Destination(cheapestUnseen), b.trackList[edge])
			if b.trackStatus[edge]==-1 {
				if cheapestVal+1 < dist[otherDest] {
					dist[otherDest] = cheapestVal+1
					set.AddOrUpdate(strings[int(otherDest)], sortedset.SCORE(cheapestVal+1), int(otherDest))
				}
			} else if b.trackStatus[edge] == b.myNumber {
				if cheapestVal < dist[otherDest] {
					dist[otherDest] = cheapestVal
					set.AddOrUpdate(strings[int(otherDest)], sortedset.SCORE(cheapestVal), int(otherDest))
				}
			}
		}
	}

	//zap.S().Debug(dist)

	edgeDistances := make([]int, b.constants.NumTracks)

	if dist[otherTarget] == 0 {
		//	we're already done with this destination ticket, return 0
		return edgeDistances, false
	}

	for i,edge := range b.trackList {
		edgeDistances[i]=min(dist[edge.d1], dist[edge.d2])
	}

	return edgeDistances, true
}




func (b *BeaverPlayer) getDTscore(dt DestinationTicket) []float64{
	edgeDistances1,ok := b.getEdgeDistancesFromTarget(dt.d1, dt.d2)
	if !ok {
		return make([]float64, b.constants.NumTracks)
	}
	edgeDistances2,ok := b.getEdgeDistancesFromTarget(dt.d2, dt.d1)
	if !ok {
		return make([]float64, b.constants.NumTracks)
	}
	edgeDistanceSum := make([]int, b.constants.NumTracks)
	sumItems := make(map[int][]int)
	for i:=0;i< b.constants.NumTracks;i++ {
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

	ans := make([]float64, b.constants.NumTracks)

	initialMultiplier := 1.0

	for _,val := range uniqueValues {
		//fmt.Println(val)
		if val == MaxInt {
			continue
		}
		for _,edge := range sumItems[val] {
			//zap.S().Debug(val,edge)
			//if val+1 == 0 {
			//	panic("WTF")
			//}
			ans[edge]=initialMultiplier/ math.Pow(float64(val+1), b.pathDenominatorPower)
			if math.IsNaN(ans[edge]) {
				//fmt.Println(initialMultiplier, val,b.pathDenominatorPower, MaxInt)
				panic("WTF")

			}
		}
		initialMultiplier *= b.longerPathMultiplier
	}

	//zap.S().Debug(dt)
	//zap.S().Debug(edgeDistanceSum)
	//zap.S().Debug(ans)
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	//fmt.Println("In func", ans)

	return ans
}

func (b * BeaverPlayer) difficultyOfGettingTrack(trid int) float64{
	ans := float64(b.trackList[trid].length)
	if b.trackList[trid].c == Other {
		maxVal := 0.0
		temp := 0.0
		for _,c := range listOfGameColors {
			temp=0
			temp += float64(b.myTrainCards[c])
			temp += float64(b.myTrainCards[Rainbow])
			temp += b.valueOfCardOnTable*float64(b.faceUpCards[c])
			temp += b.valueOfCardOnTable*float64(b.faceUpCards[Rainbow])
			maxVal = math.Max(maxVal, temp)
		}
		ans-=maxVal
	} else {
		ans -= float64(b.myTrainCards[b.trackList[trid].c])
		ans -= float64(b.myTrainCards[Rainbow])
		ans -= b.valueOfCardOnTable*float64(b.faceUpCards[b.trackList[trid].c])
		ans -= b.valueOfCardOnTable*float64(b.faceUpCards[Rainbow])
	}
	ans = math.Max(ans, 0.0)
	return ans
}

func (b * BeaverPlayer) informStatus(trackStatus []int, faceUpCards []int) {
	b.faceUpCards=faceUpCards
	b.trackStatus=trackStatus

	b.trackScores = make([]float64, b.constants.NumTracks)
	destinationTicketScores := make([]float64, b.constants.NumTracks)
	trackLengthScores := make([]float64, b.constants.NumTracks)
	trackDifficultyScores := make([]float64, b.constants.NumTracks)

	//	first, for each destination ticket, let's compute dt score
	for _,dt := range b.myDestinationTickets {
		thisDtScore := b.getDTscore(dt)
		//fmt.Println(thisDtScore)
		for i:=0;i< b.constants.NumTracks;i++ {
			destinationTicketScores[i]+=thisDtScore[i]*float64(dt.points)
		}
	}

	//fmt.Println("Finally", destinationTicketScores)

	// compute the score for building the track
	for i:=0;i< b.constants.NumTracks;i++ {
		if b.trackStatus[i] == -1 {
			trackLengthScores[i]=float64(b.constants.routeLengthScores[b.trackList[i].length])
		}
	}

	// compute the score based on difficulty of getting
	for i:=0;i< b.constants.NumTracks;i++ {
		if b.trackStatus[i]==-1 {
			trackDifficultyScores[i]= math.Pow(b.difficultyOfGettingBase, b.difficultyOfGettingTrack(i))
		}
	}


	//fmt.Println("Point 0 dt", destinationTicketScores)
	//normalize the three computed scores
	normalizeFloatSlice(&destinationTicketScores)
	normalizeFloatSlice(&trackLengthScores)
	normalizeFloatSlice(&trackDifficultyScores)

	//fmt.Println("Point 1", b.trackScores)
	//fmt.Println("Point 1 dt", destinationTicketScores)


	// merge the three computed scores
	for i:=0;i< b.constants.NumTracks;i++ {
		b.trackScores[i] = destinationTicketScores[i]*b.destinationTicketMultiplier + trackLengthScores[i]*b.trackBonusMultiplier + trackDifficultyScores[i]*b.difficultyOfGettingMultiplier
	}

	//fmt.Println("Point 2", b.trackScores)

	//for benefitting repeats
	if b.lastChosentrack != -1 && b.lastChosentrack != b.constants.NumTracks {
		b.trackScores[b.lastChosentrack]*=b.constantForRepeat
	}

	//fmt.Println("Point 3", b.trackScores)


	//set to 0 for blocked tracks
	for i:=0;i< b.constants.NumTracks;i++ {
		if b.trackStatus[i]!=-1 || b.myTrains < b.trackList[i].length {
			b.trackScores[i]=0
		}
	}
	normalizeFloatSlice(&b.trackScores)

	//fmt.Println("Point 4", b.trackScores)

	//fmt.Println("HI", b.trackScores, b.trackScores[0] == )

	//if b.myNumber == 0 {
	//	for i,score := range b.trackScores {
	//		zap.S().Debug(i,score)
	//	}
	//	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//}
}

func (b* BeaverPlayer) informCardPickup(int, GameColor) {
	//	do nothing
}

func (b* BeaverPlayer) informTrackLay(int, int) {
	//	do nothing
}

func (b* BeaverPlayer) informDestinationTicketPickup(int) {
	//	do nothing
}

func (b * BeaverPlayer) askTrackLay() (int, GameColor){
	canLay,c := b.canILayThisTrack(b.lastChosentrack)
	if !canLay {
		panic("I THOUGHT I COULD LAY THIS TRACK BUT I CANT")
	}

	if b.trackList[b.lastChosentrack].length> b.myTrainCards[c]{
		b.myTrainCards[Rainbow]-= b.trackList[b.lastChosentrack].length- b.myTrainCards[c]
		b.myTrainCards[c]=0
	} else{
		b.myTrainCards[c]-= b.trackList[b.lastChosentrack].length
	}
	b.myTrains-= b.trackList[b.lastChosentrack].length

	return b.lastChosentrack, c
} //ask this player which track he wants to lay, and with what color

func (b * BeaverPlayer) canILayThisTrack(trid int) (bool, GameColor) {
	bestColor := b.trackList[trid].c
	bestColorVal := -1
	if b.trackList[trid].c==Other {
		for _, allcolor:=range listOfGameColors{
			if allcolor!=Rainbow {
				if b.trackList[trid].length <= b.myTrainCards[allcolor]+b.myTrainCards[Rainbow]{
					return true, allcolor
				} else if b.myTrainCards[allcolor]+b.myTrainCards[Rainbow] > bestColorVal{
					bestColorVal = b.myTrainCards[allcolor]+ b.myTrainCards[Rainbow]
					bestColor = allcolor
				}
			}
		}
	} else {
		if b.myTrainCards[b.trackList[trid].c] + b.myTrainCards[Rainbow] >= b.trackList[trid].length {
			return true, b.trackList[trid].c
		}
	}
	return false, bestColor
}

func (b *BeaverPlayer) getBestMoveForTrack(trid int) int {
	canLay, c := b.canILayThisTrack(trid)
	if canLay {
		return trid
	} else {
		return b.constants.NumTracks + int(c)
	}
}

func (b * BeaverPlayer) askMove() int{
	//	in an askMove, we should have already filled trackScores, so here we just randomly sample from the distribution
	// randomly sample sampleNumber times

	sampleResults := make(map[int]int)
	bestSampleCount := 0


	for i:=0;i<b.sampleNumber;i++ {
		randomNumber := rand.Float64()
		selector := 0
		cumulativeProbability := float64(0)
		for ;selector< b.constants.NumTracks;selector++ {
			cumulativeProbability += b.trackScores[selector]
			if cumulativeProbability>=randomNumber {
				break
			}
		}

		if selector < b.constants.NumTracks && b.trackStatus[selector]!=-1 {
			//zap.S().Debug(selector,randomNumber)
			//zap.S().Debug(b.trackScores)
			panic("SOMETHING BAD HAPPENED")
		}

		theMove := 0
		if selector == b.constants.NumTracks {
			theMove = b.constants.NumTracks
			//sum := 0.0
			//for _,x := range b.trackScores {
			//	sum += x
			//}
			//fmt.Println(sum, b.myTrains)
			//panic("WUT")
		} else {
			theMove = b.getBestMoveForTrack(selector)
		}
		sampleResults[theMove]++
		bestSampleCount = max(bestSampleCount, sampleResults[theMove])
	}

	//fmt.Println(b.sampleNumber, bestSampleCount, b.myTrains)

	//if selector
	moveselectionSlice := make([]int, 0)
	for i := range sampleResults {
		if sampleResults[i] == bestSampleCount {
			moveselectionSlice = append(moveselectionSlice,i)
		}
	}

	b.chosenMove = moveselectionSlice[rand.Intn(len(moveselectionSlice))]
	if b.chosenMove < b.constants.NumTracks {
		b.lastChosentrack = b.chosenMove
		return 1
	} else {
		return 0
	}

} //Ask the player what move he wants to do: 0 is pick up cards, 1 is place Tracks, 2 is pick destination ticket

func (b * BeaverPlayer) askPickup(howManyLeft int, faceUpCards []int) GameColor {
	c := GameColor(b.chosenMove - b.constants.NumTracks)
	if b.faceUpCards[c] > 0 {
		return c
	}
	if b.faceUpCards[Rainbow] > 0 && howManyLeft > 1 {
		return Rainbow
	}
	return Other
}   //ask this player, given the gamestate, which card he wants to pick up


func (b* BeaverPlayer) giveTrainCard(card GameColor) {
	b.myTrainCards[card]++
} //tell this player he has another card of given color

func (b* BeaverPlayer) giveDestinationTicket(d DestinationTicket) {
	b.myDestinationTickets = append(b.myDestinationTickets, d)
} //tell this player has a destination card


func (b* BeaverPlayer) offerDestinationTickets(dtlist []DestinationTicket,howmany int) []int {
	//basic player doesn't care about destination tickets, so just pick the first bunch of tickets
	ret := make([]int,0)
	for i:=0;i<howmany;i++ {
		ret=append(ret,i)
	}
	return ret
}//offer a list of destination cards and tell the player to take some of them
