package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
)

type StarCheckRequest struct {
	Username string `json:"username" validate:"required"`
}

type StargazerResponse []struct {
	Login string `json:"login"`
}

var Client = http.Client{}

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

	owner := utils.Config.RepoOwner
	name := utils.Config.RepoName

	baseURL := "https://api.github.com/repos"
	u, _ := url.Parse(baseURL)
	u.Path += "/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/stargazers"
	query := u.Query()
	query.Set("per_page", "100")
	query.Set("page", fmt.Sprintf("%d", page))
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   "error creating request",
		})
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := Client.Do(req)
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
		return c.JSON(http.StatusOK, models.Response{
			Status: "success",
			Data:   "user has not starred. Please star the repository",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   "user has starred the repo",
	})
}
