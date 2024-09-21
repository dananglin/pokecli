package commands

import "os"

func ExitProgram(_ []string) error {
	os.Exit(0)

	return nil
}
