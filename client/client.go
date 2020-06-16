package client

func StartClient()  {
	SetupCui(func() {
		LoadConfiguration()
		InitManager()
		StartFileWatcher()
	})
}