package main

const (
	userID          string = "userID"
	activeReqSecret string = "activeRequestSecret"
)

type userInfo struct {
	EmailID        string         `bson:"emailId"`
	UserID         int            `bson:"_id"`
	EtsyDetails    etsyDetails    `bson:"etsyDetails"`
	TrelloDetails  trelloDetails  `bson:"trelloDetails"`
	GTasksDetails  gTasksDetails  `bson:"gTaksksDetails"`
	TodoistDetails todoistDetails `bson:"todoistDetails"`
}

type trelloDetails struct {
	TrelloAccessToken  string
	TrelloAccessSecret string
	TrelloBoards       []boardDetails
	SelectedBoardID    string   `json:"boardId"`
	SelectedListID     string   `json:"listId"`
	FieldsToUse        []string `json:"fieldsToUse"`
	IsLinked           bool     `json:"isLinked"`
	TransactionFilter  int      `json:"transactionFilter"`
	FromDate           int
}

type etsyDetails struct {
	EtsyAccessToken           string
	EtsyAccessSecret          string
	UserShopDetails           shopDetails
	UserProfileURL            string
	UserName                  string
	LastProcessedTrasactionID int
}

type gTasksDetails struct {
	Token              string
	SelectedTaskListID string `json:"listId"`
	IsLinked           bool
	TransactionFilter  int `json:"transactionFilter"`
	FromDate           int
}

type todoistDetails struct {
	Token             string
	SelectedProjectID string `json:"projectId"`
	IsLinked          bool
	TransactionFilter int `json:"transactionFilter"`
	FromDate          int
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
	IsSelected  bool `json:"isSelected"`
}

type boardList struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsSelected bool   `json:"isSelected"`
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

type etsyProfileResponse struct {
	Count   int               `json:"count"`
	Results []etsyUserProfile `json:"results"`
}

type etsyShopResponse struct {
	Count   int           `json:"count"`
	Results []shopDetails `json:"results"`
}

type etsyTransactionResponse struct {
	Count   int                      `json:"count"`
	Results []etsyTransactionDetails `json:"results"`
}

type etsyImageResponse struct {
	Count   int                `json:"count"`
	Results []etsyImageDetails `json:"results"`
}

type etsyImageDetails struct {
	ID           int    `json:"listing_image_id"`
	ImageURL     string `json:"url_570xN"`
	FullImageURL string `json:"url_fullxfull"`
}

type etsyTransactionDetails struct {
	ID             int    `json:"transaction_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	BuyerUserID    int    `json:"buyer_user_id"`
	CreationTime   int    `json:"creation_tsz"`
	PaidTime       int    `json:"paid_tsz"`
	Price          string `json:"price"`
	Currency       string `json:"currency_code"`
	ShippingPrice  string `json:"shipping_cost"`
	ImageListingID int    `json:"image_listing_id"`
	EtsyURL        string `json:"url"`
	ShippedTime    int    `json:"shipped_tsz"`
	ListingID      int    `json:"listing_id"`
}

type etsyUserProfile struct {
	EmailID        string `json:"primary_email"`
	EtsyUserID     int    `json:"user_id"`
	UserProfileURL string `json:"image_url_75x75"`
	UserName       string `json:"login_name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Region         string `json:"region"`
	City           string `json:"city"`
}

type trelloImageAttachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type todoistProject struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type todoistTask struct {
	Content   string `json:"content"`
	ProjectID int    `json:"project_id"`
	ID        int    `json:"id"`
}

func newDataStore() dataStore {
	dc := newMongoDataCache()
	return dc
}

type dataStore interface {
	saveDetailsToCache(userID int, userInfo userInfo)
	getUserInfo(userID int) (*userInfo, error)
	getUserMap() map[int]userInfo
	disconnectCache()
}
