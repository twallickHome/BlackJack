package deck

import(
	"packages/gamedata"
	
	"math/rand"	
	"time"
)

//IsAce Determine if card is an ace
func IsAce(c gamedata.Card) bool {
	out := false
	if c.Value == 0 {
		out = true
	}
	return out
}

//ShuffleDeck Shuffles the deck
func ShuffleDeck(deck  []gamedata.Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}

//DrawCard Draws a card from a shuffled deck
func DrawCard(numberOfCards int,
	deckOfCards []gamedata.Card,
	gd gamedata.GameData) []gamedata.Card {
	//fmt.Println("new hand ********")
	//.Println("Deck Before", deck[0:5])

	var c []gamedata.Card
	c = make([]gamedata.Card, numberOfCards, numberOfCards)

	for i := 0; i < numberOfCards; i++ {
		c[i] = deckOfCards[i]
	}

	//deckOfCards = append(deckOfCards[0:numberOfCards])
	//c := deckOfCards[0:numberOfCards]

	//fmt.Println("card drawn: ", c)

	//Remove card from deckOfCards
	deckOfCards = append(deckOfCards[:0], deckOfCards[numberOfCards:]...)

	//fmt.Println("deckOfCards After", deckOfCards[0:5])
	//fmt.Println("card drawn: ", c)

	//Update the size of the deckOfCards
	gd.DeckSize = len(deckOfCards)

	return c
}