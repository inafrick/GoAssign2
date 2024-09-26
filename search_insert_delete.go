package main

import (
	"GoAssign2/FPPDSemaforo"
	"fmt"
	"sync"
	"time"
)

// Definindo a estrutura do Lightswitch para controle de threads
type Lightswitch struct {
	counter int
	mutex   sync.Mutex
}

func (ls *Lightswitch) Wait(semaphore *FPPDSemaforo.Semaphore) {
	ls.mutex.Lock()
	ls.counter++
	if ls.counter == 1 {
		semaphore.Wait() // O primeiro a entrar fecha o semáforo
	}
	ls.mutex.Unlock()
}

func (ls *Lightswitch) Signal(semaphore *FPPDSemaforo.Semaphore) {
	ls.mutex.Lock()
	ls.counter--
	if ls.counter == 0 {
		semaphore.Signal() // O último a sair libera o semáforo
	}
	ls.mutex.Unlock()
}

// Declaração dos semáforos
var insertMutex = FPPDSemaforo.NewSemaphore(1) // Garante que apenas um inseridor esteja na seção crítica
var noSearcher = FPPDSemaforo.NewSemaphore(1)  // Garante que não há pesquisadores na seção crítica
var noInserter = FPPDSemaforo.NewSemaphore(1)  // Garante que não há inseridores na seção crítica
var searchSwitch Lightswitch                   // Lightswitch para pesquisadores
var insertSwitch Lightswitch                   // Lightswitch para inseridores

// Função para simular o comportamento de um pesquisador
func searcher(id int, wg *sync.WaitGroup) {
	defer wg.Done()               // Indica que a goroutine terminou
	searchSwitch.Wait(noSearcher) // Aguarda sua vez para acessar a seção crítica
	fmt.Printf("Pesquisador %d: na seção crítica\n", id)
	time.Sleep(1 * time.Second)     // Simulando trabalho na seção crítica
	searchSwitch.Signal(noSearcher) // Libera o controle quando terminar
	fmt.Printf("Pesquisador %d: saindo da seção crítica\n", id)
}

// Função para simular o comportamento de um inseridor
func inserter(id int, wg *sync.WaitGroup) {
	defer wg.Done()               // Indica que a goroutine terminou
	insertSwitch.Wait(noInserter) // Aguarda sua vez para acessar a seção crítica
	insertMutex.Wait()            // Apenas um inseridor pode entrar
	fmt.Printf("Inseridor %d: na seção crítica\n", id)
	time.Sleep(2 * time.Second)     // Simulando trabalho na seção crítica
	insertMutex.Signal()            // Libera para o próximo inseridor
	insertSwitch.Signal(noInserter) // Indica que terminou de inserir
	fmt.Printf("Inseridor %d: saindo da seção crítica\n", id)
}

// Função para simular o comportamento de um deletador
func deleter(id int, wg *sync.WaitGroup) {
	defer wg.Done()   // Indica que a goroutine terminou
	noSearcher.Wait() // Deletador bloqueia pesquisadores
	noInserter.Wait() // Deletador bloqueia inseridores
	fmt.Printf("Deletador %d: na seção crítica\n", id)
	time.Sleep(3 * time.Second) // Simulando trabalho na seção crítica
	noInserter.Signal()         // Libera inseridores
	noSearcher.Signal()         // Libera pesquisadores
	fmt.Printf("Deletador %d: saindo da seção crítica\n", id)
}

func main() {
	var wg sync.WaitGroup

	// Simulando threads pesquisadoras, inseridoras e deletadoras
	wg.Add(6) // Número de threads

	// Executando pesquisadores, inseridores e deletadores
	go searcher(1, &wg)
	go searcher(2, &wg)
	go inserter(1, &wg)
	go deleter(1, &wg)
	go inserter(2, &wg)
	go searcher(3, &wg)

	// Aguardando todas as goroutines terminarem
	wg.Wait()
}
