// Jantar dos Filósofos — V0: versão sequencial, usada como linha de base.
// Não há concorrência: um filósofo pensa e come por vez, então o tempo total é
// previsível e vale rounds * philosophers * (time_to_think + time_to_eat).
// Baseado no código de exemplo de César De Rose (PUCRS).

package main

import (
	"fmt"
	"time"
)

func Philosopher(id, iteration int) {
	time_to_think := 200
	time_to_eat := 500

	// think
	for s := 1; s <= id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" T%d\n", iteration) // Thinking
	time.Sleep(time.Duration(time_to_think) * time.Millisecond)

	// eat

	for s := 1; s <= id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" E%d\n", iteration) // Eating
	time.Sleep(time.Duration(time_to_eat) * time.Millisecond)
}

func main() {
	philosophers := 5
	rounds := 20

	fmt.Println("\n[P1]  [P2]  [P3]  [P4]  [P5]")
	fmt.Println()

	// run philosophers

	start := time.Now()

	for i := 1; i <= rounds; i++ {
		for j := 0; j < philosophers; j++ {
			Philosopher(j, i)
		}
	}
	elapsed := time.Since(start)

	fmt.Printf("\nSequential Dinner took %s\n\n", elapsed)
}
