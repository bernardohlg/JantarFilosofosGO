// Jantar dos Filósofos — V1: versão paralela ingênua, sujeita a deadlock.
// Cada filósofo pega primeiro o garfo à sua esquerda e só depois o da direita,
// sem ordenação de recursos e sem limite de concorrência. Isso satisfaz as quatro
// condições de Coffman ao mesmo tempo, entre elas a espera circular: se os cinco
// pegarem o garfo esquerdo antes de qualquer um pegar o direito, ninguém avança.
// Na prática o deadlock é raro, porque a janela entre os dois <- é de nanossegundos.
// Baseado no código de exemplo de César De Rose (PUCRS).

package main

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	mutex                            sync.Mutex
	maxNumOfPhilosophersEatingSimult = 0
	numOfPhilosophersEating          = 0
	philosophersCount                = 5
	rounds                           = 20
	forks                            []chan struct{}
)

type PhilosopherRecords struct {
	id              int
	roundsCounter   int
	timeWaitingFork []time.Duration
}

func addToSimultCounter() {
	mutex.Lock()
	numOfPhilosophersEating++
	if numOfPhilosophersEating > maxNumOfPhilosophersEatingSimult {
		maxNumOfPhilosophersEatingSimult = numOfPhilosophersEating
	}
	mutex.Unlock()
}

func decrementToSimultCounter() {
	mutex.Lock()
	numOfPhilosophersEating--
	mutex.Unlock()
}

func Philosopher(id int, recordsChan chan PhilosopherRecords) {
	record := PhilosopherRecords{id: id, roundsCounter: 0, timeWaitingFork: make([]time.Duration, rounds)}

	time_to_think := 200
	time_to_eat := 500
	espacos := strings.Repeat(" ", id*6)
	for i := 1; i <= rounds; i++ {
		// think
		fmt.Printf("%s T%d\n", espacos, i) // Thinking
		time.Sleep(time.Duration(time_to_think) * time.Millisecond)

		// eat
		startForkTime := time.Now()

		<-forks[id]
		<-forks[(id+1)%philosophersCount]
		record.timeWaitingFork[i-1] = time.Since(startForkTime)

		fmt.Printf("%s E%d\n", espacos, i) // Eating
		addToSimultCounter()
		time.Sleep(time.Duration(time_to_eat) * time.Millisecond)
		decrementToSimultCounter()

		forks[id] <- struct{}{}
		forks[(id+1)%philosophersCount] <- struct{}{}
		record.roundsCounter++

	}

	recordsChan <- record

}

func main() {
	forks = make([]chan struct{}, philosophersCount)
	recordsFromPhilosophers := make([]PhilosopherRecords, philosophersCount)
	recordsReceiver := make(chan PhilosopherRecords, 1)

	for i := 0; i < philosophersCount; i++ {
		forks[i] = make(chan struct{}, 1)
		forks[i] <- struct{}{}
	}

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]")
	fmt.Println()

	// run philosophers

	start := time.Now()

	for i := 0; i < philosophersCount; i++ {
		go Philosopher(i, recordsReceiver)
	}

	//salva dados dos filosofos
	for i := 0; i < philosophersCount; i++ {
		receivedRecord := <-recordsReceiver
		recordsFromPhilosophers[receivedRecord.id] = receivedRecord
	}

	elapsed := time.Since(start)
	close(recordsReceiver)

	var allTimesWaitedForFork []time.Duration
	var sumTimeToFork time.Duration

	for _, p := range recordsFromPhilosophers {
		for i := 0; i < p.roundsCounter; i++ {
			allTimesWaitedForFork = append(allTimesWaitedForFork, p.timeWaitingFork[i])
			sumTimeToFork += p.timeWaitingFork[i]
		}
	}

	slices.Sort(allTimesWaitedForFork)
	averageTime := sumTimeToFork / time.Duration(len(allTimesWaitedForFork))
	fmt.Printf("\nParalel Dinner took: %s\n\n", elapsed)
	fmt.Printf("average time to get fork: %s\n", averageTime)
	fmt.Printf("longest time waiting %s\n", allTimesWaitedForFork[len(allTimesWaitedForFork)-1])
	fmt.Printf("max num of philosophers eating at the same time: %d\n", maxNumOfPhilosophersEatingSimult)
	fmt.Println("\n--- Num of rounds per Philosopher ---")
	for _, p := range recordsFromPhilosophers {
		fmt.Printf("Philosopher %d: %d\n", p.id, p.roundsCounter)
	}
	fmt.Println("------------------------------------")
	fmt.Println("\n--- philosopher's waiting time ---")
	for _, p := range recordsFromPhilosophers {
		var tempoTotalEspera time.Duration
		for i := 0; i < p.roundsCounter; i++ {
			tempoTotalEspera += p.timeWaitingFork[i]
		}
		fmt.Printf("Filósofo %d esperou um total de: %s\n", p.id, tempoTotalEspera)
	}
	fmt.Println("------------------------------------")
}
