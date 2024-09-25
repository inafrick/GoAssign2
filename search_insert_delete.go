package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	sem chan struct{}
}

func NewSemaphore(initial int) *Semaphore {
	return &Semaphore{sem: make(chan struct{}, initial)}
}

func (s *Semaphore) Wait() {
	s.sem <- struct{}{}
}

func (s *Semaphore) Signal() {
	<-s.sem
}

type Lightswitch struct {
	count int
	mutex sync.Mutex
}

func (l *Lightswitch) Wait(sem *Semaphore) {
	l.mutex.Lock()
	l.count++
	if l.count == 1 {
		sem.Wait()
	}
	l.mutex.Unlock()
}

func (l *Lightswitch) Signal(sem *Semaphore) {
	l.mutex.Lock()
	l.count--
	if l.count == 0 {
		sem.Signal()
	}
	l.mutex.Unlock()
}

var (
	insertMutex  = NewSemaphore(1)
	noSearcher   = NewSemaphore(1)
	noInserter   = NewSemaphore(1)
	searchSwitch = &Lightswitch{}
	insertSwitch = &Lightswitch{}
	wg           sync.WaitGroup
)

func searcher(id int) {
	defer wg.Done()
	searchSwitch.Wait(noSearcher)
	fmt.Printf("Searcher %d is searching...\n", id)
	time.Sleep(1 * time.Second) // Simulando a seção crítica
	searchSwitch.Signal(noSearcher)
	fmt.Printf("Searcher %d finished searching.\n", id)
}

func inserter(id int) {
	defer wg.Done()
	insertSwitch.Wait(noInserter)
	insertMutex.Wait()
	fmt.Printf("Inserter %d is inserting...\n", id)
	time.Sleep(1 * time.Second) // Simulando a seção crítica
	insertMutex.Signal()
	insertSwitch.Signal(noInserter)
	fmt.Printf("Inserter %d finished inserting.\n", id)
}

func deleter(id int) {
	defer wg.Done()
	noSearcher.Wait()
	noInserter.Wait()
	fmt.Printf("Deleter %d is deleting...\n", id)
	time.Sleep(1 * time.Second) // Simulando a seção crítica
	noInserter.Signal()
	noSearcher.Signal()
	fmt.Printf("Deleter %d finished deleting.\n", id)
}

func main() {
	wg.Add(4)
	go searcher(1)
	go deleter(2)
	go inserter(1)
	go deleter(1)

	wg.Wait()
	fmt.Println("All operations completed successfully!")
}
