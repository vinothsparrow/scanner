package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/vinothsparrow/scanner/model"
)

var (
	ErrInternalServer = &model.Error{500, "Internal Server Error", "Something went wrong."}
	ErrBadRequest     = &model.Error{400, "Bad Request", "The request had bad syntax or was inherently impossible to be satisfied"}
	ErrUnauthorized   = &model.Error{401, "Unauthorized", "Login required"}
	ErrForbidden      = &model.Error{403, "Forbidden", "Access deny"}
	ErrNotFound       = &model.Error{404, "Not Found", "Not found anything matching the URI given"}
	ErrMethodNotFound = &model.Error{405, "Method Not Found", "Not found anything matching method"}
)

func SetErrors(c *gin.Context) {
	c.Set("errors", &model.Errors{})
}

func AddError(c *gin.Context, e *model.Error) {
	errors := c.MustGet("errors").(*model.Errors)
	errors.Errors = append(errors.Errors, e)
}
func GetErrors(c *gin.Context) *model.Errors {
	return c.MustGet("errors").(*model.Errors)
}
