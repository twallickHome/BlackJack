package services

import (
	"fmt"
	
	"packages/gamedata"
	"packages/deck"
	"net/http"
	
)


//NewGame Starts a new blackjack game
func NewGame(w http.ResponseWriter, 
	r *http.Request,
	deckOfCards []gamedata.Card,
	gd gamedata.GameData,
	gameOver bool){

	//Check to see if deckOfCards is large enough to play with
	if len(deckOfCards) <= 10 {
		gd.Message = "Not enough cards in deckOfCards.\nPlease shuffle the deckOfCards"
		//postJSONResponse(w, r)
		return
	}

	//Check for game over condition
	if gameOver == true {
		gameOver = false
	} else {
		gd.Message = "Game is still in progress\nCan't start a new game"
		//postJSONResponse(w, r)
		return
	}

	fmt.Println("new game")
	//Initialize a new game

	//Draw 2 cards for player 1
	player :=deck.DrawCard(2,deckOfCards,gd)
	gd.PlayerHand = append(player)
	fmt.Println("Player drew: ", gd.PlayerHand)

	//Draw 2 cards for the dealer
	//Second card is fA up
	dealer := deck.DrawCard(2,deckOfCards,gd)
	dealer[1].FaceDown = false
	gd.DealerHand = append(dealer)
	fmt.Println("Dealer drew: ", gd.DealerHand)

	// TEST:
	//need to re-engineer scoring multiple aces
	//calculateScore is off if there are 2 or more aces in hand
	//2 aces can be 2 , 12, 22
	//gd.DealerHand = []gamedata.gamedata.Card{{0, "A", "D", false}, {9, "9", "S", false}}
	//gd.DealerHand = []gamedata.gamedata.Card{{0, "A", "D", false}, {0, "A", "S", false}}
	//gd.DealerHand = []gamedata.gamedata.Card{{0, "A", "D", false}, {0, "A", "S", false}, {0, "A", "H", false}}
	//gd.DealerHand = []gamedata.gamedata.Card{{0, "A", "D", false}, {0, "A", "D", false}, {0, "A", "S", false}, {0, "A", "H", false}}

	//Get current score
	gd.DealerScore =gamedata.CalculateScore(gd.DealerHand)
	fmt.Println("DealerScore: ", gd.DealerScore)

	gd.PlayerScore =gamedata.CalculateScore(gd.PlayerHand)
	fmt.Println("PlayerScore: ", gd.PlayerScore)

	//test
	//gd.DealerScore = []int{11, 21}
	//gd.PlayerScore=[]int{19}

	playerHighestScore, dealerHighestScore := getHighestScore()

	//check for a winner
	win, p := checkForWinner(playerHighestScore, dealerHighestScore)
	if win != 4 {
		gameOver = true
		gd.Message = p

	} else {
		gd.Message = "Starting new game"
	}

	//POST back to client
	//postJSONResponse(w, r)	
}