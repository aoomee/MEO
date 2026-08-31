package handler

import "testing"

func TestReplaceManagedClientCredentialReplacesOnlyMatchingIdentity(t *testing.T) {
	config := map[string]interface{}{
		"inbounds": []interface{}{
			map[string]interface{}{
				"tag":      "vless-443",
				"protocol": "vless",
				"settings": map[string]interface{}{
					"clients": []interface{}{
						map[string]interface{}{"id": "first", "email": "same@example.com"},
						map[string]interface{}{"id": "second", "email": "same@example.com"},
					},
				},
			},
		},
	}
	replaced, err := replaceManagedClientCredential(config, "vless-443", map[string]interface{}{"id": "first"}, map[string]interface{}{"id": "new-first", "email": "same@example.com"})
	if err != nil || !replaced {
		t.Fatalf("replace failed: replaced=%v err=%v", replaced, err)
	}
	_, settings, _, _ := managedInbound(config, "vless-443")
	clients := settings["clients"].([]interface{})
	if clients[0].(map[string]interface{})["id"] != "new-first" || clients[1].(map[string]interface{})["id"] != "second" {
		t.Fatalf("replacement touched the wrong client: %#v", clients)
	}
}

func TestReplaceManagedClientCredentialDoesNotReactivateMissingClient(t *testing.T) {
	config := map[string]interface{}{
		"inbounds": []interface{}{
			map[string]interface{}{
				"tag": "hy2", "protocol": "hysteria2",
				"settings": map[string]interface{}{"clients": []interface{}{}},
			},
		},
	}
	replaced, err := replaceManagedClientCredential(config, "hy2", map[string]interface{}{"auth": "old"}, map[string]interface{}{"auth": "new"})
	if err != nil || replaced {
		t.Fatalf("missing client was reactivated: replaced=%v err=%v", replaced, err)
	}
	_, settings, _, _ := managedInbound(config, "hy2")
	if len(settings["clients"].([]interface{})) != 0 {
		t.Fatalf("missing client was appended: %#v", settings["clients"])
	}
}
