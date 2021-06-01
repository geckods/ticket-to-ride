package main

//Source for map: https://images-eu.ssl-images-amazon.com/images/I/B19d%2BVcYwWS.png

const NUMDESTINATIONS = 36
const NUMCOLORCARDS = 12
const NUMRAINBOWCARDS = 14
const NUMSTARTINGTRAINS = 48
const NUMFACEUPTRAINCARDS = 5
const NUMGAMECOLORS = 9
const NUMINITIALTRAINCARDSDEALT = 4
const NUMINITIALDESTINATIONTICKETSOFFERED = 3
const NUMINITIALDESTINATIONTICKETSPICKED = 2
const NUMDESTINATIONTICKETSOFFERED = 3
const NUMDESTINATIONTICKETSPICKED = 1
const LONGESTPATHSCORE = 10

var routeLengthScores = []int{0, 1, 2, 4, 7, 10, 15, 21}

// ['Atlanta', 'Boston', 'Calgary', 'Charleston', 'Chicago', 'Dallas', 'Denver', 'Duluth', 'El Paso', 'Helena', 'Houston', 'Kansas City', 'Las Vegas', 'Little Rock', 'Los Angeles', 'Miami', 'Montreal', 'Nashville', 'New Orleans', 'New York', 'Oklahoma City', 'Omaha', 'Phoenix', 'Pittsburgh', 'Portland', 'Raleigh', 'Saint Louis', 'Salt Lake City', 'San Francisco', 'Santa Fe', 'Sault St. Marie', 'Seattle', 'Toronto', 'Vancouver', 'Washington', 'Winnipeg']
//TODO: build enum/array for destination

const (
	Atlanta Destination = iota
)

const (
	Red GameColor = iota
	Orange
	Yellow
	Green
	Blue
	Purple
	Black
	White
	Rainbow
	Other
)

var listOfGameColors = [...]GameColor{Red, Orange, Yellow, Green, Blue, Purple, Black, White, Rainbow}

//TODO: build DestinationTicket array
var listOfDestinationTickets = []DestinationTicket{{Atlanta, Atlanta, 1}}

//TODO: build Track array
var listOfTracks = []Track{{0, Atlanta, Atlanta, Red, 3}}
