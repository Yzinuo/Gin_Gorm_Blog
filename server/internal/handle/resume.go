package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getResume(c *gin.Context) {
	doc, err := model.GetResume(GetDB(c))
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, doc)
}

func updateResume(c *gin.Context) {
	// Fail closed even when the legacy resource registry has no /setting/about entry.
	if auth, ok := c.Get(g.CTX_USER_AUTH); !ok || auth == nil {
		ReturnError(c, g.ErrPermission, nil)
		return
	}
	var doc model.ResumeDocument
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 128<<10)
	if err := c.ShouldBindJSON(&doc); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	if err := model.ValidateResume(doc); err != nil {
		ReturnError(c, g.ErrRequest, err.Error())
		return
	}
	if err := model.SaveResume(GetDB(c), doc); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, doc)
}
