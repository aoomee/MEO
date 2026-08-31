package handler

import "strings"

// isShadowsocks2022Method distinguishes the PSK-based 2022 ciphers from
// legacy Shadowsocks AEAD ciphers. The two families have different
// multi-user password semantics and must not share credential composition.
func isShadowsocks2022Method(method string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(method)), "2022-")
}

// applyShadowsocksCredentialToProxy writes one Xray Shadowsocks user's
// credential to a Clash-compatible proxy.
//
// Legacy SS ignores settings.password in multi-user mode. The client uses the
// selected user's password directly (and that user's method when present).
// SS2022 keeps the server PSK and requires ServerPassword:UserPassword.
func applyShadowsocksCredentialToProxy(proxy map[string]interface{}, credential map[string]interface{}) {
	if proxy == nil || credential == nil {
		return
	}
	method := strings.TrimSpace(shadowsocksStringValue(proxy["cipher"]))
	if userMethod := strings.TrimSpace(shadowsocksStringValue(credential["method"])); userMethod != "" && !isShadowsocks2022Method(method) {
		method = userMethod
		proxy["cipher"] = userMethod
	}
	userPassword := shadowsocksStringValue(credential["password"])
	if userPassword == "" {
		return
	}
	if !isShadowsocks2022Method(method) {
		proxy["password"] = userPassword
		return
	}

	serverPassword := shadowsocksStringValue(proxy["password"])
	if idx := strings.IndexByte(serverPassword, ':'); idx >= 0 {
		serverPassword = serverPassword[:idx]
	}
	if serverPassword == "" {
		// Keep a useful value for malformed/imported records instead of emitting
		// an empty leading segment. Valid SS2022 nodes always provide the PSK.
		proxy["password"] = userPassword
		return
	}
	proxy["password"] = serverPassword + ":" + userPassword
}

func shadowsocksStringValue(value interface{}) string {
	text, _ := value.(string)
	return text
}
