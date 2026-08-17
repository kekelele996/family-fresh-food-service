package service

import (
	"context"
	"testing"
	"time"

	"github.com/blueship581/cyfreshfood/internal/constants"
	"github.com/blueship581/cyfreshfood/internal/model"
	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
)

func TestFoodListOrderPage(t *testing.T) {
	db := newTestDB(t)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	foodRepo := repository.NewFoodItemRepository(db)
	consumeRepo := repository.NewConsumptionRecordRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	foodSvc := NewFoodItemService(foodRepo, consumeRepo, familySvc, util.NewFoodCalculator(), testLogger())
	memberSvc := NewFamilyMemberService(memberRepo, testLogger())
	statsSvc := NewStatsService(foodRepo, consumeRepo, notifyRepo, familySvc, memberSvc, util.NewFoodCalculator(), testLogger())
	ctx := context.Background()

	group, err := familySvc.Create(ctx, 1, "测试家庭")
	if err != nil {
		t.Fatalf("Create family: %v", err)
	}

	mkFood := func(name string, expiry time.Time, status string) *model.FoodItem {
		item := &model.FoodItem{FamilyID: group.ID, Name: name, Category: constants.FoodCategoryDairy, Quantity: 1, Unit: "盒", ShelfLifeDays: 5, StorageLocation: constants.StorageFridge, Status: status, CreatorID: 1, ExpiryDate: &expiry}
		if err := db.Create(item).Error; err != nil {
			t.Fatalf("seed food: %v", err)
		}
		return item
	}
	early := time.Now().AddDate(0, 0, 1)
	later := time.Now().AddDate(0, 0, 10)
	first := mkFood("早到期", early, constants.FreshnessExpiring)
	mkFood("晚到期", later, constants.FreshnessFresh)

	items, _, err := foodSvc.List(ctx, 1, group.ID, "", "", "", "", 1, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 2 || items[0].ID != first.ID {
		t.Fatalf("list order first=%v, want %d", itemIDs(items), first.ID)
	}

	pageItems, _, err := foodSvc.List(ctx, 1, group.ID, "", "", "", "", 1, 1)
	if err != nil {
		t.Fatalf("List page: %v", err)
	}
	if len(pageItems) != 1 || pageItems[0].ID != first.ID {
		t.Fatalf("page first=%v, want %d", itemIDs(pageItems), first.ID)
	}

	imported, created, err := foodSvc.ImportCSV(ctx, 1, group.ID, "牛奶,dairy\n鸡蛋,fresh\n")
	if err != nil {
		t.Fatalf("ImportCSV: %v", err)
	}
	if imported != 2 || len(created) != 2 {
		t.Fatalf("import count=%d created=%d, want 2 and 2", imported, len(created))
	}

	dash, err := statsSvc.Dashboard(ctx, 1, group.ID)
	if err != nil {
		t.Fatalf("Dashboard: %v", err)
	}
	if len(dash.ExpiringItems) != 1 {
		t.Fatalf("expiring items len=%d, want 1", len(dash.ExpiringItems))
	}
	if len(dash.ExpiredItems) != 0 {
		t.Fatalf("expired items len=%d, want 0", len(dash.ExpiredItems))
	}

	if got := util.JoinNames([]string{"张三", "李四", "王五"}); got != "张三、李四、王五" {
		t.Fatalf("JoinNames = %q, want 张三、李四、王五", got)
	}
}

func itemIDs(items []model.FoodItem) []uint {
	ids := make([]uint, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ID)
	}
	return ids
}
