package main

import (
	"flag"
	"log"
	"os"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	_ "github.com/bitrise-io/addons-ship-backend/db/migrate"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", "./migrate", "directory with migration files")
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("goose: %v", err)
	}
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	err = dataservices.InitializeConnection(dataservices.ConnectionParams{}, true)
	if err != nil {
		return
	}

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	command := args[0]
	if err := goose.Run(command, dataservices.GetDB().DB(), *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
