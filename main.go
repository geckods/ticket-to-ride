package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
)

var toLog *bool
var consoleView *bool


func main() {

	//command line flags
	toLog = flag.Bool("log", false, "Whether or not to log the operation of the engine. (default false)")
	consoleView = flag.Bool("console", true, "Whether to log the operation to console or to file. (default true, to console)")

	flag.Parse()

	var myConfig zap.Config
	if *toLog {
		if *consoleView {
			//use if you want to see events on the console
			myConfig = zap.NewDevelopmentConfig()
		} else {
			myConfig = zap.NewProductionConfig()
			myConfig.OutputPaths = append(myConfig.OutputPaths, "game.log")
			//	want to write to stderr
		}
	}

	logger, _ := myConfig.Build()
	defer logger.Sync() // flushes buffer, if any
	zap.ReplaceGlobals(logger)

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
	_ = constants
	e := Engine{}
	_ = e
	players := make([]Player, 0)
	player1 := BasicPlayer{}
	players = append(players, &player1)
	player2 := BasicPlayer{}
	players = append(players, &player2)
	player3 := BasicPlayer{}
	players = append(players, &player3)
	player4 := BasicPlayer{}
	players = append(players, &player4)

	winners := e.runGame(players, constants)

	for _,winner := range winners {
		fmt.Println("The winner was", winner)
	}

}
