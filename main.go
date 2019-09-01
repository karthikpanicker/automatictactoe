package main

import (
	"sync"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	dCache := newDataStore()
	defer dCache.disconnectCache()
	httpManager := newHTTPManager()
	go httpManager.startServer(dCache, "", 80)

	es := newEtsySynchronizer(dCache)
	go es.processOrdersForUsers()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
