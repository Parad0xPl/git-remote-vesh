package utils

import (
	"log"
	"sync"
)

// cleaningStack is a stack which register cleaning functions and
// execute them on application closing. It gives additional ways t clean
// and close external programs in case of signal like SIGINT or SIGINT.
type cleaningStack struct {
	m sync.Mutex

	stack      []cleaningRecord
	stackIndex int

	nextId int

	isInitialized bool
	isCleaned     bool
}

type cleaningRecord struct {
	id   int
	f    CleaningFunction
	name string
}

// CleaningFunction should be a funtion that close external program or
// handle logic of cleaning/restoring broken resources.
type CleaningFunction func() error

var cleaningStateGlobal cleaningStack

// InitateCleaningState initiate global cleaning state.
func InitiateCleaningState() {
	cleaningStateGlobal.m.Lock()
	defer cleaningStateGlobal.m.Unlock()

	if cleaningStateGlobal.isInitialized {
		log.Println("WARNING: Cleaning stack - double initialization call!")
		return
	}
	cleaningStateGlobal.isInitialized = true
	cleaningStateGlobal.stack = make([]cleaningRecord, 10)
}

// AddCleaning push pointer of cleaning function to the stack
func AddCleaning(f CleaningFunction, name string) {
	cleaningStateGlobal.m.Lock()
	defer cleaningStateGlobal.m.Unlock()

	if name == "" {
		name = "Unnamed"
	}

	id := cleaningStateGlobal.nextId
	cleaningStateGlobal.nextId += 1
	cleaningStateGlobal.stack[id] = cleaningRecord{
		id:   id,
		f:    f,
		name: name,
	}
}

// Set name of given id task
func SetName(id int, name string) {
	cleaningStateGlobal.stack[id].name = name
}

// CleanStack execute each registered function in LIFO order
func CleanStack() {
	cleaningStateGlobal.m.Lock()
	defer cleaningStateGlobal.m.Unlock()

	if cleaningStateGlobal.isCleaned {
		log.Println("WARNING: Cleaning stack - double clean call!")
		return
	}

	for i := cleaningStateGlobal.nextId - 1; i >= 0; i-- {
		err := cleaningStateGlobal.stack[i].f()
		if err != nil {
			log.Printf("Error while cleaning %s: %e\n", cleaningStateGlobal.stack[i].name, err)
		}
	}

	cleaningStateGlobal.isCleaned = true
}
