package client

func StartClient()  {
	SetupCui(func() {
		LoadConfiguration()
		ConnectSocket(func() {
			InitManager()
			StartFileWatcher()
		})
	})
}