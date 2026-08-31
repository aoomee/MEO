package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"time"

	"crypto/ecdh"

	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/chacha20poly1305"
)

// probeDevice POST /api/admin/wireguard/devices/{id}/probe
// 用 device.probe_*_key 对入站发一次合法 Noise_IK 握手,量 RTT。
func (h *WireGuardHandler) probeDevice(w http.ResponseWriter, r *http.Request, id int64) {
	dev, err := h.repo.GetWGDeviceByID(r.Context(), id)
	if err != nil || dev == nil {
		writeJSONError(w, http.StatusNotFound, "WG 入站不存在")
		return
	}
	if dev.ProbePrivateKey == "" || dev.ServerPublicKey == "" || dev.ListenPort <= 0 {
		writeJSONError(w, http.StatusBadRequest, "探测密钥或监听端口未就绪")
		return
	}
	endpoint := wgDeviceEndpoint(r.Context(), h.repo, dev)
	if endpoint == "" {
		writeJSONError(w, http.StatusBadRequest, "服务器没有可用的公网地址")
		return
	}
	rtt, err := wgHandshakeProbe(endpoint, dev.ProbePrivateKey, dev.ServerPublicKey, 3*time.Second)
	if err != nil {
		respondJSON(w, http.StatusOK, map[string]any{
			"ok": false, "endpoint": endpoint, "error": err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"ok": true, "endpoint": endpoint, "rtt_ms": rtt.Milliseconds(),
	})
}

func wgHandshakeProbe(endpoint, initiatorPrivB64, responderPubB64 string, timeout time.Duration) (time.Duration, error) {
	pkt, err := wgBuildInitiation(initiatorPrivB64, responderPubB64)
	if err != nil {
		return 0, err
	}
	addr, err := net.ResolveUDPAddr("udp", endpoint)
	if err != nil {
		return 0, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))
	start := time.Now()
	if _, err := conn.Write(pkt); err != nil {
		return 0, err
	}
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}
	if n < 92 || buf[0] != 2 {
		return 0, fmt.Errorf("unexpected handshake reply len=%d type=%d", n, buf[0])
	}
	return time.Since(start), nil
}

// wgBuildInitiation 构造 148 字节 Handshake Initiation(Noise_IKpsk2)。
func wgBuildInitiation(initiatorPrivB64, responderPubB64 string) ([]byte, error) {
	initPriv, err := decodeWGKey(initiatorPrivB64)
	if err != nil {
		return nil, err
	}
	respPub, err := decodeWGKey(responderPubB64)
	if err != nil {
		return nil, err
	}
	ephPrivRaw := make([]byte, 32)
	if _, err := rand.Read(ephPrivRaw); err != nil {
		return nil, err
	}
	ephPrivRaw[0] &= 248
	ephPrivRaw[31] &= 127
	ephPrivRaw[31] |= 64
	ephPriv, err := ecdh.X25519().NewPrivateKey(ephPrivRaw)
	if err != nil {
		return nil, err
	}
	ephPub := ephPriv.PublicKey().Bytes()

	respPubKey, err := ecdh.X25519().NewPublicKey(respPub)
	if err != nil {
		return nil, err
	}
	initPrivKey, err := ecdh.X25519().NewPrivateKey(initPriv)
	if err != nil {
		return nil, err
	}
	initPub := initPrivKey.PublicKey().Bytes()

	construction := hashWG([]byte("Noise_IKpsk2_25519_ChaChaPoly_BLAKE2s"))
	identifier := hashWG(append(construction, []byte("WireGuard v1 zx2c4.com")...))
	h := hashWG(append(identifier, respPub...))
	h = hashWG(append(h, ephPub...))
	ck := kdf1(construction, ephPub)

	dh1, err := ephPriv.ECDH(respPubKey)
	if err != nil {
		return nil, err
	}
	ck, k := kdf2(ck, dh1)
	encStatic, err := aeadSeal(k, 0, h, initPub)
	if err != nil {
		return nil, err
	}
	h = hashWG(append(h, encStatic...))

	dh2, err := initPrivKey.ECDH(respPubKey)
	if err != nil {
		return nil, err
	}
	ck, k = kdf2(ck, dh2)
	var ts [12]byte
	now := time.Now().UnixNano()
	binary.LittleEndian.PutUint64(ts[0:8], uint64(now))
	encTS, err := aeadSeal(k, 0, h, ts[:])
	if err != nil {
		return nil, err
	}

	msg := make([]byte, 148)
	msg[0] = 1
	binary.LittleEndian.PutUint32(msg[4:8], 1)
	copy(msg[8:40], ephPub)
	copy(msg[40:88], encStatic)
	copy(msg[88:116], encTS)

	mac1Key := hashWG(append([]byte("mac1----"), respPub...))
	mac, err := blake2s.New128(mac1Key)
	if err != nil {
		return nil, err
	}
	mac.Write(msg[:116])
	copy(msg[116:132], mac.Sum(nil))
	return msg, nil
}

func decodeWGKey(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil || len(b) != 32 {
		return nil, fmt.Errorf("invalid wireguard key")
	}
	return b, nil
}

func hashWG(b []byte) []byte {
	sum := blake2s.Sum256(b)
	return sum[:]
}

func hmacWG(key, data []byte) []byte {
	mac, err := blake2s.New256(key)
	if err != nil {
		sum := blake2s.Sum256(append(key, data...))
		return sum[:]
	}
	mac.Write(data)
	return mac.Sum(nil)
}

func kdf1(key, input []byte) []byte {
	t0 := hmacWG(key, input)
	return hmacWG(t0, []byte{0x1})
}

func kdf2(key, input []byte) (ck, k []byte) {
	t0 := hmacWG(key, input)
	t1 := hmacWG(t0, []byte{0x1})
	t2 := hmacWG(t0, append(t1, 0x2))
	return t1, t2
}

func aeadSeal(key []byte, counter uint64, ad, plain []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, chacha20poly1305.NonceSize)
	binary.LittleEndian.PutUint64(nonce[4:], counter)
	return aead.Seal(nil, nonce, plain, ad), nil
}
