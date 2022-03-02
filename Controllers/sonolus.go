package Controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RedirectSonolusUploadScript TODO: Rewrite Sonolus Uploader (mysql => meilisearch[like mongodb])
func RedirectSonolusUploadScript(ctx *gin.Context) {
	ctx.Redirect(http.StatusTemporaryRedirect, "https://service-mtxfzbvy-1300838857.gz.apigw.tencentcs.com/release/SonolusTestUpload")
}

// RedirectSonolusUploadSong TODO: Config in gate
func RedirectSonolusUploadSong(ctx *gin.Context) {
	ctx.Redirect(http.StatusTemporaryRedirect, "https://upload.ayachan.fun:24444/Sonolus")
}
