package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.SignupRequest

	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Validation errors",
			Data:    utils.FormatValidationErrors(err),
		})
	}

	existingUserByEmail, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Database error",
		})
	}

	if existingUserByEmail.ID != uuid.Nil {
		return c.JSON(http.StatusConflict, &models.Response{
			Status:  "fail",
			Message: "User with this email already exists",
		})
	}

	userId, err := uuid.NewV7()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate user ID",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to hash password",
		})
	}

	err = utils.Queries.CreateUser(ctx, db.CreateUserParams{
		ID:                userId,
		TeamID:            uuid.NullUUID{Valid: false},
		Email:             req.Email,
		Password:          string(hashedPassword),
		Role:              "student",
		IsLeader:          false,
		IsVerified:        false,
		IsBanned:          false,
		IsProfileComplete: false,
	})
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to create user",
		})
	}

	if err = utils.GenerateOTP(ctx, req.Email); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate OTP",
		})
	}

	token, err := utils.GenerateToken(&userId, false)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate token",
		})
	}

	refreshToken, err := utils.GenerateToken(&userId, true)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate refresh token",
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		Domain:   utils.Config.Domain,
		SameSite: http.SameSiteStrictMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		Domain:   utils.Config.Domain,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "User signed up successfully. OTP has been sent to email",
	})
}

func CompleteProfile(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CompleteProfileRequest
	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	fmt.Printf("Type of req: %T\n", req)

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Validation errors",
			Data:    utils.FormatValidationErrors(err),
		})
	}

	if err := middleware.TrimSpaces(&req); err != nil {
		return c.JSON(echo.ErrBadRequest.Code, models.Response{
			Status: "fail",
			Data: map[string]string{
				"error": err.Error(),
			},
		})
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.FirstName == "" || req.LastName == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "First name and last name cannot be empty",
		})
	}

	if !utils.ValidateAlphaNum(req.LastName) {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "success",
			Message: "Validation errors",
			Data: map[string]any{
				"last_name": "The name must contain only alphabetic characters and spaces. No other characters are allowed",
			},
		})
	}

	if req.Gender != "M" && req.Gender != "F" && req.Gender != "O" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Gender must be M, F or O",
		})
	}

	existingUser, err := utils.Queries.GetUserByPhoneNo(ctx, pgtype.Text{
		String: req.PhoneNo,
		Valid:  true,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Database error",
		})
	}

	if existingUser.ID != uuid.Nil && existingUser.RegNo != nil &&
		*existingUser.RegNo == req.RegNo {
		return c.JSON(http.StatusConflict, &models.Response{
			Status:  "fail",
			Message: "User with this phone number or registration number already exists",
		})
	}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to get user",
		})
	}

	if !user.IsVerified {
		err := utils.GenerateOTP(ctx, user.Email)
		if err != nil {
			logger.Errorf(logger.InternalError, err.Error())
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: "Failed to generate OTP",
			})
		}

		return c.JSON(http.StatusUnauthorized, &models.Response{
			Status:  "fail",
			Message: "User not verified. OTP has been sent to email",
		})
	}

	if user.IsProfileComplete {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Profile already completed",
		})
	}

	err = utils.Queries.CompleteProfile(ctx, db.CompleteProfileParams{
		Email:     user.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		PhoneNo: pgtype.Text{
			String: req.PhoneNo,
			Valid:  true,
		},
		Gender:        req.Gender,
		RegNo:         &req.RegNo,
		GithubProfile: &req.GithubProfile,
		HostelBlock:   &req.HostelBlock,
		RoomNo:        &req.RoomNo,
	})
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to complete profile",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Profile completed successfully",
	})
}

func VerifyOTP(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.VerifyOTPRequest

	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Validation errors",
			Data:    utils.FormatValidationErrors(err),
		})
	}

	user, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Errorf(logger.InternalError, err.Error())
			return c.JSON(http.StatusNotFound, &models.Response{
				Status:  "fail",
				Message: "User not found",
			})
		}
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to get user",
		})
	}

	if user.IsVerified {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "User already verified",
		})
	}

	otp, err := utils.RedisClient.Get(ctx, req.Email).Result()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusNotFound, &models.Response{
			Status:  "fail",
			Message: "OTP invalid/expired",
		})
	}

	if otp != req.OTP {
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Status:  "fail",
			Message: "Invalid OTP",
		})
	}

	err = utils.RedisClient.Del(ctx, req.Email).Err()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to delete OTP",
		})
	}

	err = utils.Queries.VerifyUser(ctx, req.Email)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to verify user",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "User verified successfully",
	})
}

