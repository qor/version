package publish2_test

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/l10n"
	"github.com/qor/publish2"
)

var dbGlobal, dbCN, dbEN *gorm.DB

func init() {
	dbGlobal = DB
	dbCN = dbGlobal.Set("l10n:locale", "zh")
	dbEN = dbGlobal.Set("l10n:locale", "en")
}

type L10nProduct struct {
	gorm.Model
	Name string
	Body string

	publish2.Version
	publish2.Schedule
	publish2.Visible

	l10n.Locale
}

func TestL10nWithVersions(t *testing.T) {
	name := "product 1"
	now := time.Now()
	oneDayAgo := now.Add(-24 * time.Hour)
	oneDayLater := now.Add(24 * time.Hour)

	product := L10nProduct{Name: name}
	product.SetPublishReady(true)
	dbGlobal.Create(&product)
	dbEN.Create(&product)
	dbCN.Create(&product)

	var productCN, productEN Product

	dbEN.First(&productEN, "id = ?", product.ID)
	dbCN.First(&productCN, "id = ?", product.ID)

	productEN.SetVersionName("v1-en")
	productEN.Body = name + " - v1-EN"
	productEN.SetScheduledStartAt(&oneDayAgo)
	productEN.SetScheduledEndAt(&now)
	dbEN.Save(&productEN)

	productEN.SetVersionName("v2-en")
	productEN.Body = name + " - v2-EN"
	productEN.SetPublishReady(false)
	productEN.SetScheduledStartAt(&oneDayAgo)
	productEN.SetScheduledEndAt(&oneDayLater)
	dbEN.Save(&productEN)

	productEN.SetVersionName("v3-en")
	productEN.Body = name + " - v3-EN"
	productEN.SetPublishReady(true)
	productEN.SetScheduledStartAt(&now)
	productEN.SetScheduledEndAt(&oneDayLater)
	DB.Save(&productEN)

	productCN.SetVersionName("v1-cn")
	productCN.Body = name + " - v1-CN"
	productCN.SetScheduledStartAt(&oneDayAgo)
	productCN.SetScheduledEndAt(&now)
	dbCN.Save(&productCN)

	var count int
	dbEN.Model(&L10nProduct{}).Where("id = ?", product.ID).Count(&count)
	if count != 1 {
		t.Errorf("Should only find one valid product, but got %v", count)
	}

	dbEN.Set(publish2.VersionMode, publish2.VersionMultipleMode).Set(publish2.ScheduleCurrent, now.Add(-time.Hour)).Model(&L10nProduct{}).Where("id = ?", product.ID).Count(&count)
	if count != 2 {
		t.Errorf("Should only find two valid product when scheduled time, but got %v", count)
	}

	// DB.Set(publish2.VersionMode, publish2.VersionMultipleMode).Set(publish2.ScheduleCurrent, now.Add(time.Hour)).Model(&L10nProduct{}).Where("id = ?", product.ID).Count(&count)
	// if count != 2 {
	// 	t.Errorf("Should only find two valid product when scheduled time, but got %v", count)
	// }

	// DB.Set(publish2.VersionMode, publish2.VersionMultipleMode).Set(publish2.ScheduleStart, now.Add(time.Hour)).Set(publish2.ScheduleEnd, now.Add(24*time.Hour)).Model(&L10nProduct{}).Where("id = ?", product.ID).Count(&count)
	// if count != 2 {
	// 	t.Errorf("Should only find two valid product when scheduled time, but got %v", count)
	// }

	// DB.Set(publish2.VersionMode, publish2.VersionMultipleMode).Set(publish2.ScheduleStart, now.Add(-time.Hour)).Set(publish2.ScheduleEnd, now.Add(24*time.Hour)).Model(&L10nProduct{}).Where("id = ?", product.ID).Count(&count)
	// if count != 3 {
	// 	t.Errorf("Should only find two valid product when scheduled time, but got %v", count)
	// }
}