package utils

import (
	"log"
	"sync"
)

// cleaningStack is a stack that registers cleaning functions and
// executes them on application closure. It provides a mechanism to clean
// and close external programs in case of signals like SIGINT or SIGTERM.
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

// CleaningFunction should be a function that closes external programs or
// handles the logic of cleaning/restoring broken resources.
type CleaningFunction func() error

var cleaningStateGlobal cleaningStack

// InitiateCleaningState initializes the global cleaning state.
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

// AddCleaning pushes a pointer to a cleaning function onto the stack.
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

// SetName sets the name of a given task by its ID.
func SetName(id int, name string) {
	cleaningStateGlobal.stack[id].name = name
}

// CleanStack executes each registered function in LIFO order.
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
			log.Printf("Error while cleaning %s: %v\n", cleaningStateGlobal.stack[i].name, err)
		}
	}

	cleaningStateGlobal.isCleaned = true
}
