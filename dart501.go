package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
)

type Player struct {
	Name         string
	Score        int
	ActiveRounds int
	Average      float32
}

type PlayerList []Player

func main() {
	go ctrlC()

	fmt.Printf("\nby Jon Andreas Kornberg\n\n")
	fmt.Println("-------------------")
	fmt.Println("Welcome to Dart 501")
	fmt.Println("-------------------")
	fmt.Println("Press 'Ctrl + c' at any time to quit")
	fmt.Println()

	numPlayers := 0
	for numPlayers == 0 {
		fmt.Print("Number of players: ")
		var input string
		fmt.Scanln(&input)
		numPlayers, _ = strconv.Atoi(input)
	}

	players := make(PlayerList, numPlayers)

	fmt.Println()

	players = addPlayers(players, numPlayers)

	for {
		playersInGame := make(PlayerList, numPlayers)
		copy(playersInGame, players)
		game(playersInGame)
		fmt.Print("New Game? (y/n): ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "n" {
			break
		}
	}

}

func game(players PlayerList) {
	round := 0
	finished := false
	for !finished {
		round++
		finished = true
		for pid, player := range players {
			if players[pid].Score > 0 {
				finished = false
				fmt.Print(players)
				fmt.Printf("Round %v\n", round)

				fmt.Printf("%v:\n", player.Name)
				roundScore := oneRound(player)
				if roundScore > 0 && (player.Score-roundScore) >= 0 {
					players[pid].Score -= roundScore
					players[pid].ActiveRounds++
					players[pid].Average = updateAverage(players[pid], roundScore)
				}

			} else {
				continue
			}

			if players[pid].Score == 0 {
				fmt.Printf("\n%v WINS!!\nOther players continue? (y/n): ", players[pid].Name)
				var answer string
				fmt.Scanln(&answer)
				if answer == "n" {
					finished = true
					break
				}
			}
		}

	}

}

func oneRound(p Player) int {
	roundScore := 0
	throwNum := 0
	for throwNum < 3 {
		fmt.Printf("%v. ", throwNum+1)
		var text string
		fmt.Scanln(&text)
		if len(text) == 0 {
			continue
		}

		throwScore := 0
		switch text[0] {
		case 'd':
			throwScore, _ = strconv.Atoi(text[1:])
			throwScore *= 2
		case 't':
			throwScore, _ = strconv.Atoi(text[1:])
			throwScore *= 3
		default:
			throwScore, _ = strconv.Atoi(text)
		}
		if throwScore > 60 {
			fmt.Println("Juks! Straffeslurk!")
			throwScore = 0
		} else {
			throwNum++
		}
		roundScore += throwScore
	}
	fmt.Printf("Round Score: %v\n\n", roundScore)
	return roundScore
}

func addPlayers(players PlayerList, N int) []Player {

	for player := range players {

		players[player] = Player{Name: "unnamed", Score: 501, Average: 0.0, ActiveRounds: 0}
	}

	for player := range players {
		for {
			fmt.Printf("Player %v name: ", player+1)
			var name string
			fmt.Scanln(&name)

			if len(name) < 8 {
				players[player].Name = name
				break
			}
			fmt.Println("Name too long! Try again.")
		}
	}
	return players
}

func (p PlayerList) String() string {
	out := "\n----------------"
	for range p {
		out += "----------------"
	}

	out += "\n\t\t"
	for player := range p {
		out += fmt.Sprintf("%v\t\t", p[player].Name)
	}

	out += "\n\nAverage:\t"
	for player := range p {
		out += fmt.Sprintf("(%.1f)\t\t", p[player].Average)
	}

	out += "\n\nRest:\t\t"
	for player := range p {
		out += fmt.Sprintf("%v\t\t", p[player].Score)
	}

	out += "\n----------------"
	for range p {
		out += "----------------"
	}
	out += "\n\n"
	return out
}

func ctrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Print("\n\n-------------------\n")
	fmt.Print("Goodbye Dart 501\n")
	fmt.Print("-------------------\n")
	os.Exit(1)
}

func updateAverage(player Player, roundScore int) float32 {
	round := player.ActiveRounds
	newAverage := player.Average*(float32(round-1)) + float32(roundScore)
	newAverage /= float32(round)
	return newAverage
}
