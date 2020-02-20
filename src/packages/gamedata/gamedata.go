package gamedata

import (
		
)

//Contents of this package are
//common between Client and Server

//Card : Define individual card
type Card struct {
	Value    int
	Name     string
	Suit     string
	FaceDown bool
}

//GameData : Persistant game data
type GameData struct {
	DeckSize int
	GameOver bool 
	DealerHand  []Card
	DealerScore []int
	PlayerHand  []Card
	PlayerScore []int
	Message     string
}

//********* exported functions **********

//CalculateScore Calculate score
func CalculateScore(h []Card) []int {
	//Calculate dealer score first value

	//get number of aces in hand

	cardValue := 0
	countAces := 0
	var score []int = nil

	for i := 0; i < len(h); i++ {
		if  isAce(h[i]) == true {
			countAces++
		} else {
			cardValue += h[i].Value
		}
	}

	
	if countAces > 0 {

		score = calculateAces(cardValue, score, countAces)

	} else {
		score = append(score, cardValue)

	}

	score = removeDuplicates(score)

	
	return score
}

//GetHighestScore With aces in hand, there can be more than one score. we only need the highest score 
//an ace produces without going over 21
func GetHighestScore(gd GameData) (playerHighestScore int, dealerHighestScore int ) {
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

//CheckForWinner Check for a win condition
func CheckForWinner(playerHighestScore int, dealerHighestScore int) (int, string) {
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

	
	return win, player
}






//********* Internal private functions **********

//isAce Determine if card is an ace
func isAce(c Card) bool {
	out := false
	if c.Value == 0 {
		out = true
	}
	return out
}

//calculateAces adds up all the aces in hand high and low values
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

//removeDuplicates after scores are calculated there may be duplicate scores in the score array. Don't want duplicates.
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

