package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	. "go-apiserver/handler"
	"go-apiserver/pkg/errno"
	"net/http"
)

// 新创建用户
func Create(c *gin.Context) {
	var r CreateRequest

	// Bind会检查Content-Type类型, 将消息体作为指定的格式解析到 Go struct 变量中
	if err := c.Bind(&r); err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": errno.ErrBind},
		)
		return
	}

	admin2 := c.Param("username")
	log.Infof("URL username: %s", admin2)

	// 查询url中参数
	desc := c.Query("desc")
	log.Infof("URL key param desc: %s", desc)

	contentType := c.GetHeader("Content-Type")
	log.Infof("Header Content-Type: %s", contentType)

	log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)
	if r.Username == "" {
		err := errno.New(errno.ErrUserNotFound, fmt.Errorf("username can not found in db")).Add("This is add message.")
		SendResponse(c, err, nil)
		return
	}

	if r.Password == "" {
		err := fmt.Errorf("password is empty")
		SendResponse(c, err, nil)
	}

	resp := CreateResponse{
		Username: r.Username,
	}

	SendResponse(c, nil, resp)
}
