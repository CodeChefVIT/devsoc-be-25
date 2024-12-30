// package utils

// import (
// 	"github.com/labstack/echo/v4"
// )

// func WriteError (r echo.Context,status int, err error) error {
// 	return r.JSON(status, map[string]string{"error":err.Error()})
// }

// func WriteJSON(r echo.Context, status int, v any) error {
// 	return r.JSON(status, map[string]interface{}{"message":v})
// }

package utils


func JSONResponse(status int, message any) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func ErrorResponse(status int, err error) map[string]interface{} {
	return map[string]interface{}{
		"status": status,
		"error":  err.Error(),
	}
}

func WriteJSON(status int, message any) (int, map[string]interface{}) {
	return status, JSONResponse(status, message)
}

func WriteError(status int, err error) (int, map[string]interface{}) {
	return status, ErrorResponse(status, err)
}
