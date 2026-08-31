package web

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func servedTitle(t *testing.T) string {
	t.Helper()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	Handler().ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Fatalf("GET / status = %d", recorder.Code)
	}
	return recorder.Body.String()
}

func TestSetSiteTitleInjectsEscapedInitialTitle(t *testing.T) {
	SetSiteTitle(`纠缠 & <缘>`)
	body := servedTitle(t)
	if !strings.Contains(body, `<title>纠缠 &amp; &lt;缘&gt;</title>`) {
		t.Fatalf("served index does not contain escaped custom title")
	}
	if strings.Contains(body, siteTitlePlaceholder) {
		t.Fatalf("served index leaked title placeholder")
	}

	SetSiteTitle("")
	body = servedTitle(t)
	if !strings.Contains(body, `<title></title>`) {
		t.Fatalf("empty title must stay blank until runtime branding loads")
	}
}

func TestSetSiteIconInjectsEscapedInitialIcon(t *testing.T) {
	SetSiteIcon(`/api/branding/icon?v=1&kind="custom"`)
	body := servedTitle(t)
	if !strings.Contains(body, `href="/api/branding/icon?v=1&amp;kind=&#34;custom&#34;"`) {
		t.Fatalf("served index does not contain escaped custom icon")
	}
	if strings.Contains(body, siteIconPlaceholder) {
		t.Fatalf("served index leaked icon placeholder")
	}

	SetSiteIcon("")
	body = servedTitle(t)
	if !strings.Contains(body, `href="`+defaultSiteIcon+`"`) {
		t.Fatalf("empty icon did not restore default icon")
	}
}
