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
var destinationNames = []string{"Atlanta", "Boston", "Calgary", "Charleston", "Chicago", "Dallas", "Denver", "Duluth", "El_Paso", "Helena", "Houston", "Kansas_City", "Las_Vegas", "Little_Rock", "Los_Angeles", "Miami", "Montreal", "Nashville", "New_Orleans", "New_York", "Oklahoma_City", "Omaha", "Phoenix", "Pittsburgh", "Portland", "Raleigh", "Saint_Louis", "Salt_Lake_City", "San_Francisco", "Santa_Fe", "Sault_St_Marie", "Seattle", "Toronto", "Vancouver", "Washington", "Winnipeg"}

const (
	Atlanta Destination = iota
	Boston
	Calgary
	Charleston
	Chicago
	Dallas
	Denver
	Duluth
	El_Paso
	Helena
	Houston
	Kansas_City
	Las_Vegas
	Little_Rock
	Los_Angeles
	Miami
	Montreal
	Nashville
	New_Orleans
	New_York
	Oklahoma_City
	Omaha
	Phoenix
	Pittsburgh
	Portland
	Raleigh
	Saint_Louis
	Salt_Lake_City
	San_Francisco
	Santa_Fe
	Sault_St_Marie
	Seattle
	Toronto
	Vancouver
	Washington
	Winnipeg
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

var stringColors = []string{"red", "orange", "yellow", "green", "blue", "purple", "black", "white", "rainbow", "grey"}

var listOfGameColors = [...]GameColor{Red, Orange, Yellow, Green, Blue, Purple, Black, White, Rainbow}

//TODO: build DestinationTicket array
var listOfDestinationTickets = []DestinationTicket{{Atlanta, Atlanta, 1}}

//TODO: build Track array
var listOfTracks = []Track{{0, Vancouver, Seattle, Other, 1}, {1, Seattle, Portland, Other, 1}, {2, Portland, San_Francisco, Green, 5}, {3, San_Francisco, Los_Angeles, Purple, 3}, {4, Los_Angeles, El_Paso, Black, 6},
	{5, Los_Angeles, Phoenix, Other, 3}, {6, Phoenix, El_Paso, Other, 3}, {7, Los_Angeles, Las_Vegas, White, 2}, {8, San_Francisco, Salt_Lake_City, Orange, 5}, {9, Portland, Salt_Lake_City, Blue, 6}, {10, Seattle, Helena, Yellow, 6}, {11, Seattle, Calgary, Other, 4}, {12, Vancouver, Calgary, Other, 3}, {13, Calgary, Winnipeg, White, 6},
	{14, Calgary, Helena, Other, 4}, {15, Winnipeg, Helena, Blue, 4}, {16, Helena, Salt_Lake_City, Purple, 3}, {17, Salt_Lake_City, Las_Vegas, Orange, 3}, {18, Salt_Lake_City, Denver, Red, 3}, {19, Phoenix, Denver, White, 5}, {20, Phoenix, Santa_Fe, Other, 3}, {21, El_Paso, Santa_Fe, Other, 2}, {22, Santa_Fe, Denver, Other, 2},
	{23, Helena, Duluth, Orange, 6}, {24, Winnipeg, Duluth, Black, 4}, {25, Winnipeg, Sault_St_Marie, Other, 6}, {26, Duluth, Sault_St_Marie, Other, 3}, {27, Sault_St_Marie, Montreal, Black, 5}, {28, Sault_St_Marie, Toronto, Other, 2}, {29, Montreal, Toronto, Other, 3}, {30, Duluth, Toronto, Purple, 6}, {31, Montreal, Boston, Other, 2}, {32, Boston, New_York, Red, 2},
	{33, Montreal, New_York, Blue, 3}, {34, New_York, Washington, Black, 2}, {35, Pittsburgh, New_York, Green, 2}, {36, Toronto, Pittsburgh, Other, 2}, {37, Chicago, Pittsburgh, Black, 3}, {38, Toronto, Chicago, White, 4}, {39, Duluth, Chicago, Red, 3},
	{40, Duluth, Omaha, Other, 2}, {41, Omaha, Chicago, Blue, 4}, {42, Omaha, Kansas_City, Other, 1}, {43, Denver, Omaha, Purple, 4}, {44, Helena, Denver, Green, 4}, {45, Denver, Kansas_City, Orange, 4}, {46, Denver, Oklahoma_City, Red, 4}, {47, Santa_Fe, Oklahoma_City, Blue, 3}, {48, El_Paso, Oklahoma_City, Yellow, 5}, {49, Oklahoma_City, Dallas, Other, 2}, {50, Dallas, Houston, Other, 1}, {51, El_Paso, Houston, Green, 6}, {52, El_Paso, Dallas, Red, 4},
	{53, Houston, New_Orleans, Other, 2}, {54, Oklahoma_City, Little_Rock, Other, 2}, {55, Little_Rock, Dallas, Other, 2}, {56, Kansas_City, Saint_Louis, Purple, 2}, {57, Chicago, Saint_Louis, Green, 2}, {58, Little_Rock, Saint_Louis, Other, 2}, {59, Saint_Louis, Nashville, Other, 2}, {60, Little_Rock, Nashville, White, 3}, {61, Little_Rock, New_Orleans, Green, 3}, {62, New_Orleans, Atlanta, Yellow, 4}, {63, Atlanta, Charleston, Other, 2}, {64, Charleston, Miami, Purple, 4}, {65, New_Orleans, Miami, Red, 6}, {66, Atlanta, Miami, Blue, 6},
	{67, Raleigh, Charleston, Other, 2}, {68, Nashville, Raleigh, Other, 2}, {69, Nashville, Raleigh, Black, 3}, {70, Raleigh, Washington, Other, 2}, {71, Washington, Pittsburgh, Other, 2}, {72, Pittsburgh, Raleigh, Other, 2}, {73, Pittsburgh, Saint_Louis, Yellow, 4}, {74, Pittsburgh, Saint_Louis, Green, 5}, {75, Helena, Omaha, Red, 5}, {76, Kansas_City, Oklahoma_City, Other, 2}, {77, Nashville, Atlanta, Other, 1}}
