package main

type userInfo struct {
	EmailID       string        `bson:"emailId"`
	UserID        int           `bson:"_id"`
	EtsyDetails   etsyDetails   `bson:"etsyDetails"`
	TrelloDetails trelloDetails `bson:"trelloDetails"`
	CurrentStep   int           `bson:"currentStep"`
}

type trelloDetails struct {
	TrelloAccessToken  string
	TrelloAccessSecret string
	TrelloBoards       []boardDetails
	SelectedBoardID    string `json:"boardId"`
	SelectedListID     string `json:"listId"`
}

type etsyDetails struct {
	EtsyAccessToken  string
	EtsyAccessSecret string
	UserShopDetails  shopDetails
	UserProfileURL   string
	UserName         string
}

type shopDetails struct {
	ShopID         int    `json:"shop_id"`
	ShopName       string `json:"shop_name"`
	Title          string `json:"title"`
	BannerImageURL string `json:"image_url_760x100"`
	ShopIconURL    string `json:"icon_url_fullxfull"`
	ShopFavorites  int    `json:"num_favorers"`
}

type boardDetails struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	URL         string `json:"url"`
	boardLists  []boardList
}

type boardList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type trelloCardDetails struct {
	Name       string `json:"name"`
	Descripton string `json:"desc"`
	ListID     string `json:"idList"`
	Labels     string `json:"idLabels"` //expected as comma separate strings
	URL        string `json:"urlSource"`
}

type trelloCardDetailsResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Descripton string   `json:"desc"`
	ListID     string   `json:"idList"`
	Labels     []string `json:"idLabels"` //array of labels
	URL        string   `json:"urlSource"`
}

func newDataCache() dataCache {
	dc := newMongoDataCache()
	return dc
}

type dataCache interface {
	saveDetailsToCache(userID int, userInfo userInfo)
	getUserInfo(userID int) (*userInfo, error)
	getUserMap() map[int]userInfo
	disconnectCache()
}
