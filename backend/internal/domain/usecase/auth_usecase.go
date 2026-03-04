package usecase

import (
	"context"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

// AuthUseCase defines the contract for auth business logic.
// ทุก layer ที่ต้องการใช้งาน auth ให้ขึ้นกับ interface นี้ ไม่ใช่ implementation
type AuthUseCase interface {
	GetAuthURL(state string) string
	HandleCallback(ctx context.Context, code string) (*AuthResult, error)
}

type AuthResult struct {
	Token string
	User  *entity.User
}
