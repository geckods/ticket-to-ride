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

	//GA stuff
	//toLog = new(bool)
	//consoleView = new(bool)
	//toUseVisualizer = new(bool)
	//toGenerateGraphs = false
	//optimizeBeaverParametersWithGeneticAlgorithm()


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

	//These are for BeaverPlayer OLD, without sampling code
	//{0.5, 0.5, 0.1, 0.18, 1, 0.1, 0.001, 0.01}
	//{0.5, 0.5, 0.9731249845226508, 0.44006774407186566, 0.710460877873115, 0.1, 0.001, 0.771424830503182}
	//{0.04578175501948286, 0.5, 0.1, 0.5839262959337884, 1, 0.0343, 0.0007, 0.013000000000000001}
	//{0.05951628152532772, 0.5, 0.3060318448167721, 0.21293999999999996, 1, 0.009989165164996577, 0.00091, 0.00637}
	//{0.11219467461355097, 0.5, 0.17838080992369257, 0.30255295346439276, 0.4610328365508673, 0.028145522253273242, 0.000753571, 1}
	//{0.20836153856802325 ,1 ,0.13770578824395227 ,0.9646955331718988 ,1 ,0.0021825426860652667 ,0.00107653 ,1}
	//{0.37888708948833616, 1, 0.7401533355717372, 0.06646185217789544, 0.5038589405152006, 1.5293134288138112e-05, 7.239272942909962e-06, 0.5387166810360579}
	//{0.058900132870430624, 1, 0.2895221806445169, 0.11490330624049794, 0.6407834254916002, 5.715384909614501e-06, 8.444607045954392e-07, 0.43772764436279393}
	//{0.10019273936205003 ,1 ,0.37929069417349554, 0.005788846196822391, 1 ,4.1375231334488105e-08, 1.9694065983164725e-07, 0.5796699999999999}

	//These are for BeaverPlayer NEW, with sampling code
	//This is equal to AardvardPlayer: {0.5, 0.5, 0.1, 0.18, 1, 0.1, 0.001, 0.01,0.1}
	//{0.48179801353843693, 0.7, 0.14621819593735766 ,0.030177017178349315, 0.7362757503240934, 0.8041609718224826, 0.031657202511261404, 0.4788249209903217, 1} : This one seems pretty strong, it accumulates cards early and aggressively builds midgame, and focuses on long roads when it fails to connect its destinations. I think it's a strat that can be countered by a bunch of bots that spam roads early. Strongest strat I have ATM
	//	{0.65, 0.5, 0.5613935345478703, 0.014786738417391164, 0.2636930125458408, 0.048999999999999995, 0.000637, 0.7, 0.7835672711986463} //this one also seems pretty high level
	// {0.44454935033352205 ,0.5 ,0.107653, 0.010350716892173813, 0.8450304277220828, 0.08914744929999999, 0.00013917876699999997, 0.2525800021749184, 1} //this one is even stronger, and it starts early as well, doesn't hoard as long
	e := Engine{}
	players := make([]Player, 0)
	player1 := BeaverPlayer{}
	player1.setScoringParameters([]float64{0.44454935033352205 ,0.5 ,0.107653, 0.010350716892173813, 0.8450304277220828, 0.08914744929999999, 0.00013917876699999997, 0.2525800021749184, 1})
	players = append(players, &player1)
	player2 := BeaverPlayer{}
	player2.setScoringParameters([]float64{0.65, 0.5, 0.5613935345478703, 0.014786738417391164, 0.2636930125458408, 0.048999999999999995, 0.000637, 0.7, 0.7835672711986463})
	players = append(players, &player2)
	player3 := AardvarkPlayer{}
	//player3.setScoringParameters([]float64{0.65, 0.5, 0.5613935345478703, 0.014786738417391164, 0.2636930125458408, 0.048999999999999995, 0.000637, 0.7, 0.7835672711986463})
	players = append(players, &player3)
	player4 := AardvarkPlayer{}
	//player4.setScoringParameters([]float64{0.65, 0.5, 0.5613935345478703, 0.014786738417391164, 0.2636930125458408, 0.048999999999999995, 0.000637, 0.7, 0.7835672711986463})
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
