package main

import (
	"encoding/json"
	"fmt"
	"os"

	"net/http"
	"packages/gamedata"
	"packages/services"
	
)

var deckOfCards []gamedata.Card //Holds an unshuffled deck of cards
var cardsInPlay []gamedata.Card //hold a record of cards removed from deck
var gd = gamedata.GameData{}    //Persistant game data
var port string                 //Holds port #

//Sends game data back to client
func postJSONResponse(w http.ResponseWriter, r *http.Request) {

	
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
	gd, deckOfCards = services.ShuffleDeck(w, r, deckOfCards, gd)
	postJSONResponse(w, r)

}

//Player draws a card
func hitHandler(w http.ResponseWriter, r *http.Request) {
	gd, deckOfCards = services.HitMe(w, r, deckOfCards, gd)
	postJSONResponse(w, r)
}

//player stays
func stayHandler(w http.ResponseWriter, r *http.Request) {

	gd, deckOfCards = services.Stay(w, r, gd, deckOfCards)
	postJSONResponse(w, r)
}

//Show player hand
func showHandler(w http.ResponseWriter, r *http.Request) {
	gd = services.ShowHand(w, r, gd)
	postJSONResponse(w, r)
}

//Starts a new game from an allready shuffled deck
func newHandler(w http.ResponseWriter, r *http.Request) {
	gd,deckOfCards = services.NewGame(w, r, deckOfCards, gd)
	postJSONResponse(w, r)
}

func main() {

	if len(os.Args) >= 2 {
		port = os.Args[1]
	} else {
		port = "8081"
	}

	gd.GameOver = true

	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/hit", hitHandler)
	http.HandleFunc("/stay", stayHandler)
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/shuffle", shuffleHandler)

	err := http.ListenAndServe(":"+port, nil)
	
	if err != nil {
		fmt.Println(err) // Ugly debug output

	}

}


