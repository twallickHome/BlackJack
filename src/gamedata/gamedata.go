package gamedata
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
	DealerHand  []Card
	DealerScore []int
	PlayerHand  []Card
	PlayerScore []int
	Message     string
}