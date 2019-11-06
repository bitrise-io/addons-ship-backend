package main

import (
	"fmt"

	"github.com/bitrise-io/addons-ship-backend/tasks/lib"
)

func main() {
	fmt.Println(lib.MigrateSelectedProvisioningProfileSlugToArray())
}
