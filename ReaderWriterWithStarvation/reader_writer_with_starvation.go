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

var readLightSwitch Lightswitch
var roomEmpty = FPPDSemaforo.NewSemaphore(1)
var mutex = FPPDSemaforo.NewSemaphore(1) // protege o contador de leitores

func reader(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	mutex.Wait()
	readLightSwitch.Wait(roomEmpty)
	mutex.Signal()

	fmt.Printf("Leitor %d: na seção crítica\n", id)
	time.Sleep(3 * time.Second) // simulando trabalho na seção crítica

	mutex.Wait()
	readLightSwitch.Signal(roomEmpty)
	mutex.Signal()

	fmt.Printf("Leitor %d: saindo da seção crítica\n", id)
}

func writer(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	roomEmpty.Wait()
	fmt.Printf("Escritor %d: na seção crítica\n", id)
	time.Sleep(3 * time.Second) // simulando trabalho na seção crítica
	fmt.Printf("Escritor %d: saindo da seção crítica\n", id)
	roomEmpty.Signal()

}

func main() {
	var wg sync.WaitGroup

	// simulando threads leitores e escritores
	wg.Add(4) // número de threads

	go writer(1, &wg)
	go reader(1, &wg)
	go reader(2, &wg)
	go reader(3, &wg)

	// aguardando todas as goroutines terminarem
	wg.Wait()
}
