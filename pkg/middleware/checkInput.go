package middleware

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/labstack/echo/v4"

	"regexp"
)

func RejectSpecialMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		
		allowed := regexp.MustCompile(`[^a-zA-Z0-9]`)

		for k, values := range c.QueryParams() {
			for _, value := range values {
				if allowed.MatchString(value) {
					return c.JSON(echo.ErrBadRequest.Code, &models.Response{
						Status: "fail",
						Data: "Special characters in "+ k + "are not allowed",
					})
				}
			}
		}

		if err := c.Request().ParseForm(); err == nil {
			for k, values := range c.Request().PostForm {
				for _, value := range values {
					if allowed.MatchString(value) {
						return c.JSON(echo.ErrBadRequest.Code, &models.Response{
							Status: "fail",
							Data: "Special characters in "+ k + "are not allowed",
						})
					}
				}
			}
		}
		return next(c)

	}
}

