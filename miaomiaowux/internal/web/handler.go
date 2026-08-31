package web

import (
	"bytes"
	"embed"
	"html"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

// Embed the directory as a single tree. Using dist/* makes each immediate
// child an individual match, so an empty directory copied by Vite (for
// example dist/twemoji) causes "contains no embeddable files" at build time.
// Empty directories are irrelevant to the runtime file server and are safely
// ignored when the whole tree is embedded.
//
//go:embed dist
var embeddedFiles embed.FS

const siteTitlePlaceholder = "__MMW_SITE_TITLE__"
const siteIconPlaceholder = "__MMW_SITE_ICON__"
const defaultSiteTitle = "MEO"
const defaultSiteIcon = "/images/meo-mark.png"

var (
	initOnce    sync.Once
	staticFS    fs.FS
	staticFiles http.Handler
	indexBytes  []byte
	indexMod    time.Time

	themeMu          sync.RWMutex
	servedIndex      []byte // indexBytes 替换占位符后的实际下发内容
	currentSiteTitle string
	currentSiteIcon  = defaultSiteIcon
)

func rebuildServedIndexLocked() {
	servedIndex = bytes.ReplaceAll(
		indexBytes,
		[]byte(siteTitlePlaceholder),
		[]byte(html.EscapeString(currentSiteTitle)),
	)
	servedIndex = bytes.ReplaceAll(
		servedIndex,
		[]byte(siteIconPlaceholder),
		[]byte(html.EscapeString(currentSiteIcon)),
	)
	indexMod = time.Now()
}

// SetSiteTitle 更新 HTML 首屏的 <title>。React 启动后仍会通过 /api/branding 热更新，
// 但首次解析 HTML 时浏览器已经拿到正确标题，不再先显示内置名称再跳变。
func SetSiteTitle(title string) {
	initOnce.Do(initialize)
	title = strings.TrimSpace(title)
	themeMu.Lock()
	defer themeMu.Unlock()
	currentSiteTitle = title
	rebuildServedIndexLocked()
}

// SetSiteIcon 更新 HTML 首屏的 favicon。空值使用内置图标。
func SetSiteIcon(iconURL string) {
	initOnce.Do(initialize)
	iconURL = strings.TrimSpace(iconURL)
	if iconURL == "" {
		iconURL = defaultSiteIcon
	}
	themeMu.Lock()
	defer themeMu.Unlock()
	currentSiteIcon = iconURL
	rebuildServedIndexLocked()
}

func initialize() {
	sub, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		panic(err)
	}

	staticFS = sub
	staticFiles = http.FileServer(http.FS(sub))

	indexBytes, err = fs.ReadFile(sub, "index.html")
	if err != nil {
		panic(err)
	}
	rebuildServedIndexLocked()

	if info, err := fs.Stat(sub, "index.html"); err == nil {
		indexMod = info.ModTime()
	} else {
		indexMod = time.Now()
	}
}

// 返回一个为嵌入式前端 SPA 提供服务的 HTTP 处理程序。
func Handler() http.Handler {
	initOnce.Do(initialize)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/traffic/") {
			http.NotFound(w, r)
			return
		}

		cleaned := path.Clean(r.URL.Path)
		if cleaned == "." {
			cleaned = "/"
		}

		if cleaned == "/" {
			serveIndex(w, r)
			return
		}

		resource := strings.TrimPrefix(cleaned, "/")
		if resource == "" {
			serveIndex(w, r)
			return
		}

		if fileExists(resource) {
			staticFiles.ServeHTTP(w, r)
			return
		}

		serveIndex(w, r)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	initOnce.Do(initialize)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// index.html 引用的是带内容哈希的 JS 资源,本身不能被浏览器长缓存,
	// 否则发布新版本后浏览器仍加载旧 bundle(导致动态菜单等新功能不生效)。
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	themeMu.RLock()
	content := servedIndex
	mod := indexMod
	themeMu.RUnlock()
	http.ServeContent(w, r, "index.html", mod, bytes.NewReader(content))
}

func fileExists(name string) bool {
	initOnce.Do(initialize)

	info, err := fs.Stat(staticFS, name)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
