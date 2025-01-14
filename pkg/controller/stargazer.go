package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/labstack/echo/v4"
)

type StarCheckRequest struct {
	Username string `json:"username" validate:"required"`
}

type StargazerResponse []struct {
	Login string `json:"login"`
}

const (
	REPO_OWNER = "CodeChefVIT"
	REPO_NAME  = "cookoff-9.0-backend"
)

func CheckStarred(c echo.Context) error {
	var request StarCheckRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   "invalid body",
		})
	}

	hasStarred := false
	page := 1

	client := &http.Client{}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/stargazers?page=%d&per_page=100", REPO_OWNER, REPO_NAME, page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   "error creating request",
		})
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   "error",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return c.JSON(http.StatusNotFound, models.Response{
			Status: "fail",
			Data:   "repo not found",
		})
	}

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   "error from github",
		})
	}

	var stargazers StargazerResponse
	if err := json.NewDecoder(resp.Body).Decode(&stargazers); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   "parsing error",
		})
	}

	for _, user := range stargazers {
		if user.Login == request.Username {
			hasStarred = true
			break
		}
	}

	if !hasStarred {
		return c.JSON(http.StatusNotFound, models.Response{
			Status: "fail",
			Data:   "user has not starred",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   "user has starred the repo",
	})
}
