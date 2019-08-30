package main

import (
	"encoding/json"
	"time"
)

type etsySynchronizer struct {
	dCache dataCache
}

func newEtsySynchronizer(cache dataCache) *etsySynchronizer {
	es := new(etsySynchronizer)
	es.dCache = cache
	return es
}

func (es *etsySynchronizer) processOrdersForUsers() {
	for {
		edm := newEtsyDataManager()
		userList := es.dCache.getUserMap()
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
	tdm.addCard(info, card, nil)
}

func logUserInfo(info *userInfo) {
	userInfoBytes, _ := json.Marshal(info)
	Info(string(userInfoBytes))
}
