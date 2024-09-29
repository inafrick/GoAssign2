package main

import (
	"GoAssign2/FPPDSemaforo"
	"fmt"
	"sync"
	"time"
)

type Lightswitch struct {
	counter int
	mutex   sync.Mutex
}

func (ls *Lightswitch) Wait(semaphore *FPPDSemaforo.Semaphore) {
	ls.mutex.Lock()
	ls.counter++
	if ls.counter == 1 {
		semaphore.Wait() // o primeiro a entrar fecha o semáforo
	}
	ls.mutex.Unlock()
}

func (ls *Lightswitch) Signal(semaphore *FPPDSemaforo.Semaphore) {
	ls.mutex.Lock()
	ls.counter--
	if ls.counter == 0 {
		semaphore.Signal() // o último a sair libera o semáforo
	}
	ls.mutex.Unlock()
}

var servings = 0
var mutex = FPPDSemaforo.NewSemaphore(1)
var emptyPot = FPPDSemaforo.NewSemaphore(1)
var fullPot = FPPDSemaforo.NewSemaphore(0)

func cook(M int) {
	for{
		emptyPot.Wait()
		putServingsInPot(M)
		fullPot.Signal()
	}
}

func savage(id int){
	for{
		mutex.Wait()
		if servings == 0{
			emptyPot.Signal()
			fullPot.Wait()
		}
		getServingFromPot(id)
		mutex.Signal()

		eat(id)
	}
}

func putServingsInPot(M int){
	servings = M
	fmt.Println("Cook is refilling the pot")
	time.Sleep(1 * time.Second)
}

func getServingFromPot(id int){
	servings--
	fmt.Println("Savage", id, "is getting a serving")
	time.Sleep(1 * time.Second)
}

func eat(id int){
	fmt.Println("Savage", id, "is eating")
	time.Sleep(1 * time.Second)
}

func main() {
	var wg sync.WaitGroup

	// simulando threads pesquisadoras, inseridoras e deletadoras
	wg.Add(7) // número de threads

	go cook(10) // cozinheiro - com número máximo do pote

	for i := 0; i < 6; i++ {
		go func(id int) {
			defer wg.Done()
			savage(id)
		}(i)
	}

	// aguardando todas as goroutines terminarem
	wg.Wait()
}
