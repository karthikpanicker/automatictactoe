package main

import (
	"etsello/apps"
	"etsello/common"
	"strings"
	"time"

	"google.golang.org/api/tasks/v1"
)

type etsySynchronizer struct {
	dCache common.DataStore
}

func newEtsySynchronizer(cache common.DataStore) *etsySynchronizer {
	es := new(etsySynchronizer)
	es.dCache = cache
	return es
}

func (es *etsySynchronizer) processOrdersForUsers() {
	for {
		edm := apps.GetAppManager(apps.Etsy)
		userList := es.dCache.GetUserMap()
		for _, userDetails := range userList {
			//Need to fetch transactions for the user only if he has linked any of the apps
			if userDetails.TrelloDetails.IsLinked == false &&
				userDetails.GTasksDetails.IsLinked == false &&
				userDetails.TodoistDetails.IsLinked == false {
				continue
			}
			response, err := edm.GetAppData(&userDetails, apps.EtsyTransactionListRequest, nil)
			if err != nil {
				common.Error(err)
				continue
			}
			orderList := response.(*common.EtsyTransactionResponse)
			for i := len(orderList.Results) - 1; i >= 0; i-- {
				etsyTransaction := orderList.Results[i]
				if etsyTransaction.ShippedTime == 0 {
					reqParamsMap := make(map[string]interface{})
					reqParamsMap[apps.EtsyUserIDKey] = etsyTransaction.BuyerUserID
					response, _ := edm.GetAppData(&userDetails, apps.ProfileDetailsForUserRequest, reqParamsMap)
					buyerProfile := response.(*common.EtsyUserProfile)
					if userDetails.TrelloDetails.IsLinked {
						es.postTransactionToTrello(edm, etsyTransaction, &userDetails, buyerProfile)
					}
					if userDetails.GTasksDetails.IsLinked {
						es.postTransactionToGTasks(etsyTransaction, &userDetails, buyerProfile)
					}
					if userDetails.TodoistDetails.IsLinked {
						es.postTransactionToTodoist(etsyTransaction, &userDetails, buyerProfile)
					}
				}
			}
			es.dCache.SaveDetailsToCache(userDetails.UserID, userDetails)
		}
		time.Sleep(time.Minute * 2)
	}
}

func (es *etsySynchronizer) postTransactionToTrello(edm apps.AppDataManager, tranDetails common.EtsyTransactionDetails,
	info *common.UserInfo, buyerProfile *common.EtsyUserProfile) {
	// checking the criteria for selecting a transaction based on user preference
	// Second condition prevents a transaction from processed twice
	if tranDetails.PaidTime < info.TrelloDetails.FromDate ||
		tranDetails.ID <= info.TrelloDetails.LastProcessedTrasactionID {
		return
	}
	tdm := apps.GetAppManager(apps.Trello)
	card := common.TrelloCardDetails{
		Name:   tranDetails.Title,
		ListID: info.TrelloDetails.SelectedListID,
	}
	if contains(info.TrelloDetails.FieldsToUse, "listing_desc") {
		card.Descripton = es.formattedDescriptionWithMarkDown(tranDetails, buyerProfile, info)
	}
	if contains(info.TrelloDetails.FieldsToUse, "listing_link") {
		card.URL = tranDetails.EtsyURL
	}
	var resultCard common.TrelloCardDetailsResponse
	trelloReqParamsMap := make(map[string]interface{})
	if contains(info.TrelloDetails.FieldsToUse, "listing_image") {
		reqParamsMap := make(map[string]interface{})
		reqParamsMap[apps.EtsyTranDetailsKey] = tranDetails
		response, _ := edm.GetAppData(info, apps.EtsyImageDetailsRequest, reqParamsMap)
		trelloReqParamsMap[apps.TrelloShouldAttachImage] = true
		trelloReqParamsMap[apps.EtsyImageDetailsKey] = response
		tdm.AddItem(info, card, trelloReqParamsMap, &resultCard)
	} else {
		trelloReqParamsMap[apps.TrelloShouldAttachImage] = false
		tdm.AddItem(info, card, trelloReqParamsMap, &resultCard)
	}
	info.TrelloDetails.LastProcessedTrasactionID = tranDetails.ID
}

func (es *etsySynchronizer) postTransactionToGTasks(tranDetails common.EtsyTransactionDetails,
	info *common.UserInfo, buyerProfile *common.EtsyUserProfile) {
	if tranDetails.PaidTime < info.GTasksDetails.FromDate ||
		tranDetails.ID <= info.GTasksDetails.LastProcessedTrasactionID {
		return
	}
	todoItem := &tasks.Task{
		Title: tranDetails.Title,
		Notes: tranDetails.Description,
	}
	gtm := apps.GetAppManager(apps.Gtask)
	err := gtm.AddItem(info, todoItem, nil, nil)
	if err != nil {
		common.Error(err)
	}
	info.GTasksDetails.LastProcessedTrasactionID = tranDetails.ID
}

func (es *etsySynchronizer) postTransactionToTodoist(tranDetails common.EtsyTransactionDetails,
	info *common.UserInfo, buyerProfile *common.EtsyUserProfile) {
	if tranDetails.PaidTime < info.GTasksDetails.FromDate ||
		tranDetails.ID <= info.TodoistDetails.LastProcessedTrasactionID {
		return
	}
	tdm := apps.GetAppManager(apps.Todoist)
	task := common.TodoistTask{
		Content:   tranDetails.Title + " " + tranDetails.EtsyURL,
		ProjectID: info.TodoistDetails.SelectedProjectID,
	}
	err := tdm.AddItem(info, task, nil, task)
	if err != nil {
		common.Error(err)
	}
	info.TodoistDetails.LastProcessedTrasactionID = tranDetails.ID
}

func (es *etsySynchronizer) formattedDescriptionWithMarkDown(tranDetails common.EtsyTransactionDetails,
	buyerProfile *common.EtsyUserProfile, info *common.UserInfo) string {
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
