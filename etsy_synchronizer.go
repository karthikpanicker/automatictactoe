package main

import "time"

type etsySynchronizer struct {
	userCache *userCache
}

func newEtsySynchronizer(cache *userCache) *etsySynchronizer {
	es := new(etsySynchronizer)
	es.userCache = cache
	return es
}

func (es *etsySynchronizer) processOrdersForUsers() {
	for {
		edm := newEtsyDataManager()
		userList := es.userCache.getUserMap()
		for _, value := range userList {
			orderList, err := edm.getTransactionList(value)
			if err != nil {
				Error(err)
				continue
			}
			Info(orderList)
		}
		time.Sleep(time.Second * 300)
	}
}
