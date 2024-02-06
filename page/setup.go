package page

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, path string) {
    r.LoadHTMLGlob("./templates/*.html")
    r.Static("/static", "./static")

    g := r.Group(path)
    for _, page := range pages {
        enablePage(g, page.Template, page.Data)
    }
}

func enablePage(
        r    *gin.RouterGroup,
        tpl  string,
        data map[string]string,
    ) {

    handler := func (c *gin.Context) {
        c.HTML(http.StatusOK, tpl, gin.H{
            "common": nil,
            "page": data,
        })
    }
    r.GET("/" + tpl, handler)
    if tpl == "index.html" {
        r.GET("/", handler)
    }
}
