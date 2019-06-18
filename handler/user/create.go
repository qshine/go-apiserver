package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	. "go-apiserver/handler"
	"go-apiserver/model"
	"go-apiserver/pkg/errno"
	"go-apiserver/util"
)

// @Summary Add new user to the database
// @Description Add a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.CreateRequest true "Create a new user"
// @Success 200 {object} user.CreateResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /user [post]
func Create(c *gin.Context) {
	log.Info("User create function called", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest

	// Bind会检查Content-Type类型和参数, 将消息体作为指定的格式解析到 Go struct 变量中
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// 检验参数
	if err := r.checkParam(); err != nil {
		SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// 加密密码
	if err := u.Encrypt(); err != nil {
		SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	// 保存数据到db
	if err := u.Create(); err != nil {
		SendResponse(c, errno.ErrDatabase, nil)
		return
	}
	rsp := CreateResponse{
		Username: r.Username,
	}
	SendResponse(c, nil, rsp)

}

func (r *CreateRequest) checkParam() error {
	if r.Username == "" {
		return errno.New(errno.ErrValidation, nil).Add("username is empty.")
	}

	if r.Password == "" {
		return errno.New(errno.ErrValidation, nil).Add("password is empty.")
	}

	return nil
}
