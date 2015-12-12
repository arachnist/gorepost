package bot

import (
	"github.com/arachnist/gorepost/config"
)

var cfg *config.Config

func Initialize(config *config.Config) {
	cfg = config
}
