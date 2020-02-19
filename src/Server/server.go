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


//Dealer draws
func dealerDraw() {
	dealerDraw := deck.DrawCard(1,deckOfCards,gd)
	gd.DealerHand = append(dealerDraw, gd.PlayerHand...)

	//Get current score
	gd.DealerScore =gamedata.CalculateScore(gd.DealerHand)
	fmt.Println("DealerScore: ", gd.DealerScore)

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
func shuffleHandler(w http.ResponseWriter, r *http.Request){
	gd,deckOfCards= services.ShuffleDeck(w,r,deckOfCards,gd)
	postJSONResponse(w, r)

}
	

func hitHandler(w http.ResponseWriter, r *http.Request) {
	gd= services.HitMe(w,r,gd)
	postJSONResponse(w, r)
}

func stayHandler(w http.ResponseWriter, r *http.Request) {
	if gd.GameOver == true {
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

	playerHighestScore, dealerHighestScore :=gamedata.GetHighestScore(gd)

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
			gd.GameOver = true
			condition = false
			break
		} else if dealerHighestScore < playerHighestScore {
			dealerDraw()

		} else if dealerHighestScore <= 15 {
			dealerDraw()
		}

		playerHighestScore, dealerHighestScore =gamedata.GetHighestScore(gd)

		//check for a winner
		win, p :=gamedata.CheckForWinner(playerHighestScore, dealerHighestScore)
		if win != 4 {
			gd.GameOver = true
			gd.Message = p
			gd.GameOver = true
			condition = false

		}
		//Check for push only on dealers turn
		if playerHighestScore == dealerHighestScore {
			gd.GameOver = true
			gd.Message = "It's a push"
			gd.GameOver = true
			condition = false

		}
	} //for ok := true; ok; ok = condition

	deck.ShowAllCards(gd)

	postJSONResponse(w, r)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	gd=services.ShowHand(w,r,gd)
	postJSONResponse(w, r)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	gd = services.NewGame(w,r,deckOfCards,gd)
	postJSONResponse(w, r)
}

func main() {

	if len(os.Args) >= 2 {
		port = os.Args[1]
	} else {
		port = "8090"
	}

	gd.GameOver = true

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

