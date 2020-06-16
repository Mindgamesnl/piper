package client

func StartClient()  {
	go SetupCui()
	LoadConfiguration()
	StartFileWatcher()
}