package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"strconv"
	"gamedata"
	// "math/rand"
	// "time"
	//"html/template"
	//"regexp"
	//"errors"
)

// //Definition of a single card
// type card struct {
// 	Value    int
// 	Name     string
// 	Suit     string
// 	FaceDown bool
// }

// type gameData struct {
// 	DeckSize int
// 	DealerHand  []card
// 	DealerScore []int
// 	PlayerHand  []card
// 	PlayerScore []int
// 	Message     string
// }

var deck []gamedata.Card            //Holds a deck of cards
var gd = gamedata.GameData{}         //Persistant game data
var port string             //Holds port #
var clear map[string]func() //create a map for storing clear funcs
var condition bool          //Check for exit program condition
var screenWidth int         //Width of play area
var commandline bool=false 	//Use command line UI or built in UI

func displayScore(s []int) string {
	//fmt.Println("s = ",s)
	if (len(s)<=0){ return ""}

	out := strconv.Itoa(s[0])

	for i := 1; i < len(s); i++ {
		out += " or " + strconv.Itoa(s[i])
	}
	return out
}

//Support for call_clear()
func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//Clear Console screen
func callClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func center(s string, n int, fill string) string {
	div := n / 2

	return strings.Repeat(fill, div) + s + strings.Repeat(fill, div)
}

func displayHand(h []gamedata.Card) string {
	out := ""

	//fmt.Println("Card: ", h)

	for i := 0; i < len(h); i++ {
		if h[i].FaceDown == true {
			out += "[@] "
		} else {
			out += "[" + h[i].Name + " " + h[i].Suit + "] "
		}
	}
	return out
}

func updateDisplay() {
	callClear()

	//fmt.Println("Game Data: ", gd)

	//fmt.Println("update display gd=")
	//fmt.Println(gd)
	//Display game information
	
	fmt.Println(center("Play area", screenWidth, " "))
	fmt.Println("Deck Size: ",gd.DeckSize)
	fmt.Println("\nDealer Says... ")
	fmt.Println(gd.Message)
	//fmt.Println(center("", screenWidth, "*"))

	//Play area
	fmt.Println(center("", screenWidth, " "))
	fmt.Println(center(" Dealer's Hand ", screenWidth, "*"))
	fmt.Println(center(displayHand(gd.DealerHand), screenWidth, " "))
	//fmt.Println(center(" Score: "+displayScore(gd.DealerScore), screenWidth, " "))

	fmt.Println(center("", screenWidth, " "))
	fmt.Println(center(" Players's Hand ", screenWidth, "*"))
	fmt.Println(center(displayHand(gd.PlayerHand), screenWidth, " "))
	fmt.Println(center(" Score: "+displayScore(gd.PlayerScore), screenWidth, " "))

	if commandline==false{
	//Display menu
	fmt.Println("")
	fmt.Println(center("Main Menu", screenWidth, " "))
	fmt.Println("1 - New Game")
	fmt.Println("2 - Hit")
	fmt.Println("3 - Stay")
	fmt.Println("4 - Show")
	fmt.Println("5 - Shuffle")

	fmt.Println("")
	fmt.Println("X - Exit Game")
	}
}

func transmitCommand(inp string) {
	command := "error"

	fmt.Println(inp)
	switch inp {
	case "new":
		command = inp
		break
	case "shuffle":
		command = inp
		break
	case "hit":
		command = inp
		break

	case "stay":
		command = inp
		break
	case "show":
		command = inp
		break

	}

	//Set up output to server
	//in this case there is no need to send any data
	jOut, err := json.Marshal(nil)
	//Then we set up a new HTTP request for posting the JSON data to local port 8080.
	strRequest := "http://localhost:" + port + "/" + command
	fmt.Println(strRequest)
	req, err := http.NewRequest("POST", strRequest, bytes.NewBuffer(jOut))
	req.Header.Set("Content-Type", "application/json")
	//An HTTP client will send our HTTP request to the server and collect the response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	jsonResponse, err := ioutil.ReadAll(resp.Body)
	//Finally, we print the received response and close the response body.

	//fmt.Println("Response: ", jsonResponse)
	json.Unmarshal([]byte(jsonResponse), &gd)

	//fmt.Println("Game Data: ", gd)
	resp.Body.Close()
}

func mainGUI(inp rune) {

	out := "error"
	switch inp {
	case 'x', 'X':
		exitGame()
		break
	case '1':
		out = "new"
		break
	case '2':
		out = "hit"
		break
	case '3':
		out = "stay"
		break
	case '4':
		out = "show"
		break
	case '5':
		out = "shuffle"
		break

	}
	transmitCommand(out)
}

func exitGame() {
	condition = false
	fmt.Println("Thank you for playing.")
}

func getJSONResponse(body []byte) (*gamedata.GameData, error) {
	var s = new(gamedata.GameData)

	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func main() {
	//initialize main loop
	condition = true
	screenWidth = 50
	port = "8081"

	if len(os.Args) >= 2 {
		commandline=true
		transmitCommand(os.Args[1])
		updateDisplay()
	} else {
		commandline=false
		
		for ok := true; ok; ok = condition {
			updateDisplay()

			//Read the command line
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println(err)
			} else {
				mainGUI(char)
			}

		}

	}

}
