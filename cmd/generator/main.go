package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Character struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Lives int    `json:"lives"`
}

// Функция для генерации случайных значений
func generateRandomCharacter(id int) Character {
	names := []string{"zombie", "human", "alien", "robot"}
	colors := []string{"red", "green", "blue", "yellow"}
	return Character{
		ID:    id,
		Name:  names[rand.Intn(len(names))],
		Color: colors[rand.Intn(len(colors))],
		Lives: rand.Intn(10) + 1,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	characters := make([]Character, 100)
	for i := 0; i < 100; i++ {
		characters[i] = generateRandomCharacter(i + 1)
	}

	jsonData, err := json.MarshalIndent(characters, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
