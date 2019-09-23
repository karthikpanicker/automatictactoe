package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

func TestSaveUserInfo(t *testing.T) {
	gotenv.Apply(strings.NewReader("MONGO_URL=mongodb://localhost:27017"))
	mdc := newMongoDataCache()
	defer mdc.DisconnectCache()
	info := buildDummyUserInfo()
	info.EtsyDetails.UserShopDetails.ShopName = "Whatay Shop"
	info.EtsyDetails.UserShopDetails.ShopID = 54321
	mdc.SaveDetailsToCache(info.UserID, *info)
}

func TestGetUserInfo(t *testing.T) {
	gotenv.Apply(strings.NewReader("MONGO_URL=mongodb://localhost:27017"))
	mdc := newMongoDataCache()
	defer mdc.DisconnectCache()
	info := buildDummyUserInfo()
	info.EtsyDetails.UserShopDetails.ShopName = "Whatay Shop"
	info.EtsyDetails.UserShopDetails.ShopID = 54321
	mdc.SaveDetailsToCache(info.UserID, *info)
	savedInfo, err := mdc.GetUserInfo(info.UserID)
	assert.Nil(t, err)
	assert.EqualValues(t, savedInfo, info)
}

func TestGetUserMap(t *testing.T) {
	gotenv.Apply(strings.NewReader("MONGO_URL=mongodb://localhost:27017"))
	mdc := newMongoDataCache()
	defer mdc.DisconnectCache()
	info := buildDummyUserInfo()
	info.EtsyDetails.UserShopDetails.ShopName = "Whatay Shop"
	info.EtsyDetails.UserShopDetails.ShopID = 54321
	mdc.SaveDetailsToCache(info.UserID, *info)
	usersMap := mdc.GetUserMap()
	assert.EqualValues(t, usersMap[info.UserID], *info)
}

func buildDummyUserInfo() *UserInfo {
	info := &UserInfo{
		UserID:  1234,
		EmailID: "karthik.panicker@gmail.com",
	}
	return info
}