func ResendOTP(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.ResendOTP
	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Validation errors",
			Data:    utils.FormatValidationErrors(err),
		})
	}

	err := utils.GenerateOTP(ctx, req.Email)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())

		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate OTP",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "OTP has been sent to email",
	})
}

func Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.LoginRequest

	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Validation errors",
			Data:    utils.FormatValidationErrors(err),
		})
	}

	user, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Errorf(logger.InternalError, err.Error())
			return c.JSON(http.StatusNotFound, &models.Response{
				Status:  "fail",
				Message: "User not found",
			})
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to get user",
		})
	}

	if user.IsBanned {
		return c.JSON(http.StatusTeapot, &models.Response{
			Status:  "fail",
			Message: "User banned",
		})
	}

	if !user.IsVerified {
		err := utils.GenerateOTP(ctx, req.Email)
		if err != nil {
			logger.Errorf(logger.InternalError, err.Error())

			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: "Failed to generate OTP",
			})
		}

		return c.JSON(http.StatusExpectationFailed, &models.Response{
			Status:  "fail",
			Message: "User not verified. OTP has been sent to email",
			Data: map[string]any{
				"is_verified": false,
			},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Status:  "fail",
			Message: "Invalid password",
		})
	}

	token, err := utils.GenerateToken(&user.ID, false)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate token",
		})
	}

	refreshToken, err := utils.GenerateToken(&user.ID, true)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate refresh token",
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		Domain:   utils.Config.Domain,
		SameSite: http.SameSiteStrictMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		Domain:   utils.Config.Domain,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "User logged in successfully",
		Data: map[string]interface{}{
			"is_profile_complete": user.IsProfileComplete,
			"is_starred":          user.IsStarred,
		},
	})
}

func UpdatePassword(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.UpdatePasswordRequest

	if err := c.Bind(&req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Validation errors",
			Data:    utils.FormatValidationErrors(err),
		})
	}

	_, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Errorf(logger.InternalError, err.Error())
			return c.JSON(http.StatusNotFound, &models.Response{
				Status:  "fail",
				Message: "User not found",
			})
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to get user",
		})
	}

	storedOTP, err := utils.RedisClient.Get(ctx, req.Email).Result()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "OTP expired or invalid",
		})
	}

	if storedOTP != req.OTP {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid OTP",
		})
	}

	err = utils.RedisClient.Del(ctx, req.Email).Err()
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to delete OTP",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to hash password",
		})
	}

	err = utils.Queries.UpdatePassword(ctx, db.UpdatePasswordParams{
		Email:    req.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to update password",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Password updated successfully",
	})
}

func RefreshToken(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return c.JSON(http.StatusForbidden, &models.Response{
				Status:  "success",
				Message: "please login again",
			})
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Status:  "fail",
			Message: "Refresh token not found",
		})
	}

	claims, err := utils.ValidateRefreshToken(refreshToken.Value)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Status:  "fail",
			Message: "Invalid refresh token",
		})
	}

	token, err := utils.GenerateToken(&claims.UserID, false)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate token",
		})
	}

	newRefreshToken, err := utils.GenerateToken(&claims.UserID, true)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to generate refresh token",
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		Domain:   utils.Config.Domain,
		SameSite: http.SameSiteStrictMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		Domain:   utils.Config.Domain,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Token refreshed successfully",
	})
}

func Logout(c echo.Context) error {
	access, err := c.Cookie("jwt")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			access = &http.Cookie{
				Name:     "jwt",
				MaxAge:   -1,
				Value:    "",
				Path:     "/",
				Domain:   utils.Config.Domain,
				SameSite: http.SameSiteStrictMode,
				Secure:   true,
				HttpOnly: true,
			}
		}
	}

	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			refresh = &http.Cookie{
				Name:     "refresh_token",
				MaxAge:   -1,
				Value:    "",
				Path:     "/",
				Domain:   utils.Config.Domain,
				SameSite: http.SameSiteStrictMode,
				Secure:   true,
				HttpOnly: true,
			}
		}
	}

	access.MaxAge = -1
	access.Value = ""
	access.Path = "/"
	access.SameSite = http.SameSiteStrictMode
	access.Secure = true
	access.HttpOnly = true
	c.SetCookie(access)

	refresh.MaxAge = -1
	refresh.Value = ""
	refresh.Path = "/"
	refresh.SameSite = http.SameSiteStrictMode
	refresh.Secure = true
	refresh.HttpOnly = true
	c.SetCookie(refresh)

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Logged out successfully",
	})
}
