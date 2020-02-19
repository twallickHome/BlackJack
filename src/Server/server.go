package main

import (
	
	"encoding/json"
	"fmt"
	"os"

	//"io/ioutil"
	//"log"
	"packages/gamedata"
	"packages/deck"
	"packages/services"
	"net/http"
	//"html/template"
	//"regexp"
	//"errors"
)

var deckOfCards []gamedata.Card //Holds an unshuffled deck of cards
var cardsInPlay []gamedata.Card //hold a record of cards removed from deck
var gd = gamedata.GameData{}    //Persistant game data
var port string                 //Holds port #
var gameOver bool = true


func removeDuplicates(elements []int) []int {

	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// ace calculation
func calculateAces(cardValue int, score []int, aceCount int) []int {

	score = append(score, 1*aceCount+cardValue)
	calc := 11 * aceCount
	if calc <= 21 {
		score = append(score, calc+cardValue)
	}
	calc = 11 + (aceCount - 1)
	if calc <= 21 {
		score = append(score, calc+cardValue)
	}

	return score
}

//Shows the hand of both players
func showAllCards() {
	//show all cards
	for i := 0; i < len(gd.DealerHand); i++ {
		gd.DealerHand[i].FaceDown = false
	}

	for i := 0; i < len(gd.PlayerHand); i++ {
		gd.PlayerHand[i].FaceDown = false
	}
}
func getHighestScore() (playerHighestScore int, dealerHighestScore int) {
	//Get player highest score
	playerHighestScore = gd.PlayerScore[0]
	for i := 1; i < len(gd.PlayerScore); i++ {
		if (gd.PlayerScore[i] > playerHighestScore) && (gd.PlayerScore[i] <= 21) {
			playerHighestScore = gd.PlayerScore[i]
		}
	}

	dealerHighestScore = gd.DealerScore[0]
	for i := 1; i < len(gd.DealerScore); i++ {
		if (gd.DealerScore[i] > dealerHighestScore) && (gd.DealerScore[i] <= 21) {
			dealerHighestScore = gd.DealerScore[i]
		}
	}

	return playerHighestScore, dealerHighestScore
}

//Dealer draws
func dealerDraw() {
	dealerDraw := drawCard(1)
	gd.DealerHand = append(dealerDraw, gd.PlayerHand...)

	//Get current score
	gd.DealerScore = calculateScore(gd.DealerHand)
	fmt.Println("DealerScore: ", gd.DealerScore)

}

//Check for a win condition
func checkForWinner(playerHighestScore int, dealerHighestScore int) (int, string) {
	//1=win
	//2=bust
	//3=Push
	//4=none
	win := 4
	player := ""

	//test
	//gd.DealerScore = []int{11, 21}
	//gd.PlayerScore=[]int{19}

	//Dealer instant win/loss
	if dealerHighestScore == 21 {
		win = 1
		player = "Dealer WINS!!!"

	} else if dealerHighestScore > 21 {
		win = 2
		player = "Dealer is BUST!!!"

	}

	//Player instant win/loss

	if playerHighestScore == 21 {
		win = 1
		player = "Player WINS!!!"

	} else if playerHighestScore > 21 {
		win = 2
		player = "Player is BUST!!!"

	}

	if win != 4 {
		showAllCards()
	}
	return win, player
}

//Sends game data back to client
func postJSONResponse(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Transmit Game Data: ", gd)
	//Now it’s time to prepare our response by setting up a weatherData structure.
	//We could try fetching the data from a weather service, but for the purpose of demonstrating JSON handling, let’s just use some mock-up data.

	//For encoding the Go struct as JSON, we use the Marshal function from encoding/json.

	outJSON, err := json.Marshal(gd)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}
	//We send a JSON response, so we need to set the Con10t-Type header accordingly.
	w.Header().Set("Content-Type", "application/json")
	//Sending the response is as easy as writing to the ResponseWriter object.
	//fmt.Println("Output JSON: ",outJSON)
	w.Write(outJSON)
}

//Shuffles a new deckOfCards of cards
func shuffleHandler(w http.ResponseWriter, r *http.Request) {

	if gameOver == false {
		gd.Message = "Game is still in progress\nCan't shuffle the deckOfCards"
		postJSONResponse(w, r)
		return
	}

	//Reset Game data
	gd = gamedata.GameData{}

	//Shuffle a new deckOfCards
	deckOfCards = nil
	buildDeck()

	
	deck.ShuffleDeck(deckOfCards)
	
	
	gd.DeckSize = len(deckOfCards)
	gd.Message = "Deck is shuffled"
	postJSONResponse(w, r)

}

