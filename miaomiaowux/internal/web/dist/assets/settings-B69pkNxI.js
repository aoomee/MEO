import {
  u as _0x26a6ee,
  aQ as _0x34937c,
  q as _0x459e31,
  r as _0x380420,
  v as _0x2396d8,
  aR as _0x23341c,
  a4 as _0x176a5e,
  t as _0x22679b,
  j as _0x208f5b,
  ax as _0x26e7dd,
  aS as _0x32f3cf,
} from "./vendor-modules-0UUaSA6d.js";
import {
  u as _0x4dbf5b,
  h as _0x682145,
  a as _0x16cfcb,
  L as _0x15e369,
  I as _0x35b87f,
  B as _0x3caa21,
  m as _0xb0a545,
  s as _0x41a457,
  D as _0xc729b2,
  d as _0x53abab,
  e as _0x53cd1c,
  f as _0x1e47f6,
  g as _0x1fdbdd,
  p as _0xbd2056,
} from "./index-0oY9qUmNNK.js";
import { a as _0x52ffac, c as _0x1ec1d9 } from "./use-capabilities.js";
import {
  T as _0x4245ed,
  A as _0x2070fd,
  b as _0x486c48,
  c as _0x127458,
} from "./topbar-C6fNrllX.js";
import {
  C as _0x1508a3,
  a as _0x3fba8f,
  b as _0x426f96,
  c as _0x42b0a0,
  d as _0x47b978,
} from "./card-cr-ClpdW.js";
import {
  I as _0x16a78e,
  a as _0x57094b,
  b as _0x52e504,
} from "./input-otp-BV4jMh9b.js";
import "./dropdown-menu-BpG341YX.js";
import "./minimal-badge.js";
import "./checkbox-C1SgKqW3.js";
import "./progress-xgHFkzjD.js";
function Se() {
  const { t: _0x146fac } = _0x26a6ee("settings"),
    _0x210c72 = _0x34937c(),
    _0x6e2750 = _0x459e31(),
    { auth: _0x752f4 } = _0x4dbf5b(),
    { data: _0x549ab9 } = _0x52ffac(),
    _0x8bd5b3 = _0x1ec1d9(_0x549ab9),
    [_0x344f70, _0x25cab9] = _0x380420["useState"](!0x1),
    { data: _0x5ee4c1, isLoading: _0x1fd676 } = _0x2396d8({
      queryKey: ["profile"],
      queryFn: _0xbd2056,
      enabled: !!_0x752f4["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    { data: _0x17236c, isLoading: _0x374159 } = _0x2396d8({
      queryKey: ["user-token"],
      queryFn: async () => (await _0x16cfcb["get"]("/api/user/token"))["data"],
      enabled: !!_0x752f4["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x2bdbd8 = _0x23341c({
      defaultValues: { username: "", nickname: "", email: "", avatar_url: "" },
    });
  _0x380420["useEffect"](() => {
    _0x5ee4c1 &&
      _0x2bdbd8["reset"]({
        username: _0x5ee4c1["username"],
        nickname: _0x5ee4c1["nickname"],
        email: _0x5ee4c1["email"],
        avatar_url: _0x5ee4c1["avatar_url"],
      });
  }, [_0x5ee4c1, _0x2bdbd8]);
  const _0x448514 = _0x176a5e({
      mutationFn: async (_0x27d9f5) => {
        const _0x52060a = {
          username: _0x27d9f5["username"]["trim"](),
          nickname: _0x27d9f5["nickname"]["trim"](),
          email: _0x27d9f5["email"]["trim"](),
          avatar_url: _0x27d9f5["avatar_url"]["trim"](),
        };
        return (await _0x16cfcb["put"]("/api/user/settings", _0x52060a))[
          "data"
        ];
      },
      onSuccess: () => {
        (_0x22679b["success"](_0x146fac("profile.updated")),
          _0x6e2750["invalidateQueries"]({ queryKey: ["profile"] }));
      },
      onError: (_0x432e31) => {
        (_0x682145(_0x432e31),
          _0x22679b["error"](_0x146fac("profile.updateFailed")));
      },
    }),
    _0x109027 = _0x176a5e({
      mutationFn: async () =>
        (await _0x16cfcb["post"]("/api/user/token"))["data"],
      onSuccess: (_0x2c8224) => {
        (_0x6e2750["setQueryData"](["user-token"], _0x2c8224),
          _0x6e2750["invalidateQueries"]({
            queryKey: ["user-package-assignments"],
          }),
          _0x6e2750["invalidateQueries"]({ queryKey: ["admin-nodes"] }),
          _0x25cab9(!0x1),
          _0x22679b["success"](
            _0x146fac("token.reset", {
              count: _0x2c8224["credentials_updated"],
            }),
          ));
      },
      onError: (_0x4e4ccc) => {
        (_0x682145(_0x4e4ccc),
          _0x22679b["error"](_0x146fac("token.resetFailed")));
      },
    }),
    _0x595040 = _0x23341c({
      defaultValues: {
        current_password: "",
        new_password: "",
        confirm_password: "",
      },
    }),
    _0x264197 = _0x176a5e({
      mutationFn: async (_0x546240) =>
        (
          await _0x16cfcb["post"]("/api/user/password", {
            current_password: _0x546240["current_password"],
            new_password: _0x546240["new_password"],
          })
        )["data"],
      onSuccess: () => {
        (_0x22679b["success"](_0x146fac("password.updated")),
          _0x595040["reset"](),
          _0x752f4["reset"](),
          _0x210c72({ to: "/", replace: !0x0 }));
      },
      onError: (_0x3b84ba) => {
        (_0x682145(_0x3b84ba),
          _0x22679b["error"](_0x146fac("password.changeFailed")));
      },
    }),
    _0x2c574a = _0x2bdbd8["handleSubmit"]((_0x492169) => {
      if (!_0x492169["username"]["trim"]()) {
        _0x22679b["error"](_0x146fac("profile.usernameEmpty"));
        return;
      }
      if (
        _0x5ee4c1?.["is_admin"] &&
        _0x492169["username"]["trim"]() !== _0x5ee4c1["username"]
      ) {
        _0x22679b["error"](_0x146fac("profile.adminUsernameImmutable"));
        return;
      }
      if (!/^[a-zA-Z0-9-]{3,20}$/["test"](_0x492169["username"]["trim"]())) {
        _0x22679b["error"](
          _0x146fac("profile.usernameInvalid", {
            defaultValue:
              "用户名只能包含字母、数字、短横线,长度\x203-20,不能包含下划线",
          }),
        );
        return;
      }
      _0x448514["mutate"](_0x492169);
    }),
    _0x3cee94 = _0x595040["handleSubmit"]((_0x1acf40) => {
      if (_0x1acf40["new_password"]["trim"]()["length"] < 0x8) {
        _0x22679b["error"](_0x146fac("password.minLength"));
        return;
      }
      if (_0x1acf40["new_password"] !== _0x1acf40["confirm_password"]) {
        _0x22679b["error"](_0x146fac("password.mismatch"));
        return;
      }
      _0x264197["mutate"](_0x1acf40);
    }),
    _0x5186ab =
      _0x5ee4c1?.["nickname"] ||
      _0x5ee4c1?.["username"] ||
      _0x146fac("defaultUser"),
    _0x269952 = _0x5ee4c1?.["is_admin"]
      ? "/images/meo-mark.png"
      : "/images/meo-mark.png",
    _0x1f0a90 = _0x5ee4c1?.["avatar_url"]?.["trim"]()
      ? _0x5ee4c1["avatar_url"]["trim"]()
      : _0x269952,
    _0x178bf3 = _0x5186ab["slice"](0x0, 0x2) || _0x146fac("defaultUser"),
    _0x2c72d9 = _0x17236c?.["token"] ?? "";
  return _0x208f5b["jsxs"]("div", {
    className: "bg-background\x20min-h-svh",
    children: [
      _0x208f5b["jsx"](_0x4245ed, {}),
      _0x208f5b["jsxs"]("main", {
        className:
          "mx-auto\x20w-full\x20max-w-4xl\x20px-4\x20py-8\x20pt-24\x20sm:px-6",
        children: [
          _0x208f5b["jsx"]("section", {
            className: "space-y-2",
            children: _0x208f5b["jsx"]("h1", {
              className: "text-3xl\x20font-semibold\x20tracking-tight",
              children: _0x146fac("title"),
            }),
          }),
          _0x208f5b["jsxs"]("div", {
            className: "mt-8\x20grid\x20gap-6\x20lg:grid-cols-2",
            children: [
              _0x208f5b["jsxs"]("div", {
                className: "space-y-6",
                children: [
                  _0x208f5b["jsxs"](_0x1508a3, {
                    children: [
                      _0x208f5b["jsxs"](_0x3fba8f, {
                        children: [
                          _0x208f5b["jsx"](_0x426f96, {
                            children: _0x146fac("profile.title"),
                          }),
                          _0x208f5b["jsx"](_0x42b0a0, {
                            children: _0x146fac("profile.description"),
                          }),
                        ],
                      }),
                      _0x208f5b["jsx"](_0x47b978, {
                        children: _0x208f5b["jsxs"]("form", {
                          className: "space-y-5",
                          onSubmit: _0x2c574a,
                          children: [
                            _0x208f5b["jsxs"]("div", {
                              className: "flex\x20items-center\x20gap-4",
                              children: [
                                _0x208f5b["jsxs"](_0x2070fd, {
                                  className: "size-12",
                                  children: [
                                    _0x208f5b["jsx"](_0x486c48, {
                                      src: _0x1f0a90,
                                      alt: _0x5186ab,
                                    }),
                                    _0x208f5b["jsx"](_0x127458, {
                                      children: _0x178bf3,
                                    }),
                                  ],
                                }),
                                _0x208f5b["jsx"]("div", {
                                  className: "text-muted-foreground\x20text-sm",
                                  children: _0x5ee4c1?.["is_admin"]
                                    ? _0x146fac("profile.adminAvatarHint")
                                    : _0x146fac("profile.avatarHint"),
                                }),
                              ],
                            }),
                            _0x208f5b["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x208f5b["jsx"](_0x15e369, {
                                  htmlFor: "username",
                                  children: _0x146fac("profile.username"),
                                }),
                                _0x208f5b["jsx"](_0x35b87f, {
                                  id: "username",
                                  placeholder: _0x146fac(
                                    "profile.usernamePlaceholder",
                                  ),
                                  disabled:
                                    _0x1fd676 || _0x5ee4c1?.["is_admin"],
                                  ..._0x2bdbd8["register"]("username", {
                                    required: !0x0,
                                  }),
                                }),
                                _0x5ee4c1?.["is_admin"]
                                  ? _0x208f5b["jsx"]("p", {
                                      className:
                                        "text-muted-foreground\x20text-xs",
                                      children: _0x146fac(
                                        "profile.adminUsernameDisabled",
                                      ),
                                    })
                                  : null,
                              ],
                            }),
                            _0x208f5b["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x208f5b["jsx"](_0x15e369, {
                                  htmlFor: "nickname",
                                  children: _0x146fac("profile.nickname"),
                                }),
                                _0x208f5b["jsx"](_0x35b87f, {
                                  id: "nickname",
                                  placeholder: _0x146fac(
                                    "profile.nicknamePlaceholder",
                                  ),
                                  disabled: _0x1fd676,
                                  ..._0x2bdbd8["register"]("nickname"),
                                }),
                              ],
                            }),
                            _0x208f5b["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x208f5b["jsx"](_0x15e369, {
                                  htmlFor: "email",
                                  children: _0x146fac("profile.email"),
                                }),
                                _0x208f5b["jsx"](_0x35b87f, {
                                  id: "email",
                                  type: "email",
                                  placeholder: _0x146fac(
                                    "profile.emailPlaceholder",
                                  ),
                                  disabled: _0x1fd676,
                                  ..._0x2bdbd8["register"]("email"),
                                }),
                              ],
                            }),
                            _0x208f5b["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x208f5b["jsx"](_0x15e369, {
                                  htmlFor: "avatar_url",
                                  children: _0x146fac("profile.avatarUrl"),
                                }),
                                _0x208f5b["jsx"](_0x35b87f, {
                                  id: "avatar_url",
                                  placeholder: "https://example.com/avatar.png",
                                  disabled: _0x1fd676,
                                  ..._0x2bdbd8["register"]("avatar_url"),
                                }),
                              ],
                            }),
                            _0x208f5b["jsx"](_0x3caa21, {
                              type: "submit",
                              className: "w-full",
                              disabled: _0x448514["isPending"],
                              children: _0x448514["isPending"]
                                ? _0x146fac("actions.saving", { ns: "common" })
                                : _0x146fac("profile.saveButton"),
                            }),
                          ],
                        }),
                      }),
                    ],
                  }),
                  _0x208f5b["jsxs"](_0x1508a3, {
                    children: [
                      _0x208f5b["jsxs"](_0x3fba8f, {
                        children: [
                          _0x208f5b["jsx"](_0x426f96, {
                            children: _0x146fac("themeStyle.title"),
                          }),
                          _0x208f5b["jsx"](_0x42b0a0, {
                            children: _0x146fac("themeStyle.description"),
                          }),
                        ],
                      }),
                      _0x208f5b["jsxs"](_0x47b978, {
                        children: [
                          _0x208f5b["jsx"]("div", {
                            className: "flex\x20gap-2",
                            children: [
                              {
                                value: "flat",
                                label: "MEO 简约",
                              },
                            ]["map"]((_0x306c7a) =>
                              _0x208f5b["jsx"](
                                "button",
                                {
                                  type: "button",
                                  onClick: () => {
                                    (_0xb0a545("mmw-theme-style") ||
                                      "flat") !== _0x306c7a["value"] &&
                                      (_0x41a457(
                                        "mmw-theme-style",
                                        _0x306c7a["value"],
                                        0xe10 * 0x18 * 0x16d,
                                      ),
                                      window["location"]["reload"]());
                                  },
                                  className:
                                    "flex\x20items-center\x20gap-2\x20rounded-md\x20border\x20px-3\x20py-1.5\x20text-sm\x20transition-colors\x20disabled:cursor-not-allowed\x20disabled:opacity-40\x20" +
                                    ((_0xb0a545("mmw-theme-style") ||
                                      "flat") === _0x306c7a["value"]
                                      ? "bg-primary\x20text-primary-foreground\x20border-primary"
                                      : "bg-background\x20hover:bg-muted\x20border-border"),
                                  children: _0x306c7a["label"],
                                },
                                _0x306c7a["value"],
                              ),
                            ),
                          }),
                        ],
                      }),
                    ],
                  }),
                  _0x208f5b["jsx"](he, {}),
                ],
              }),
              _0x208f5b["jsxs"]("div", {
                className: "space-y-6",
                children: [
                  _0x208f5b["jsxs"](_0x1508a3, {
                    children: [
                      _0x208f5b["jsxs"](_0x3fba8f, {
                        children: [
                          _0x208f5b["jsx"](_0x426f96, {
                            children: _0x146fac("password.title"),
                          }),
                          _0x208f5b["jsx"](_0x42b0a0, {
                            children: _0x146fac("password.description"),
                          }),
                        ],
                      }),
                      _0x208f5b["jsx"](_0x47b978, {
                        children: _0x208f5b["jsxs"]("form", {
                          className: "space-y-4",
                          onSubmit: _0x3cee94,
                          children: [
                            _0x208f5b["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x208f5b["jsx"](_0x15e369, {
                                  htmlFor: "current_password",
                                  children: _0x146fac(
                                    "password.currentPassword",
                                  ),
                                }),
                                _0x208f5b["jsx"](_0x35b87f, {
                                  id: "current_password",
                                  type: "password",
                                  autoComplete: "current-password",
                                  placeholder: _0x146fac(
                                    "password.currentPasswordPlaceholder",
                                  ),
                                  ..._0x595040["register"]("current_password", {
                                    required: !0x0,
                                  }),
                                }),
                              ],
                            }),
                            _0x208f5b["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x208f5b["jsx"](_0x15e369, {
                                  htmlFor: "new_password",
                                  children: _0x146fac("password.newPassword"),
                                }),
                                _0x208f5b["jsx"](_0x35b87f, {
                                  id: "new_password",
                                  type: "password",
                                  autoComplete: "new-password",
                                  placeholder: _0x146fac(
                                    "password.newPasswordHint",
                                  ),
                                  ..._0x595040["register"]("new_password", {
                                    required: !0x0,
                                  }),
                                }),
                              ],
                            }),
                            _0x208f5b["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x208f5b["jsx"](_0x15e369, {
                                  htmlFor: "confirm_password",
                                  children: _0x146fac(
                                    "password.confirmPassword",
                                  ),
                                }),
                                _0x208f5b["jsx"](_0x35b87f, {
                                  id: "confirm_password",
                                  type: "password",
                                  autoComplete: "new-password",
                                  placeholder: _0x146fac(
                                    "password.confirmPasswordPlaceholder",
                                  ),
                                  ..._0x595040["register"]("confirm_password", {
                                    required: !0x0,
                                  }),
                                }),
                              ],
                            }),
                            _0x208f5b["jsx"](_0x3caa21, {
                              type: "submit",
                              className: "w-full",
                              disabled: _0x264197["isPending"],
                              children: _0x264197["isPending"]
                                ? _0x146fac("password.changing")
                                : _0x146fac("password.updateButton"),
                            }),
                          ],
                        }),
                      }),
                    ],
                  }),
                  _0x208f5b["jsxs"](_0x1508a3, {
                    children: [
                      _0x208f5b["jsxs"](_0x3fba8f, {
                        children: [
                          _0x208f5b["jsx"](_0x426f96, {
                            children: _0x146fac("token.title"),
                          }),
                          _0x208f5b["jsx"](_0x42b0a0, {
                            children: _0x208f5b["jsx"]("p", {
                              className:
                                "text-destructive\x20mt-2\x20text-sm\x20font-semibold",
                              children: _0x146fac("token.warning"),
                            }),
                          }),
                        ],
                      }),
                      _0x208f5b["jsxs"](_0x47b978, {
                        className: "space-y-4",
                        children: [
                          _0x208f5b["jsx"]("div", {
                            className:
                              "bg-muted/40\x20rounded-md\x20border\x20p-3\x20font-mono\x20text-xs\x20break-all\x20shadow-inner\x20sm:text-sm",
                            children: _0x374159
                              ? _0x146fac("actions.loading", { ns: "common" })
                              : _0x2c72d9 || _0x146fac("token.notGenerated"),
                          }),
                          _0x208f5b["jsxs"]("div", {
                            className: "flex\x20flex-wrap\x20gap-2",
                            children: [
                              _0x208f5b["jsx"](_0x3caa21, {
                                size: "sm",
                                variant: "secondary",
                                disabled: !_0x2c72d9 || _0x109027["isPending"],
                                onClick: async () => {
                                  if (_0x2c72d9) {
                                    if (
                                      typeof navigator < "u" &&
                                      navigator["clipboard"]?.["writeText"]
                                    )
                                      try {
                                        (await navigator["clipboard"][
                                          "writeText"
                                        ](_0x2c72d9),
                                          _0x22679b["success"](
                                            _0x146fac("token.copied"),
                                          ));
                                        return;
                                      } catch (_0x5ddfc6) {
                                        console["error"](
                                          "copy\x20token\x20failed",
                                          _0x5ddfc6,
                                        );
                                      }
                                    _0x22679b["error"](
                                      _0x146fac("actions.copyFailed", {
                                        ns: "common",
                                      }),
                                    );
                                  }
                                },
                                children: _0x146fac("token.copyButton"),
                              }),
                              _0x208f5b["jsx"](_0x3caa21, {
                                size: "sm",
                                variant: "outline",
                                disabled: _0x109027["isPending"],
                                onClick: () => _0x25cab9(!0x0),
                                children: _0x109027["isPending"]
                                  ? _0x146fac("actions.resetting", {
                                      ns: "common",
                                    })
                                  : _0x146fac("token.resetButton"),
                              }),
                            ],
                          }),
                        ],
                      }),
                    ],
                  }),
                  _0x208f5b["jsx"](_0xc729b2, {
                    open: _0x344f70,
                    onOpenChange: (_0x1819ee) =>
                      !_0x109027["isPending"] && _0x25cab9(_0x1819ee),
                    children: _0x208f5b["jsxs"](_0x53abab, {
                      className: "sm:max-w-md",
                      children: [
                        _0x208f5b["jsxs"](_0x53cd1c, {
                          children: [
                            _0x208f5b["jsx"](_0x1e47f6, {
                              children: _0x146fac("token.confirmTitle"),
                            }),
                            _0x208f5b["jsx"](_0x1fdbdd, {
                              children: _0x146fac("token.confirmDescription"),
                            }),
                          ],
                        }),
                        _0x208f5b["jsxs"]("div", {
                          className: "flex\x20justify-end\x20gap-2",
                          children: [
                            _0x208f5b["jsx"](_0x3caa21, {
                              variant: "outline",
                              disabled: _0x109027["isPending"],
                              onClick: () => _0x25cab9(!0x1),
                              children: _0x146fac("actions.cancel", {
                                ns: "common",
                              }),
                            }),
                            _0x208f5b["jsx"](_0x3caa21, {
                              variant: "destructive",
                              disabled: _0x109027["isPending"],
                              onClick: () => _0x109027["mutate"](),
                              children: _0x109027["isPending"]
                                ? _0x146fac("token.resetting")
                                : _0x146fac("token.confirmButton"),
                            }),
                          ],
                        }),
                      ],
                    }),
                  }),
                  _0x208f5b["jsx"](je, {}),
                ],
              }),
            ],
          }),
        ],
      }),
    ],
  });
}
function he() {
  const { t: _0x392956 } = _0x26a6ee("settings"),
    _0x26e785 = _0x459e31(),
    { data: _0x5a4562 } = _0x2396d8({
      queryKey: ["profile"],
      queryFn: _0xbd2056,
      staleTime: 0x12c * 0x3e8,
    }),
    [_0x8a9047, _0x29d88e] = _0x380420["useState"](!0x1),
    [_0x561d3e, _0x213a21] = _0x380420["useState"](!0x1),
    [_0x2e9352, _0x42dc64] = _0x380420["useState"]("password"),
    [_0x2a9127, _0x45725d] = _0x380420["useState"](""),
    [_0x2049fd, _0x3f0e40] = _0x380420["useState"](""),
    [_0x1c01cb, _0x4c3b08] = _0x380420["useState"](""),
    [_0x23d6f5, _0x396929] = _0x380420["useState"](""),
    [_0x5a12d8, _0xd918c3] = _0x380420["useState"]([]),
    [_0x9ea63c, _0x3cadde] = _0x380420["useState"](""),
    { data: _0x365482 } = _0x2396d8({
      queryKey: ["2fa-status"],
      queryFn: async () =>
        (await _0x16cfcb["get"]("/api/user/2fa/status"))["data"],
      staleTime: 0x7530,
    }),
    _0x57aaf3 = _0x176a5e({
      mutationFn: async (_0x2dc944) =>
        (
          await _0x16cfcb["post"]("/api/user/2fa/setup", {
            password: _0x2dc944,
          })
        )["data"],
      onSuccess: (_0x658e86) => {
        (_0x4c3b08(_0x658e86["secret"]),
          _0x3f0e40(_0x658e86["url"]),
          _0x42dc64("qr"));
      },
      onError: (_0x549697) => {
        (_0x682145(_0x549697),
          _0x22679b["error"](_0x392956("twoFactor.passwordFailed")));
      },
    }),
    _0x546aa2 = _0x176a5e({
      mutationFn: async (_0xb08bdc) =>
        (
          await _0x16cfcb["post"]("/api/user/2fa/verify-setup", {
            code: _0xb08bdc,
          })
        )["data"],
      onSuccess: (_0x44e697) => {
        (_0xd918c3(_0x44e697["recovery_codes"]),
          _0x42dc64("recovery"),
          _0x26e785["invalidateQueries"]({ queryKey: ["2fa-status"] }));
      },
      onError: (_0x2ea6af) => {
        (_0x682145(_0x2ea6af),
          _0x22679b["error"](_0x392956("twoFactor.invalidCode")),
          _0x396929(""));
      },
    }),
    _0x3763be = _0x176a5e({
      mutationFn: async (_0x42f9f2) => {
        await _0x16cfcb["post"]("/api/user/2fa/disable", { code: _0x42f9f2 });
      },
      onSuccess: () => {
        (_0x22679b["success"](_0x392956("twoFactor.disabled")),
          _0x213a21(!0x1),
          _0x3cadde(""),
          _0x26e785["invalidateQueries"]({ queryKey: ["2fa-status"] }));
      },
      onError: (_0x1605c5) => {
        (_0x682145(_0x1605c5),
          _0x22679b["error"](_0x392956("twoFactor.invalidCode")),
          _0x3cadde(""));
      },
    }),
    _0x15c7ec = () => {
      (_0x42dc64("password"),
        _0x45725d(""),
        _0x3f0e40(""),
        _0x4c3b08(""),
        _0x396929(""),
        _0xd918c3([]));
    };
  return _0x208f5b["jsxs"](_0x208f5b["Fragment"], {
    children: [
      _0x208f5b["jsxs"](_0x1508a3, {
        children: [
          _0x208f5b["jsxs"](_0x3fba8f, {
            children: [
              _0x208f5b["jsx"](_0x426f96, {
                children: _0x392956("twoFactor.title"),
              }),
              _0x208f5b["jsx"](_0x42b0a0, {
                children: _0x365482?.["enabled"]
                  ? _0x392956("twoFactor.enabledDesc")
                  : _0x392956("twoFactor.disabledDesc"),
              }),
            ],
          }),
          _0x208f5b["jsx"](_0x47b978, {
            children: _0x365482?.["enabled"]
              ? _0x208f5b["jsx"](_0x3caa21, {
                  variant: "destructive",
                  className: "w-full",
                  onClick: () => _0x213a21(!0x0),
                  children: _0x392956("twoFactor.disableButton"),
                })
              : _0x208f5b["jsx"](_0x3caa21, {
                  className: "w-full",
                  onClick: () => {
                    (_0x15c7ec(), _0x29d88e(!0x0));
                  },
                  children: _0x392956("twoFactor.enableButton"),
                }),
          }),
        ],
      }),
      _0x208f5b["jsx"](_0xc729b2, {
        open: _0x8a9047,
        onOpenChange: (_0x4dafd7) => {
          !_0x4dafd7 &&
            _0x2e9352 !== "recovery" &&
            (_0x29d88e(!0x1), _0x15c7ec());
        },
        children: _0x208f5b["jsxs"](_0x53abab, {
          className: "sm:max-w-md",
          onInteractOutside: (_0x52dcef) => {
            _0x2e9352 === "recovery" && _0x52dcef["preventDefault"]();
          },
          children: [
            _0x208f5b["jsxs"](_0x53cd1c, {
              children: [
                _0x208f5b["jsxs"](_0x1e47f6, {
                  children: [
                    _0x2e9352 === "password" &&
                      _0x392956("twoFactor.steps.password"),
                    _0x2e9352 === "qr" && _0x392956("twoFactor.steps.qrcode"),
                    _0x2e9352 === "verify" &&
                      _0x392956("twoFactor.steps.verify"),
                    _0x2e9352 === "recovery" &&
                      _0x392956("twoFactor.steps.recovery"),
                  ],
                }),
                _0x208f5b["jsxs"](_0x1fdbdd, {
                  children: [
                    _0x2e9352 === "password" &&
                      _0x392956("twoFactor.passwordDesc"),
                    _0x2e9352 === "qr" && _0x392956("twoFactor.qrcodeDesc"),
                    _0x2e9352 === "verify" && _0x392956("twoFactor.verifyDesc"),
                    _0x2e9352 === "recovery" &&
                      _0x392956("twoFactor.recoveryDesc"),
                  ],
                }),
              ],
            }),
            _0x2e9352 === "password" &&
              _0x208f5b["jsxs"]("div", {
                className: "space-y-4",
                children: [
                  _0x208f5b["jsx"](_0x35b87f, {
                    type: "password",
                    placeholder: _0x392956("twoFactor.passwordPlaceholder"),
                    value: _0x2a9127,
                    onChange: (_0x1cc74d) =>
                      _0x45725d(_0x1cc74d["target"]["value"]),
                    onKeyDown: (_0x3d30f1) => {
                      _0x3d30f1["key"] === "Enter" &&
                        _0x2a9127 &&
                        _0x57aaf3["mutate"](_0x2a9127);
                    },
                    autoFocus: !0x0,
                  }),
                  _0x208f5b["jsx"](_0x3caa21, {
                    className: "w-full",
                    disabled: !_0x2a9127 || _0x57aaf3["isPending"],
                    onClick: () => _0x57aaf3["mutate"](_0x2a9127),
                    children: _0x57aaf3["isPending"]
                      ? _0x392956("actions.verifying", { ns: "common" })
                      : _0x392956("actions.next", { ns: "common" }),
                  }),
                ],
              }),
            _0x2e9352 === "qr" &&
              _0x208f5b["jsxs"]("div", {
                className: "space-y-4",
                children: [
                  _0x208f5b["jsx"]("div", {
                    className:
                      "flex\x20justify-center\x20rounded-lg\x20border\x20bg-white\x20p-4",
                    children: _0x208f5b["jsx"](_0x26e7dd, {
                      value: _0x2049fd,
                      size: 0xc8,
                    }),
                  }),
                  _0x208f5b["jsxs"]("div", {
                    className: "space-y-1",
                    children: [
                      _0x208f5b["jsx"](_0x15e369, {
                        className: "text-muted-foreground\x20text-xs",
                        children: _0x392956("twoFactor.manualKey"),
                      }),
                      _0x208f5b["jsx"]("div", {
                        className:
                          "bg-muted/40\x20rounded-md\x20border\x20p-2\x20font-mono\x20text-xs\x20break-all\x20select-all",
                        children: _0x1c01cb,
                      }),
                    ],
                  }),
                  _0x208f5b["jsx"](_0x3caa21, {
                    className: "w-full",
                    onClick: () => _0x42dc64("verify"),
                    children: _0x392956("actions.next", { ns: "common" }),
                  }),
                ],
              }),
            _0x2e9352 === "verify" &&
              _0x208f5b["jsxs"]("div", {
                className: "space-y-4",
                children: [
                  _0x208f5b["jsx"]("div", {
                    className: "flex\x20justify-center",
                    children: _0x208f5b["jsxs"](_0x16a78e, {
                      maxLength: 0x6,
                      value: _0x23d6f5,
                      onChange: _0x396929,
                      onComplete: (_0x3548bd) => _0x546aa2["mutate"](_0x3548bd),
                      autoFocus: !0x0,
                      children: [
                        _0x208f5b["jsxs"](_0x57094b, {
                          children: [
                            _0x208f5b["jsx"](_0x52e504, { index: 0x0 }),
                            _0x208f5b["jsx"](_0x52e504, { index: 0x1 }),
                            _0x208f5b["jsx"](_0x52e504, { index: 0x2 }),
                          ],
                        }),
                        _0x208f5b["jsxs"](_0x57094b, {
                          children: [
                            _0x208f5b["jsx"](_0x52e504, { index: 0x3 }),
                            _0x208f5b["jsx"](_0x52e504, { index: 0x4 }),
                            _0x208f5b["jsx"](_0x52e504, { index: 0x5 }),
                          ],
                        }),
                      ],
                    }),
                  }),
                  _0x208f5b["jsx"](_0x3caa21, {
                    className: "w-full",
                    disabled:
                      _0x23d6f5["length"] !== 0x6 || _0x546aa2["isPending"],
                    onClick: () => _0x546aa2["mutate"](_0x23d6f5),
                    children: _0x546aa2["isPending"]
                      ? _0x392956("actions.verifying", { ns: "common" })
                      : _0x392956("twoFactor.verifyAndEnable"),
                  }),
                ],
              }),
            _0x2e9352 === "recovery" &&
              _0x208f5b["jsxs"]("div", {
                className: "space-y-4",
                children: [
                  _0x208f5b["jsx"]("div", {
                    className:
                      "bg-muted/40\x20grid\x20grid-cols-2\x20gap-2\x20rounded-lg\x20border\x20p-3",
                    children: _0x5a12d8["map"]((_0x180afb) =>
                      _0x208f5b["jsx"](
                        "div",
                        {
                          className: "text-center\x20font-mono\x20text-sm",
                          children: _0x180afb,
                        },
                        _0x180afb,
                      ),
                    ),
                  }),
                  _0x208f5b["jsxs"]("div", {
                    className: "grid\x20grid-cols-2\x20gap-2",
                    children: [
                      _0x208f5b["jsx"](_0x3caa21, {
                        variant: "outline",
                        onClick: async () => {
                          try {
                            (await navigator["clipboard"]["writeText"](
                              _0x5a12d8["join"]("\x0a"),
                            ),
                              _0x22679b["success"](
                                _0x392956("twoFactor.recoveryCodesCopied"),
                              ));
                          } catch {
                            _0x22679b["error"](
                              _0x392956("twoFactor.recoveryCodesCopyFailed"),
                            );
                          }
                        },
                        children: _0x392956("twoFactor.copyRecoveryCodes"),
                      }),
                      _0x208f5b["jsxs"](_0x3caa21, {
                        variant: "outline",
                        onClick: () => {
                          const _0x11fbd7 = _0x5a12d8["join"]("\x0a"),
                            _0x5ba845 = new Blob([_0x11fbd7], {
                              type: "text/plain",
                            }),
                            _0x40a797 = URL["createObjectURL"](_0x5ba845),
                            _0x562bf9 = document["createElement"]("a");
                          ((_0x562bf9["href"] = _0x40a797),
                            (_0x562bf9["download"] =
                              _0x392956("brand", { ns: "common" }) +
                              "-" +
                              (_0x5a4562?.["username"] || "user") +
                              ".txt"),
                            _0x562bf9["click"](),
                            URL["revokeObjectURL"](_0x40a797));
                        },
                        children: [
                          _0x208f5b["jsx"](_0x32f3cf, {
                            className: "mr-1\x20size-4",
                          }),
                          _0x392956("twoFactor.downloadRecoveryCodes"),
                        ],
                      }),
                    ],
                  }),
                  _0x208f5b["jsx"](_0x3caa21, {
                    className: "w-full",
                    onClick: () => {
                      (_0x29d88e(!0x1), _0x15c7ec());
                    },
                    children: _0x392956("twoFactor.recoveryCodesSaved"),
                  }),
                ],
              }),
          ],
        }),
      }),
      _0x208f5b["jsx"](_0xc729b2, {
        open: _0x561d3e,
        onOpenChange: (_0x4259e3) => {
          _0x4259e3 || (_0x213a21(!0x1), _0x3cadde(""));
        },
        children: _0x208f5b["jsxs"](_0x53abab, {
          className: "sm:max-w-md",
          children: [
            _0x208f5b["jsxs"](_0x53cd1c, {
              children: [
                _0x208f5b["jsx"](_0x1e47f6, {
                  children: _0x392956("twoFactor.disableTitle"),
                }),
                _0x208f5b["jsx"](_0x1fdbdd, {
                  children: _0x392956("twoFactor.disableDesc"),
                }),
              ],
            }),
            _0x208f5b["jsxs"]("div", {
              className: "space-y-4",
              children: [
                _0x208f5b["jsx"]("div", {
                  className: "flex\x20justify-center",
                  children: _0x208f5b["jsxs"](_0x16a78e, {
                    maxLength: 0x6,
                    value: _0x9ea63c,
                    onChange: _0x3cadde,
                    onComplete: (_0x1249b5) => _0x3763be["mutate"](_0x1249b5),
                    autoFocus: !0x0,
                    children: [
                      _0x208f5b["jsxs"](_0x57094b, {
                        children: [
                          _0x208f5b["jsx"](_0x52e504, { index: 0x0 }),
                          _0x208f5b["jsx"](_0x52e504, { index: 0x1 }),
                          _0x208f5b["jsx"](_0x52e504, { index: 0x2 }),
                        ],
                      }),
                      _0x208f5b["jsxs"](_0x57094b, {
                        children: [
                          _0x208f5b["jsx"](_0x52e504, { index: 0x3 }),
                          _0x208f5b["jsx"](_0x52e504, { index: 0x4 }),
                          _0x208f5b["jsx"](_0x52e504, { index: 0x5 }),
                        ],
                      }),
                    ],
                  }),
                }),
                _0x208f5b["jsx"](_0x3caa21, {
                  variant: "destructive",
                  className: "w-full",
                  disabled:
                    _0x9ea63c["length"] !== 0x6 || _0x3763be["isPending"],
                  onClick: () => _0x3763be["mutate"](_0x9ea63c),
                  children: _0x3763be["isPending"]
                    ? _0x392956("actions.disabling", { ns: "common" })
                    : _0x392956("twoFactor.confirmDisable"),
                }),
              ],
            }),
          ],
        }),
      }),
    ],
  });
}
function je() {
  const { t: _0x46aaff } = _0x26a6ee("settings"),
    _0x4f1d7e = _0x459e31(),
    [_0x3048cf, _0x336e52] = _0x380420["useState"](""),
    [_0x130abd, _0x134fbb] = _0x380420["useState"](""),
    { data: _0x7f2461, isLoading: _0x27da56 } = _0x2396d8({
      queryKey: ["api-tokens"],
      queryFn: async () =>
        (await _0x16cfcb["get"]("/api/user/api-tokens"))["data"],
    }),
    _0x2db673 = _0x7f2461?.["tokens"] || [],
    _0x2ccddd = _0x176a5e({
      mutationFn: async () =>
        (
          await _0x16cfcb["post"]("/api/user/api-tokens", {
            name: _0x3048cf["trim"](),
          })
        )["data"],
      onSuccess: (_0x4bfab9) => {
        (_0x134fbb(_0x4bfab9["token"]),
          _0x336e52(""),
          _0x4f1d7e["invalidateQueries"]({ queryKey: ["api-tokens"] }),
          _0x22679b["success"](_0x46aaff("apiToken.created")));
      },
      onError: () => _0x22679b["error"](_0x46aaff("apiToken.createFailed")),
    }),
    _0x5c2616 = _0x176a5e({
      mutationFn: async (_0x54d41d) =>
        _0x16cfcb["delete"]("/api/user/api-tokens/" + _0x54d41d),
      onSuccess: () => {
        (_0x4f1d7e["invalidateQueries"]({ queryKey: ["api-tokens"] }),
          _0x22679b["success"](_0x46aaff("apiToken.revoked")));
      },
      onError: () => _0x22679b["error"](_0x46aaff("apiToken.revokeFailed")),
    }),
    _0x27f63f = (_0x454388) =>
      navigator["clipboard"]?.["writeText"](_0x454388)["then"](
        () => _0x22679b["success"](_0x46aaff("token.copied")),
        () => {},
      ),
    _0x2c5ed7 =
      "{\x0a\x20\x20\x22mcp\x22:\x20{\x0a\x20\x20\x20\x20\x22servers\x22:\x20{\x0a\x20\x20\x20\x20\x20\x20\x22miaomiaowux\x22:\x20{\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x22url\x22:\x20\x22" +
      (typeof window < "u"
        ? window["location"]["origin"]
        : "https://your-mmwx") +
      "/mcp\x22,\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x22transport\x22:\x20\x22streamable-http\x22,\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x22headers\x22:\x20{\x20\x22Authorization\x22:\x20\x22Bearer\x20" +
      (_0x130abd || "<your-api-token>") +
      "\x22\x20}\x0a\x20\x20\x20\x20\x20\x20}\x0a\x20\x20\x20\x20}\x0a\x20\x20}\x0a}";
  return _0x208f5b["jsxs"](_0x1508a3, {
    children: [
      _0x208f5b["jsxs"](_0x3fba8f, {
        children: [
          _0x208f5b["jsx"](_0x426f96, {
            children: _0x46aaff("apiToken.title"),
          }),
          _0x208f5b["jsx"](_0x42b0a0, {
            children: _0x46aaff("apiToken.description"),
          }),
        ],
      }),
      _0x208f5b["jsxs"](_0x47b978, {
        className: "space-y-4",
        children: [
          _0x208f5b["jsxs"]("div", {
            className: "flex\x20items-end\x20gap-2",
            children: [
              _0x208f5b["jsxs"]("div", {
                className: "flex-1\x20space-y-1",
                children: [
                  _0x208f5b["jsx"](_0x15e369, {
                    className: "text-xs",
                    children: _0x46aaff("apiToken.nameLabel"),
                  }),
                  _0x208f5b["jsx"](_0x35b87f, {
                    value: _0x3048cf,
                    onChange: (_0x4ee76b) =>
                      _0x336e52(_0x4ee76b["target"]["value"]),
                    placeholder: _0x46aaff("apiToken.namePlaceholder"),
                  }),
                ],
              }),
              _0x208f5b["jsx"](_0x3caa21, {
                size: "sm",
                onClick: () => _0x2ccddd["mutate"](),
                disabled: _0x2ccddd["isPending"],
                children: _0x46aaff("apiToken.createButton"),
              }),
            ],
          }),
          _0x130abd &&
            _0x208f5b["jsxs"]("div", {
              className:
                "border-primary/40\x20bg-primary/5\x20space-y-2\x20rounded-md\x20border\x20p-3",
              children: [
                _0x208f5b["jsx"](_0x15e369, {
                  className: "text-primary\x20text-xs",
                  children: _0x46aaff("apiToken.tokenOnce"),
                }),
                _0x208f5b["jsxs"]("div", {
                  className: "flex\x20gap-2",
                  children: [
                    _0x208f5b["jsx"](_0x35b87f, {
                      readOnly: !0x0,
                      value: _0x130abd,
                      className: "font-mono\x20text-xs",
                    }),
                    _0x208f5b["jsx"](_0x3caa21, {
                      variant: "outline",
                      size: "sm",
                      className: "shrink-0",
                      onClick: () => _0x27f63f(_0x130abd),
                      children: _0x46aaff("token.copyButton"),
                    }),
                  ],
                }),
                _0x208f5b["jsx"](_0x15e369, {
                  className: "text-xs",
                  children: _0x46aaff("apiToken.openclawSnippet"),
                }),
                _0x208f5b["jsx"]("pre", {
                  className:
                    "bg-muted/60\x20max-h-48\x20overflow-auto\x20rounded\x20p-2\x20text-[11px]",
                  children: _0x2c5ed7,
                }),
                _0x208f5b["jsx"](_0x3caa21, {
                  variant: "ghost",
                  size: "sm",
                  onClick: () => _0x27f63f(_0x2c5ed7),
                  children: _0x46aaff("apiToken.copySnippet"),
                }),
              ],
            }),
          _0x208f5b["jsxs"]("div", {
            className: "space-y-1.5",
            children: [
              _0x208f5b["jsx"](_0x15e369, {
                className: "text-xs",
                children: _0x46aaff("apiToken.listTitle"),
              }),
              _0x27da56
                ? _0x208f5b["jsx"]("p", {
                    className: "text-muted-foreground\x20text-xs",
                    children: _0x46aaff("actions.loading", { ns: "common" }),
                  })
                : _0x2db673["length"] === 0x0
                  ? _0x208f5b["jsx"]("p", {
                      className: "text-muted-foreground\x20text-xs",
                      children: _0x46aaff("apiToken.empty"),
                    })
                  : _0x2db673["map"]((_0x310070) =>
                      _0x208f5b["jsxs"](
                        "div",
                        {
                          className:
                            "flex\x20items-center\x20justify-between\x20rounded-md\x20border\x20px-3\x20py-1.5\x20text-xs",
                          children: [
                            _0x208f5b["jsxs"]("div", {
                              className: "min-w-0",
                              children: [
                                _0x208f5b["jsx"]("span", {
                                  className: "font-medium",
                                  children:
                                    _0x310070["name"] || "#" + _0x310070["id"],
                                }),
                                _0x208f5b["jsxs"]("span", {
                                  className: "text-muted-foreground\x20ml-2",
                                  children: [
                                    new Date(_0x310070["created_at"])[
                                      "toLocaleDateString"
                                    ](),
                                    _0x310070["last_used_at"]
                                      ? "\x20·\x20" +
                                        _0x46aaff("apiToken.lastUsed") +
                                        "\x20" +
                                        new Date(_0x310070["last_used_at"])[
                                          "toLocaleDateString"
                                        ]()
                                      : "\x20·\x20" +
                                        _0x46aaff("apiToken.neverUsed"),
                                  ],
                                }),
                              ],
                            }),
                            _0x208f5b["jsx"](_0x3caa21, {
                              variant: "ghost",
                              size: "sm",
                              className:
                                "h-6\x20shrink-0\x20text-red-600\x20hover:text-red-700",
                              disabled: _0x5c2616["isPending"],
                              onClick: () =>
                                _0x5c2616["mutate"](_0x310070["id"]),
                              children: _0x46aaff("apiToken.revoke"),
                            }),
                          ],
                        },
                        _0x310070["id"],
                      ),
                    ),
            ],
          }),
        ],
      }),
    ],
  });
}
export { Se as component };
