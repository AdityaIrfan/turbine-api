package services

import (
	"errors"
	"net/http"
	"time"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/phuslu/log"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo      contract.IUserRepository
	authRedisRepo contract.IAuthRedisRepository
	divisionRepo  contract.IDivisionRepository
}

func NewAuthService(
	userRepo contract.IUserRepository,
	authRedisRepo contract.IAuthRedisRepository,
	divisionRepo contract.IDivisionRepository) contract.IAuthService {
	return &authService{
		userRepo:      userRepo,
		authRedisRepo: authRedisRepo,
		divisionRepo:  divisionRepo,
	}
}

func (a *authService) Register(c echo.Context, in *models.Register) error {
	// check username
	exist, err := a.userRepo.IsUsernameExist(in.Username)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "username already in use")
	}

	// check division
	division, err := a.divisionRepo.GetByIdWithSelectedFields(in.DivisionId, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if division.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "division not found")
	}

	user, err := in.ToModel()
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}
	if err := a.userRepo.Create(user); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "register success, waiting for admin permission")
}

func (a *authService) Login(c echo.Context, in *models.Login) error {
	user, err := a.userRepo.GetByUsernameWithSelectedFields(in.Username, "*", "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "credential not found")
	}

	isBlocked, err := a.authRedisRepo.IsLoginBlocked(user.Id)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if isBlocked {
		return helpers.Response(c, http.StatusForbidden, "you were blocked for 10 minutes due to invalid credential 3 times")
	}

	if user.IsInActive() {
		return helpers.Response(c, http.StatusForbidden, "your account is inactive, waiting for admin permission")
	} else if user.IsBlockedByAdmin() {
		return helpers.Response(c, http.StatusForbidden, "your account has been blocked, contact admin for more information")
	}

	hash, _ := helpers.Decrypt(user.PasswordHash)
	salt, _ := helpers.Decrypt(user.PasswordSalt)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Password+salt)); err != nil {
		go a.authRedisRepo.IncLoginFailedCounter(user.Id)
		return helpers.Response(c, http.StatusBadRequest, "wrong password")
	}

	tokenExpiration := time.Now().Add(helpers.LoginExpiration)
	token, err := helpers.GenerateToken(user.Id, tokenExpiration)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	refreshTokenExpiration := time.Now().Add(helpers.RefreshTokenExpiration)
	refreshToken, err := helpers.GenerateRefreshToken(user.Id, refreshTokenExpiration)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	go a.authRedisRepo.SaveRefreshToken(user.Id, &models.RefreshTokenRedis{
		RefreshToken: refreshToken,
		Exp:          refreshTokenExpiration.Unix(),
		Active:       tokenExpiration.Add(time.Second * 1).Unix(),
	}, helpers.RefreshTokenExpiration)

	response := &models.AuthResponse{
		Name:         user.Name,
		Division:     string(user.Division.Type),
		Token:        token,
		RefreshToken: refreshToken,
	}

	return helpers.Response(c, http.StatusOK, "login success", response)
}

func (a *authService) RefreshToken(c echo.Context, in *models.RefreshTokenRequest) error {
	refrehToken, err := helpers.VerifyRefreshToken(in.RefreshToken)
	if err != nil {
		log.Error().Err(errors.New("ERROR VERIFY REFRESH TOKEN : " + err.Error())).Msg("")
		return helpers.Response(c, http.StatusBadRequest, "invalid refresh token")
	} else if !refrehToken.Valid {
		return helpers.Response(c, http.StatusBadRequest, "invalid refresh token")
	}

	var userId string
	if value, ok := refrehToken.Claims.(jwt.MapClaims)["Id"].(string); !ok {
		log.Error().Err(errors.New("ERROR GETTING USER ID FROM CLAIMS : USER ID IS EMPTY OR VALUE IS NOT STRING")).Msg("")
		return helpers.Response(c, http.StatusBadRequest, "invalid refresh token")
	} else {
		userId = value
	}

	if helpers.IsTokenExpired(refrehToken) {
		go a.authRedisRepo.DeleteRefreshToken(userId)
		return helpers.Response(c, http.StatusBadRequest, "refresh token expired")
	}

	refreshTokenRedis, err := a.authRedisRepo.GetRefreshToken(userId)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if refreshTokenRedis == nil {
		return helpers.Response(c, http.StatusBadRequest, "refres token expired")
	}

	if refreshTokenRedis.RefreshToken != in.RefreshToken {
		return helpers.Response(c, http.StatusBadRequest, "invalid refresh token")
	}

	if !refreshTokenRedis.IsActive() {
		return helpers.Response(c, http.StatusBadRequest, "token still active")
	}

	user, err := a.userRepo.GetByIdWithSelectedFields(userId, "*", "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "invalid credential")
	}

	tokenExpiration := time.Now().Add(helpers.LoginExpiration)
	token, err := helpers.GenerateToken(user.Id, tokenExpiration)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	refreshTokenExpiration := time.Now().Add(helpers.RefreshTokenExpiration)
	refreshToken, err := helpers.GenerateRefreshToken(user.Id, refreshTokenExpiration)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	go a.authRedisRepo.SaveRefreshToken(user.Id, &models.RefreshTokenRedis{
		RefreshToken: refreshToken,
		Exp:          refreshTokenExpiration.Unix(),
		Active:       tokenExpiration.Add(time.Second * 1).Unix(),
	}, helpers.RefreshTokenExpiration)

	response := &models.AuthResponse{
		Name:         user.Name,
		Division:     string(user.Division.Type),
		Token:        token,
		RefreshToken: refreshToken,
	}

	return helpers.Response(c, http.StatusOK, "refresh token success", response)
}
