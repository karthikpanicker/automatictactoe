package main

import (
	"etsello/common"
	"etsello/web"
	"sync"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	dCache := common.NewDataStore()
	defer dCache.DisconnectCache()
	httpManager := web.NewHTTPManager()
	go httpManager.StartServer(dCache, "", 80)

	es := newEtsySynchronizer(dCache)
	go es.processOrdersForUsers()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
