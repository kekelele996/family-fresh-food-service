package service

import (
	"context"
	"testing"

	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
)


func codeOfFreshNil(t *testing.T, err error) int {
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

func TestFoodNilComposite(t *testing.T) {
	db := newTestDB(t)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	notifySvc := NewNotificationService(notifyRepo, familySvc, testLogger())
	ctx := context.Background()

	if err := notifySvc.MarkRead(ctx, 1, 9999); codeOfFreshNil(t, err) != 1003 {
		t.Fatalf("mark missing notification code=%d, want 1003", codeOfFreshNil(t, err))
	}
}
