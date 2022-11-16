package main

import (
	"bufio"
	"fmt"
	"github.com/RyabovNick/databasecourse_2/golang/tasks/console_game/creature"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	time := "day"
	gameOver := false
	cycles := 0

	player := creature.New(10, 100, 20, 30)

	for !gameOver {

		if cycles == 5 {
			cycles = 0
			timeChange(&time, player)
		}

		showPlayerStats(player, time)
		actionDialog(scanner, player)
		cycles++

		checkGameOver(player, &gameOver)

	}
}

func checkGameOver(player *creature.Creature, gameOver *bool) {
	if !player.IsStateOk() {
		*gameOver = true
		fmt.Println("Game over! You lose")
	} else if player.GetRespect() >= 100 {
		*gameOver = true
		fmt.Println("Game over. You win!")
	}
}

func actionDialog(scanner *bufio.Scanner, player *creature.Creature) {
	fmt.Println("Actions: \n" +
		"---------------------")
	fmt.Println("1. Dig\n" +
		"2. Eat\n" +
		"3. Fight\n" +
		"4. Sleep\n" +
		"---------------------")
	fmt.Println("")
	fmt.Print("Enter action number: ")

	scanner.Scan()

	switch scanner.Text() {
	case "1":
		digDialog(scanner, player)
	case "2":
		eatDialog(scanner, player)
	case "3":
		fightDialog(scanner, player)
	case "4":
		player.Sleep()
		fmt.Println("")
		fmt.Println("You has slept. Hole length - 2, health + 20, respect - 2, weight - 5")
	default:
		fmt.Println("Unknown command")
	}
}

func digDialog(scanner *bufio.Scanner, player *creature.Creature) {
	fmt.Println("")
	fmt.Println("Dig intensively?")
	fmt.Print("Enter y or n: ")

	scanner.Scan()

	switch scanner.Text() {
	case "y":
		player.Dig(true)
	case "n":
		player.Dig(false)
	default:
		fmt.Println("Unknown command")
	}
}

func eatDialog(scanner *bufio.Scanner, player *creature.Creature) {
	fmt.Println("")
	fmt.Println("Eat: green(1) or withered(2) grass?")
	fmt.Print("Enter 1 or 2: ")

	scanner.Scan()
	ans := scanner.Text()

	switch ans {
	case "1":
		player.Eat(true)
	case "2":
		player.Eat(false)
	default:
		fmt.Println("Unknown command")
	}
}

func fightDialog(scanner *bufio.Scanner, player *creature.Creature) {
	fmt.Println("")
	fmt.Println("Choose your enemy.")
	fmt.Println("Fight with weak(1), average(2) or strong(3) creature?")
	fmt.Print("Enter: ")

	r := player.GetRespect()
	h := player.GetHealth()
	enemy := creature.New(0, 100, 0, 0)

	scanner.Scan()

	switch scanner.Text() {
	case "1":
		enemy.SetWeight(30)
	case "2":
		enemy.SetWeight(50)
	case "3":
		enemy.SetWeight(70)
	default:
		fmt.Println("Unknown command")
	}

	if player.FightWith(enemy) {
		fmt.Println("")
		fmt.Printf("You win. respect + %d, health - %d\n", player.GetRespect()-r, h-player.GetHealth())
	} else {
		fmt.Println("")
		fmt.Printf("You lose. health - %d\n", h-player.GetHealth())
	}
}

func showPlayerStats(player *creature.Creature, time string) {
	fmt.Println("")
	fmt.Println("Player stats: ")
	fmt.Printf("------------------------------------------------------------\n"+
		"| Hole length: %d | health: %d | respect: %d | weight: %d |\n"+
		"Time: %s\n"+
		"------------------------------------------------------------\n",
		player.GetHoleLength(), player.GetHealth(), player.GetRespect(), player.GetWeight(), time)
	fmt.Println()

}

func timeChange(time *string, player *creature.Creature) {
	if *time == "day" {
		fmt.Println("")
		fmt.Println("Night has come. Hole length - 2, health + 20, respect - 2, weight - 5")
		*time = "night"
		player.Sleep()
	} else {
		fmt.Println("")
		fmt.Println("Day has come")
		*time = "day"
	}
}
