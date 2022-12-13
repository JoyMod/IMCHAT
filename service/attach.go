package service

import (
	"IMCHAT/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Upload(c *gin.Context) {
	w := c.Writer
	r := c.Request

	src, head, err := r.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
	}

	sux := ".png"
	of := head.Filename
	tem := strings.Split(of, ".")

	if len(tem) > 1 {
		sux = "." + tem[len(tem)-1]
	}

	fileName := fmt.Sprintf("%d%04%d%s", time.Now().Unix(), rand.Int31(), sux)
	createRc, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(w, err.Error())
	}

	_, err = io.Copy(createRc, src)
	if err != nil {
		return
	}
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	url := "./asset/upload/" + fileName
	utils.RespOK(w, url, "发送图片成功！")
}
