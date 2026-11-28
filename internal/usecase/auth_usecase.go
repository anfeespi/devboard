package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anfeespi/devboard/internal/auth"
	"github.com/anfeespi/devboard/internal/domain"
)

// authUserRepository defines what the use case needs
type authUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

// AuthTokens tokens
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// AuthUseCase use case of the authentication
type AuthUseCase struct {
	userRepo             authUserRepository
	jwtManager           *auth.JWTManager
	tokenStore           *auth.TokenStore
	refreshTokenDuration time.Duration
}

// NewAuthUseCase constructor
func NewAuthUseCase(userRepo authUserRepository, jwtManager *auth.JWTManager, tokenStore *auth.TokenStore, refreshTokenDuration time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:             userRepo,
		jwtManager:           jwtManager,
		tokenStore:           tokenStore,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// Register register a new user creating tokens
func (uc *AuthUseCase) Register(ctx context.Context, email, name, password string) (*AuthTokens, error) {
	passwordHash, err := auth.HashPassword(password)

	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := domain.NewUser(email, name, passwordHash)
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	tokens, err := uc.generateTokenPair(ctx, user.ID, user.Email, "member")
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// generateTokenPair generate refresh and access token
func (uc *AuthUseCase) generateTokenPair(ctx context.Context, userID, email, role string) (*AuthTokens, error) {
	accessToken, err := uc.jwtManager.GenerateAccessToken(userID, email, role)

	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(userID, email, role)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	refreshClaims, err := uc.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("validating claims on the refresh token: %w", err)
	}

	if err := uc.tokenStore.StoreRefreshToken(
		ctx,
		refreshClaims.ID,
		userID,
		uc.refreshTokenDuration,
	); err != nil {
		return nil, fmt.Errorf("storing refresh token: %w", err)
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * time.Minute,
	}, nil
}

// Login verify credentials of the user to generate the tokens
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*AuthTokens, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("invalid credentials: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("searching user: %w", err)
	}

	if err := auth.VerifyPassword(password, user.PasswordHash); err != nil {
		if errors.Is(err, auth.ErrPasswordMismatch) {
			return nil, fmt.Errorf("invalid credentials: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("verifying password: %w", err)
	}

	tokens, err := uc.generateTokenPair(ctx, user.ID, user.Email, "member")
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// RefreshTokens create the new tokens based if the refresh token still valid
func (uc *AuthUseCase) RefreshTokens(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	claims, err := uc.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", domain.ErrUnauthorized)
	}

	_, err = uc.tokenStore.ValidateRefreshToken(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, auth.ErrTokenRevoked) || errors.Is(err, auth.ErrTokenNotFound) {
			return nil, fmt.Errorf("invalid refresh token: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("validating refresh token: %w", err)
	}

	if err := uc.tokenStore.RevokeRefreshToken(ctx, claims.ID, uc.refreshTokenDuration); err != nil {
		return nil, fmt.Errorf("revoking token: %w", err)
	}

	tokens, err := uc.generateTokenPair(ctx, claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Logout revoke the user refresh token
func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	claims, err := uc.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil
	}

	if err := uc.tokenStore.RevokeRefreshToken(ctx, claims.ID, uc.refreshTokenDuration); err != nil {
		return fmt.Errorf("revoking refresh token in logout: %w", err)
	}

	return nil
}
