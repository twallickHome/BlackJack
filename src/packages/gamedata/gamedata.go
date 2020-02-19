package gamedata
//Contents of this package are 
//common between Client and Server

import(
	"packages/deck"	
)

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
	DealerHand  []Card
	DealerScore []int
	PlayerHand  []Card
	PlayerScore []int
	Message     string
}

//CalculateScore Calculate score
func CalculateScore(h []Card) []int {
	//Calculate dealer score first value

	//get number of aces in hand

	cardValue := 0
	countAces := 0
	var score []int = nil

	for i := 0; i < len(h); i++ {
		if deck.IsAce(h[i]) == true {
			countAces++
		} else {
			cardValue += h[i].Value
		}
	}

	fmt.Println("Card Value: ", cardValue)
	if countAces > 0 {

		score = calculateAces(cardValue, score, countAces)

	} else {
		score = append(score, cardValue)

	}

	score = removeDuplicates(score)

	fmt.Println("score: ", score)
