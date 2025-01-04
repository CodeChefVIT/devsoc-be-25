package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

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

	if err := utils.Validate.Struct(req); err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	existingUserByEmail, err := utils.Queries.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUserByEmail.ID != uuid.Nil {
		return c.JSON(http.StatusConflict, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User with this email already exists"},
		})
	}

	existingUserByRegNo, err := utils.Queries.GetUserByRegNo(ctx, req.RegNo)
	if err == nil && existingUserByRegNo.ID != uuid.Nil {
		return c.JSON(http.StatusConflict, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User with this registration number already exists"},
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

	userId, _ := uuid.NewV7()
	user := &db.User{
		ID:   userId,
		Name: req.UserName,
		TeamID: uuid.NullUUID{
			Valid: false,
		},
		Email:      req.Email,
		IsVitian:   true,
		RegNo:      req.RegNo,
		Password:   string(hashedPassword),
		PhoneNo:    req.PhoneNo,
		Role:       "student",
		IsLeader:   false,
		College:    "VIT",
		IsVerified: false,
		IsBanned:   false,
	}

	err = utils.Queries.CreateUser(ctx, db.CreateUserParams{
		ID:         user.ID,
		Name:       user.Name,
		TeamID:     user.TeamID,
		Email:      user.Email,
		IsVitian:   user.IsVitian,
		RegNo:      user.RegNo,
		Password:   user.Password,
		PhoneNo:    user.PhoneNo,
		Role:       user.Role,
		IsLeader:   user.IsLeader,
		College:    user.College,
		IsVerified: user.IsVerified,
		IsBanned:   user.IsBanned,
	})
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to create user"},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "User signed up successfully",
		},
	})
}

func SendOTP(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.SendOTPRequest

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

	err = utils.GenerateOTP(ctx, req.Email)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to generate OTP"},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "OTP sent successfully",
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
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User not verified"},
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid password"},
		})
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to generate token"},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "User logged in successfully",
			"token":   token,
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
