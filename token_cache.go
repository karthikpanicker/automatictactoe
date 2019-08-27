package main

type userInfo struct {
	EmailID       string
	UserID        int
	EtsyDetails   etsyDetails
	TrelloDetails trelloDetails
	CurrentStep   int
}

type trelloDetails struct {
	trelloAccessToken  string
	trelloAccessSecret string
	TrelloBoards       []boardDetails
	SelectedBoardID    string `json:"boardId"`
	SelectedListID     string `json:"listId"`
}

type etsyDetails struct {
	etsyAccessToken  string
	etsyAccessSecret string
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

type userCache struct {
	userMap map[int]userInfo
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
	Labels     string `json:"idLabels"`
	URL        string `json:"urlSource"`
}

func newUserCache() *userCache {
	uc := new(userCache)
	uc.userMap = make(map[int]userInfo)
	return uc
}

func (uc *userCache) saveDetailsToCache(userID int, userInfo userInfo) {
	uc.userMap[userID] = userInfo
}

func (uc *userCache) getUserInfo(userID int) userInfo {
	return uc.userMap[userID]
}

func (uc *userCache) getUserMap() map[int]userInfo {
	return uc.userMap
}
