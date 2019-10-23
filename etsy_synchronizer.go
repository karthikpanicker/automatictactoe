package main

import (
	"etsello/apps"
	"etsello/common"
	"sort"
	"strconv"
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
			//
			orderList := response.(*common.EtsyTransactionResponse)
			sort.Slice(orderList.Results[:], func(i, j int) bool {
				return orderList.Results[i].PaidTime < orderList.Results[j].PaidTime
			})
			for i := 0; i < len(orderList.Results); i++ {
				etsyTransaction := orderList.Results[i]
				if etsyTransaction.ShippedTime == 0 {
					reqParamsMap := make(map[string]interface{})
					reqParamsMap[apps.EtsyReceiptIdKey] = strconv.Itoa(etsyTransaction.ReceiptID)
					response, _ := edm.GetAppData(&userDetails, apps.EtsyReceiptDetailsRequest, reqParamsMap)
					receiptDetails := response.(*common.EtsyReceiptDetails)
					time.Sleep(time.Second * 1)
					var lastTransactionIndex = 0
					if i > 0 {
						lastTransactionIndex = i - 1
					}
					if userDetails.TrelloDetails.IsLinked {
						es.postTransactionToTrello(edm, etsyTransaction, &userDetails, receiptDetails, orderList.Results[lastTransactionIndex])
					}
					if userDetails.GTasksDetails.IsLinked {
						es.postTransactionToGTasks(etsyTransaction, &userDetails, receiptDetails)
					}
					if userDetails.TodoistDetails.IsLinked {
						es.postTransactionToTodoist(etsyTransaction, &userDetails, receiptDetails)
					}
				}
			}
			es.dCache.SaveDetailsToCache(userDetails.UserID, userDetails)
		}
		time.Sleep(time.Minute * 15)
	}
}

func (es *etsySynchronizer) postTransactionToTrello(edm apps.AppDataManager, tranDetails common.EtsyTransactionDetails,
	info *common.UserInfo, receiptDetails *common.EtsyReceiptDetails, lastTransaction common.EtsyTransactionDetails) {
	// checking the criteria for selecting a transaction based on user preference
	// if the transaction was paid before the user selected date or if the transaction was paid the paid time for
	// last processed transaction ignore the transaction
	if tranDetails.PaidTime < info.TrelloDetails.FromDate ||
		tranDetails.PaidTime < info.TrelloDetails.LastProcessedTransactionPaidTime {
		return
		// Condition to process orders if two orders have the same paid time (happens when there are multiple
		//transactions in an order) Paid time would be the same but transaction id would vary.
	} else if tranDetails.PaidTime == lastTransaction.PaidTime {
		if tranDetails.ID != lastTransaction.ID && tranDetails.ID != info.TrelloDetails.LastProcessedTransactionID{
		}else {
			return
		}
		// Ignore if the paid time and last processed transaction paid time are the same.
	} else if tranDetails.PaidTime == info.TrelloDetails.LastProcessedTransactionPaidTime {
		return
	}
	// If details are not configured skip the transaction
	if info.TrelloDetails.SelectedBoardID == "" || info.TrelloDetails.SelectedListID == "" {
		return
	}
	tdm := apps.GetAppManager(apps.Trello)
	card := common.TrelloCardDetails{
		Name:   tranDetails.Title,
		ListID: info.TrelloDetails.SelectedListID,
	}

	card.Descripton = es.formattedDescriptionWithMarkDown(tranDetails, receiptDetails, info)
	if contains(info.TrelloDetails.FieldsToUse, "listing_link") {
		card.URL = tranDetails.EtsyURL
	}
	var resultCard common.TrelloCardDetailsResponse
	trelloReqParamsMap := make(map[string]interface{})
	if contains(info.TrelloDetails.FieldsToUse, "listing_image") {
		reqParamsMap := make(map[string]interface{})
		reqParamsMap[apps.EtsyTranDetailsKey] = tranDetails
		response, err := edm.GetAppData(info, apps.EtsyImageDetailsRequest, reqParamsMap)
		if err != nil {
			common.Error("Error getting image details.", err)
		}
		trelloReqParamsMap[apps.TrelloShouldAttachImage] = true
		trelloReqParamsMap[apps.EtsyImageDetailsKey] = response
		tdm.AddItem(info, card, trelloReqParamsMap, &resultCard)
	} else {
		trelloReqParamsMap[apps.TrelloShouldAttachImage] = false
		tdm.AddItem(info, card, trelloReqParamsMap, &resultCard)
	}
	common.Info("Last processed transaction id: " + strconv.Itoa(tranDetails.ID))
	info.TrelloDetails.LastProcessedTransactionPaidTime = tranDetails.PaidTime
	info.TrelloDetails.LastProcessedTransactionID = tranDetails.ID
}

func (es *etsySynchronizer) postTransactionToGTasks(tranDetails common.EtsyTransactionDetails,
	info *common.UserInfo, receiptDetails *common.EtsyReceiptDetails) {
	if tranDetails.PaidTime < info.GTasksDetails.FromDate ||
		tranDetails.PaidTime <= info.GTasksDetails.LastProcessedTransactionPaidTime {
		return
	}
	if info.GTasksDetails.SelectedTaskListID == "" {
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
	info.GTasksDetails.LastProcessedTransactionPaidTime = tranDetails.PaidTime
}

func (es *etsySynchronizer) postTransactionToTodoist(tranDetails common.EtsyTransactionDetails,
	info *common.UserInfo, receiptDetails *common.EtsyReceiptDetails) {
	if tranDetails.PaidTime < info.GTasksDetails.FromDate ||
		tranDetails.PaidTime <= info.TodoistDetails.LastProcessedTransactionPaidTime {
		return
	}
	if info.TodoistDetails.SelectedProjectID == 0 {
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
	info.TodoistDetails.LastProcessedTransactionPaidTime = tranDetails.PaidTime
}

func (es *etsySynchronizer) formattedDescriptionWithMarkDown(tranDetails common.EtsyTransactionDetails,
	receiptDetails *common.EtsyReceiptDetails, info *common.UserInfo) string {
	var sb strings.Builder
	if contains(info.TrelloDetails.FieldsToUse, "listing_desc") {
		sb.WriteString(tranDetails.Description)
		sb.WriteString("\n\n")
	}
	if contains(info.TrelloDetails.FieldsToUse, "listing_buy_profile") && receiptDetails != nil {
		sb.WriteString("Buyer Details\n")
		sb.WriteString("--------------\n")
		sb.WriteString(receiptDetails.FormattedAddress)
		sb.WriteString("\n\n")
	}
	if contains(info.TrelloDetails.FieldsToUse, "listing_buyer_variations") && len(tranDetails.Variations) > 0 {
		sb.WriteString("Variations\n")
		sb.WriteString("--------------\n")
		for _, variation := range tranDetails.Variations {
			sb.WriteString(variation.Name)
			sb.WriteString(": ")
			sb.WriteString(variation.Value)
			sb.WriteString("\n")
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
