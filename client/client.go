package client

import (
	"os"
)

func StartClient()  {
	SetupCui(func() {
		LoadConfiguration()
		ConnectSocket(func() {
			InitManager()
			StartFileWatcher(func() {
				// other routines
				args := os.Args
				for i := range args {
					arg := args[i]

					if arg == "--upload-all" {
						SyncAllAndClose()
					}
				}
			})
		})
	})
}