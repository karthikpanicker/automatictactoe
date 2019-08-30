package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

func TestSaveUserInfo(t *testing.T) {
	gotenv.Apply(strings.NewReader("MONGO_URL=mongodb://localhost:27017"))
	mdc := newMongoDataCache()
	info := buildDummyUserInfo()
	info.EtsyDetails.UserShopDetails.ShopName = "Whatay Shop"
	info.EtsyDetails.UserShopDetails.ShopID = 54321
	mdc.saveDetailsToCache(info.UserID, *info)
}

func TestGetUserInfo(t *testing.T) {
	gotenv.Apply(strings.NewReader("MONGO_URL=mongodb://localhost:27017"))
	mdc := newMongoDataCache()
	info := buildDummyUserInfo()
	info.EtsyDetails.UserShopDetails.ShopName = "Whatay Shop"
	info.EtsyDetails.UserShopDetails.ShopID = 54321
	mdc.saveDetailsToCache(info.UserID, *info)
	savedInfo, err := mdc.getUserInfo(info.UserID)
	assert.Nil(t, err)
	assert.EqualValues(t, savedInfo, info)
}

func TestGetUserMap(t *testing.T) {
	gotenv.Apply(strings.NewReader("MONGO_URL=mongodb://localhost:27017"))
	mdc := newMongoDataCache()
	info := buildDummyUserInfo()
	info.EtsyDetails.UserShopDetails.ShopName = "Whatay Shop"
	info.EtsyDetails.UserShopDetails.ShopID = 54321
	mdc.saveDetailsToCache(info.UserID, *info)
	usersMap := mdc.getUserMap()
	assert.EqualValues(t, usersMap[info.UserID], *info)
}
