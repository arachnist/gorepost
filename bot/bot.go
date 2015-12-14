package bot

import (
	"github.com/arachnist/dyncfg"
	"sync"
)

var cfg *dyncfg.Dyncfg
var cfgLock sync.Mutex

func Initialize(config *dyncfg.Dyncfg) {
	cfg = config
	cfgLock.Unlock()
}

func init() {
	cfgLock.Lock()
}
