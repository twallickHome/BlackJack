package deck

import(
	"packages/gamedata"
	
	"math/rand"	
	"time"
)

//BuildDeck Creates an unshuffled deck of cards
func BuildDeck() []gamedata.Card {
	deckOfCards:= []gamedata.Card{}
	
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

	return deckOfCards
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

//ShowAllCards Shows the hand of both players
func ShowAllCards(gd gamedata.GameData) {
	//show all cards
	for i := 0; i < len(gd.DealerHand); i++ {
		gd.DealerHand[i].FaceDown = false
	}

	for i := 0; i < len(gd.PlayerHand); i++ {
		gd.PlayerHand[i].FaceDown = false
	}
}
