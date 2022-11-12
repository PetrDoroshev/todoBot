package main

import (
	"bufio"
	"fmt"
	"github.com/RyabovNick/databasecourse_2/golang/tasks/console_game/creature"
	"os"
)

var player creature.Creature
var scanner *bufio.Scanner
var time string

func main() {

	scanner = bufio.NewScanner(os.Stdin)

	time = "day"
	gameOver := false
	cycles := 0
	player = creature.Creature{HoleLength: 10,
		Health:  100,
		Respect: 20,
		Weight:  30}

	for !gameOver {

		if cycles == 5 {
			cycles = 0
			timeChange()
		}

		showPlayerStats()
		actionDialog()
		cycles++

		if player.Health <= 0 || player.HoleLength <= 0 || player.Respect <= 0 || player.Weight <= 0 {
			gameOver = true
			fmt.Println("Game over! You lose")
		} else if player.Respect >= 100 {
			gameOver = true
			fmt.Println("Game over. You win!")
		}

	}
}

func actionDialog() {
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
		digDialog()
	case "2":
		eatDialog()
	case "3":
		fightDialog()
	case "4":
		player.Sleep()
		fmt.Println("")
		fmt.Println("You has slept. Hole length - 2, Health + 20, Respect - 2, Weight - 5")
	default:
		fmt.Println("Unknown command")
	}
}

func digDialog() {
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

func eatDialog() {
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

func fightDialog() {
	fmt.Println("")
	fmt.Println("Choose your enemy.")
	fmt.Println("Fight with weak(1), average(2) or strong(3) creature?")
	fmt.Print("Enter: ")

	r := player.Respect
	h := player.Health
	enemy := creature.Creature{HoleLength: 0, Health: 100, Respect: 0, Weight: 0}

	scanner.Scan()

	switch scanner.Text() {
	case "1":
		enemy.Weight = 30
	case "2":
		enemy.Weight = 50
	case "3":
		enemy.Weight = 70
	default:
		fmt.Println("Unknown command")
	}

	if player.FightWith(enemy) {
		fmt.Println("")
		fmt.Printf("You win. Respect + %d, Health - %d\n", player.Respect-r, h-player.Health)
	} else {
		fmt.Println("")
		fmt.Printf("You lose. Health - %d\n", h-player.Health)
	}
}

func showPlayerStats() {
	fmt.Println("")
	fmt.Println("Player stats: ")
	fmt.Printf("------------------------------------------------------------\n"+
		"| Hole length: %d | Health: %d | Respect: %d | Weight: %d |\n"+
		"Time: %s\n"+
		"------------------------------------------------------------\n",
		player.HoleLength, player.Health, player.Respect, player.Weight, time)
	fmt.Println()

}

func timeChange() {
	if time == "day" {
		fmt.Println("")
		fmt.Println("Night has come. Hole length - 2, Health + 20, Respect - 2, Weight - 5")
		time = "night"
		player.Sleep()
	} else {
		fmt.Println("")
		fmt.Println("Day has come")
		time = "day"
	}
}
