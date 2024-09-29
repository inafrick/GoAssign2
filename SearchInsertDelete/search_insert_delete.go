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

var insertMutex = FPPDSemaforo.NewSemaphore(1) // garante que apenas um inseridor esteja na seção crítica
var noSearcher = FPPDSemaforo.NewSemaphore(1)  // garante que não há pesquisadores na seção crítica
var noInserter = FPPDSemaforo.NewSemaphore(1)  // garante que não há inseridores na seção crítica
var searchSwitch Lightswitch                   // lightswitch para pesquisadores
var insertSwitch Lightswitch                   // lightswitch para inseridores

// função para simular o comportamento de um pesquisador
func searcher(id int, wg *sync.WaitGroup) {
	defer wg.Done() // indica que a goroutine terminou

	searchSwitch.Wait(noSearcher) // aguarda sua vez para acessar a seção crítica

	fmt.Printf("Searcher %d: na seção crítica\n", id)
	time.Sleep(1 * time.Second) // simulando trabalho na seção crítica

	searchSwitch.Signal(noSearcher) // libera o controle quando terminar
	fmt.Printf("Searcher %d: saindo da seção crítica\n", id)
}

// função para simular o comportamento de um inseridor
func inserter(id int, wg *sync.WaitGroup) {
	defer wg.Done() // indica que a goroutine terminou

	insertSwitch.Wait(noInserter) // aguarda sua vez para acessar a seção crítica
	insertMutex.Wait()            // apenas um inseridor pode entrar

	fmt.Printf("Inserter %d: na seção crítica\n", id)
	time.Sleep(2 * time.Second) // simulando trabalho na seção crítica

	insertMutex.Signal()            // libera para o próximo inseridor
	insertSwitch.Signal(noInserter) // indica que terminou de inserir
	fmt.Printf("Inserter %d: saindo da seção crítica\n", id)
}

// função para simular o comportamento de um deletador
func deleter(id int, wg *sync.WaitGroup) {
	defer wg.Done() // indica que a goroutine terminou

	noSearcher.Wait() // deletador bloqueia pesquisadores
	noInserter.Wait() // deletador bloqueia inseridores

	fmt.Printf("Deleter %d: na seção crítica\n", id)
	time.Sleep(3 * time.Second) // simulando trabalho na seção crítica

	noInserter.Signal() // libera inseridores
	noSearcher.Signal() // libera pesquisadores
	fmt.Printf("Deleter %d: saindo da seção crítica\n", id)
}

func main() {
	var wg sync.WaitGroup

	// simulando threads pesquisadoras, inseridoras e deletadoras
	wg.Add(5) // número de threads

	go searcher(1, &wg)
	go searcher(2, &wg)
	go inserter(1, &wg)
	go inserter(2, &wg)
	go deleter(1, &wg)

	// aguardando todas as goroutines terminarem
	wg.Wait()
}
