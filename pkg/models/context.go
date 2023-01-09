package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/magicLian/gostarter/pkg/log"
)

const (
	CTX_KEY = "go-starter-context"
)

type ReqContext struct {
	*gin.Context
	*SignInUser
	IsSignedIn        bool
	IsMaintainer      bool
	IsAuditSupervisor bool
	IsSystemAdmin     bool
	CurrentOrgId      string
	CurrentOrgName    string
	Log               log.Logger
}

func CtxGetByGin(c *gin.Context) *ReqContext {
	res, has := c.Get("ReqContext")
	if !has {
		return &ReqContext{
			Context:    c,
			SignInUser: &SignInUser{},
		}
	}
	ctx := res.(*ReqContext)
	if ctx == nil {
		c.JSON(500, "System error")
		return nil
	}
	return ctx
}

func CtxGet(ctx context.Context) (*ReqContext, error) {
	c := ctx.Value(CTX_KEY).(*ReqContext)
	if c == nil {
		return nil, errors.New("cloud not retrieve context")
	}
	return c, nil
}

func CtxSet(ctx *ReqContext) context.Context {
	return context.WithValue(ctx.Request.Context(), CTX_KEY, ctx)
}

func (c *ReqContext) Success(message string) {
	resp := make(map[string]interface{})
	resp["message"] = message
	c.JSON(200, resp)
}

func (c *ReqContext) Error(statusCode int, message string, err error) {
	switch err.(type) {
	case *HttpErrorResponse:
		errInfo := err.(*HttpErrorResponse)
		c.HttpError(errInfo.Message, errInfo.StatusCode, errInfo)
	default:
		c.HttpError(message, statusCode, err)
	}
}

func (c *ReqContext) Response(statusCode int, data interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	var b []byte
	var err error

	switch t := data.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		if b, err = json.Marshal(data); err != nil {
			c.Log.Errorf(fmt.Sprintf("body json marshal error:%s", err.Error()))
			b = []byte("body json marshal error")
		}
	}
	c.Writer.Write(b)
}

func (c *ReqContext) HttpError(message string, status int, err error) {
	data := make(map[string]interface{})
	switch status {
	case 404:
		data["message"] = "Not Found"
	case 500:
		data["message"] = "Internal Server Error"
	}

	if message != "" {
		data["message"] = message
	}
	if os.Getenv("env") == "development" {
		if err != nil {
			data["error"] = err.Error()
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		c.Log.Errorf("body json marshal")
		c.Response(status, []byte("body json marshal"))
	}
	c.Response(status, b)
	c.Abort()
}
