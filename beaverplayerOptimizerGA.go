package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
)

const beaverIndividualSize = 9
const fileWriteFrequency = 1

type individual []float64

type scorestruct struct {
	idx int
	score float64
}

type GA_Beaver struct {
	populationSize int
	crossoverRate,mutateRate float64
	otherMutateRate, powerOfScore float64
	elitismCount int
	numGamesInTournament int

	population []individual
	popscores []scorestruct

}

func (g* GA_Beaver) mutate(x *individual){
	for i:=0;i<beaverIndividualSize;i++ {
		if rand.Float64() < g.mutateRate {
			(*x)[i] = rand.Float64()
		}
	}
}

func (g* GA_Beaver) mutate2(x *individual){
	for i:=0;i<beaverIndividualSize;i++ {
		if rand.Float64() < g.otherMutateRate {
			if rand.Intn(2)==0 {
				(*x)[i]*=1.3
				(*x)[i] = math.Min((*x)[i], 1.0)
			} else {
				(*x)[i]*=0.7
			}
		}
	}
}


func (g* GA_Beaver) crossover(x,y *individual){
	for i:=0;i<beaverIndividualSize;i++ {
		if rand.Float64() < g.crossoverRate {
			(*x)[i],(*y)[i] = (*y)[i],(*x)[i]
		}
	}
}

func (g GA_Beaver) randomIndividual() individual {
	var x individual
	for i:=0;i<beaverIndividualSize;i++ {
		x=append(x,rand.Float64())
	}
	return x
}

func (g* GA_Beaver) fillInParameters(filename string) {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	fmt.Fscanf(file, "%d\n%f\n%f\n%f\n%f\n%d\n%d", &g.populationSize, &g.crossoverRate, &g.mutateRate, &g.otherMutateRate, &g.powerOfScore, &g.elitismCount, &g.numGamesInTournament)
}

func (g* GA_Beaver) fillInPopulationScores(){
	g.popscores = make([]scorestruct, g.populationSize)
	for i:=0;i<g.populationSize;i++ {
		g.popscores[i].idx=i
	}

	totSum := 0.0

	for numIter:=0;numIter<g.numGamesInTournament;numIter++ {

		channelArray := make([][]chan bool, g.populationSize)
		for i:=0;i<g.populationSize;i++ {
			channelArray[i] = make([]chan bool, g.populationSize)
			for j:=0;j<g.populationSize;j++ {
				channelArray[i][j] = make(chan bool)
			}
		}

		for i:=0;i<g.populationSize;i++ {
			for j:=0;j<g.populationSize;j++ {
				if i==j {
					continue
				}
				go g.twoWayTourney(g.population[i],g.population[j],channelArray[i][j])
			}
		}

		for i:=0;i<g.populationSize;i++ {
			for j:=0;j<g.populationSize;j++ {
				if i==j {
					continue
				}
				res := <- channelArray[i][j]
				//fmt.Println("RES", res)
				if res {
					g.popscores[j].score+=1.0
				} else {
					g.popscores[i].score += 1.0
				}
				totSum += 1.0
			}
		}
	}

	sort.Slice(g.popscores, func(i, j int) bool {
		return g.popscores[i].score>g.popscores[j].score
	})

	newTotSum := 0.0
	for i:=0;i<g.populationSize;i++ {
		g.popscores[i].score/=totSum
		g.popscores[i].score = math.Pow(g.popscores[i].score,g.powerOfScore)
		newTotSum+=g.popscores[i].score
		//fmt.Println(i, g.popscores[i])
	}

	for i:=0;i<g.populationSize;i++ {
		g.popscores[i].score/=newTotSum
	}

	sort.Slice(g.popscores, func(i, j int) bool {
		return g.popscores[i].score>g.popscores[j].score
	})


}

//func (g* GA_Beaver) safeGame(e* Engine, c *GameConstants, players *[]Player) []int{
//	//fmt.Println("AM HERE")
//	//fmt.Println(e.OptimizerMode)
//	toReturn := []int{-1}
//	//defer func() {
//	//	if r:= recover();r!=nil {
//	//		//fmt.Println("Recovering")
//	//		toReturn = []int{-1}
//	//	}
//	//}()
//	toReturn = e.runGame(*players,*c)
//	return toReturn
//}

func (g* GA_Beaver) twoWayTourney(a,b individual, ch chan bool){
	constants := GameConstants{
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
		LongestPathScore:                    LONGESTPATHSCORE,
		NumPlayers:                          0,
		NumTracks:                           0,
		NumDestinations:                     NUMDESTINATIONS,
		routeLengthScores:                   routeLengthScores,
	}
	e := Engine{}
	e.OptimizerMode = true
	players := make([]Player, 0)
	player1 := BeaverPlayer{}
	player1.setScoringParameters(a)
	players = append(players, &player1)
	player2 := BeaverPlayer{}
	player2.setScoringParameters(b)
	players = append(players, &player2)
	player3 := BeaverPlayer{}
	player3.setScoringParameters(a)
	players = append(players, &player3)
	player4 := BeaverPlayer{}
	player4.setScoringParameters(b)
	players = append(players, &player4)

	winners := e.runGame(players, constants)
	//fmt.Println(winners)
	if winners[0] == 0 || winners[0]==2 {
		//fmt.Println("Writing To Channel")
		ch <- false
	} else {
		//fmt.Println("Writing To Channel")
		ch <- true
	}

}

