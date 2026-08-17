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

func TestFoodRecommendStats(t *testing.T) {
	db := newTestDB(t)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	foodRepo := repository.NewFoodItemRepository(db)
	consumeRepo := repository.NewConsumptionRecordRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)
	familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	recipeSvc := NewRecipeService(recipeRepo, foodRepo, familySvc, util.NewFoodCalculator(), testLogger())
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
	mkFood("临期牛奶", time.Now().AddDate(0, 0, 1), constants.FreshnessExpiring)
	mkFood("过期酸奶", time.Now().AddDate(0, 0, -1), constants.FreshnessExpired)

	rec, err := recipeSvc.Recommend(ctx, 1, group.ID)
	if err != nil {
		t.Fatalf("Recommend: %v", err)
	}
	if len(rec.FoodItems) != 1 || rec.FoodItems[0].Name != "临期牛奶" {
		t.Fatalf("recommend foods=%v, want only 临期牛奶", rec.FoodItems)
	}
	if rec.Summary.ExpiringCount != 1 {
		t.Fatalf("recommend expiring count=%d, want 1", rec.Summary.ExpiringCount)
	}

	st, err := statsSvc.Statistics(ctx, 1, group.ID, "")
	if err != nil {
		t.Fatalf("Statistics: %v", err)
	}
	if st.WasteAmount != 15 || len(st.TopWasted) != 1 {
		t.Fatalf("waste amount=%v topWasted=%d, want 15 and 1", st.WasteAmount, len(st.TopWasted))
	}

	if got := util.JoinNames([]string{"张三", "李四", "王五"}); got != "张三、李四、王五" {
		t.Fatalf("JoinNames = %q, want 张三、李四、王五", got)
	}
}
