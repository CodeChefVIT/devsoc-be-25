package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.SignupRequest

	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid request body"},
		})
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.FirstName == "" || req.LastName == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "First name and last name cannot be empty"},
		})
	}

	if req.Gender != "M" && req.Gender != "F" && req.Gender != "O" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Gender must be M, F or O"},
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	existingUserByEmail, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Database error"},
		})
	}
	if existingUserByEmail.ID != uuid.Nil {
		return c.JSON(http.StatusConflict, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User with this email already exists"},
		})
	}

	existingUserByVitEmail, err := utils.Queries.GetUserByVitEmail(ctx, req.VitEmail)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Database error"},
		})
	}
	if existingUserByVitEmail.ID != uuid.Nil {
		return c.JSON(http.StatusConflict, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User with this VIT email already exists"},
		})
	}

	existingUserByPhoneNo, err := utils.Queries.GetUserByPhoneNo(ctx, req.PhoneNo)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Database error"},
		})
	}
	if existingUserByPhoneNo.ID != uuid.Nil {
		return c.JSON(http.StatusConflict, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User with this phone number already exists"},
		})
	}

	existingUserByRegNo, err := utils.Queries.GetUserByRegNo(ctx, req.RegNo)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Database error"},
		})
	}
	if existingUserByRegNo.ID != uuid.Nil {
		return c.JSON(http.StatusConflict, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User with this registration number already exists"},
		})
	}

	userId, err := uuid.NewV7()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to generate user ID"},
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to hash password"},
		})
	}

	err = utils.Queries.CreateUser(ctx, db.CreateUserParams{
		ID:            userId,
		TeamID:        uuid.NullUUID{Valid: false},
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		PhoneNo:       req.PhoneNo,
		Gender:        req.Gender,
		RegNo:         req.RegNo,
		VitEmail:      req.VitEmail,
		HostelBlock:   req.HostelBlock,
		RoomNo:        int32(req.RoomNumber),
		GithubProfile: req.GithubProfile,
		Password:      string(hashedPassword),
		Role:          "student",
		IsLeader:      false,
		IsVerified:    false,
		IsBanned:      false,
	})
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to create user"},
		})
	}

	if err = utils.GenerateOTP(ctx, req.Email); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to generate OTP"},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "User signed up successfully. OTP has been sent to email",
		},
	})
}

func VerifyOTP(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.VerifyOTPRequest

	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid request body"},
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	_, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf(logger.InternalError, err.Error())
			return c.JSON(http.StatusNotFound, models.Response{
				Status: "fail",
				Data:   map[string]string{"error": "User not found"},
			})
		}
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to get user"},
		})
	}

	otp, err := utils.RedisClient.Get(ctx, req.Email).Result()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusNotFound, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "OTP invalid/expired"},
		})
	}

	if otp != req.OTP {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid OTP"},
		})
	}

	err = utils.RedisClient.Del(ctx, req.Email).Err()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to delete OTP"},
		})
	}

	err = utils.Queries.VerifyUser(ctx, req.Email)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to verify user"},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "User verified successfully",
		},
	})
}

func Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.LoginRequest

	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid request body"},
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	user, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf(logger.InternalError, err.Error())
			return c.JSON(http.StatusNotFound, models.Response{
				Status: "fail",
				Data:   map[string]string{"error": "User not found"},
			})
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to get user"},
		})
	}

	if !user.IsVerified {
		err := utils.GenerateOTP(ctx, req.Email)
		if err != nil {
			logger.Errorf(logger.InternalError, err.Error())
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status: "fail",
				Data:   map[string]string{"error": "Failed to generate OTP"},
			})
		}

		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User not verified. OTP has been sent to email"},
		})
	}

	if user.IsBanned {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User banned"},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid password"},
		})
	}

	token, err := utils.GenerateToken(&user, false)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to generate token"},
		})
	}

	refreshToken, err := utils.GenerateToken(&user, true)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to generate refresh token"},
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "User logged in successfully",
		},
	})
}

func UpdatePassword(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.UpdatePasswordRequest

	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid request body"},
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	_, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf(logger.InternalError, err.Error())
			return c.JSON(http.StatusNotFound, models.Response{
				Status: "fail",
				Data:   map[string]string{"error": "User not found"},
			})
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to get user"},
		})
	}

	storedOTP, err := utils.RedisClient.Get(ctx, req.Email).Result()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "OTP expired or invalid"},
		})
	}

	if storedOTP != req.OTP {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid OTP"},
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to hash password"},
		})
	}

	err = utils.Queries.UpdatePassword(ctx, db.UpdatePasswordParams{
		Email:    req.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to update password"},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "Password updated successfully",
		},
	})
}

func RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Refresh token not found"},
		})
	}

	refreshClaims, err := utils.ValidateRefreshToken(refreshToken.Value)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid refresh token"},
		})
	}

	dbUser, err := utils.Queries.GetUserByID(ctx, refreshClaims.UserID)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusNotFound, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User not found"},
		})
	}

	token, err := utils.GenerateToken(&dbUser, false)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to generate token"},
		})
	}

	newRefreshToken, err := utils.GenerateToken(&dbUser, true)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to generate refresh token"},
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "Token refreshed successfully",
		},
	})
}
