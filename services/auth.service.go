package services

import (
	"errors"
	"net/http"
	"time"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

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
		return helpers.Response(c, http.StatusBadRequest, "username sudah digunakan")
	}

	// check email
	exist, err = a.userRepo.IsEmailExist(in.Email)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "email sudah digunakan")
	}

	// check phone
	exist, err = a.userRepo.IsPhoneExist(in.Phone)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if exist {
		return helpers.Response(c, http.StatusBadRequest, "nomor telefon sudah digunakan")
	}

	// check division
	division, err := a.divisionRepo.GetByIdWithSelectedFields(in.DivisionId, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if division.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "divisi tidak ditemukan")
	}

	user, err := in.ToModel()
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	if err := a.userRepo.Create(user); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil register, menunggu persetujuan admin")
}

func (a *authService) Login(c echo.Context, in *models.Login) error {
	user, err := a.userRepo.GetByUsernameWithSelectedFields(in.Username, "*", "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "username atau password salah")
	}

	isBlocked, err := a.authRedisRepo.IsLoginBlocked(user.Id)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if isBlocked {
		return helpers.Response(c, http.StatusForbidden, "akun anda diblokir, silahkan tunggu 10 menit untuk mencoba lagi")
	}

	if user.IsInActive() {
		return helpers.Response(c, http.StatusForbidden, "akun anda tidak aktif, silahkan hubungi admin untuk informasi lebih lanjut")
	} else if user.IsBlockedByAdmin() {
		return helpers.Response(c, http.StatusForbidden, "akun anda diblokir, silahkan hubungi admin untuk informasi lebih lanjut")
	}

	if user.PasswordHash == "" || user.PasswordSalt == "" {
		return helpers.Response(c, http.StatusBadRequest, "akun anda tidak memiliki password, silahkan hubungi admin")
	}

	hash, err := helpers.Decrypt(user.PasswordHash)
	if err != nil {
		log.Error().Err(errors.New("ERROR DECRYPT PASSWORD HASH : " + err.Error()))
		return helpers.ResponseUnprocessableEntity(c)
	}
	salt, err := helpers.Decrypt(user.PasswordSalt)
	if err != nil {
		log.Error().Err(errors.New("ERROR DECRYPT PASSWORD SALT : " + err.Error()))
		return helpers.ResponseUnprocessableEntity(c)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Password+salt)); err != nil {
		go a.authRedisRepo.IncLoginFailedCounter(user.Id)
		return helpers.Response(c, http.StatusBadRequest, "username atau password salah")
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

	// go a.authRedisRepo.SaveToken(user.Id, token, helpers.LoginExpiration)

	go a.authRedisRepo.SaveRefreshToken(user.Id, &models.RefreshTokenRedis{
		RefreshToken: refreshToken,
		Exp:          refreshTokenExpiration.Unix(),
		Active:       tokenExpiration.Add(time.Second * 1).Unix(),
	}, helpers.RefreshTokenExpiration)

	response := &models.AuthResponse{
		Name:         user.Name,
		Division:     string(user.Division.Name),
		Source:       user.GetSource(),
		Token:        token,
		RefreshToken: refreshToken,
	}

	return helpers.Response(c, http.StatusOK, "berhasil login", response)
}

func (a *authService) RefreshToken(c echo.Context, in *models.RefreshTokenRequest) error {
	refrehToken, err := helpers.VerifyRefreshToken(in.RefreshToken)
	if err != nil {
		log.Error().Err(errors.New("ERROR VERIFY REFRESH TOKEN : " + err.Error())).Msg("")
		return helpers.Response(c, http.StatusBadRequest, "invalid refresh token")
	} else if !refrehToken.Valid {
		return helpers.Response(c, http.StatusBadRequest, "refresh token tidak valid")
	}

	var userId string
	if value, ok := refrehToken.Claims.(jwt.MapClaims)["Id"].(string); !ok {
		log.Error().Err(errors.New("ERROR GETTING USER ID FROM CLAIMS : USER ID IS EMPTY OR VALUE IS NOT STRING")).Msg("")
		return helpers.Response(c, http.StatusBadRequest, "refresh token tidak valid")
	} else {
		userId = value
	}

	if helpers.IsTokenExpired(refrehToken) {
		go a.authRedisRepo.DeleteRefreshToken(userId)
		return helpers.Response(c, http.StatusBadRequest, "refresh token kedaluarsa")
	}

	refreshTokenRedis, err := a.authRedisRepo.GetRefreshToken(userId)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if refreshTokenRedis == nil {
		return helpers.Response(c, http.StatusBadRequest, "refresh token kedaluarsa")
	}

	if refreshTokenRedis.RefreshToken != in.RefreshToken {
		return helpers.Response(c, http.StatusBadRequest, "refresh token tidak valid")
	}

	if !refreshTokenRedis.IsActive() {
		return helpers.Response(c, http.StatusBadRequest, "token masih aktif")
	}

	user, err := a.userRepo.GetByIdWithSelectedFields(userId, "*", "Division")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.Response(c, http.StatusForbidden, "akun tidak valid")
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

	// go a.authRedisRepo.SaveToken(user.Id, token, helpers.LoginExpiration)

	go a.authRedisRepo.SaveRefreshToken(user.Id, &models.RefreshTokenRedis{
		RefreshToken: refreshToken,
		Exp:          refreshTokenExpiration.Unix(),
		Active:       tokenExpiration.Add(time.Second * 1).Unix(),
	}, helpers.RefreshTokenExpiration)

	response := &models.AuthResponse{
		Name:         user.Name,
		Division:     string(user.Division.Name),
		Source:       user.GetSource(),
		Token:        token,
		RefreshToken: refreshToken,
	}

	return helpers.Response(c, http.StatusOK, "refresh token berhasil", response)
}

func (a *authService) Logout(c echo.Context, token string) error {
	// userId, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string)
	// if !ok {
	// 	return helpers.ResponseForbiddenAccess(c)
	// }

	// user, err := a.userRepo.GetByIdWithSelectedFields(userId, "id")
	// if err != nil {
	// 	return helpers.ResponseUnprocessableEntity(c)
	// } else if user.IsEmpty() {
	// 	return helpers.ResponseForbiddenAccess(c)
	// }

	// existingToken, err := a.authRedisRepo.GetToken(userId)
	// if err != nil {
	// 	return helpers.ResponseUnprocessableEntity(c)
	// } else if existingToken == "" || existingToken != token {
	// 	return helpers.Response(c, http.StatusBadRequest, "token tidak valid")
	// }

	// go a.authRedisRepo.DeleteToken(userId)

	return helpers.Response(c, http.StatusOK, "berhasil logout")
}
