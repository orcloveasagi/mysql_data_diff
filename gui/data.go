package gui

import (
	"db_diff/db"
	"fmt"
	"fyne.io/fyne/v2"
	"strconv"
	"strings"
)

type MenuData struct {
	Title, Intro string
	Onselect     func(ctx *NavCtx) fyne.CanvasObject
	Children     []string
	UID          string
}

type NavCtx struct {
	uid string
}

var (
	// Menus defines the metadata for each tutorial
	Menus = make(map[string]*MenuData)
)

const (
	CreateMenuUID = "create"
	ProfileUID    = "profile"
	RootUID       = ""
)

func init() {
	refreshMenuData()
}

func refreshMenuData() {
	data := db.LoadAll()
	allMenus := []*MenuData{
		{"Root", "", nil, nil, ""},
		{"Create", "", createCompareProfile, nil, CreateMenuUID},
		{"AllProfile", "", nil, nil, ProfileUID},
	}
	profileUIDs := make([]string, len(data))
	for i, value := range data {
		subMenuUID := profileUID(ProfileUID, strconv.FormatInt(value.Id, 10))
		profileUIDs[i] = subMenuUID
		allMenus = append(allMenus, &MenuData{value.Common.Name, "", selectCompareProfile, nil, subMenuUID})
	}

	for _, v := range allMenus {
		Menus[v.UID] = v
	}

	Menus[RootUID].Children = []string{CreateMenuUID, ProfileUID}
	Menus[ProfileUID].Children = profileUIDs
}

func profileUID(parentUID, childrenId string) string {
	return fmt.Sprintf("%s_%s", parentUID, childrenId)
}
func getDataByProfileUID(uid string) *db.CompareData {
	index := strings.Index(uid, "_")
	dataId, err := strconv.ParseInt(uid[index+1:], 10, 64)
	if err != nil {
		fyne.LogError("parse profileUID error!", err)
	}
	return db.Load(dataId)
}

func selectCompareProfile(ctx *NavCtx) fyne.CanvasObject {
	data := getDataByProfileUID(ctx.uid)
	return createForm(data)
}

func createCompareProfile(_ *NavCtx) fyne.CanvasObject {
	data := db.DefaultData()
	return createForm(&data)
}
