package main

import (
	"encoding/json"
	"time"
)

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
			es.postTransactionToTrello(orderList.Results[0], &value)
		}
		time.Sleep(time.Second * 30)
	}
}

func (es *etsySynchronizer) postTransactionToTrello(tranDetails transactionDetails, info *userInfo) {
	if info.TrelloDetails.SelectedBoardID == "" {
		return
	}
	logUserInfo(info)
	tdm := newTrelloDataManager()
	card := trelloCardDetails{
		Name:       tranDetails.Title,
		Descripton: tranDetails.Description,
		ListID:     info.TrelloDetails.SelectedListID,
		//Labels:     info.EtsyDetails.UserShopDetails.ShopName,
		URL: tranDetails.EtsyURL,
	}
	tdm.addCard(info, card)
}

func logUserInfo(info *userInfo) {
	userInfoBytes, _ := json.Marshal(info)
	Info(string(userInfoBytes))
}
