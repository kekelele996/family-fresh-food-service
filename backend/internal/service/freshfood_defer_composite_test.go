package service

import (
	"errors"
	"testing"

	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
)

func TestFoodDeferSwallow(t *testing.T) {
	db := newTestDB(t)
	foodRepo := repository.NewFoodItemRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	userRepo := repository.NewUserRepository(db)

	if _, err := foodRepo.FindByID(9999); !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("food FindByID err=%v, want ErrNotFound", err)
	}
	if _, err := memberRepo.FindByID(9999); !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("member FindByID err=%v, want ErrNotFound", err)
	}
	if _, err := notifyRepo.FindByID(9999); !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("notification FindByID err=%v, want ErrNotFound", err)
	}
	if _, err := userRepo.FindByID(9999); !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("user FindByID err=%v, want ErrNotFound", err)
	}
}
