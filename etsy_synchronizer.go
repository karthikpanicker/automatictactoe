package main

import (
	"strings"
	"time"

	"google.golang.org/api/tasks/v1"
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
			//Need to fetch transactions for the user only if he has linked any of the apps
			if !userDetails.TrelloDetails.IsLinked || !userDetails.GTasksDetails.IsLinked {
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
					if userDetails.TrelloDetails.IsLinked {
						es.postTransactionToTrello(etsyTransaction, &userDetails, buyerProfile)
					}
					if userDetails.GTasksDetails.IsLinked {
						es.postTransactionToGTasks(etsyTransaction, &userDetails, buyerProfile)
					}
					lptID = etsyTransaction.ID
				}
			}
			userDetails.EtsyDetails.LastProcessedTrasactionID = lptID
			es.dCache.saveDetailsToCache(userDetails.UserID, userDetails)
		}
		time.Sleep(time.Minute * 30)
	}
}

func (es *etsySynchronizer) postTransactionToTrello(tranDetails etsyTransactionDetails,
	info *userInfo, buyerProfile *etsyUserProfile) {
	if tranDetails.PaidTime < info.TrelloDetails.FromDate {
		return
	}
	tdm := newTrelloDataManager()
	card := trelloCardDetails{
		Name:   tranDetails.Title,
		ListID: info.TrelloDetails.SelectedListID,
	}
	if contains(info.TrelloDetails.FieldsToUse, "listing_desc") {
		card.Descripton = es.formattedDescriptionWithMarkDown(tranDetails, buyerProfile, info)
	}
	if contains(info.TrelloDetails.FieldsToUse, "listing_link") {
		card.URL = tranDetails.EtsyURL
	}
	var resultCard trelloCardDetailsResponse
	tdm.addCard(info, card, &resultCard)
	if contains(info.TrelloDetails.FieldsToUse, "listing_image") {
		edm := newEtsyDataManager()
		imageDetails, _ := edm.getImageDetails(info, tranDetails)
		tdm.attachImage(info, &resultCard, imageDetails)
	}
}

func (es *etsySynchronizer) postTransactionToGTasks(tranDetails etsyTransactionDetails,
	info *userInfo, buyerProfile *etsyUserProfile) {
	if tranDetails.PaidTime < info.GTasksDetails.FromDate {
		return
	}
	todoItem := &tasks.Task{
		Title: tranDetails.Title,
		Notes: tranDetails.Description,
	}
	gtm := newGTasksDataManager()
	_, err := gtm.addToDoItem(info, todoItem, nil)
	if err != nil {
		Error(err)
	}
}

func (es *etsySynchronizer) formattedDescriptionWithMarkDown(tranDetails etsyTransactionDetails,
	buyerProfile *etsyUserProfile, info *userInfo) string {
	var sb strings.Builder
	sb.WriteString(tranDetails.Description)
	sb.WriteString("\n\n")
	if contains(info.TrelloDetails.FieldsToUse, "listing_buy_profile") && buyerProfile != nil {
		sb.WriteString("Buyer Details\n")
		sb.WriteString("--------------\n")
		sb.WriteString(buyerProfile.FirstName)
		sb.WriteString(" ")
		sb.WriteString(buyerProfile.LastName)
		sb.WriteString("\n")
		if buyerProfile.Region == "" || buyerProfile.City == "" {
			sb.WriteString(buyerProfile.Region)
			sb.WriteString(" ")
			sb.WriteString(buyerProfile.City)
		}
	}
	return sb.String()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
