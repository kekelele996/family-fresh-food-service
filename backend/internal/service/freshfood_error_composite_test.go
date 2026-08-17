package service

import (
	"context"
	"testing"

	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
)

func codeOfFresh(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		return 0
	}
	appErr, ok := err.(*util.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	return appErr.Code
}

func TestFoodErrorComposite(t *testing.T) {
	db := newTestDB(t)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	foodRepo := repository.NewFoodItemRepository(db)
	consumeRepo := repository.NewConsumptionRecordRepository(db)
	userRepo := repository.NewUserRepository(db)
	familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	foodSvc := NewFoodItemService(foodRepo, consumeRepo, familySvc, util.NewFoodCalculator(), testLogger())
	userSvc := NewUserService(userRepo, "test-secret", 24, testLogger())
	ctx := context.Background()

	if _, _, err := userSvc.Login(ctx, "13800000000", "pass123"); codeOfFresh(t, err) != 1001 {
		t.Fatalf("login missing user code=%d, want 1001", codeOfFresh(t, err))
	}
	if _, err := userSvc.GetByID(ctx, 9999); codeOfFresh(t, err) != 1003 {
		t.Fatalf("get missing user code=%d, want 1003", codeOfFresh(t, err))
	}

	group, err := familySvc.Create(ctx, 1, "测试家庭")
	if err != nil {
		t.Fatalf("Create family: %v", err)
	}
	if _, err := familySvc.GetByID(ctx, 9999); codeOfFresh(t, err) != 1003 {
		t.Fatalf("get missing group code=%d, want 1003", codeOfFresh(t, err))
	}
	if _, err := familySvc.InviteMember(ctx, group.ID, 1, 2); err != nil {
		t.Fatalf("InviteMember new user: %v", err)
	}
	if _, err := foodSvc.GetByID(ctx, 1, 9999); codeOfFresh(t, err) != 1003 {
		t.Fatalf("get missing food code=%d, want 1003", codeOfFresh(t, err))
	}
}
