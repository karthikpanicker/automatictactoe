package main

import (
	"sync"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	userCache := newUserCache()
	httpManager := newHTTPManager()
	go httpManager.startServer(userCache, "localhost", 8900)

	es := newEtsySynchronizer(userCache)
	go es.processOrdersForUsers()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
