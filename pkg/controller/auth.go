package controller

import (
	"fmt"
	"errors"
	"net/http"
	"strings"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5"
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
		fmt.Println(err.Error());
		logger.Errorf(logger.InternalError, err.Error())
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
		SameSite: http.SameSiteNoneMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
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

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Validation errors",
			Data:    utils.FormatValidationErrors(err),
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

	if req.Gender != "M" && req.Gender != "F" && req.Gender != "O" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Gender must be M, F or O",
		})
	}

	existingUserByVitEmail, err := utils.Queries.GetUserByVitEmail(ctx, &req.VitEmail)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Database error",
		})
	}
	if existingUserByVitEmail.ID != uuid.Nil {
		return c.JSON(http.StatusConflict, &models.Response{
			Status:  "fail",
			Message: "User with this VIT email already exists",
		})
	}

	existingUser, err := utils.Queries.GetUserByPhoneNo(ctx, pgtype.Text{
		String: req.PhoneNo,
		Valid:  true,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Database error",
		})
	}

	if existingUser.ID != uuid.Nil && existingUser.RegNo != nil && *existingUser.RegNo == req.RegNo {
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
		VitEmail:      &req.VitEmail,
		HostelBlock:   req.HostelBlock,
		RoomNo:        int32(req.RoomNumber),
		GithubProfile: req.GithubProfile,
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
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
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
		SameSite: http.SameSiteNoneMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "User logged in successfully",
		Data: map[string]interface{}{
			"is_profile_complete": user.IsProfileComplete,
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
		if errors.Is(err, sql.ErrNoRows) {
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
		SameSite: http.SameSiteNoneMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Token refreshed successfully",
	})
}
