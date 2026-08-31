package bot

import (
	_ "embed"
	"net/http"
)

// Mini App 顶部使用与主面板一致的 MEO 黑白标记。
// 不依赖外部 URL;路由放在 /api/tg-webapp/ 下,确保被现有 nginx 反代覆盖。
//
//go:embed assets/meo-mark.png
var meoMark []byte

func (s *Service) webAppLogoLight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(meoMark)
}

func (s *Service) webAppLogoDark(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(meoMark)
}
