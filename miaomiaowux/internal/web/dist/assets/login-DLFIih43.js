import {
  u as _0x417ac5,
  j as _0x52eab8,
  aZ as _0x11fe01,
  t as _0x4dd190,
  v as _0x2db853,
  aQ as _0x47ed37,
  q as _0x4b8fe5,
  r as _0x11aa79,
  aR as _0x55913b,
  a4 as _0x5890e1,
  a_ as _0x371828,
  a5 as _0x5e0695,
  a$ as _0x9baad4,
  aK as _0xaf4184,
  aA as _0x901ad0,
  b0 as _0x18c6a5,
} from "./vendor-modules-0UUaSA6d.js";
import {
  l as _0x4dc602,
  n as _0xc8a919,
  u as _0x4438a,
  h as _0x102abc,
  a as _0x45cc8,
  L as _0x43a322,
  I as _0x14b462,
  B as _0x148670,
} from "./index-0oY9qUmNNK.js";
import { u as _0x4e7b14 } from "./useTurnstile-Bm81HBF_.js";
import {
  C as _0x321023,
  a as _0x39adc6,
  b as _0x478de2,
  c as _0x68dba6,
  d as _0x298ce4,
} from "./card-cr-ClpdW.js";
import { C as _0x5adb55 } from "./checkbox-C1SgKqW3.js";
import {
  I as _0xa23836,
  a as _0x2ab460,
  b as _0x13ec89,
} from "./input-otp-BV4jMh9b.js";
const pe = [
  { value: "zh-CN", label: "中文" },
  { value: "en", label: "English" },
];
function xe({ className: _0x1a2442 }) {
  const { i18n: _0x35c44b } = _0x417ac5();
  return _0x52eab8["jsxs"]("div", {
    className: _0x4dc602(
      "border-border\x20bg-background/80\x20inline-flex\x20items-center\x20gap-1\x20rounded-md\x20border\x20px-1.5\x20py-1\x20backdrop-blur",
      _0x1a2442,
    ),
    children: [
      _0x52eab8["jsx"](_0x11fe01, {
        className: "text-muted-foreground\x20size-4\x20shrink-0",
      }),
      pe["map"]((_0x29cfd1) =>
        _0x52eab8["jsx"](
          "button",
          {
            type: "button",
            onClick: () => {
              _0x35c44b["language"] !== _0x29cfd1["value"] &&
                _0x35c44b["changeLanguage"](_0x29cfd1["value"]);
            },
            className: _0x4dc602(
              "border\x20px-2\x20py-0.5\x20text-xs\x20transition-colors",
              _0x35c44b["language"] === _0x29cfd1["value"]
                ? "bg-primary\x20text-primary-foreground\x20border-primary"
                : "bg-background\x20hover:bg-muted\x20border-border",
            ),
            children: _0x29cfd1["label"],
          },
          _0x29cfd1["value"],
        ),
      ),
    ],
  });
}
function he() {
  const { data: _0x3684ae } = _0x2db853({
      queryKey: ["login-wallpaper"],
      queryFn: async () =>
        (await _0x45cc8["get"]("/api/public/login-wallpaper"))["data"],
      staleTime: 0x1 / 0x0,
    }),
    _0x3355b6 = (_0x3684ae?.["login_wallpaper"] || "")["trim"](),
    _0x2390fd = false;
  return { custom: _0x3355b6, isAnime: _0x2390fd };
}
function Fe() {
  const { t: _0x7a510e } = _0x417ac5("auth"),
    { data: _0x5ed860, isLoading: _0x152b5e } = _0x2db853({
      queryKey: ["setup-status"],
      queryFn: async () => (await _0x45cc8["get"]("/api/setup/status"))["data"],
      staleTime: 0x1 / 0x0,
    });
  let _0x4ce37f;
  return (
    _0x152b5e
      ? (_0x4ce37f = _0x52eab8["jsx"]("div", {
          className:
            "login-pixel-bg\x20flex\x20min-h-svh\x20items-center\x20justify-center\x20px-4\x20py-12",
          children: _0x52eab8["jsx"](_0x321023, {
            className: "w-full\x20max-w-sm",
            children: _0x52eab8["jsxs"](_0x39adc6, {
              className: "space-y-2\x20text-center",
              children: [
                _0x52eab8["jsx"](_0x478de2, {
                  children: _0x7a510e("login.loading"),
                }),
                _0x52eab8["jsx"](_0x68dba6, {
                  children: _0x7a510e("login.checkingStatus"),
                }),
              ],
            }),
          }),
        }))
      : _0x5ed860?.["needs_setup"]
        ? (_0x4ce37f = _0x52eab8["jsx"](fe, {}))
        : (_0x4ce37f = _0x52eab8["jsx"](ge, {})),
    _0x52eab8["jsxs"](_0x52eab8["Fragment"], {
      children: [
        _0x52eab8["jsx"]("div", {
          className: "fixed\x20top-4\x20right-4\x20z-50",
          children: _0x52eab8["jsx"](xe, {}),
        }),
        _0x4ce37f,
      ],
    })
  );
}
function G(_0x4691ea, _0x4aba23, _0x16544c, _0x4c9402, _0xed37b8) {
  (_0x4aba23["setAccessToken"](_0x4691ea["token"]),
    _0x16544c["invalidateQueries"]({ queryKey: ["traffic-summary"] }),
    _0x16544c["setQueryData"](["profile"], {
      username: _0x4691ea["username"],
      email: _0x4691ea["email"],
      nickname: _0x4691ea["nickname"],
      role: _0x4691ea["role"],
      is_admin: _0x4691ea["is_admin"],
    }),
    _0x4dd190["success"](_0xed37b8("login.success")),
    _0x4c9402({ to: "/" }));
}
function ge() {
  const { t: _0x3c5416 } = _0x417ac5("auth"),
    { siteTitle: _0x10b728, logoUrl: _0x56e45e } = _0xc8a919(),
    _0x16ca56 = _0x10b728 || _0x3c5416("login.title"),
    _0x91446e = _0x47ed37(),
    _0x37c9ac = _0x4b8fe5(),
    { auth: _0x3bd1a2 } = _0x4438a(),
    [_0x1250cb, _0x597cbd] = _0x11aa79["useState"](null),
    _0xbf39d1 = _0x55913b({
      defaultValues: { username: "", password: "", remember_me: !0x1 },
    }),
    { data: _0x24f890 } = _0x2db853({
      queryKey: ["captcha-config"],
      queryFn: async () =>
        (await _0x45cc8["get"]("/api/captcha/config"))["data"],
      staleTime: 0x1 / 0x0,
    }),
    _0xa8ba28 = _0x24f890?.["site_key"] || "",
    {
      containerRef: _0x45145e,
      token: _0x211528,
      reset: _0x33ef6d,
    } = _0x4e7b14(_0xa8ba28),
    { custom: _0xf8a00e } = he();
  let _0x35b653,
    _0x574840,
    _0x886261 = null;
  _0xf8a00e
    ? ((_0x35b653 =
        "relative\x20flex\x20min-h-svh\x20items-center\x20justify-center\x20lg:justify-end\x20px-4\x20py-12\x20lg:pr-[7vw]\x20bg-cover\x20bg-center"),
      (_0x574840 = { backgroundImage: "url(\x22" + _0xf8a00e + "\x22)" }),
      (_0x886261 = _0x52eab8["jsx"]("div", {
        className:
          "pointer-events-none\x20absolute\x20inset-0\x20bg-black/25\x20dark:bg-black/45",
      })))
    : (_0x35b653 =
        "login-pixel-bg\x20flex\x20min-h-svh\x20items-center\x20justify-center\x20px-4\x20py-12");
  const _0x5312d0 = _0x5890e1({
      mutationFn: async (_0xc3bb7f) =>
        (await _0x45cc8["post"]("/api/login", _0xc3bb7f))["data"],
      onSuccess: (_0x51c371) => {
        if (_0x51c371["requires_2fa"] && _0x51c371["two_factor_token"]) {
          _0x597cbd(_0x51c371["two_factor_token"]);
          return;
        }
        (G(_0x51c371, _0x3bd1a2, _0x37c9ac, _0x91446e, _0x3c5416),
          _0xbf39d1["reset"]());
      },
      onError: (_0x50e908) => {
        (_0x102abc(_0x50e908),
          _0x4dd190["error"](_0x3c5416("login.failed")),
          _0x33ef6d());
      },
    }),
    _0x3c2226 = _0xbf39d1["handleSubmit"]((_0x4fef61) => {
      if (_0xa8ba28 && !_0x211528) {
        _0x4dd190["error"](_0x3c5416("login.captchaRequired"));
        return;
      }
      _0x5312d0["mutate"]({ ..._0x4fef61, turnstile_token: _0x211528 });
    });
  if (_0x1250cb)
    return _0x52eab8["jsx"](je, {
      twoFactorToken: _0x1250cb,
      onBack: () => _0x597cbd(null),
      onSuccess: (_0x40e537) =>
        G(_0x40e537, _0x3bd1a2, _0x37c9ac, _0x91446e, _0x3c5416),
    });
  const _0x48ce07 = _0x52eab8["jsxs"](_0x321023, {
    className: "relative\x20z-10\x20w-full\x20max-w-sm\x20shadow-lg",
    children: [
      _0x52eab8["jsxs"](_0x39adc6, {
        className: "space-y-2\x20text-center",
        children: [
          _0x52eab8["jsxs"](_0x478de2, {
            className:
              "flex\x20items-center\x20justify-center\x20gap-2\x20text-2xl\x20font-semibold",
            children: [
              _0x56e45e
                ? _0x52eab8["jsx"]("img", {
                    src: _0x56e45e,
                    alt: _0x16ca56 + "\x20Logo",
                    className: "h-9\x20w-9\x20shrink-0\x20object-contain",
                  })
                : null,
              _0x10b728
                ? _0x52eab8["jsx"]("span", {
                    className: "min-w-0\x20truncate",
                    children: _0x10b728,
                  })
                : _0x52eab8["jsx"]("span", {
                    children: _0x3c5416("login.title"),
                  }),
            ],
          }),
        ],
      }),
      _0x52eab8["jsx"](_0x298ce4, {
        children: _0x52eab8["jsxs"]("form", {
          className: "space-y-6",
          onSubmit: _0x3c2226,
          children: [
            _0x52eab8["jsxs"]("div", {
              className: "space-y-2",
              children: [
                _0x52eab8["jsx"](_0x43a322, {
                  htmlFor: "username",
                  children: _0x3c5416("login.username"),
                }),
                _0x52eab8["jsx"](_0x14b462, {
                  id: "username",
                  name: "username",
                  type: "text",
                  autoCapitalize: "none",
                  autoComplete: "username",
                  autoFocus: !0x0,
                  placeholder: _0x3c5416("login.usernamePlaceholder"),
                  ..._0xbf39d1["register"]("username", { required: !0x0 }),
                }),
              ],
            }),
            _0x52eab8["jsxs"]("div", {
              className: "space-y-2",
              children: [
                _0x52eab8["jsx"](_0x43a322, {
                  htmlFor: "password",
                  children: _0x3c5416("login.password"),
                }),
                _0x52eab8["jsx"](_0x14b462, {
                  id: "password",
                  name: "password",
                  type: "password",
                  autoComplete: "current-password",
                  placeholder: _0x3c5416("login.passwordPlaceholder"),
                  ..._0xbf39d1["register"]("password", { required: !0x0 }),
                }),
              ],
            }),
            _0x52eab8["jsxs"]("div", {
              className: "flex\x20items-center\x20space-x-2",
              children: [
                _0x52eab8["jsx"](_0x5adb55, {
                  id: "remember_me",
                  checked: _0xbf39d1["watch"]("remember_me"),
                  onCheckedChange: (_0x8b066d) =>
                    _0xbf39d1["setValue"]("remember_me", _0x8b066d === !0x0),
                }),
                _0x52eab8["jsx"](_0x43a322, {
                  htmlFor: "remember_me",
                  className: "cursor-pointer\x20text-sm\x20font-normal",
                  children: _0x3c5416("login.rememberMe"),
                }),
              ],
            }),
            _0xa8ba28 &&
              _0x52eab8["jsx"]("div", {
                ref: _0x45145e,
                className: "flex\x20justify-center",
              }),
            _0x52eab8["jsx"](_0x148670, {
              type: "submit",
              className: "w-full",
              disabled: _0x5312d0["isPending"],
              children: _0x5312d0["isPending"]
                ? _0x3c5416("login.loggingIn")
                : _0x3c5416("login.loginButton"),
            }),
          ],
        }),
      }),
    ],
  });
  return _0x52eab8["jsxs"]("div", {
    className: _0x35b653,
    style: _0x574840,
    children: [_0x886261, _0x48ce07],
  });
}
function je({
  twoFactorToken: _0x1ac6eb,
  onBack: _0x2a62b4,
  onSuccess: _0x223252,
}) {
  const { t: _0x5d85f4 } = _0x417ac5("auth"),
    [_0x17b704, _0x22e749] = _0x11aa79["useState"](""),
    [_0x7d30a6, _0x111808] = _0x11aa79["useState"](!0x1),
    [_0xa9870b, _0xe54db1] = _0x11aa79["useState"](""),
    _0x4d5d51 = _0x5890e1({
      mutationFn: async (_0x8a4c12) =>
        (
          await _0x45cc8["post"]("/api/login/2fa", {
            two_factor_token: _0x1ac6eb,
            code: _0x8a4c12,
          })
        )["data"],
      onSuccess: (_0xdc4266) => _0x223252(_0xdc4266),
      onError: (_0x1e8789) => {
        (_0x102abc(_0x1e8789),
          _0x4dd190["error"](_0x5d85f4("twoFactor.invalidCode")),
          _0x22e749(""));
      },
    }),
    _0x40f6b4 = _0x5890e1({
      mutationFn: async (_0x3f0fb4) =>
        (
          await _0x45cc8["post"]("/api/login/recovery", {
            two_factor_token: _0x1ac6eb,
            recovery_code: _0x3f0fb4,
          })
        )["data"],
      onSuccess: (_0x1390e3) => {
        (_0x4dd190["success"](_0x5d85f4("twoFactor.recoverySuccess")),
          _0x223252(_0x1390e3));
      },
      onError: (_0x3c0501) => {
        (_0x102abc(_0x3c0501),
          _0x4dd190["error"](_0x5d85f4("twoFactor.invalidRecovery")));
      },
    });
  return _0x52eab8["jsx"]("div", {
    className:
      "login-pixel-bg\x20flex\x20min-h-svh\x20items-center\x20justify-center\x20px-4\x20py-12",
    children: _0x52eab8["jsxs"](_0x321023, {
      className: "w-full\x20max-w-sm\x20shadow-lg",
      children: [
        _0x52eab8["jsxs"](_0x39adc6, {
          className: "space-y-2\x20text-center",
          children: [
            _0x52eab8["jsx"](_0x478de2, {
              className: "text-2xl\x20font-semibold",
              children: _0x5d85f4("twoFactor.title"),
            }),
            _0x52eab8["jsx"](_0x68dba6, {
              children: _0x5d85f4(
                _0x7d30a6 ? "twoFactor.recoveryDesc" : "twoFactor.codeDesc",
              ),
            }),
          ],
        }),
        _0x52eab8["jsxs"](_0x298ce4, {
          className: "space-y-6",
          children: [
            _0x7d30a6
              ? _0x52eab8["jsxs"]("div", {
                  className: "space-y-4",
                  children: [
                    _0x52eab8["jsx"](_0x14b462, {
                      value: _0xa9870b,
                      onChange: (_0x44da64) =>
                        _0xe54db1(_0x44da64["target"]["value"]),
                      placeholder: _0x5d85f4("twoFactor.recoveryPlaceholder"),
                      autoFocus: !0x0,
                      onKeyDown: (_0x2e5665) => {
                        _0x2e5665["key"] === "Enter" &&
                          _0xa9870b["trim"]() &&
                          _0x40f6b4["mutate"](_0xa9870b["trim"]());
                      },
                    }),
                    _0x52eab8["jsx"](_0x148670, {
                      className: "w-full",
                      onClick: () => _0x40f6b4["mutate"](_0xa9870b["trim"]()),
                      disabled: !_0xa9870b["trim"]() || _0x40f6b4["isPending"],
                      children: _0x40f6b4["isPending"]
                        ? _0x5d85f4("twoFactor.verifying")
                        : _0x5d85f4("twoFactor.useRecoveryLogin"),
                    }),
                  ],
                })
              : _0x52eab8["jsxs"]("div", {
                  className: "space-y-4",
                  children: [
                    _0x52eab8["jsx"]("div", {
                      className: "flex\x20justify-center",
                      children: _0x52eab8["jsxs"](_0xa23836, {
                        maxLength: 0x6,
                        value: _0x17b704,
                        onChange: _0x22e749,
                        onComplete: (_0x4d7a13) =>
                          _0x4d5d51["mutate"](_0x4d7a13),
                        autoFocus: !0x0,
                        children: [
                          _0x52eab8["jsxs"](_0x2ab460, {
                            children: [
                              _0x52eab8["jsx"](_0x13ec89, { index: 0x0 }),
                              _0x52eab8["jsx"](_0x13ec89, { index: 0x1 }),
                              _0x52eab8["jsx"](_0x13ec89, { index: 0x2 }),
                            ],
                          }),
                          _0x52eab8["jsxs"](_0x2ab460, {
                            children: [
                              _0x52eab8["jsx"](_0x13ec89, { index: 0x3 }),
                              _0x52eab8["jsx"](_0x13ec89, { index: 0x4 }),
                              _0x52eab8["jsx"](_0x13ec89, { index: 0x5 }),
                            ],
                          }),
                        ],
                      }),
                    }),
                    _0x52eab8["jsx"](_0x148670, {
                      className: "w-full",
                      onClick: () => _0x4d5d51["mutate"](_0x17b704),
                      disabled:
                        _0x17b704["length"] !== 0x6 || _0x4d5d51["isPending"],
                      children: _0x4d5d51["isPending"]
                        ? _0x5d85f4("twoFactor.verifying")
                        : _0x5d85f4("twoFactor.verify"),
                    }),
                  ],
                }),
            _0x52eab8["jsxs"]("div", {
              className: "flex\x20items-center\x20justify-between\x20text-sm",
              children: [
                _0x52eab8["jsxs"]("button", {
                  type: "button",
                  className:
                    "text-muted-foreground\x20hover:text-foreground\x20flex\x20items-center\x20gap-1\x20transition-colors",
                  onClick: _0x2a62b4,
                  children: [
                    _0x52eab8["jsx"](_0x371828, { className: "size-3" }),
                    _0x5d85f4("twoFactor.back"),
                  ],
                }),
                _0x52eab8["jsx"]("button", {
                  type: "button",
                  className:
                    "text-muted-foreground\x20hover:text-foreground\x20transition-colors",
                  onClick: () => {
                    (_0x111808(!_0x7d30a6), _0x22e749(""), _0xe54db1(""));
                  },
                  children: _0x5d85f4(
                    _0x7d30a6
                      ? "twoFactor.useVerificationCode"
                      : "twoFactor.useRecoveryCode",
                  ),
                }),
              ],
            }),
          ],
        }),
      ],
    }),
  });
}
function fe() {
  const { t: _0x93cd55 } = _0x417ac5("auth"),
    _0x98d5b3 = _0x4b8fe5(),
    [_0x42e89a, _0x5d3d59] = _0x11aa79["useState"](null),
    [_0x262c5e, _0x270622] = _0x11aa79["useState"](""),
    [_0x48b9e4, _0x1f2d5d] = _0x11aa79["useState"](""),
    [_0x43a296, _0x3af403] = _0x11aa79["useState"](!0x1),
    [_0x503027, _0x5d84da] = _0x11aa79["useState"](!0x1),
    [_0x2fbc26, _0x5349ac] = _0x11aa79["useState"]({
      host: "127.0.0.1",
      port: 0x1538,
      database: "mmwx",
      username: "mmwx",
      password: "",
      ssl_mode: "prefer",
    }),
    [_0x4a9d8b, _0x48635b] = _0x11aa79["useState"](null),
    _0x324f1c = _0x55913b({
      defaultValues: {
        username: "",
        password: "",
        nickname: "",
        email: "",
        avatar_url: "",
      },
    }),
    _0x84e907 = _0x5890e1({
      mutationFn: async (_0x405eb9) =>
        (
          await _0x45cc8["post"]("/api/setup/verify-domain", {
            domain: _0x405eb9,
          })
        )["data"],
      onSuccess: (_0x4f0b05) => {
        (_0x48635b(_0x4f0b05),
          _0x3af403(_0x4f0b05["match"]),
          _0x4f0b05["match"]
            ? _0x4dd190["success"](_0x93cd55("setup.domainVerified"))
            : _0x4dd190["error"](
                _0x4f0b05["message"] || _0x93cd55("setup.domainMismatch"),
              ));
      },
      onError: (_0x4988a6) => {
        (_0x102abc(_0x4988a6), _0x3af403(!0x1), _0x48635b(null));
      },
    });
  _0x11aa79["useEffect"](() => {
    const _0x5909bf = window["location"]["hostname"];
    _0x5909bf &&
      _0x5909bf !== "localhost" &&
      !/^\d+\.\d+\.\d+\.\d+$/["test"](_0x5909bf) &&
      !_0x5909bf["includes"](":") &&
      (_0x1f2d5d(_0x5909bf), _0x84e907["mutate"](_0x5909bf));
  }, []);
  const _0x2cfe5e = _0x5890e1({
      mutationFn: async (_0x262f2d) =>
        (
          await _0x45cc8["post"]("/api/setup/init", {
            ..._0x262f2d,
            domain: _0x43a296 ? _0x48b9e4 : "",
            database: _0x503027
              ? {
                  driver: "postgres",
                  ..._0x2fbc26,
                  max_open_conns: 0x1e,
                  max_idle_conns: 0xa,
                }
              : void 0x0,
          })
        )["data"],
      onSuccess: (_0x136780) => {
        _0x98d5b3["invalidateQueries"]({ queryKey: ["setup-status"] });
        let _0x762eb2 = _0x93cd55("setup.success");
        if (
          (_0x136780["nginx_setup"] &&
            (_0x762eb2 += "\x20" + _0x93cd55("setup.nginxConfigured")),
          _0x4dd190["success"](_0x762eb2),
          _0x324f1c["reset"](),
          _0x136780["restarting"])
        ) {
          setTimeout(() => {
            _0x136780["redirect_url"]
              ? (window["location"]["href"] =
                  _0x136780["redirect_url"] + "/login")
              : window["location"]["reload"]();
          }, 0x9c4);
          return;
        }
        _0x136780["redirect_url"] &&
          setTimeout(() => {
            window["location"]["href"] = _0x136780["redirect_url"] + "/login";
          }, 0x5dc);
      },
      onError: (_0x3b193c) => {
        (_0x102abc(_0x3b193c), _0x4dd190["error"](_0x93cd55("setup.failed")));
      },
    }),
    _0x587525 = _0x5890e1({
      mutationFn: async (_0x465275) => {
        const _0x160a22 = new FormData();
        (_0x160a22["append"]("backup", _0x465275),
          _0x262c5e && _0x160a22["append"]("passphrase", _0x262c5e));
        const _0x5965a8 =
            (_0x45cc8["defaults"]["baseURL"] ?? "") +
            "/api/setup/restore-backup",
          _0x17dde4 = await fetch(_0x5965a8, {
            method: "POST",
            body: _0x160a22,
          }),
          _0x57c07e = await _0x17dde4["text"]();
        if (!_0x17dde4["ok"]) {
          let _0x191549 = _0x57c07e;
          try {
            const _0x51a3f9 = JSON["parse"](_0x57c07e);
            _0x191549 = _0x51a3f9["error"] || _0x51a3f9["message"] || _0x57c07e;
          } catch {}
          throw new Error(_0x191549 || "HTTP\x20" + _0x17dde4["status"]);
        }
        return _0x57c07e ? JSON["parse"](_0x57c07e) : {};
      },
      onSuccess: () => {
        (_0x98d5b3["invalidateQueries"]({ queryKey: ["setup-status"] }),
          _0x4dd190["success"](_0x93cd55("setup.restoreSuccess")),
          _0x5d3d59(null),
          _0x270622(""),
          setTimeout(() => {
            window["location"]["reload"]();
          }, 0x5dc));
      },
      onError: (_0x443271) => {
        _0x4dd190["error"](_0x93cd55("setup.restoreFailed"), {
          description: _0x443271["message"],
        });
      },
    }),
    _0x37a66a = _0x324f1c["handleSubmit"]((_0x45104f) => {
      _0x2cfe5e["mutate"](_0x45104f);
    }),
    _0x175d9c = _0x48b9e4["trim"]()["length"] > 0x0,
    _0xc84fc4 = _0x2cfe5e["isPending"] || (_0x175d9c && !_0x43a296);
  return _0x52eab8["jsx"]("div", {
    className:
      "login-pixel-bg\x20flex\x20min-h-svh\x20items-center\x20justify-center\x20px-4\x20py-12",
    children: _0x52eab8["jsxs"](_0x321023, {
      className: "w-full\x20max-w-2xl\x20shadow-lg",
      children: [
        _0x52eab8["jsxs"](_0x39adc6, {
          className: "space-y-2\x20text-center",
          children: [
            _0x52eab8["jsx"](_0x478de2, {
              className: "text-2xl\x20font-semibold",
              children: _0x93cd55("setup.welcome"),
            }),
            _0x52eab8["jsx"](_0x68dba6, {
              children: _0x93cd55("setup.firstAdminDesc"),
            }),
          ],
        }),
        _0x52eab8["jsxs"](_0x298ce4, {
          children: [
            _0x52eab8["jsxs"]("form", {
              className: "space-y-4",
              onSubmit: _0x37a66a,
              children: [
                _0x52eab8["jsxs"]("div", {
                  className: "space-y-2",
                  children: [
                    _0x52eab8["jsxs"](_0x43a322, {
                      htmlFor: "setup-username",
                      children: [
                        _0x93cd55("setup.username"),
                        "\x20",
                        _0x52eab8["jsx"]("span", {
                          className: "text-destructive",
                          children: "*",
                        }),
                      ],
                    }),
                    _0x52eab8["jsx"](_0x14b462, {
                      id: "setup-username",
                      name: "username",
                      type: "text",
                      autoCapitalize: "none",
                      autoComplete: "username",
                      autoFocus: !0x0,
                      placeholder: _0x93cd55("setup.usernamePlaceholder"),
                      ..._0x324f1c["register"]("username", {
                        required: !0x0,
                        pattern: /^[a-zA-Z0-9-]{3,20}$/,
                      }),
                    }),
                    _0x324f1c["formState"]["errors"]["username"]?.["type"] ===
                      "pattern" &&
                      _0x52eab8["jsx"]("p", {
                        className: "text-destructive\x20text-xs",
                        children: _0x93cd55("setup.usernameInvalid", {
                          defaultValue:
                            "用户名只能包含字母、数字、短横线,长度\x203-20,不能包含下划线",
                        }),
                      }),
                  ],
                }),
                _0x52eab8["jsxs"]("div", {
                  className: "space-y-2",
                  children: [
                    _0x52eab8["jsxs"](_0x43a322, {
                      htmlFor: "setup-password",
                      children: [
                        _0x93cd55("setup.password"),
                        "\x20",
                        _0x52eab8["jsx"]("span", {
                          className: "text-destructive",
                          children: "*",
                        }),
                      ],
                    }),
                    _0x52eab8["jsx"](_0x14b462, {
                      id: "setup-password",
                      name: "password",
                      type: "password",
                      autoComplete: "new-password",
                      placeholder: _0x93cd55("setup.passwordPlaceholder"),
                      ..._0x324f1c["register"]("password", { required: !0x0 }),
                    }),
                  ],
                }),
                _0x52eab8["jsxs"]("div", {
                  className: "space-y-2",
                  children: [
                    _0x52eab8["jsx"](_0x43a322, {
                      htmlFor: "setup-domain",
                      children: _0x93cd55("setup.domainLabel"),
                    }),
                    _0x52eab8["jsxs"]("div", {
                      className: "flex\x20gap-2",
                      children: [
                        _0x52eab8["jsx"](_0x14b462, {
                          id: "setup-domain",
                          type: "text",
                          placeholder: _0x93cd55("setup.domainPlaceholder"),
                          value: _0x48b9e4,
                          onChange: (_0x251f64) => {
                            (_0x1f2d5d(_0x251f64["target"]["value"]),
                              _0x3af403(!0x1),
                              _0x48635b(null));
                          },
                        }),
                        _0x52eab8["jsx"](_0x148670, {
                          type: "button",
                          variant: "outline",
                          disabled: !_0x175d9c || _0x84e907["isPending"],
                          onClick: () =>
                            _0x84e907["mutate"](_0x48b9e4["trim"]()),
                          children: _0x84e907["isPending"]
                            ? _0x52eab8["jsx"](_0x5e0695, {
                                className: "size-4\x20animate-spin",
                              })
                            : _0x93cd55("setup.verifyButton"),
                        }),
                      ],
                    }),
                    _0x4a9d8b &&
                      _0x52eab8["jsxs"]("div", {
                        className:
                          "flex\x20items-start\x20gap-2\x20text-xs\x20" +
                          (_0x4a9d8b["match"]
                            ? "text-green-600"
                            : "text-destructive"),
                        children: [
                          _0x4a9d8b["match"]
                            ? _0x52eab8["jsx"](_0x9baad4, {
                                className: "mt-0.5\x20size-4\x20shrink-0",
                              })
                            : _0x52eab8["jsx"](_0xaf4184, {
                                className: "mt-0.5\x20size-4\x20shrink-0",
                              }),
                          _0x52eab8["jsx"]("span", {
                            children: _0x4a9d8b["match"]
                              ? _0x93cd55("setup.domainCorrect", {
                                  serverIp: _0x4a9d8b["server_ip"],
                                })
                              : _0x93cd55("setup.domainMismatchDetailed", {
                                  domainIp:
                                    _0x4a9d8b["domain_ips"]?.["join"](
                                      ",\x20",
                                    ) || _0x93cd55("setup.none"),
                                  serverIp:
                                    _0x4a9d8b["server_ip"] ||
                                    _0x93cd55("setup.unknown"),
                                }),
                          }),
                        ],
                      }),
                  ],
                }),
                _0x52eab8["jsxs"]("div", {
                  className: "space-y-2",
                  children: [
                    _0x52eab8["jsx"](_0x43a322, {
                      htmlFor: "setup-nickname",
                      children: _0x93cd55("setup.nickname"),
                    }),
                    _0x52eab8["jsx"](_0x14b462, {
                      id: "setup-nickname",
                      name: "nickname",
                      type: "text",
                      autoComplete: "name",
                      placeholder: _0x93cd55("setup.nicknamePlaceholder"),
                      ..._0x324f1c["register"]("nickname"),
                    }),
                  ],
                }),
                _0x52eab8["jsxs"]("div", {
                  className: "space-y-2",
                  children: [
                    _0x52eab8["jsx"](_0x43a322, {
                      htmlFor: "setup-email",
                      children: _0x93cd55("setup.email"),
                    }),
                    _0x52eab8["jsx"](_0x14b462, {
                      id: "setup-email",
                      name: "email",
                      type: "email",
                      autoComplete: "email",
                      placeholder: _0x93cd55("setup.emailPlaceholder"),
                      ..._0x324f1c["register"]("email"),
                    }),
                  ],
                }),
                _0x52eab8["jsxs"]("div", {
                  className: "space-y-2",
                  children: [
                    _0x52eab8["jsx"](_0x43a322, {
                      htmlFor: "setup-avatar",
                      children: _0x93cd55("setup.avatarUrl"),
                    }),
                    _0x52eab8["jsx"](_0x14b462, {
                      id: "setup-avatar",
                      name: "avatar_url",
                      type: "url",
                      autoComplete: "url",
                      placeholder: _0x93cd55("setup.avatarPlaceholder"),
                      ..._0x324f1c["register"]("avatar_url"),
                    }),
                  ],
                }),
                _0x52eab8["jsxs"]("div", {
                  className: "space-y-4\x20rounded-lg\x20border\x20p-4",
                  children: [
                    _0x52eab8["jsxs"]("label", {
                      className:
                        "flex\x20cursor-pointer\x20items-center\x20gap-2",
                      children: [
                        _0x52eab8["jsx"](_0x5adb55, {
                          checked: _0x503027,
                          onCheckedChange: (_0x25a34c) =>
                            _0x5d84da(_0x25a34c === !0x0),
                        }),
                        _0x52eab8["jsx"]("span", {
                          className: "font-medium",
                          children: "使用\x20PostgreSQL\x20数据库",
                        }),
                      ],
                    }),
                    _0x52eab8["jsx"]("p", {
                      className: "text-muted-foreground\x20text-xs",
                      children:
                        "默认使用\x20data/mmwx.db。高并发或多服务器部署可在首次初始化时直接连接\x20PostgreSQL。",
                    }),
                    _0x503027 &&
                      _0x52eab8["jsxs"]("div", {
                        className: "grid\x20gap-3\x20sm:grid-cols-2",
                        children: [
                          _0x52eab8["jsxs"]("div", {
                            className: "space-y-1",
                            children: [
                              _0x52eab8["jsx"](_0x43a322, { children: "主机" }),
                              _0x52eab8["jsx"](_0x14b462, {
                                value: _0x2fbc26["host"],
                                onChange: (_0x260817) =>
                                  _0x5349ac((_0x22415d) => ({
                                    ..._0x22415d,
                                    host: _0x260817["target"]["value"],
                                  })),
                              }),
                            ],
                          }),
                          _0x52eab8["jsxs"]("div", {
                            className: "space-y-1",
                            children: [
                              _0x52eab8["jsx"](_0x43a322, { children: "端口" }),
                              _0x52eab8["jsx"](_0x14b462, {
                                type: "number",
                                value: _0x2fbc26["port"],
                                onChange: (_0x22dc98) =>
                                  _0x5349ac((_0xd48306) => ({
                                    ..._0xd48306,
                                    port: Number(_0x22dc98["target"]["value"]),
                                  })),
                              }),
                            ],
                          }),
                          _0x52eab8["jsxs"]("div", {
                            className: "space-y-1",
                            children: [
                              _0x52eab8["jsx"](_0x43a322, {
                                children: "数据库名",
                              }),
                              _0x52eab8["jsx"](_0x14b462, {
                                value: _0x2fbc26["database"],
                                onChange: (_0x3569cb) =>
                                  _0x5349ac((_0x1a2cc9) => ({
                                    ..._0x1a2cc9,
                                    database: _0x3569cb["target"]["value"],
                                  })),
                              }),
                            ],
                          }),
                          _0x52eab8["jsxs"]("div", {
                            className: "space-y-1",
                            children: [
                              _0x52eab8["jsx"](_0x43a322, {
                                children: "用户名",
                              }),
                              _0x52eab8["jsx"](_0x14b462, {
                                value: _0x2fbc26["username"],
                                onChange: (_0x1f5877) =>
                                  _0x5349ac((_0x3fe973) => ({
                                    ..._0x3fe973,
                                    username: _0x1f5877["target"]["value"],
                                  })),
                              }),
                            ],
                          }),
                          _0x52eab8["jsxs"]("div", {
                            className: "space-y-1",
                            children: [
                              _0x52eab8["jsx"](_0x43a322, { children: "密码" }),
                              _0x52eab8["jsx"](_0x14b462, {
                                type: "password",
                                autoComplete: "new-password",
                                value: _0x2fbc26["password"],
                                onChange: (_0x2a7c99) =>
                                  _0x5349ac((_0x12791c) => ({
                                    ..._0x12791c,
                                    password: _0x2a7c99["target"]["value"],
                                  })),
                              }),
                            ],
                          }),
                          _0x52eab8["jsxs"]("div", {
                            className: "space-y-1",
                            children: [
                              _0x52eab8["jsx"](_0x43a322, {
                                children: "SSL\x20模式",
                              }),
                              _0x52eab8["jsx"](_0x14b462, {
                                value: _0x2fbc26["ssl_mode"],
                                onChange: (_0x15d299) =>
                                  _0x5349ac((_0x3a2add) => ({
                                    ..._0x3a2add,
                                    ssl_mode: _0x15d299["target"]["value"],
                                  })),
                                placeholder: "prefer",
                              }),
                            ],
                          }),
                        ],
                      }),
                  ],
                }),
                _0x52eab8["jsx"](_0x148670, {
                  type: "submit",
                  className: "w-full",
                  disabled: _0xc84fc4,
                  children: _0x2cfe5e["isPending"]
                    ? _0x93cd55("setup.creating")
                    : _0x93cd55("setup.createAdmin"),
                }),
              ],
            }),
            _0x52eab8["jsxs"]("div", {
              className: "relative\x20my-6",
              children: [
                _0x52eab8["jsx"]("div", {
                  className: "absolute\x20inset-0\x20flex\x20items-center",
                  children: _0x52eab8["jsx"]("span", {
                    className: "w-full\x20border-t",
                  }),
                }),
                _0x52eab8["jsx"]("div", {
                  className:
                    "relative\x20flex\x20justify-center\x20text-xs\x20uppercase",
                  children: _0x52eab8["jsx"]("span", {
                    className: "bg-card\x20text-muted-foreground\x20px-2",
                    children: _0x93cd55("setup.or"),
                  }),
                }),
              ],
            }),
            _0x52eab8["jsxs"]("div", {
              className: "space-y-3",
              children: [
                _0x52eab8["jsx"](_0x43a322, {
                  children: _0x93cd55("setup.restoreFromBackup"),
                }),
                _0x52eab8["jsx"](_0x14b462, {
                  type: "file",
                  accept: ".zip,.enc",
                  onChange: (_0x1c785f) =>
                    _0x5d3d59(_0x1c785f["target"]["files"]?.[0x0] || null),
                  className: "cursor-pointer",
                }),
                _0x52eab8["jsx"](_0x14b462, {
                  type: "password",
                  value: _0x262c5e,
                  onChange: (_0xd02218) =>
                    _0x270622(_0xd02218["target"]["value"]),
                  placeholder: _0x93cd55(
                    "setup.legacyBackupPassphrasePlaceholder",
                  ),
                  autoComplete: "off",
                }),
                _0x52eab8["jsx"]("p", {
                  className: "text-muted-foreground\x20text-xs",
                  children: _0x93cd55("setup.legacyBackupPassphraseHint"),
                }),
                _0x52eab8["jsxs"](_0x148670, {
                  type: "button",
                  onClick: () => _0x42e89a && _0x587525["mutate"](_0x42e89a),
                  disabled: !_0x42e89a || _0x587525["isPending"],
                  variant: "outline",
                  className: "w-full",
                  children: [
                    _0x52eab8["jsx"](_0x901ad0, {
                      className: "mr-2\x20size-4",
                    }),
                    _0x587525["isPending"]
                      ? _0x93cd55("setup.restoring")
                      : _0x93cd55("setup.restoreFromBackup"),
                  ],
                }),
                _0x52eab8["jsxs"]("div", {
                  className:
                    "text-muted-foreground\x20flex\x20items-start\x20gap-2\x20text-xs",
                  children: [
                    _0x52eab8["jsx"](_0x18c6a5, {
                      className: "size-4\x20shrink-0\x20text-amber-500",
                    }),
                    _0x52eab8["jsx"]("span", {
                      children: _0x93cd55("setup.backupHint"),
                    }),
                  ],
                }),
              ],
            }),
          ],
        }),
      ],
    }),
  });
}
export { Fe as component, G as handleLoginSuccess };
