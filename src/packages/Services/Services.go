package services

import (
	"net/http"
	"packages/deck"
	"packages/gamedata"
)

//NewGame Starts a new blackjack game
func NewGame(w http.ResponseWriter,
	r *http.Request,
	deckOfCards []gamedata.Card,
	gd gamedata.GameData) (gamedata.GameData,[]gamedata.Card){

	//Check to see if deckOfCards is large enough to play with
	if len(deckOfCards) <= 10 {
		gd.Message = "Not enough cards in deckOfCards.\nPlease shuffle the deckOfCards"
		return gd,deckOfCards
	}

	//Check for game over condition
	if gd.GameOver == true {
		gd.GameOver = false
	} else {
		gd.Message = "Game is still in progress\nCan't start a new game"
		return gd,deckOfCards
	}

	
	//Initialize a new game

	//Draw 2 cards for player 1
	gd,player,deckOfCards := deck.DrawCard(2, deckOfCards, gd)
	gd.PlayerHand = append(player)
	
	//Draw 2 cards for the dealer
	//Second card is face up
	gd,dealer,deckOfCards := deck.DrawCard(2, deckOfCards, gd)
	dealer[1].FaceDown = false
	gd.DealerHand = append(dealer)
	
	// TEST:
	//need to re-engineer scoring multiple aces
	//calculateScore is off if there are 2 or more aces in hand
	//2 aces can be 2 , 12, 22
	//gd.DealerHand = []gamedata.gamedata.Card{{0, "A", "D", false}, {9, "9", "S", false}}
	//gd.DealerHand = []gamedata.gamedata.Card{{0, "A", "D", false}, {0, "A", "S", false}}
	//gd.DealerHand = []gamedata.gamedata.Card{{0, "A", "D", false}, {0, "A", "S", false}, {0, "A", "H", false}}
	//gd.DealerHand = []gamedata.gamedata.Card{{0, "A", "D", false}, {0, "A", "D", false}, {0, "A", "S", false}, {0, "A", "H", false}}

	//Get current score
	gd.DealerScore = gamedata.CalculateScore(gd.DealerHand)
	gd.PlayerScore = gamedata.CalculateScore(gd.PlayerHand)
	
	//test
	//gd.DealerScore = []int{11, 21}
	//gd.PlayerScore=[]int{19}

	playerHighestScore, dealerHighestScore := gamedata.GetHighestScore(gd)

	//check for a winner
	win, p := gamedata.CheckForWinner(playerHighestScore, dealerHighestScore)
	if win != 4 {
		gd.GameOver = true
		gd.Message = p
		deck.ShowAllCards(gd)

	} else {
		gd.Message = "Starting new game"
	}
	return gd,deckOfCards
}

//ShuffleDeck Create deck and shuffle it
func ShuffleDeck(w http.ResponseWriter,
	r *http.Request,
	deckOfCards []gamedata.Card,
	gd gamedata.GameData) (gamedata.GameData,[]gamedata.Card) {

	if gd.GameOver == false {
		gd.Message = "Game is still in progress\nCan't shuffle the deckOfCards"
		return gd,deckOfCards
	}

	//Reset Game data
	gd = gamedata.GameData{}
	gd.GameOver=true

	//Shuffle a new deckOfCards
	deckOfCards = nil
 	deckOfCards=deck.BuildDeck()

	deck.ShuffleDeck(deckOfCards)

	gd.DeckSize = len(deckOfCards)
	gd.Message = "Deck is shuffled"

	return gd,deckOfCards
}

//ShowHand Show the players hand
func ShowHand(w http.ResponseWriter,
	r *http.Request,
	gd gamedata.GameData) (gamedata.GameData){

	gd.Message = "showing cards"
	
	//show all cards
	for i := 0; i < len(gd.PlayerHand); i++ {
		gd.PlayerHand[i].FaceDown = false
	}

	return gd

}

//HitMe give the player a card
func HitMe(w http.ResponseWriter,
	r *http.Request,
	deckOfCards []gamedata.Card,
	gd gamedata.GameData) (gamedata.GameData,[]gamedata.Card){

	if gd.GameOver == true {
		gd.Message = "Game is over... Start a new game"
		return gd,deckOfCards
	}

	//Draw a card
	gd,playerDraw,deckOfCards :=deck.DrawCard(1,deckOfCards,gd)
	gd.PlayerHand = append(playerDraw, gd.PlayerHand...)
	
	//Get current score
	gd.DealerScore =gamedata.CalculateScore(gd.DealerHand)
	gd.PlayerScore =gamedata.CalculateScore(gd.PlayerHand)
	
	playerHighestScore, dealerHighestScore :=gamedata.GetHighestScore(gd)

	//check for a winner
	win, p :=gamedata.CheckForWinner(playerHighestScore, dealerHighestScore)
	if win != 4 {
		gd.GameOver = true
		gd.Message = p
		deck.ShowAllCards(gd)
	} else {
		gd.Message = "Here is your card"
	}
	return gd,deckOfCards
}

//Stay Player stays
func Stay(w http.ResponseWriter,
	r *http.Request,
	gd gamedata.GameData,
	deckOfCards []gamedata.Card) (gamedata.GameData,[]gamedata.Card){
	
		if gd.GameOver == true {
		gd.Message = "Game is over...Start a new game"
		
		return gd,deckOfCards
	}
	
	gd.Message = "You are staying"
	

	//test
	//gd.DealerScore = []int{19, 29}
	//gd.PlayerScore=[]int{19}

	playerHighestScore, dealerHighestScore :=gamedata.GetHighestScore(gd)

	//Dealers turn
	condition := true
	for ok := true; ok; ok = condition {
		//determine dealers highest score

		//a little AI
		//Check if dealer should draw or stay

		//compare player score with dealer
		//Whos ahead

		if dealerHighestScore >= playerHighestScore {
			gd.Message = "Dealer WINS!!!"
			gd.GameOver = true
			condition = false
			break
		} else if dealerHighestScore < playerHighestScore {
			gd,deckOfCards=deck.DealerDraw(gd,deckOfCards)

		} else if dealerHighestScore <= 15 {
			gd,deckOfCards=deck.DealerDraw(gd,deckOfCards)
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

	gd=deck.ShowAllCards(gd)
	return gd,deckOfCards
	}