func (g* GA_Beaver) tournament(inds [4]individual) int{

	scores := make([]int, 4)
	for i:=0;i<g.numGamesInTournament;i++ {

		constants := GameConstants{
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
			LongestPathScore:                    LONGESTPATHSCORE,
			NumPlayers:                          0,
			NumTracks:                           0,
			NumDestinations:                     NUMDESTINATIONS,
			routeLengthScores:                   routeLengthScores,
		}
		e := Engine{}
		e.OptimizerMode = true
		players := make([]Player, 0)
		player1 := BeaverPlayer{}
		player1.setScoringParameters(inds[0])
		players = append(players, &player1)
		player2 := BeaverPlayer{}
		player2.setScoringParameters(inds[1])
		players = append(players, &player2)
		player3 := BeaverPlayer{}
		player3.setScoringParameters(inds[2])
		players = append(players, &player3)
		player4 := BeaverPlayer{}
		player4.setScoringParameters(inds[3])
		players = append(players, &player4)

		winners := e.runGame(players, constants)
		for _,winner := range winners {
			scores[winner]++
		}
	}

	maxWinner := -1
	maxWins := 0

	for i,score := range scores {
		if score > maxWins {
			maxWins = score
			maxWinner = i
		}
	}

	return maxWinner
}

func (g* GA_Beaver) selectRandomFromFourWayTournament() int{
	ints := make([]int,4)
	ints[0] = rand.Intn(g.populationSize)
	ints[1] = rand.Intn(g.populationSize)
	ints[2] = rand.Intn(g.populationSize)
	ints[3] = rand.Intn(g.populationSize)
	inds := [4]individual{g.population[ints[0]],g.population[ints[1]],g.population[ints[2]],g.population[ints[3]]}
	winner := g.tournament(inds)
	if winner==-1 {
		return 0
	}
	return ints[winner]
}

func (g* GA_Beaver) rouletteWheelSelection() int {
	randomNumber := rand.Float64()
	for i:=0;i<g.populationSize;i++ {
		randomNumber -= g.popscores[i].score
		if randomNumber <= 0 {
			return g.popscores[i].idx
		}
	}
	return -1
}

func optimizeBeaverParametersWithGeneticAlgorithm() {
	g := GA_Beaver{}
	g.fillInParameters("gaparams.txt")

	//TODO: seeding
	g.population = append(g.population, individual{0.5,0.5,0.1,0.18,1,0.1,0.001,0.01,0.00})
	g.population = append(g.population, individual{0.5,0.5,0.1,0.18,1,0.1,0.001,0.01,0.10})


	for len(g.population) < g.populationSize {
		g.population = append(g.population,g.randomIndividual())
	}

	iterationNumber := 0

	for {
		//fetch params
		g.fillInParameters("gaparams.txt")
		fmt.Println("ITER NO:", iterationNumber)
		fmt.Println("PARAMS ARE:", g.populationSize, g.crossoverRate, g.mutateRate, g.otherMutateRate, g.powerOfScore, g.elitismCount, g.numGamesInTournament)
		fmt.Println("BEST MEMBER IS:", g.population[0])

		//score each individual
		g.fillInPopulationScores()

		//fmt.Println(g.popscores)

		newPopulation := make([]individual, 0)
		//elitism
		for i:=0;i<g.elitismCount;i++ {
			newPopulation = append(newPopulation, g.population[g.popscores[i].idx])
		}

		for len(newPopulation)<g.populationSize {
			//fmt.Println(len(newPopulation))
			parent1 := make(individual, beaverIndividualSize)
			parent2 := make(individual, beaverIndividualSize)
			copy(parent1,g.population[g.rouletteWheelSelection()])
			copy(parent2,g.population[g.rouletteWheelSelection()])
			//copy(parent2,g.population[g.selectRandomFromFourWayTournament()])
			g.crossover(&parent1, &parent2)
			g.mutate(&parent1)
			g.mutate(&parent2)
			g.mutate2(&parent1)
			g.mutate2(&parent2)
			newPopulation = append(newPopulation, parent1)
			newPopulation = append(newPopulation, parent2)
		}
		g.population = newPopulation

		if iterationNumber%fileWriteFrequency == 0 {
			file, err := os.OpenFile("garesults.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			//defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(file,"Iteration Number", iterationNumber)
			fmt.Fprintln(file,g.population)
			fmt.Fprintln(file)
			file.Close()
		}
		iterationNumber++
	}

}
