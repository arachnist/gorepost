package bot

import (
	"github.com/arachnist/gorepost/config"
	"sync"
)

var cfg *config.Config
var cfgLock sync.Mutex

func Initialize(config *config.Config) {
	cfg = config
	cfgLock.Unlock()
}

func init() {
	cfgLock.Lock()
}
