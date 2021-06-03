package main

import (
	"fmt"
	strcstrconv "strconv"
)

func main() {
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
		fmt.Println("The winner was" + strcstrconv.Itoa(winner))
	}

}
