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

	//TOTAL SEEDING:
	//g.population = []individual{{0.44454935033352205, 0.5, 0.107653, 0.010350716892173813, 0.8450304277220828, 0.08914744929999999, 0.00013917876699999997, 0.2525800021749184, 1}, {0.6350705004764602, 0.5, 0.7432586245994012, 0.010350716892173813, 0.5379226477374233, 0.1459111877410948, 0.00024009999999999998, 0.4458999999999999, 1}, {0.21483373423678898, 0.31849999999999995, 0.107653, 1, 0.5379226477374233, 0.1459111877410948, 0.00010706058999999998, 0.4458999999999999, 1}, {0.6232157208577027, 0.35, 0.5202810372195807, 0.005071851277165168, 0.8450304277220828, 0.12236608647264025, 0.00019882680999999997, 1, 0.8281}, {0.42100618271204354, 0.31849999999999995, 0.37328426390196345, 0.027461085632297877, 1, 0.12120828043789884, 8.235429999999999e-05, 0.4458999999999999, 1}, {0.34196103871809386, 0.5, 0.107653, 0.010350716892173813, 0.8450304277220828, 0.08914744929999999, 0.00013917876699999997, 0.2525800021749184, 1}, {0.5779141554335787, 1, 0.0753571, 0.010350716892173813, 0.8450304277220828, 0.2084445539158497, 0.00013917876699999997, 0.2525800021749184, 1}, {0.6232157208577027, 0.65, 0.04800247269999999, 0.8419108093375564, 0.844462555317776, 0.08914744929999999, 0.00028403829999999996, 0.44534175986638375, 1}, {0.488515769597277, 0.5, 0.18193357000000002, 0.45333658964329954, 0.6992994420586502, 0.062403214509999985, 0.00031213, 0.7219137896444351, 0.4375231749678653}, {0.42100618271204354, 0.5914999999999999, 0.107653, 0.017492711547773744, 0.31829742469670014, 0.2084445539158497, 0.0001529437, 0.5534071433793584, 0.48999999999999994}, {0.21483373423678898, 0.45499999999999996, 0.10546143252190002, 0.48999999999999994, 0.5379226477374233, 0.48224371107540437, 0.00013917876699999997, 0.17680600152244286, 1}, {0.6350705004764602, 0.65, 0.7432586245994012, 0.09494994600000001, 0.31829742469670014, 0.18968454406342325, 0.00025847485299999997, 0.5796699999999999, 1}, {0.8101804371150135, 0.35, 0.5202810372195807, 0.005071851277165168, 0.8450304277220828, 0.12236608647264025, 0.00019882680999999997, 1, 0.7}, {0.3757813612286746, 0.6364049956294002, 0.5202810372195807, 0.005071851277165168, 0.5561939824459268, 0.04368225015699999, 0.00031213, 0.0536599210648232, 0.8281}, {0.34196103871809386, 0.5, 0.5202810372195807, 0.0072455018245216685, 0.5379226477374233, 0.062403214509999985, 0.00031213, 0.4458999999999999, 1}, {0.6350705004764602, 0.5, 0.107653, 0.1448366087894187, 1, 0.1459111877410948, 0.00018093239709999996, 0.2525800021749184, 1}, {0.6350705004764602, 0.5914999999999999, 0.9662362119792216, 0.31733561275030964, 0.4137866521057102, 0.10213783141876635, 0.00031213, 0.5796699999999999, 0.41209331509738134}, {0.42100618271204354, 0.7, 0.18193357000000002, 0.0072455018245216685, 0.3765458534161963, 0.2084445539158497, 0.00010706059, 0.9384879265377657, 0.9430479823280116}, {0.09674009053919595, 0.31849999999999995, 0.13994890000000001, 0.008571428658409133, 0.7, 0.4180448855919369, 0.2664507153186618, 0.515469392193711, 0.6761203588006081}, {0.488515769597277, 0.40641843454115995, 0.107653, 0.01922275994260851, 0.844462555317776, 0.42650801894646384, 0.00018093239709999996, 0.6701102098518242, 0.9099999999999999}, {0.8255916506193981, 0.40641843454115995, 0.107653, 0.013455931959825957, 0.844462555317776, 0.22510109224181216, 0.00028403829999999996, 0.46907714689627694, 0.6761203588006081}, {0.34196103871809386, 0.132784546047367, 0.107653, 0.01922275994260851, 0.48999999999999994, 0.42650801894646384, 0.5521203209528786, 0.15294369999999993, 0.9658862268580115}, {0.42100618271204354, 0.65, 0.9662362119792216, 0.027461085632297877, 0.45471060670957164, 0.12120828043789884, 5.764800999999999e-05, 0.5796699999999999, 1}, {0.5736849062358256, 0.31849999999999995, 0.2612989847313744, 0.013455931959825957, 1, 0.015620761659111707, 0.00019882680999999997, 0.3608285745355977, 1}, {0.5779141554335788, 1, 0.08321795108585943, 0.010350716892173813, 0.8450304277220828, 0.08914744929999999, 0.000405769, 0.4256978025995064, 1}, {0.45499999999999996, 0.35, 0.2422354741113431, 0.02112391202484452, 0.8450304277220828, 0.6252051730108312, 0.00043682250156999996, 1, 0.4732842511604256}, {0.21483373423678898, 0.8450000000000001, 0.10546143252190002, 0.48999999999999994, 0.31829742469670014, 0.1459111877410948, 0.00025847485299999997, 0.34257058451260286, 0.7}, {0.6232157208577027, 0.45499999999999996, 0.5332632341456621, 0.014786738417391164, 0.3765458534161963, 0.48224371107540437, 0.0005274997, 0.22984780197917573, 0.6250331070969505}, {0.5779141554335788, 1, 0.04800247269999999, 0.0072455018245216685, 1, 0.11589168408999999, 0.4465379883817542, 0.4256978025995064, 1}, {0.5779141554335788, 1, 0.04800247269999999, 0.010350716892173813, 0.8450304277220828, 0.08914744929999999, 0.000405769, 0.4256978025995064, 0.7}, {0.6350705004764602, 0.5, 0.7432586245994012, 0.010350716892173813, 0.5202712699125625, 0.1459111877410948, 0.00016806999999999998, 0.4458999999999999, 1}, {0.3757813612286746, 0.20765701469221365, 0.04800247269999999, 0.9099999999999999, 0.3666625470543794, 0.12236608647264025, 0.555134866673812, 0.2476164298480412, 0.9099999999999999}, {0.4362510046003919, 0.35, 0.5332632341456621, 0.027461085632297877, 0.26358209739133737, 0.22510109224181216, 0.00030577575109899994, 0.34257058451260286, 1}, {0.5914999999999999, 0.45499999999999996, 0.4521328063216694, 0.005071851277165168, 0.8450304277220828, 0.1459111877410948, 0.00028403829999999996, 1, 0.6250331070969505}, {0.5699670500159709, 0.5, 0.7432586245994012, 0.010350716892173813, 0.5379226477374233, 0.22510109224181216, 0.00025847485299999997, 0.13287305512538392, 0.7}, {0.6350705004764602, 0.885565691820244, 0.16555954870000003, 0.7, 0.6495865810136738, 0.1459111877410948, 0.00016806999999999998, 0.4458999999999999, 0.48999999999999994}, {0.44454935033352205, 0.5, 0.107653, 0.010350716892173813, 1, 0.08914744929999999, 0.00013917876699999997, 0.2525800021749184, 1}, {0.44454935033352205, 0.5, 0.0753571, 0.010350716892173813, 0.8450304277220828, 0.08914744929999999, 0.00013917876699999997, 0.2525800021749184, 1}, {0.488515769597277, 0.5805977636302285, 0.107653, 0.027461085632297877, 0.6369004380736037, 0.554460424630403, 0.9043574131772127, 0.21849099999999994, 0.9658862268580115}, {0.22669563684494648, 0.35, 0.5332632341456621, 0.027461085632297877, 0.45471060670957164, 0.38812229724128205, 0.00021849099999999997, 0.14253325029259156, 0.9658862268580115}, {0.6350705004764602, 0.7, 0.7432586245994012, 0.010350716892173813, 0.5379226477374233, 0.10213783141876635, 0.00025847485299999997, 0.5796699999999999, 0.9658862268580115}, {0.488515769597277, 0.31849999999999995, 0.4230659320703461, 0.0072455018245216685, 0.7, 0.062403214509999985, 0.00024009999999999998, 0.5053396527511045, 1}, {0.3757813612286746, 0.5, 0.107653, 0.7, 0.19743367918312735, 0.08914744929999999, 0.00013917876699999997, 0.2476164298480412, 1}, {0.34196103871809386, 0.37678549999999994, 0.04800247269999999, 0.010350716892173813, 0.8450304277220828, 0.3215729889168745, 0.00030577575109899994, 0.2525800021749184, 0.7}, {0.488515769597277, 0.5891273193885694, 0.13994890000000001, 0.004615384662220302, 1, 0.18968454406342325, 0.00013917876699999997, 0.515469392193711, 0.3429999999999999}, {0.1107951529471101, 0.31849999999999995, 0.2612989847313744, 0.647623699490428, 0.5911237887224432, 0.22510109224181216, 0.2664507153186618, 0.4458999999999999, 0.7}, {0.6350705004764602, 0.9099999999999999, 0.7432586245994012, 0.010350716892173813, 0.5379226477374233, 0.07957750387456933, 0.00024009999999999998, 0.5796699999999999, 1}, {0.6168388863819738, 0.5, 0.11589168409000002, 0.013455931959825957, 0.3666625470543794, 0.1459111877410948, 0.00030577575109899994, 0.14601434629163068, 0.6369999999999999}, {0.42100618271204354, 0.37678549999999994, 0.9662362119792216, 0.006593406660314718, 0.45471060670957164, 0.3215729889168745, 0.0005274997, 0.3608285745355977, 0.9099999999999999}, {0.4001093022439462, 0.5, 0.5332632341456621, 0.010350716892173813, 0.9607234207297121, 0.08914744929999999, 8.235429999999999e-05, 0.3608285745355977, 1}}


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