func hitHandler(w http.ResponseWriter, r *http.Request) {

	if gameOver == true {
		gd.Message = "Game is over... Start a new game"
		postJSONResponse(w, r)
		return
	}

	fmt.Println("Hit")

	fmt.Println("Game Data: ", gd)

	//Draw a card
	playerDraw := drawCard(1)
	gd.PlayerHand = append(playerDraw, gd.PlayerHand...)
	fmt.Println("Player drew: ", gd.PlayerHand)

	//Get current score
	gd.DealerScore = calculateScore(gd.DealerHand)
	fmt.Println("DealerScore: ", gd.DealerScore)

	gd.PlayerScore = calculateScore(gd.PlayerHand)
	fmt.Println("PlayerScore: ", gd.PlayerScore)

	playerHighestScore, dealerHighestScore := getHighestScore()

	//check for a winner
	win, p := checkForWinner(playerHighestScore, dealerHighestScore)
	if win != 4 {
		gameOver = true
		gd.Message = p

	} else {
		gd.Message = "Here is your card"
	}

	postJSONResponse(w, r)

}

func stayHandler(w http.ResponseWriter, r *http.Request) {
	if gameOver == true {
		gd.Message = "Game is over...Start a new game"
		postJSONResponse(w, r)
		return
	}
	//fmt.Println("stay")
	gd.Message = "You are staying"
	//fmt.Println("Game Data: ", gd)

	//test
	//gd.DealerScore = []int{19, 29}
	//gd.PlayerScore=[]int{19}

	playerHighestScore, dealerHighestScore := getHighestScore()

	//Dealers turn
	condition := true
	for ok := true; ok; ok = condition {
		//determine dealers highest score

		//a little PAI
		//Check if dealer should draw or stay

		//compare player score with dealer
		//Whos ahead

		if dealerHighestScore > playerHighestScore {
			gd.Message = "Dealer WINS!!!"
			gameOver = true
			condition = false
			break
		} else if dealerHighestScore < playerHighestScore {
			dealerDraw()

		} else if dealerHighestScore <= 15 {
			dealerDraw()
		}

		playerHighestScore, dealerHighestScore = getHighestScore()

		//check for a winner
		win, p := checkForWinner(playerHighestScore, dealerHighestScore)
		if win != 4 {
			gameOver = true
			gd.Message = p
			gameOver = true
			condition = false

		}
		//Check for push only on dealers turn
		if playerHighestScore == dealerHighestScore {
			gameOver = true
			gd.Message = "It's a push"
			gameOver = true
			condition = false

		}
	} //for ok := true; ok; ok = condition

	showAllCards()

	postJSONResponse(w, r)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Show")
	gd.Message = "showing cards"
	//fmt.Println("Game Data: ", gd)

	//show all cards

	for i := 0; i < len(gd.PlayerHand); i++ {
		gd.PlayerHand[i].FaceDown = false
	}

	postJSONResponse(w, r)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	services.NewGame(w,r,deckOfCards,gd,gameOver)
	postJSONResponse(w, r)
}

func main() {

	if len(os.Args) >= 2 {
		port = os.Args[1]
	} else {
		port = "8090"
	}

	gameOver = true

	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/hit", hitHandler)
	http.HandleFunc("/stay", stayHandler)
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/shuffle", shuffleHandler)

	//err := http.ListenAndServe(":"+port, nil)
	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		fmt.Println(err) // Ugly debug output

	}

	

}

//********************************************************************
//    support functions
//********************************************************************

//Creates an unshuffled deck of cards
func buildDeck() {
	//Add s
	deckOfCards = append(deckOfCards, gamedata.Card{0, "A", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{2, "2", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{3, "3", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{4, "4", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{5, "5", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{6, "6", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{7, "7", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{8, "8", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{9, "9", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "10", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "J", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "Q", "S", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "K", "S", true})

	//Add H
	deckOfCards = append(deckOfCards, gamedata.Card{0, "A", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{2, "2", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{3, "3", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{4, "4", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{5, "5", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{6, "6", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{7, "7", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{8, "8", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{9, "9", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "10", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "J", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "Q", "H", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "K", "H", true})

	//Add C
	deckOfCards = append(deckOfCards, gamedata.Card{0, "A", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{2, "2", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{3, "3", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{4, "4", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{5, "5", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{6, "6", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{7, "7", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{8, "8", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{9, "9", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "10", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "J", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "Q", "C", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "K", "C", true})

	//Add D
	deckOfCards = append(deckOfCards, gamedata.Card{0, "A", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{2, "2", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{3, "3", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{4, "4", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{5, "5", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{6, "6", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{7, "7", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{8, "8", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{9, "9", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "10", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "J", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "Q", "D", true})
	deckOfCards = append(deckOfCards, gamedata.Card{10, "K", "D", true})

}
