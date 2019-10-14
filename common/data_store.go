package common

const (
	// UserID is the key used to store userid information in session
	UserID string = "userID"
	// ActiveReqSecret is the key used to store requets secret in session
	ActiveReqSecret string = "activeRequestSecret"
)

// UserInfo is the root document to store user details
type UserInfo struct {
	EmailID        string         `bson:"emailId"`
	UserID         int            `bson:"_id"`
	EtsyDetails    EtsyDetails    `bson:"etsyDetails"`
	TrelloDetails  TrelloDetails  `bson:"trelloDetails"`
	GTasksDetails  GTasksDetails  `bson:"gTaksksDetails"`
	TodoistDetails TodoistDetails `bson:"todoistDetails"`
}

// TrelloDetails is a struct used to store trello realted information.
type TrelloDetails struct {
	TrelloAccessToken                string
	TrelloAccessSecret               string
	TrelloBoards                     []BoardDetails
	SelectedBoardID                  string   `json:"boardId"`
	SelectedListID                   string   `json:"listId"`
	FieldsToUse                      []string `json:"fieldsToUse"`
	IsLinked                         bool     `json:"isLinked"`
	TransactionFilter                int      `json:"transactionFilter"`
	FromDate                         int
	LastProcessedTransactionPaidTime int
}

// EtsyDetails is a struct used to store etsy realted information.
type EtsyDetails struct {
	EtsyAccessToken  string
	EtsyAccessSecret string
	UserShopDetails  ShopDetails
	UserProfileURL   string
	UserName         string
}

// GTasksDetails is a struct to store google tasks details
type GTasksDetails struct {
	Token                            string
	SelectedTaskListID               string `json:"listId"`
	IsLinked                         bool
	TransactionFilter                int `json:"transactionFilter"`
	FromDate                         int
	LastProcessedTransactionPaidTime int
}

// GTasksListDetails is a struct to store gtask list details
type GTasksListDetails struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	IsSelected bool   `json:"isSelected"`
}

// TodoistDetails is a struct to store todoist details
type TodoistDetails struct {
	Token                            string
	SelectedProjectID                int `json:"projectId"`
	IsLinked                         bool
	TransactionFilter                int `json:"transactionFilter"`
	FromDate                         int
	LastProcessedTransactionPaidTime int
}

// ShopDetails is a struct to store etsy shop details
type ShopDetails struct {
	ShopID         int    `json:"shop_id"`
	ShopName       string `json:"shop_name"`
	Title          string `json:"title"`
	BannerImageURL string `json:"image_url_760x100"`
	ShopIconURL    string `json:"icon_url_fullxfull"`
	ShopFavorites  int    `json:"num_favorers"`
}

// BoardDetails is a struct to store details of a single board in trello
type BoardDetails struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	URL         string `json:"url"`
	boardLists  []BoardList
	IsSelected  bool `json:"isSelected"`
}

// BoardList is a struct to store the details of a list within a board.
type BoardList struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsSelected bool   `json:"isSelected"`
}

// TrelloCardDetails is a struct to store card details.
type TrelloCardDetails struct {
	Name       string `json:"name"`
	Descripton string `json:"desc"`
	ListID     string `json:"idList"`
	Labels     string `json:"idLabels"` //expected as comma separate strings
	URL        string `json:"urlSource"`
}

// TrelloCardDetailsResponse struct is used to store request response from trello for card
// details request
type TrelloCardDetailsResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Descripton string   `json:"desc"`
	ListID     string   `json:"idList"`
	Labels     []string `json:"idLabels"` //array of labels
	URL        string   `json:"urlSource"`
}

// EtsyProfileResponse is a struct to store the results of a request to
// get etsy profile information.
type EtsyProfileResponse struct {
	Count   int               `json:"count"`
	Results []EtsyUserProfile `json:"results"`
}

// EtsyShopResponse is struct to store response from the api to get etsy shop details
type EtsyShopResponse struct {
	Count   int           `json:"count"`
	Results []ShopDetails `json:"results"`
}

// EtsyTransactionResponse is struct to store response from the api to get etsy transaction list
type EtsyTransactionResponse struct {
	Count   int                      `json:"count"`
	Results []EtsyTransactionDetails `json:"results"`
}

// EtsyImageResponse is struct to store response from the api to get etsy image list
type EtsyImageResponse struct {
	Count   int                `json:"count"`
	Results []EtsyImageDetails `json:"results"`
}

// EtsyImageDetails is used to store the detailed information of an etsy image.
type EtsyImageDetails struct {
	ID           int    `json:"listing_image_id"`
	ImageURL     string `json:"url_570xN"`
	FullImageURL string `json:"url_fullxfull"`
}

// EtsyTransactionDetails is a struct to store a new transaction in  etsy.
type EtsyTransactionDetails struct {
	ID             int         `json:"transaction_id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	BuyerUserID    int         `json:"buyer_user_id"`
	CreationTime   int         `json:"creation_tsz"`
	PaidTime       int         `json:"paid_tsz"`
	Price          string      `json:"price"`
	Currency       string      `json:"currency_code"`
	ShippingPrice  string      `json:"shipping_cost"`
	ImageListingID int         `json:"image_listing_id"`
	EtsyURL        string      `json:"url"`
	ShippedTime    int         `json:"shipped_tsz"`
	ListingID      int         `json:"listing_id"`
	Variations     []Variation `json:"variations"`
}

// Variation is a struct to store variations from original Etsy listing as requested by the buyer
type Variation struct {
	Name  string `json:"formatted_name"`
	Value string `json:"formatted_value"`
}

// EtsyUserProfile is struct to store etsy user profile details.
type EtsyUserProfile struct {
	EmailID        string `json:"primary_email"`
	EtsyUserID     int    `json:"user_id"`
	UserProfileURL string `json:"image_url_75x75"`
	UserName       string `json:"login_name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Region         string `json:"region"`
	City           string `json:"city"`
}

// TrelloImageAttachment is a struct to store image details in trello
type TrelloImageAttachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// TodoistProject is a struct to store todoist project details.
type TodoistProject struct {
	Name       string `json:"name"`
	ID         int    `json:"id"`
	IsSelected bool   `json:"isSelected"`
}

// TodoistTask is struct to store tasks details in todoist.
type TodoistTask struct {
	Content   string `json:"content"`
	ProjectID int    `json:"project_id"`
	ID        int    `json:"id"`
}

// NewDataStore creates a new datastore with an encapsulated implemetation of data store.
func NewDataStore() DataStore {
	dc := newMongoDataCache()
	return dc
}

// DataStore is an interface to abstract implementation of datastore
type DataStore interface {
	SaveDetailsToCache(userID int, userInfo UserInfo)
	GetUserInfo(userID int) (*UserInfo, error)
	GetUserMap() map[int]UserInfo
	DisconnectCache()
}

// IsFieldSelected is used to check current selections
func (td *TrelloDetails) IsFieldSelected(fieldValue string) bool {
	return contains(td.FieldsToUse, fieldValue)
}

func contains(valueArray []string, value string) bool {
	for _, field := range valueArray {
		if field == value {
			return true
		}
	}
	return false
}
