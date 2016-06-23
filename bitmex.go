package bitmex

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func init() {
	log.SetHandler(cli.New(os.Stderr))
}
