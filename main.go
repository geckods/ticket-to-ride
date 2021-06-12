package main

import (
	"flag"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var toLog *bool
var consoleView *bool
var toUseVisualizer *bool
var toGenerateGraphs bool

var server *socketio.Server //may be required globally

func main() {

	//seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// GA stuff
	toLog = new(bool)
	consoleView = new(bool)
	toUseVisualizer = new(bool)
	toGenerateGraphs = false
	optimizeBeaverParametersWithGeneticAlgorithm()


	//command line flags
	toLog = flag.Bool("log", false, "Whether or not to log the operation of the engine. (default false)")
	consoleView = flag.Bool("console", true, "Whether to log the operation to console or to file. (default true, to console)")
	toUseVisualizer = flag.Bool("visualize", false, "Whether or not to send data on a socket for visualization")
	flag.Parse()

	toGenerateGraphs = (*toLog && !(*consoleView))||(*toUseVisualizer)

	//logging related code
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

		logger, _ := myConfig.Build()
		defer logger.Sync() // flushes buffer, if any
		zap.ReplaceGlobals(logger)
	}

	//visualizer related code
	var wg *sync.WaitGroup
	numConnections := 0

	if *toUseVisualizer {
		server = socketio.NewServer(nil)

		server.OnConnect("/", func(s socketio.Conn) error {
			s.SetContext("")
			zap.S().Info("connected:", s.ID(), s.URL(), s.LocalAddr(), s.RemoteAddr(), s.Namespace())
			numConnections++
			return nil
		})

		server.OnError("/", func(s socketio.Conn, e error) {
			zap.S().Error("meet error:", e)
		})

		//
		server.OnDisconnect("/", func(s socketio.Conn, reason string) {
			zap.S().Debug("closed", reason)
		})

		go server.Serve()
		defer server.Close()

		http.Handle("/socket.io/", server)
		http.Handle("/", http.FileServer(http.Dir("./visualizer")))
		zap.S().Infof("Serving at localhost:8000...")

		wg = new(sync.WaitGroup)
		wg.Add(1)

		go func() {
			log.Fatal(http.ListenAndServe(":8000", nil))
			wg.Done()
		}()


	}

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

	//{0.046012701338462726, 0.8234722509035575, 0.8661244141006816, 0.5790593293832864, 0.9077195653391886, 0.8283225125248498, 0.18343344230131683, 0.19529733461220988}
	//{0.18292491645390843, 0.4283570818068078, 0.8969919575618727, 0.6826534880132438, 0.9789293555766876, 0.9222122589217269, 0.09083727535388708, 0.4931419977048804}

	e := Engine{}
	players := make([]Player, 0)
	player1 := BeaverPlayer{}
	player1.setScoringParameters([]float64{0.18292491645390843, 0.4283570818068078, 0.8969919575618727, 0.6826534880132438, 0.9789293555766876, 0.9222122589217269, 0.09083727535388708, 0.4931419977048804})
	players = append(players, &player1)
	player2 := AardvarkPlayer{}
	players = append(players, &player2)
	player3 := BasicPlayer{}
	players = append(players, &player3)
	player4 := BasicPlayer{}
	players = append(players, &player4)

	if *toUseVisualizer {
		zap.S().Infof("Waiting for visualizer to connect via socket.")
		for numConnections == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	winners := e.runGame(players, constants)

	for _,winner := range winners {
		fmt.Println("The winner was", winner)
	}

	if *toUseVisualizer {
		for _,winner := range winners {
			server.BroadcastToNamespace("/", "ENGINE_UPDATE", "The winner was " + strconv.Itoa(winner))
		}
		wg.Wait()
	}
}
