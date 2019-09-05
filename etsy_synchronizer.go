package main

import (
	"strings"
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
					buyerProfile, _ := edm.getProfileDetails(etsyTransaction.BuyerUserID, &userDetails)
					es.postTransactionToTrello(etsyTransaction, &userDetails, buyerProfile)
					lptID = etsyTransaction.ID
				}
			}
			userDetails.EtsyDetails.LastProcessedTrasactionID = lptID
			es.dCache.saveDetailsToCache(userDetails.UserID, userDetails)
		}
		time.Sleep(time.Second * 30)
	}
}

func (es *etsySynchronizer) postTransactionToTrello(tranDetails etsyTransactionDetails,
	info *userInfo, buyerProfile *etsyUserProfile) {
	tdm := newTrelloDataManager()
	edm := newEtsyDataManager()
	imageDetails, _ := edm.getImageDetails(info, tranDetails)
	card := trelloCardDetails{
		Name:       tranDetails.Title,
		Descripton: es.formattedDescriptionWithMarkDown(tranDetails, buyerProfile),
		ListID:     info.TrelloDetails.SelectedListID,
		URL:        tranDetails.EtsyURL,
	}
	var resultCard trelloCardDetailsResponse
	tdm.addCard(info, card, &resultCard)
	tdm.attachImage(info, &resultCard, imageDetails)
}

func (es *etsySynchronizer) formattedDescriptionWithMarkDown(tranDetails etsyTransactionDetails,
	buyerProfile *etsyUserProfile) string {
	var sb strings.Builder
	sb.WriteString(tranDetails.Description)
	sb.WriteString("\n\n")
	sb.WriteString("Buyer Details\n")
	sb.WriteString("--------------\n")
	sb.WriteString(buyerProfile.FirstName)
	sb.WriteString(" ")
	sb.WriteString(buyerProfile.LastName)
	sb.WriteString("\n")
	sb.WriteString(buyerProfile.Region)
	sb.WriteString(", ")
	sb.WriteString(buyerProfile.City)
	return sb.String()
}
