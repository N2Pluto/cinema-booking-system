package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
	jwtpkg "github.com/n2pluto/cinema-booking-system/pkg/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Ensure UseCase implements the domain interface at compile time.
var _ domainusecase.AuthUseCase = (*UseCase)(nil)

type UseCase struct {
	userRepo    repository.UserRepository
	oauthConfig *oauth2.Config
}

func NewUseCase(userRepo repository.UserRepository, clientID, clientSecret, callbackURL string) *UseCase {
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return &UseCase{userRepo: userRepo, oauthConfig: cfg}
}

func (uc *UseCase) GetAuthURL(state string) string {
	return uc.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (uc *UseCase) HandleCallback(ctx context.Context, code string) (*domainusecase.AuthResult, error) {
	token, err := uc.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	googleUser, err := fetchGoogleUser(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("fetch google user: %w", err)
	}

	user, err := uc.userRepo.Upsert(ctx, &entity.User{
		GoogleID:    googleUser.ID,
		Email:       googleUser.Email,
		DisplayName: googleUser.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}

	jwtToken, err := jwtpkg.Generate(user.ID.Hex(), user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("generate jwt: %w", err)
	}

	return &domainusecase.AuthResult{Token: jwtToken, User: user}, nil
}

type googleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func fetchGoogleUser(accessToken string) (*googleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info googleUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
