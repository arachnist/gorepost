package bot

import (
	"github.com/arachnist/dyncfg"
	"sync"
)

var cfg *dyncfg.Dyncfg
var initLock sync.Mutex
var initList []func()

func addInit(f func()) {
	initLock.Lock()
	defer initLock.Unlock()

	initList = append(initList, f)
}

func Initialize(config *dyncfg.Dyncfg) {
	cfg = config

	for _, f := range initList {
		f()
	}
}
