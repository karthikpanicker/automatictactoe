package main

import (
	"time"
)

type etsySynchronizer struct {
	dCache dataStore
}

func newEtsySynchronizer(cache dataStore) *etsySynchronizer {
	es := new(etsySynchronizer)
	es.dCache = cache
	return es
}

func (es *etsySynchronizer) processOrdersForUsers() {
	for {
		edm := newEtsyDataManager()
		userList := es.dCache.getUserMap()
		for _, userDetails := range userList {
			if userDetails.TrelloDetails.SelectedBoardID == "" {
				continue
			}
			orderList, err := edm.getTransactionList(userDetails)
			if err != nil {
				Error(err)
				continue
			}
			lptID := userDetails.EtsyDetails.LastProcessedTrasactionID
			for i := len(orderList.Results) - 1; i >= 0; i-- {
				etsyTransaction := orderList.Results[i]
				if etsyTransaction.ID > lptID && etsyTransaction.ShippedTime == 0 {
					es.postTransactionToTrello(etsyTransaction, &userDetails)
					lptID = etsyTransaction.ID
				}
			}
			userDetails.EtsyDetails.LastProcessedTrasactionID = lptID
			es.dCache.saveDetailsToCache(userDetails.UserID, userDetails)
		}
		time.Sleep(time.Second * 30)
	}
}

func (es *etsySynchronizer) postTransactionToTrello(tranDetails etsyTransactionDetails, info *userInfo) {
	tdm := newTrelloDataManager()
	edm := newEtsyDataManager()
	imageDetails, _ := edm.getImageDetails(info, tranDetails)
	card := trelloCardDetails{
		Name:       tranDetails.Title,
		Descripton: tranDetails.Description,
		ListID:     info.TrelloDetails.SelectedListID,
		//Labels:     info.EtsyDetails.UserShopDetails.ShopName,
		URL: tranDetails.EtsyURL,
	}
	var resultCard trelloCardDetailsResponse
	tdm.addCard(info, card, &resultCard)
	tdm.attachImage(info, &resultCard, imageDetails)
}
