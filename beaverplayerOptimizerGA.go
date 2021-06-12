package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

const beaverIndividualSize = 8
const numGamesInTournament = 10

type individual []float64

type GA_Beaver struct {
	populationSize int
	crossoverRate,mutateRate float64

	population []individual

}

func (g* GA_Beaver) mutate(x *individual){
	for i:=0;i<beaverIndividualSize;i++ {
		if rand.Float64() < g.crossoverRate {
			(*x)[i] = rand.Float64()
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

	fmt.Fscanf(file, "%d\n%lf\n%lf", &g.populationSize, &g.crossoverRate, &g.mutateRate)
}

func (g* GA_Beaver) safeGame(e* Engine, c *GameConstants, players *[]Player) []int{
	//fmt.Println("AM HERE")
	//fmt.Println(e.OptimizerMode)
	toReturn := []int{-1}
	//defer func() {
	//	if r:= recover();r!=nil {
	//		//fmt.Println("Recovering")
	//		toReturn = []int{-1}
	//	}
	//}()
	toReturn = e.runGame(*players,*c)
	return toReturn
}

func (g* GA_Beaver) tournament(inds [4]individual) int{

	scores := make([]int, 4)
	for i:=0;i<numGamesInTournament;i++ {

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

		winners := g.safeGame(&e,&constants,&players)
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

func optimizeBeaverParametersWithGeneticAlgorithm() {
	g := GA_Beaver{}
	g.fillInParameters("gaparams.txt")
	//TODO: seeding
	g.population = append(g.population, individual{0.5,0.5,0.1,0.18,1,0.1,0,0.01})
	g.population = append(g.population, individual{0.5,0.5,0.1,0.18,1,0.1,0,0.01})

	for len(g.population) < g.populationSize {
		g.population = append(g.population,g.randomIndividual())
	}

	iterationNumber := 0

	for {
		//TODO: add elitism
		newPopulation := make([]individual, 0)
		for len(newPopulation)<g.populationSize {
			//fmt.Println(len(newPopulation))
			parent1 := make(individual, beaverIndividualSize)
			parent2 := make(individual, beaverIndividualSize)
			copy(parent1,g.population[g.selectRandomFromFourWayTournament()])
			copy(parent2,g.population[g.selectRandomFromFourWayTournament()])
			g.crossover(&parent1, &parent2)
			g.mutate(&parent1)
			g.mutate(&parent2)
			newPopulation = append(newPopulation, parent1)
			newPopulation = append(newPopulation, parent2)
		}
		g.population = newPopulation

		iterationNumber++
		fmt.Println(iterationNumber)
		if iterationNumber%20 == 0 {
			file, err := os.OpenFile("garesults.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			//defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			for _,ind := range g.population {
				fmt.Fprint(file,ind)
				fmt.Fprint(file, "\n")
			}
			file.Close()
		}
	}

}
