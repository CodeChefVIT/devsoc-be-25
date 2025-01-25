package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
)

type StarCheckRequest struct {
	Github_link string `json:"github_link" validate:"required"`
}

var Client = http.Client{}

func CheckStarred(c echo.Context) error {
	//var request StarCheckRequest
	//if err := c.Bind(&request); err != nil {
	//	return c.JSON(http.StatusBadRequest, models.Response{
	//		Status: "fail",
	//		Data:   "invalid body",
	//	})
	//}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status:  "fail",
			Message: "User not found",
		})
	}
	github_link := user.GithubProfile

	hasStarred := false

	owner := utils.Config.RepoOwner
	name := utils.Config.RepoName

	github_user := strings.Split(github_link, "github.com/")

	if len(github_user) != 2 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: "error invalid github link",
		})
	}

	baseURL := fmt.Sprintf("https://api.github.com/users/%s/starred", github_user[1])
	u, _ := url.Parse(baseURL)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "fail",
			Message: "error creating request",
		})
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", utils.Config.GithubPAT))

	resp, err := Client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return c.JSON(http.StatusNotFound, models.Response{
			Status:  "fail",
			Message: "repo not found",
		})
	}

	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "fail",
			Message: "error from github",
		})
	}

	bytes, _ := io.ReadAll(resp.Body)

	var repos []struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	}

	err = json.Unmarshal(bytes, &repos)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "fail",
			Message: "parsing error",
		})
	}

	fmt.Println(repos)

	for _, user := range repos {
		fmt.Println(owner, name)
		fmt.Println(user.Owner.Login, user.Name)
		if user.Owner.Login == name && user.Name == owner {
			hasStarred = true
			break
		}
	}

	if !hasStarred {
		return c.JSON(http.StatusExpectationFailed, models.Response{
			Status:  "fail",
			Message: "user has not starred. Please star the repository",
		})
	}

	err = utils.Queries.UpdateStarred(c.Request().Context(), db.UpdateStarredParams{IsStarred: true, Email: user.Email})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to update the user",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   "user has starred the repo",
	})
}
