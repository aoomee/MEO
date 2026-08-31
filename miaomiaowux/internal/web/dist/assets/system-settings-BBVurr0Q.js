import {
  v as _0x37d641,
  r as _0x46e13c,
  j as _0x2ddfe2,
  t as _0x54e43f,
  q as _0x29acb4,
  a4 as _0x144b3f,
  ay as _0x32387d,
  a5 as _0x4e8a70,
  ao as _0x14801b,
  az as _0x18725b,
  aA as _0xa3252a,
  aB as _0x38a01f,
  aC as _0x18a4f9,
  aD as _0x4fa2aa,
  H as _0x195e80,
  u as _0x2443c3,
  aE as _0x1e002a,
  Q as _0x51d33c,
  J as _0x6a84da,
  aF as _0x1618b9,
  aG as _0x18ddc9,
  aH as _0x2d8607,
  aI as _0x1814ad,
  av as _0x3061be,
  aJ as _0x157104,
  aK as _0x5749cc,
  aL as _0x4701f1,
  aM as _0x555bf5,
  aN as _0x3039d8,
  aO as _0x56abb7,
  aP as _0x28e2d0,
} from "./vendor-modules-0UUaSA6d.js";
import {
  a as _0x495bb8,
  h as _0x4deb0e,
  L as _0x34df34,
  B as _0x5185a8,
  I as _0x549353,
  C as _0x1e06f2,
  D as _0x1cc136,
  k as _0x3f9d9d,
  d as _0x3e8fab,
  e as _0x333ff1,
  f as _0x4b37d4,
  g as _0x1e6ae9,
  i as _0x5881af,
  u as _0x235baf,
  l as _0x7633f5,
} from "./index-0oY9qUmNNK.js";
import {
  g as _0xe15cc6,
  c as _0x3f64f4,
  f as _0x7c97c,
} from "./country-flag-CF94k7yy.js";
import {
  a as _0x26349a,
  c as _0x24539f,
  u as _0x3910d6,
} from "./use-capabilities.js";
import { u as _0xd1c38e } from "./use-proxy-groups-DnqyXqIF.js";
import { u as _0x472a2b } from "./useTurnstile-Bm81HBF_.js";
import {
  C as _0x2aeb40,
  a as _0x1db8ce,
  b as _0x30c50e,
  c as _0x54661d,
  d as _0x42cb32,
} from "./card-cr-ClpdW.js";
import { C as _0x77832a } from "./checkbox-C1SgKqW3.js";
import { R as _0x2db305, a as _0x356e54 } from "./radio-group-B32CHMgi.js";
import {
  S as _0x3ba40b,
  a as _0x3d4a68,
  b as _0x1929c7,
  c as _0x51ae8f,
  d as _0x1260f9,
} from "./select-CHQxHHZU.js";
import { S as _0x543b01 } from "./switch-BmK7fVZr.js";
import {
  T as _0x2b16ce,
  a as _0x383cce,
  b as _0x3102ee,
  c as _0x550af3,
} from "./tabs-DlVWTTSI.js";
import { T as _0x5d9440 } from "./textarea-BJ6734aT.js";
import {
  a as _0x30391b,
  b as _0x2b1216,
  c as _0x42f438,
} from "./tooltip-BGx-1utF.js";
import { F as _0x3e16bd } from "./flag-emoji-picker-j7cGeacV.js";
import { T as _0x5d8e9a } from "./topbar-C6fNrllX.js";
import { L as _0x4c44aa } from "./minimal-badge.js";
import { P as _0x23ae0c } from "./feature-wrapper.js";
import { P as _0x4c3ba9 } from "./progress-xgHFkzjD.js";
import "./popover-DFYIlHRV.js";
import "./twemoji-Dd73OA7b.js";
import "./dropdown-menu-BpG341YX.js";
const ps = { unicom: "联通", mobile: "移动", telecom: "电信" },
  Nn = 0x1e;
function Cr(_0x2536bf) {
  return _0x37d641({
    queryKey: ["probe-cdn-regions"],
    queryFn: async () =>
      (await _0x495bb8["get"]("/api/admin/probe/regions"))["data"],
    enabled: _0x2536bf,
    staleTime: 0xe10 * 0x3e8,
  });
}
function wn({
  regions: _0x37b9de,
  source: _0x45e0e6,
  selected: _0x52037e,
  onChange: _0x1a4d5e,
  max: _0x1def2c = Nn,
  emptyHint: _0x406ba7,
}) {
  const [_0x537007, _0x1a786f] = _0x46e13c["useState"](!0x1),
    [_0x2e1373, _0x236bda] = _0x46e13c["useState"](!0x1),
    [_0x2715dc, _0x473077] = _0x46e13c["useState"](""),
    [_0xd5cac6, _0x5809bf] = _0x46e13c["useState"](""),
    [_0x245db0, _0x4883c4] = _0x46e13c["useState"]("443"),
    [_0x3b0a23, _0x50fcad] = _0x46e13c["useState"]("tcp"),
    _0x441359 = new Set(_0x52037e["map"]((_0x825a84) => _0x825a84["key"])),
    _0x21186e = (_0x32c168, _0x421777) => {
      if (_0x421777) {
        if (_0x441359["has"](_0x32c168["key"])) return;
        if (_0x52037e["length"] >= _0x1def2c) {
          _0x54e43f["error"]("最多选\x20" + _0x1def2c + "\x20个目标");
          return;
        }
        _0x1a4d5e([..._0x52037e, _0x32c168]);
      } else
        _0x1a4d5e(
          _0x52037e["filter"](
            (_0x57b040) => _0x57b040["key"] !== _0x32c168["key"],
          ),
        );
    },
    _0x4f8d25 = _0x52037e["filter"]((_0x12df94) =>
      _0x12df94["key"]["startsWith"]("custom-"),
    ),
    _0x145301 = (_0x1e9a27) =>
      (_0x1e9a27["match"](/:/g)?.["length"] ?? 0x0) >= 0x2 &&
      /^[0-9a-fA-F:.]+$/["test"](_0x1e9a27),
    _0x382bb3 = () => {
      const _0x3e4005 = _0xd5cac6["trim"](),
        _0x262e62 = _0x2715dc["trim"]() || _0x3e4005;
      if (!_0x3e4005) {
        _0x54e43f["error"]("请填写域名或\x20IP");
        return;
      }
      if (_0x3e4005["includes"]("/")) {
        _0x54e43f["error"]("地址不要带协议或路径");
        return;
      }
      if (_0x3e4005["includes"](":") && !_0x145301(_0x3e4005)) {
        _0x54e43f["error"]("地址不要带端口,IPv6\x20请直接填地址(不加方括号)");
        return;
      }
      const _0x2d5ef7 = _0x3b0a23 === "tcp" ? Number(_0x245db0) : 0x0;
      if (
        _0x3b0a23 === "tcp" &&
        (!Number["isInteger"](_0x2d5ef7) ||
          _0x2d5ef7 < 0x1 ||
          _0x2d5ef7 > 0xffff)
      ) {
        _0x54e43f["error"]("端口无效");
        return;
      }
      if (_0x52037e["length"] >= _0x1def2c) {
        _0x54e43f["error"]("最多选\x20" + _0x1def2c + "\x20个目标");
        return;
      }
      const _0x43d401 = ("custom-" + _0x3e4005 + "-" + _0x2d5ef7)["replace"](
        /[^a-zA-Z0-9.-]/g,
        "_",
      );
      if (_0x441359["has"](_0x43d401)) {
        _0x54e43f["error"]("该目标已添加");
        return;
      }
      (_0x1a4d5e([
        ..._0x52037e,
        {
          key: _0x43d401,
          label: _0x262e62,
          isp: "custom",
          host: _0x3e4005,
          port: _0x2d5ef7,
          type: _0x3b0a23,
        },
      ]),
        _0x473077(""),
        _0x5809bf(""),
        _0x4883c4("443"));
    },
    _0xc8901b = [];
  for (const _0x434f63 of _0x37b9de?.["international"] ?? []) {
    const _0x429197 = _0xc8901b["find"](
      ([_0x14cae1]) => _0x14cae1 === _0x434f63["group"],
    );
    _0x429197
      ? _0x429197[0x1]["push"](_0x434f63)
      : _0xc8901b["push"]([_0x434f63["group"], [_0x434f63]]);
  }
  return _0x37b9de
    ? _0x2ddfe2["jsxs"]("div", {
        className: "space-y-2",
        children: [
          _0x2ddfe2["jsxs"]("div", {
            className:
              "text-muted-foreground\x20flex\x20items-center\x20gap-2\x20text-xs",
            children: [
              _0x2ddfe2["jsxs"]("span", {
                children: [
                  "Ping\x20目标(已选\x20",
                  _0x52037e["length"],
                  "/",
                  _0x1def2c,
                  ")",
                ],
              }),
              _0x45e0e6 === "embedded" &&
                _0x2ddfe2["jsx"]("span", {
                  className: "text-yellow-600",
                  children: "·\x20使用内置快照(CDN\x20不可达)",
                }),
            ],
          }),
          _0x52037e["length"] === 0x0 && _0x406ba7,
          _0xc8901b["length"] > 0x0 &&
            _0x2ddfe2["jsx"]("div", {
              className: "divide-y\x20rounded-md\x20border",
              children: _0xc8901b["map"](([_0x438eaa, _0x29fd5e]) =>
                _0x2ddfe2["jsxs"](
                  "div",
                  {
                    className: "px-3\x20py-2",
                    children: [
                      _0x2ddfe2["jsx"]("div", {
                        className:
                          "text-muted-foreground\x20mb-1.5\x20text-xs\x20font-medium",
                        children: _0x438eaa,
                      }),
                      _0x2ddfe2["jsx"]("div", {
                        className: "flex\x20flex-wrap\x20gap-x-4\x20gap-y-1.5",
                        children: _0x29fd5e["map"]((_0x177a33) => {
                          const _0x27100d = {
                            key: _0x177a33["key"],
                            label: _0x177a33["label"],
                            isp: "intl",
                            host: _0x177a33["host"],
                            port: _0x177a33["port"],
                            type: _0x177a33["type"],
                          };
                          return _0x2ddfe2["jsxs"](
                            "label",
                            {
                              className:
                                "flex\x20cursor-pointer\x20items-center\x20gap-1.5\x20text-xs",
                              children: [
                                _0x2ddfe2["jsx"](_0x77832a, {
                                  checked: _0x441359["has"](_0x177a33["key"]),
                                  onCheckedChange: (_0x7e59a4) =>
                                    _0x21186e(_0x27100d, _0x7e59a4 === !0x0),
                                }),
                                _0x2ddfe2["jsx"]("span", {
                                  children: _0x177a33["label"],
                                }),
                                _0x2ddfe2["jsx"]("span", {
                                  className:
                                    "text-muted-foreground\x20text-[10px]",
                                  children:
                                    _0x177a33["type"] === "icmp"
                                      ? "ICMP"
                                      : ":" + _0x177a33["port"],
                                }),
                              ],
                            },
                            _0x177a33["key"],
                          );
                        }),
                      }),
                    ],
                  },
                  _0x438eaa,
                ),
              ),
            }),
          _0x2ddfe2["jsxs"]("div", {
            className: "space-y-2\x20rounded-md\x20border\x20px-3\x20py-2",
            children: [
              _0x2ddfe2["jsxs"]("div", {
                className: "flex\x20items-center\x20justify-between",
                children: [
                  _0x2ddfe2["jsx"]("span", {
                    className:
                      "text-muted-foreground\x20text-xs\x20font-medium",
                    children: "自定义目标",
                  }),
                  _0x2ddfe2["jsx"]("button", {
                    type: "button",
                    className: "text-primary\x20text-xs\x20hover:underline",
                    onClick: () => _0x236bda((_0x2407d1) => !_0x2407d1),
                    children: _0x2e1373 ? "收起" : "添加",
                  }),
                ],
              }),
              _0x4f8d25["length"] > 0x0 &&
                _0x2ddfe2["jsx"]("div", {
                  className: "flex\x20flex-wrap\x20gap-1.5",
                  children: _0x4f8d25["map"]((_0x51f8bc) =>
                    _0x2ddfe2["jsxs"](
                      "span",
                      {
                        className:
                          "inline-flex\x20items-center\x20gap-1\x20rounded-md\x20border\x20px-2\x20py-0.5\x20text-xs",
                        children: [
                          _0x51f8bc["label"],
                          _0x2ddfe2["jsxs"]("span", {
                            className: "text-muted-foreground\x20text-[10px]",
                            children: [
                              _0x51f8bc["host"],
                              _0x51f8bc["type"] === "icmp"
                                ? "\x20·\x20ICMP"
                                : ":" + _0x51f8bc["port"],
                            ],
                          }),
                          _0x2ddfe2["jsx"]("button", {
                            type: "button",
                            className:
                              "text-muted-foreground\x20hover:text-destructive",
                            onClick: () => _0x21186e(_0x51f8bc, !0x1),
                            children: "×",
                          }),
                        ],
                      },
                      _0x51f8bc["key"],
                    ),
                  ),
                }),
              _0x2e1373 &&
                _0x2ddfe2["jsxs"]("div", {
                  className: "space-y-2",
                  children: [
                    _0x2ddfe2["jsxs"]("div", {
                      className: "flex\x20flex-wrap\x20gap-2",
                      children: [
                        _0x2ddfe2["jsx"]("input", {
                          className:
                            "h-8\x20min-w-[120px]\x20flex-1\x20rounded-md\x20border\x20bg-transparent\x20px-2\x20text-xs",
                          placeholder: "名称,如\x20我的服务器",
                          value: _0x2715dc,
                          onChange: (_0xe22114) =>
                            _0x473077(_0xe22114["target"]["value"]),
                        }),
                        _0x2ddfe2["jsx"]("input", {
                          className:
                            "h-8\x20min-w-[140px]\x20flex-1\x20rounded-md\x20border\x20bg-transparent\x20px-2\x20text-xs",
                          placeholder: "域名或\x20IP(不带端口)",
                          value: _0xd5cac6,
                          onChange: (_0x14af87) =>
                            _0x5809bf(_0x14af87["target"]["value"]),
                        }),
                        _0x2ddfe2["jsxs"]("select", {
                          className:
                            "h-8\x20rounded-md\x20border\x20bg-transparent\x20px-2\x20text-xs",
                          value: _0x3b0a23,
                          onChange: (_0x31af47) =>
                            _0x50fcad(_0x31af47["target"]["value"]),
                          children: [
                            _0x2ddfe2["jsx"]("option", {
                              value: "tcp",
                              children: "TCP",
                            }),
                            _0x2ddfe2["jsx"]("option", {
                              value: "icmp",
                              children: "ICMP",
                            }),
                          ],
                        }),
                        _0x3b0a23 === "tcp" &&
                          _0x2ddfe2["jsx"]("input", {
                            className:
                              "h-8\x20w-20\x20rounded-md\x20border\x20bg-transparent\x20px-2\x20text-xs",
                            placeholder: "端口",
                            value: _0x245db0,
                            onChange: (_0x452262) =>
                              _0x4883c4(_0x452262["target"]["value"]),
                          }),
                        _0x2ddfe2["jsx"]("button", {
                          type: "button",
                          className:
                            "hover:bg-accent\x20h-8\x20rounded-md\x20border\x20px-3\x20text-xs",
                          onClick: _0x382bb3,
                          children: "添加",
                        }),
                      ],
                    }),
                    _0x2ddfe2["jsx"]("p", {
                      className: "text-muted-foreground\x20text-xs",
                      children:
                        "不支持内网/环回地址\x20——\x20所有\x20agent\x20都会去拨测这个地址,指向内网等于让面板变成内网扫描器",
                    }),
                  ],
                }),
            ],
          }),
          _0x2ddfe2["jsx"]("div", {
            className:
              "max-h-64\x20divide-y\x20overflow-y-auto\x20rounded-md\x20border",
            children: _0x37b9de["provinces"]["map"]((_0x126b82) =>
              _0x2ddfe2["jsx"](
                "div",
                {
                  className: "px-3\x20py-2",
                  children: _0x2ddfe2["jsxs"]("div", {
                    className: "flex\x20items-center\x20gap-2\x20text-sm",
                    children: [
                      _0x2ddfe2["jsx"]("span", {
                        className: "w-16\x20shrink-0\x20font-medium",
                        children: _0x126b82["province"],
                      }),
                      _0x2ddfe2["jsx"]("div", {
                        className: "flex\x20flex-wrap\x20gap-x-3\x20gap-y-1",
                        children: _0x126b82["targets"]["map"]((_0x153d38) => {
                          const _0xaf82d6 = {
                            key: _0x153d38["key"],
                            label:
                              "" +
                              _0x126b82["province"] +
                              (ps[_0x153d38["isp"]] ?? _0x153d38["isp"]),
                            isp: _0x153d38["isp"],
                            host: _0x153d38["host"],
                            port: _0x153d38["port"],
                          };
                          return _0x2ddfe2["jsxs"](
                            "label",
                            {
                              className:
                                "flex\x20cursor-pointer\x20items-center\x20gap-1\x20text-xs",
                              children: [
                                _0x2ddfe2["jsx"](_0x77832a, {
                                  checked: _0x441359["has"](_0x153d38["key"]),
                                  onCheckedChange: (_0x5d8378) =>
                                    _0x21186e(_0xaf82d6, _0x5d8378 === !0x0),
                                }),
                                _0x2ddfe2["jsx"]("span", {
                                  children:
                                    ps[_0x153d38["isp"]] ?? _0x153d38["isp"],
                                }),
                              ],
                            },
                            _0x153d38["key"],
                          );
                        }),
                      }),
                    ],
                  }),
                },
                _0x126b82["province"],
              ),
            ),
          }),
          _0x2ddfe2["jsx"]("button", {
            type: "button",
            className: "text-primary\x20text-xs\x20hover:underline",
            onClick: () => _0x1a786f((_0x201a51) => !_0x201a51),
            children: _0x537007
              ? "收起"
              : "展开市级节点(" + _0x37b9de["cities"]["length"] + ")",
          }),
          _0x537007 &&
            _0x2ddfe2["jsx"]("div", {
              className:
                "max-h-64\x20divide-y\x20overflow-y-auto\x20rounded-md\x20border",
              children: _0x37b9de["cities"]["map"]((_0x597cfe) => {
                const _0x18efb = {
                  key: _0x597cfe["key"],
                  label:
                    "" +
                    _0x597cfe["label"] +
                    (ps[_0x597cfe["isp"]] ?? _0x597cfe["isp"]),
                  isp: _0x597cfe["isp"],
                  host: _0x597cfe["host"],
                  port: _0x597cfe["port"],
                };
                return _0x2ddfe2["jsxs"](
                  "label",
                  {
                    className:
                      "hover:bg-accent/40\x20flex\x20cursor-pointer\x20items-center\x20gap-2\x20px-3\x20py-1.5\x20text-xs",
                    children: [
                      _0x2ddfe2["jsx"](_0x77832a, {
                        checked: _0x441359["has"](_0x597cfe["key"]),
                        onCheckedChange: (_0x160fab) =>
                          _0x21186e(_0x18efb, _0x160fab === !0x0),
                      }),
                      _0x2ddfe2["jsxs"]("span", {
                        className: "truncate",
                        children: [
                          _0x597cfe["label"],
                          "\x20·\x20",
                          ps[_0x597cfe["isp"]] ?? _0x597cfe["isp"],
                        ],
                      }),
                    ],
                  },
                  _0x597cfe["key"],
                );
              }),
            }),
        ],
      })
    : _0x2ddfe2["jsx"]("div", {
        className: "text-muted-foreground\x20text-xs",
        children: "加载省市列表中…",
      });
}
const pn = [
  { key: "node_blocked", label: "节点被墙(自动)", auto: !0x0 },
  { key: "node_recovered", label: "节点恢复(自动)", auto: !0x0 },
  { key: "maintenance", label: "系统维护", auto: !0x1 },
  { key: "sub_update", label: "订阅更新提醒", auto: !0x1 },
  { key: "general", label: "通用公告", auto: !0x1 },
];
function Sr() {
  const _0x45d617 = _0x29acb4(),
    { data: _0x1d8907 } = _0x37d641({
      queryKey: ["announcement-config"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/announcements"))[
          "data"
        ],
      staleTime: 0x12c * 0x3e8,
    }),
    [_0x3a4bc3, _0x225230] = _0x46e13c["useState"]({}),
    [_0x236747, _0x280f88] = _0x46e13c["useState"]([]),
    [_0x132f02, _0x1471cb] = _0x46e13c["useState"](!0x1),
    _0x21dc9d = !!_0x1d8907?.["official_probe_available"];
  _0x46e13c["useEffect"](() => {
    (_0x1d8907?.["config"]?.["types"] &&
      _0x225230(_0x1d8907["config"]["types"]),
      Array["isArray"](_0x1d8907?.["probe_tester_ids"]) &&
        _0x280f88(_0x1d8907["probe_tester_ids"]),
      _0x1471cb(!!_0x1d8907?.["official_probe"]));
  }, [_0x1d8907]);
  const { data: _0x715d4c } = _0x37d641({
      queryKey: ["speed-testers-for-announce"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/speedtest/testers"))["data"],
      staleTime: 0xea60,
    }),
    _0x569aa6 = _0x715d4c?.["testers"] || [],
    _0x59d0b2 = (_0x10ed3f) =>
      _0x280f88((_0x4b9a71) =>
        _0x4b9a71["includes"](_0x10ed3f)
          ? _0x4b9a71["filter"]((_0x2dc196) => _0x2dc196 !== _0x10ed3f)
          : [..._0x4b9a71, _0x10ed3f],
      ),
    _0x288923 = _0x144b3f({
      mutationFn: async (_0x15f4a0) => {
        await _0x495bb8["put"]("/api/admin/system-settings/announcements", {
          config: { types: _0x15f4a0 },
          probe_tester_ids: _0x236747,
          official_probe: _0x132f02,
        });
      },
      onSuccess: () => {
        (_0x45d617["invalidateQueries"]({ queryKey: ["announcement-config"] }),
          _0x54e43f["success"]("公告配置已保存"));
      },
      onError: _0x4deb0e,
    }),
    _0x10e58a = (_0x90197c, _0x4ccc0c) =>
      _0x225230((_0x2cbb89) => ({
        ..._0x2cbb89,
        [_0x90197c]: { ..._0x2cbb89[_0x90197c], ..._0x4ccc0c },
      })),
    { data: _0x2532cb } = _0x37d641({
      queryKey: ["announcements-active"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/announcements"))["data"],
    }),
    _0x14e00e = _0x2532cb?.["announcements"] || [],
    _0x1e453a = _0x144b3f({
      mutationFn: async (_0x58460d) => {
        await _0x495bb8["delete"]("/api/admin/announcements?id=" + _0x58460d);
      },
      onSuccess: () => {
        (_0x45d617["invalidateQueries"]({ queryKey: ["announcements-active"] }),
          _0x54e43f["success"]("已删除"));
      },
      onError: _0x4deb0e,
    }),
    [_0x450d42, _0x4c09bb] = _0x46e13c["useState"]("general"),
    [_0x21c4c2, _0x11c4ab] = _0x46e13c["useState"](""),
    [_0xc27c71, _0x3731af] = _0x46e13c["useState"](""),
    [_0x59fc7a, _0x1bf270] = _0x46e13c["useState"](""),
    [_0x118d10, _0x206b3f] = _0x46e13c["useState"](!0x0),
    [_0x467a8c, _0x5c4494] = _0x46e13c["useState"](!0x0),
    _0xf8e82a = _0x144b3f({
      mutationFn: async () => {
        await _0x495bb8["post"]("/api/admin/announcements", {
          type: _0x450d42,
          title: _0x21c4c2["trim"](),
          body: _0xc27c71["trim"](),
          expires_minutes: (parseInt(_0x59fc7a, 0xa) || 0x0) * 0x3c,
          via_bot: _0x118d10,
          via_miniapp: _0x467a8c,
        });
      },
      onSuccess: () => {
        (_0x45d617["invalidateQueries"]({ queryKey: ["announcements-active"] }),
          _0x54e43f["success"]("公告已发布(bot\x20将在\x201\x20分钟内广播)"),
          _0x11c4ab(""),
          _0x3731af(""),
          _0x1bf270(""));
      },
      onError: _0x4deb0e,
    }),
    _0xb38c53 = (_0x34aeb8) => {
      _0x4c09bb(_0x34aeb8);
      const _0x15d548 = _0x3a4bc3[_0x34aeb8];
      _0x15d548 &&
        (_0x11c4ab(_0x15d548["title"] || ""),
        _0x3731af(_0x15d548["template"] || ""));
    };
  return _0x2ddfe2["jsxs"](_0x2aeb40, {
    children: [
      _0x2ddfe2["jsxs"](_0x1db8ce, {
        className: "pb-4",
        children: [
          _0x2ddfe2["jsxs"](_0x30c50e, {
            className: "flex\x20items-center\x20gap-2",
            children: [
              _0x2ddfe2["jsx"](_0x32387d, { className: "h-5\x20w-5" }),
              "公告",
            ],
          }),
          _0x2ddfe2["jsx"](_0x54661d, {
            children:
              "配置各类公告的文案与推送渠道;可手动发布公告,广播给所有绑定\x20Telegram\x20的用户并显示在\x20Mini\x20App。",
          }),
        ],
      }),
      _0x2ddfe2["jsxs"](_0x42cb32, {
        className: "space-y-4",
        children: [
          _0x2ddfe2["jsxs"]("div", {
            className: "space-y-3",
            children: [
              pn["map"]((_0x3e0a33) => {
                const _0x3e1759 = _0x3a4bc3[_0x3e0a33["key"]] || {
                  enabled: !0x0,
                  template: "",
                  via_bot: !0x0,
                  via_miniapp: !0x0,
                };
                return _0x2ddfe2["jsxs"](
                  "div",
                  {
                    className: "space-y-2\x20rounded-lg\x20border\x20p-3",
                    children: [
                      _0x2ddfe2["jsxs"]("div", {
                        className: "flex\x20items-center\x20justify-between",
                        children: [
                          _0x2ddfe2["jsxs"]("div", {
                            className:
                              "flex\x20items-center\x20gap-2\x20text-sm\x20font-medium",
                            children: [
                              _0x3e0a33["label"],
                              _0x3e0a33["auto"] &&
                                _0x2ddfe2["jsx"]("span", {
                                  className: "text-muted-foreground\x20text-xs",
                                  children: "·\x20探测到自动触发",
                                }),
                            ],
                          }),
                          _0x2ddfe2["jsx"](_0x543b01, {
                            checked: _0x3e1759["enabled"],
                            onCheckedChange: (_0x1eb4ac) =>
                              _0x10e58a(_0x3e0a33["key"], {
                                enabled: _0x1eb4ac,
                              }),
                          }),
                        ],
                      }),
                      _0x3e1759["enabled"] &&
                        _0x2ddfe2["jsxs"](_0x2ddfe2["Fragment"], {
                          children: [
                            _0x3e0a33["key"] !== "general" &&
                              _0x2ddfe2["jsx"](_0x5d9440, {
                                rows: 0x2,
                                value: _0x3e1759["template"],
                                placeholder: _0x3e0a33["auto"]
                                  ? "文案,可用占位符\x20{node}\x20{time}"
                                  : "默认文案",
                                onChange: (_0x1dd5c1) =>
                                  _0x10e58a(_0x3e0a33["key"], {
                                    template: _0x1dd5c1["target"]["value"],
                                  }),
                              }),
                            _0x2ddfe2["jsxs"]("div", {
                              className:
                                "flex\x20flex-wrap\x20items-center\x20gap-4\x20text-sm",
                              children: [
                                _0x2ddfe2["jsxs"]("label", {
                                  className: "flex\x20items-center\x20gap-2",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x543b01, {
                                      checked: _0x3e1759["via_bot"],
                                      onCheckedChange: (_0x117ac7) =>
                                        _0x10e58a(_0x3e0a33["key"], {
                                          via_bot: _0x117ac7,
                                        }),
                                    }),
                                    "发送\x20bot\x20消息",
                                  ],
                                }),
                                _0x2ddfe2["jsxs"]("label", {
                                  className: "flex\x20items-center\x20gap-2",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x543b01, {
                                      checked: _0x3e1759["via_miniapp"],
                                      onCheckedChange: (_0x7e91d5) =>
                                        _0x10e58a(_0x3e0a33["key"], {
                                          via_miniapp: _0x7e91d5,
                                        }),
                                    }),
                                    "显示\x20Mini\x20App",
                                  ],
                                }),
                              ],
                            }),
                          ],
                        }),
                    ],
                  },
                  _0x3e0a33["key"],
                );
              }),
              _0x2ddfe2["jsxs"]("div", {
                className: "space-y-2\x20rounded-lg\x20border\x20p-3",
                children: [
                  _0x2ddfe2["jsx"]("div", {
                    className: "text-sm\x20font-medium",
                    children: "被墙探测源",
                  }),
                  _0x2ddfe2["jsxs"]("p", {
                    className: "text-muted-foreground\x20text-xs",
                    children: [
                      "从选中的",
                      _0x2ddfe2["jsx"]("span", {
                        className: "font-medium",
                        children: "家用测速端",
                      }),
                      "视角\x20TCP\x20拨测各节点判断是否被墙。测速端部署在国内家庭网络,才能反映真实可达性——",
                      _0x2ddfe2["jsx"]("span", {
                        className: "text-amber-600\x20dark:text-amber-400",
                        children: "主控和机房\x20agent\x20都探不准",
                      }),
                      "(机房能连的节点国内未必能连)。测速端与下方官方探测都不启用时,不做被墙探测。",
                    ],
                  }),
                  _0x569aa6["length"] === 0x0
                    ? _0x2ddfe2["jsx"]("div", {
                        className: "text-muted-foreground\x20text-xs",
                        children:
                          "暂无测速端。请先在「节点测速」里添加家用测速端。",
                      })
                    : _0x2ddfe2["jsx"]("div", {
                        className: "flex\x20flex-wrap\x20gap-2",
                        children: _0x569aa6["map"]((_0x2bcbe9) => {
                          const _0x5e8c96 =
                            Array["isArray"](_0x2bcbe9["caps"]) &&
                            _0x2bcbe9["caps"]["includes"]("probe");
                          return _0x2ddfe2["jsxs"](
                            "button",
                            {
                              type: "button",
                              disabled: !_0x5e8c96,
                              title: _0x5e8c96
                                ? void 0x0
                                : "该测速端版本过旧,不支持可达性探测,请升级",
                              onClick: () => _0x59d0b2(_0x2bcbe9["id"]),
                              className:
                                "rounded-md\x20border\x20px-2\x20py-1\x20text-xs\x20transition-colors\x20" +
                                (_0x5e8c96
                                  ? _0x236747["includes"](_0x2bcbe9["id"])
                                    ? "border-primary\x20bg-primary/10\x20text-primary"
                                    : "bg-card\x20hover:bg-accent/50"
                                  : "cursor-not-allowed\x20opacity-50"),
                              children: [
                                _0x2bcbe9["name"],
                                !_0x5e8c96 &&
                                  _0x2ddfe2["jsx"]("span", {
                                    className: "text-muted-foreground\x20ml-1",
                                    children: "(需升级)",
                                  }),
                              ],
                            },
                            _0x2bcbe9["id"],
                          );
                        }),
                      }),
                  _0x2ddfe2["jsxs"]("div", {
                    className:
                      "flex\x20items-center\x20justify-between\x20gap-3\x20border-t\x20pt-2",
                    children: [
                      _0x2ddfe2["jsxs"]("div", {
                        className: "space-y-0.5",
                        children: [
                          _0x2ddfe2["jsx"](_0x34df34, {
                            className: "text-xs",
                            children: "官方探测",
                          }),
                          _0x2ddfe2["jsx"]("p", {
                            className: "text-muted-foreground\x20text-xs",
                            children: _0x21dc9d
                              ? "额外使用官方部署在国内的探测端,与上面选中的测速端结果取并集(任一可达即可达)。没有自建测速端时也可单独使用。"
                              : "请先配置可用的远程测速端。",
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsx"](_0x543b01, {
                        checked: _0x132f02,
                        disabled: !_0x21dc9d,
                        onCheckedChange: _0x1471cb,
                      }),
                    ],
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"](_0x5185a8, {
                size: "sm",
                disabled: _0x288923["isPending"],
                onClick: () => _0x288923["mutate"](_0x3a4bc3),
                children: [
                  _0x288923["isPending"] &&
                    _0x2ddfe2["jsx"](_0x4e8a70, {
                      className: "mr-2\x20h-4\x20w-4\x20animate-spin",
                    }),
                  "保存配置",
                ],
              }),
            ],
          }),
          _0x2ddfe2["jsxs"]("div", {
            className: "space-y-3\x20rounded-lg\x20border\x20p-3",
            children: [
              _0x2ddfe2["jsx"]("div", {
                className: "text-sm\x20font-medium",
                children: "立即发布公告",
              }),
              _0x2ddfe2["jsxs"]("div", {
                className: "grid\x20grid-cols-2\x20gap-2",
                children: [
                  _0x2ddfe2["jsxs"]("div", {
                    className: "space-y-1",
                    children: [
                      _0x2ddfe2["jsx"](_0x34df34, {
                        className: "text-xs",
                        children: "类型",
                      }),
                      _0x2ddfe2["jsxs"](_0x3ba40b, {
                        value: _0x450d42,
                        onValueChange: _0xb38c53,
                        children: [
                          _0x2ddfe2["jsx"](_0x3d4a68, {
                            children: _0x2ddfe2["jsx"](_0x1929c7, {}),
                          }),
                          _0x2ddfe2["jsx"](_0x51ae8f, {
                            children: pn["filter"](
                              (_0x3da4f1) => !_0x3da4f1["auto"],
                            )["map"]((_0x57a15e) =>
                              _0x2ddfe2["jsx"](
                                _0x1260f9,
                                {
                                  value: _0x57a15e["key"],
                                  children: _0x57a15e["label"],
                                },
                                _0x57a15e["key"],
                              ),
                            ),
                          }),
                        ],
                      }),
                    ],
                  }),
                  _0x2ddfe2["jsxs"]("div", {
                    className: "space-y-1",
                    children: [
                      _0x2ddfe2["jsx"](_0x34df34, {
                        className: "text-xs",
                        children: "生效时长(小时,0=永久)",
                      }),
                      _0x2ddfe2["jsx"](_0x549353, {
                        type: "number",
                        value: _0x59fc7a,
                        onChange: (_0x51f842) =>
                          _0x1bf270(_0x51f842["target"]["value"]),
                        placeholder: "0",
                      }),
                    ],
                  }),
                ],
              }),
              _0x2ddfe2["jsx"](_0x549353, {
                value: _0x21c4c2,
                onChange: (_0x506ac7) =>
                  _0x11c4ab(_0x506ac7["target"]["value"]),
                placeholder: "标题(可空)",
              }),
              _0x2ddfe2["jsx"](_0x5d9440, {
                rows: 0x3,
                value: _0xc27c71,
                onChange: (_0x3d2f07) =>
                  _0x3731af(_0x3d2f07["target"]["value"]),
                placeholder: "公告正文",
              }),
              _0x2ddfe2["jsxs"]("div", {
                className:
                  "flex\x20flex-wrap\x20items-center\x20gap-4\x20text-sm",
                children: [
                  _0x2ddfe2["jsxs"]("label", {
                    className: "flex\x20items-center\x20gap-2",
                    children: [
                      _0x2ddfe2["jsx"](_0x543b01, {
                        checked: _0x118d10,
                        onCheckedChange: _0x206b3f,
                      }),
                      "\x20发送\x20bot\x20消息",
                    ],
                  }),
                  _0x2ddfe2["jsxs"]("label", {
                    className: "flex\x20items-center\x20gap-2",
                    children: [
                      _0x2ddfe2["jsx"](_0x543b01, {
                        checked: _0x467a8c,
                        onCheckedChange: _0x5c4494,
                      }),
                      "\x20",
                      "显示\x20Mini\x20App",
                    ],
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"](_0x5185a8, {
                size: "sm",
                disabled: _0xf8e82a["isPending"] || !_0xc27c71["trim"](),
                onClick: () => _0xf8e82a["mutate"](),
                children: [
                  _0xf8e82a["isPending"]
                    ? _0x2ddfe2["jsx"](_0x4e8a70, {
                        className: "mr-2\x20h-4\x20w-4\x20animate-spin",
                      })
                    : _0x2ddfe2["jsx"](_0x14801b, {
                        className: "mr-2\x20h-4\x20w-4",
                      }),
                  "发布",
                ],
              }),
            ],
          }),
          _0x14e00e["length"] > 0x0 &&
            _0x2ddfe2["jsxs"]("div", {
              className: "space-y-2",
              children: [
                _0x2ddfe2["jsxs"]("div", {
                  className: "text-sm\x20font-medium",
                  children: ["当前生效公告(", _0x14e00e["length"], ")"],
                }),
                _0x14e00e["map"]((_0x44abdf) =>
                  _0x2ddfe2["jsxs"](
                    "div",
                    {
                      className:
                        "flex\x20items-start\x20justify-between\x20gap-2\x20rounded-md\x20border\x20p-2\x20text-sm",
                      children: [
                        _0x2ddfe2["jsxs"]("div", {
                          className: "min-w-0\x20flex-1",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className: "font-medium",
                              children: [
                                _0x44abdf["title"] || "公告",
                                "\x20",
                                _0x2ddfe2["jsxs"]("span", {
                                  className: "text-muted-foreground\x20text-xs",
                                  children: ["·\x20", _0x44abdf["type"]],
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsx"]("div", {
                              className:
                                "text-muted-foreground\x20text-xs\x20break-words\x20whitespace-pre-wrap",
                              children: _0x44abdf["body"],
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsx"](_0x5185a8, {
                          variant: "ghost",
                          size: "icon",
                          className: "text-destructive\x20shrink-0",
                          onClick: () => _0x1e453a["mutate"](_0x44abdf["id"]),
                          children: _0x2ddfe2["jsx"](_0x18725b, {
                            className: "h-4\x20w-4",
                          }),
                        }),
                      ],
                    },
                    _0x44abdf["id"],
                  ),
                ),
              ],
            }),
        ],
      }),
    ],
  });
}
function Pr() {
  const _0x50e643 = _0x29acb4(),
    [_0x23baf3, _0x4eb364] = _0x46e13c["useState"](""),
    [_0x395c54, _0x58c2f6] = _0x46e13c["useState"](""),
    [_0x14a01f, _0x18581a] = _0x46e13c["useState"](""),
    [_0x4ea065, _0x346413] = _0x46e13c["useState"](""),
    _0x131559 = _0x46e13c["useRef"](null),
    _0x919ca4 = _0x46e13c["useRef"](null),
    { data: _0x1c8a14, isLoading: _0x38cd27 } = _0x37d641({
      queryKey: ["admin-branding"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/branding"))["data"],
    });
  _0x46e13c["useEffect"](() => {
    _0x1c8a14?.["branding"] &&
      (_0x4eb364(_0x1c8a14["branding"]["site_title"] || ""),
      _0x58c2f6(_0x1c8a14["branding"]["brand_title"] || ""),
      _0x18581a(_0x1c8a14["branding"]["logo_url"] || ""),
      _0x346413(_0x1c8a14["branding"]["icon_url"] || ""));
  }, [_0x1c8a14]);
  const _0x1d78aa = () => {
      (_0x50e643["invalidateQueries"]({ queryKey: ["admin-branding"] }),
        _0x50e643["invalidateQueries"]({ queryKey: ["branding"] }));
    },
    _0x5194c5 = _0x144b3f({
      mutationFn: async () =>
        _0x495bb8["post"]("/api/admin/system-settings/branding", {
          site_title: _0x23baf3,
          brand_title: _0x395c54,
          logo_url: _0x14a01f,
          icon_url: _0x4ea065,
        }),
      onSuccess: () => {
        (_0x54e43f["success"]("品牌设置已保存"), _0x1d78aa());
      },
      onError: (_0x5eb0d5) =>
        _0x54e43f["error"](
          _0x5eb0d5?.["response"]?.["data"]?.["message"] ||
            _0x5eb0d5?.["response"]?.["data"]?.["error"] ||
            "保存失败",
        ),
    }),
    _0x4d64a6 = _0x144b3f({
      mutationFn: async (_0x53fec7) => {
        const _0x53db08 = new FormData();
        return (
          _0x53db08["append"]("logo", _0x53fec7),
          (
            await _0x495bb8["post"](
              "/api/admin/system-settings/branding/logo",
              _0x53db08,
              { headers: { "Content-Type": "multipart/form-data" } },
            )
          )["data"]
        );
      },
      onSuccess: (_0x44030f) => {
        (_0x18581a(_0x44030f["logo_url"]),
          _0x54e43f["success"]("logo\x20已上传"),
          _0x1d78aa());
      },
      onError: (_0x5ee0b7) =>
        _0x54e43f["error"](
          _0x5ee0b7?.["response"]?.["data"]?.["error"] || "logo\x20上传失败",
        ),
    }),
    _0x4dd692 = _0x144b3f({
      mutationFn: async (_0x1f51d7) => {
        const _0x2e6aa7 = new FormData();
        return (
          _0x2e6aa7["append"]("icon", _0x1f51d7),
          (
            await _0x495bb8["post"](
              "/api/admin/system-settings/branding/icon",
              _0x2e6aa7,
              { headers: { "Content-Type": "multipart/form-data" } },
            )
          )["data"]
        );
      },
      onSuccess: (_0x4f57c6) => {
        (_0x346413(_0x4f57c6["icon_url"]),
          _0x54e43f["success"]("浏览器图标已上传"),
          _0x1d78aa());
      },
      onError: (_0x32d088) =>
        _0x54e43f["error"](
          _0x32d088?.["response"]?.["data"]?.["error"] || "浏览器图标上传失败",
        ),
    });
  return _0x2ddfe2["jsxs"](_0x2aeb40, {
    children: [
      _0x2ddfe2["jsxs"](_0x1db8ce, {
        className: "pb-4",
        children: [
          _0x2ddfe2["jsx"](_0x30c50e, { children: "自定义品牌" }),
          _0x2ddfe2["jsx"](_0x54661d, {
            children:
              "自定义站点标题、左上角标题、Logo\x20与浏览器图标。",
          }),
        ],
      }),
      _0x2ddfe2["jsx"](_0x42cb32, {
        children: _0x2ddfe2["jsx"](_0x23ae0c, {
          feature: "custom_branding",
          children: _0x2ddfe2["jsxs"]("div", {
            className: "space-y-4",
            children: [
              _0x2ddfe2["jsxs"]("div", {
                className: "space-y-1.5",
                children: [
                  _0x2ddfe2["jsx"](_0x34df34, {
                    children: "站点标题(浏览器标签页)",
                  }),
                  _0x2ddfe2["jsx"](_0x549353, {
                    value: _0x23baf3,
                    onChange: (_0xbb2b5) =>
                      _0x4eb364(_0xbb2b5["target"]["value"]),
                    placeholder: "留空用默认",
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"]("div", {
                className: "space-y-1.5",
                children: [
                  _0x2ddfe2["jsx"](_0x34df34, { children: "左上角标题" }),
                  _0x2ddfe2["jsx"](_0x549353, {
                    value: _0x395c54,
                    onChange: (_0x562a0) =>
                      _0x58c2f6(_0x562a0["target"]["value"]),
                    placeholder: "留空用默认",
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"]("div", {
                className: "space-y-1.5",
                children: [
                  _0x2ddfe2["jsx"](_0x34df34, {
                    children: "Logo(填图片\x20URL,或上传)",
                  }),
                  _0x2ddfe2["jsxs"]("div", {
                    className: "flex\x20gap-2",
                    children: [
                      _0x2ddfe2["jsx"](_0x549353, {
                        value: _0x14a01f,
                        onChange: (_0x3e2f6e) =>
                          _0x18581a(_0x3e2f6e["target"]["value"]),
                        placeholder: "https://.../logo.png,留空用默认",
                        className: "flex-1",
                      }),
                      _0x2ddfe2["jsx"]("input", {
                        ref: _0x131559,
                        type: "file",
                        accept:
                          "image/png,image/jpeg,image/webp,image/gif,image/svg+xml,image/x-icon",
                        className: "hidden",
                        onChange: (_0x542a50) => {
                          const _0x4d2954 = _0x542a50["target"]["files"]?.[0x0];
                          (_0x4d2954 && _0x4d64a6["mutate"](_0x4d2954),
                            (_0x542a50["target"]["value"] = ""));
                        },
                      }),
                      _0x2ddfe2["jsxs"](_0x5185a8, {
                        type: "button",
                        variant: "outline",
                        onClick: () => _0x131559["current"]?.["click"](),
                        disabled: _0x4d64a6["isPending"],
                        children: [
                          _0x4d64a6["isPending"]
                            ? _0x2ddfe2["jsx"](_0x4e8a70, {
                                className: "h-4\x20w-4\x20animate-spin",
                              })
                            : _0x2ddfe2["jsx"](_0xa3252a, {
                                className: "h-4\x20w-4",
                              }),
                          _0x2ddfe2["jsx"]("span", {
                            className: "ml-1",
                            children: "上传",
                          }),
                        ],
                      }),
                    ],
                  }),
                  _0x14a01f &&
                    _0x2ddfe2["jsx"]("img", {
                      src: _0x14a01f,
                      alt: "logo\x20预览",
                      className:
                        "mt-1\x20h-12\x20w-12\x20rounded\x20border\x20object-contain",
                    }),
                ],
              }),
              _0x2ddfe2["jsxs"]("div", {
                className: "space-y-1.5",
                children: [
                  _0x2ddfe2["jsx"](_0x34df34, {
                    children: "浏览器图标(填图片\x20URL,或上传)",
                  }),
                  _0x2ddfe2["jsxs"]("div", {
                    className: "flex\x20gap-2",
                    children: [
                      _0x2ddfe2["jsx"](_0x549353, {
                        value: _0x4ea065,
                        onChange: (_0x34624a) =>
                          _0x346413(_0x34624a["target"]["value"]),
                        placeholder: "https://.../favicon.ico,留空用默认",
                        className: "flex-1",
                      }),
                      _0x2ddfe2["jsx"]("input", {
                        ref: _0x919ca4,
                        type: "file",
                        accept:
                          "image/png,image/jpeg,image/webp,image/gif,image/svg+xml,image/x-icon",
                        className: "hidden",
                        onChange: (_0x57dd98) => {
                          const _0x2132f3 = _0x57dd98["target"]["files"]?.[0x0];
                          (_0x2132f3 && _0x4dd692["mutate"](_0x2132f3),
                            (_0x57dd98["target"]["value"] = ""));
                        },
                      }),
                      _0x2ddfe2["jsxs"](_0x5185a8, {
                        type: "button",
                        variant: "outline",
                        onClick: () => _0x919ca4["current"]?.["click"](),
                        disabled: _0x4dd692["isPending"],
                        children: [
                          _0x4dd692["isPending"]
                            ? _0x2ddfe2["jsx"](_0x4e8a70, {
                                className: "h-4\x20w-4\x20animate-spin",
                              })
                            : _0x2ddfe2["jsx"](_0xa3252a, {
                                className: "h-4\x20w-4",
                              }),
                          _0x2ddfe2["jsx"]("span", {
                            className: "ml-1",
                            children: "上传",
                          }),
                        ],
                      }),
                    ],
                  }),
                  _0x4ea065 &&
                    _0x2ddfe2["jsx"]("img", {
                      src: _0x4ea065,
                      alt: "浏览器图标预览",
                      className:
                        "mt-1\x20h-8\x20w-8\x20rounded\x20border\x20object-contain",
                    }),
                ],
              }),
              _0x2ddfe2["jsxs"](_0x5185a8, {
                onClick: () => _0x5194c5["mutate"](),
                disabled: _0x5194c5["isPending"] || _0x38cd27,
                children: [
                  _0x5194c5["isPending"] &&
                    _0x2ddfe2["jsx"](_0x4e8a70, {
                      className: "mr-1\x20h-4\x20w-4\x20animate-spin",
                    }),
                  "保存",
                ],
              }),
            ],
          }),
        }),
      }),
    ],
  });
}
const Tr = {
    driver: "postgres",
    host: "127.0.0.1",
    port: 0x1538,
    database: "mmwx",
    username: "mmwx",
    password: "",
    ssl_mode: "prefer",
    max_open_conns: 0x1e,
    max_idle_conns: 0xa,
  },
  hn = (_0x371c19) => ({
    driver: _0x371c19["driver"],
    path: _0x371c19["path"],
    host: _0x371c19["host"],
    port: _0x371c19["port"],
    database: _0x371c19["database"],
    username: _0x371c19["username"],
    password: _0x371c19["password"],
    ssl_mode: _0x371c19["ssl_mode"],
    max_open_conns: _0x371c19["max_open_conns"],
    max_idle_conns: _0x371c19["max_idle_conns"],
  }),
  gn = (_0x3481c2 = 0x0) => {
    if (_0x3481c2 < 0x400) return _0x3481c2 + "\x20B";
    const _0x404c75 = ["KB", "MB", "GB", "TB"];
    let _0x322fef = _0x3481c2 / 0x400,
      _0x4d062c = 0x0;
    for (; _0x322fef >= 0x400 && _0x4d062c < _0x404c75["length"] - 0x1; )
      ((_0x322fef /= 0x400), _0x4d062c++);
    return (
      _0x322fef["toFixed"](_0x322fef >= 0xa ? 0x1 : 0x2) +
      "\x20" +
      _0x404c75[_0x4d062c]
    );
  };
function Fr() {
  const [_0x165033, _0x1ff594] = _0x46e13c["useState"](Tr),
    [_0xf55247, _0x103a7b] = _0x46e13c["useState"](!0x1),
    [_0x4953a1, _0x50eed6] = _0x46e13c["useState"](!0x1),
    _0xc3b86 = _0x37d641({
      queryKey: ["database-status"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/database/status"))["data"][
          "status"
        ],
    });
  _0x46e13c["useEffect"](() => {
    const _0x499185 = _0xc3b86["data"];
    _0x499185?.["driver"] === "postgres" &&
      _0x1ff594((_0x366a86) => ({
        ..._0x366a86,
        ..._0x499185["config"],
        driver: "postgres",
        password: "",
      }));
  }, [_0xc3b86["data"]]);
  const _0x162b8f = _0x144b3f({
      mutationFn: () =>
        _0x495bb8["post"]("/api/admin/database/test", hn(_0x165033)),
      onSuccess: (_0x382428) =>
        _0x54e43f["success"](_0x382428["data"]["message"]),
      onError: _0x4deb0e,
    }),
    _0x51d1b4 = _0x37d641({
      queryKey: ["database-migration-progress"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/database/migration-progress"))[
          "data"
        ],
      enabled: _0x4953a1,
      refetchInterval: _0x4953a1 ? 0x1f4 : !0x1,
    }),
    _0x5de050 = _0x144b3f({
      mutationFn: () =>
        _0x495bb8["post"]("/api/admin/database/migrate", hn(_0x165033)),
      onMutate: () => {
        (_0x103a7b(!0x1), _0x50eed6(!0x0));
      },
      onSuccess: (_0x54a4b9) => {
        (_0x54e43f["success"](_0x54a4b9["data"]["message"], {
          description:
            "已迁移\x20" +
            _0x54a4b9["data"]["report"]["rows"] +
            "\x20行数据" +
            (_0x54a4b9["data"]["report"]["skipped"]
              ? "，跳过\x20" +
                _0x54a4b9["data"]["report"]["skipped"] +
                "\x20条孤儿数据"
              : "") +
            "。主控正在自动重启，页面即将刷新。",
        }),
          window["setTimeout"](() => window["location"]["reload"](), 0x1f40));
      },
      onError: _0x4deb0e,
      onSettled: () => _0x50eed6(!0x1),
    }),
    _0x28bc27 = _0x51d1b4["data"]?.["progress"],
    _0xa96a5b = _0x28bc27?.["total_tables"]
      ? Math["round"](
          (_0x28bc27["completed_tables"] / _0x28bc27["total_tables"]) * 0x64,
        )
      : 0x0,
    _0x14dfce = (_0x10a90c, _0x4decac) =>
      _0x1ff594((_0x3658ea) => ({ ..._0x3658ea, [_0x10a90c]: _0x4decac }));
  return _0x2ddfe2["jsxs"](_0x2aeb40, {
    children: [
      _0x2ddfe2["jsxs"](_0x1db8ce, {
        children: [
          _0x2ddfe2["jsxs"](_0x30c50e, {
            className: "flex\x20items-center\x20gap-2",
            children: [
              _0x2ddfe2["jsx"](_0x38a01f, { className: "h-5\x20w-5" }),
              "数据库",
            ],
          }),
          _0x2ddfe2["jsx"](_0x54661d, {
            children:
              "查看当前数据库占用，测试\x20PostgreSQL，并将现有\x20SQLite\x20数据一致迁移后自动重启切换。",
          }),
        ],
      }),
      _0x2ddfe2["jsxs"](_0x42cb32, {
        className: "space-y-6",
        children: [
          _0x2ddfe2["jsxs"]("div", {
            className:
              "grid\x20gap-3\x20rounded-lg\x20border\x20p-4\x20sm:grid-cols-4",
            children: [
              _0x2ddfe2["jsxs"]("div", {
                children: [
                  _0x2ddfe2["jsx"]("div", {
                    className: "text-muted-foreground\x20text-sm",
                    children: "当前类型",
                  }),
                  _0x2ddfe2["jsx"]("div", {
                    className: "font-medium",
                    children:
                      _0xc3b86["data"]?.["driver"] === "postgres"
                        ? "PostgreSQL"
                        : "SQLite",
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"]("div", {
                children: [
                  _0x2ddfe2["jsx"]("div", {
                    className: "text-muted-foreground\x20text-sm",
                    children: "数据库大小",
                  }),
                  _0x2ddfe2["jsx"]("div", {
                    className: "font-medium",
                    children: gn(_0xc3b86["data"]?.["size"]),
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"]("div", {
                children: [
                  _0x2ddfe2["jsx"]("div", {
                    className: "text-muted-foreground\x20text-sm",
                    children: "WAL",
                  }),
                  _0x2ddfe2["jsx"]("div", {
                    className: "font-medium",
                    children:
                      _0xc3b86["data"]?.["driver"] === "sqlite"
                        ? gn(_0xc3b86["data"]?.["wal_size"])
                        : "由\x20PostgreSQL\x20管理",
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"]("div", {
                children: [
                  _0x2ddfe2["jsx"]("div", {
                    className: "text-muted-foreground\x20text-sm",
                    children: "连接池",
                  }),
                  _0x2ddfe2["jsxs"]("div", {
                    className: "font-medium",
                    children: [
                      _0xc3b86["data"]?.["in_use_connections"] ?? 0x0,
                      "\x20使用中\x20/",
                      "\x20",
                      _0xc3b86["data"]?.["open_connections"] ?? 0x0,
                      "\x20已打开",
                    ],
                  }),
                ],
              }),
            ],
          }),
          _0x2ddfe2["jsxs"]("div", {
            className: "grid\x20gap-4\x20sm:grid-cols-2\x20lg:grid-cols-3",
            children: [
              _0x2ddfe2["jsx"](be, {
                label: "主机",
                children: _0x2ddfe2["jsx"](_0x549353, {
                  value: _0x165033["host"],
                  onChange: (_0x5b6c03) =>
                    _0x14dfce("host", _0x5b6c03["target"]["value"]),
                  placeholder: "127.0.0.1",
                }),
              }),
              _0x2ddfe2["jsx"](be, {
                label: "端口",
                children: _0x2ddfe2["jsx"](_0x549353, {
                  type: "number",
                  min: 0x1,
                  max: 0xffff,
                  value: _0x165033["port"],
                  onChange: (_0x462d8f) =>
                    _0x14dfce("port", Number(_0x462d8f["target"]["value"])),
                }),
              }),
              _0x2ddfe2["jsx"](be, {
                label: "数据库名",
                children: _0x2ddfe2["jsx"](_0x549353, {
                  value: _0x165033["database"],
                  onChange: (_0xddb3fa) =>
                    _0x14dfce("database", _0xddb3fa["target"]["value"]),
                }),
              }),
              _0x2ddfe2["jsx"](be, {
                label: "用户名",
                children: _0x2ddfe2["jsx"](_0x549353, {
                  value: _0x165033["username"],
                  onChange: (_0x58c90f) =>
                    _0x14dfce("username", _0x58c90f["target"]["value"]),
                }),
              }),
              _0x2ddfe2["jsx"](be, {
                label: "密码",
                children: _0x2ddfe2["jsx"](_0x549353, {
                  type: "password",
                  value: _0x165033["password"],
                  onChange: (_0x3b5478) =>
                    _0x14dfce("password", _0x3b5478["target"]["value"]),
                  autoComplete: "new-password",
                  placeholder: _0xc3b86["data"]?.["config"][
                    "password_configured"
                  ]
                    ? "留空表示不修改"
                    : "PostgreSQL\x20密码",
                }),
              }),
              _0x2ddfe2["jsx"](be, {
                label: "SSL\x20模式",
                children: _0x2ddfe2["jsxs"](_0x3ba40b, {
                  value: _0x165033["ssl_mode"],
                  onValueChange: (_0x4337ea) =>
                    _0x14dfce("ssl_mode", _0x4337ea),
                  children: [
                    _0x2ddfe2["jsx"](_0x3d4a68, {
                      children: _0x2ddfe2["jsx"](_0x1929c7, {}),
                    }),
                    _0x2ddfe2["jsx"](_0x51ae8f, {
                      children: [
                        "disable",
                        "prefer",
                        "require",
                        "verify-ca",
                        "verify-full",
                      ]["map"]((_0x22cad4) =>
                        _0x2ddfe2["jsx"](
                          _0x1260f9,
                          { value: _0x22cad4, children: _0x22cad4 },
                          _0x22cad4,
                        ),
                      ),
                    }),
                  ],
                }),
              }),
              _0x2ddfe2["jsx"](be, {
                label: "最大连接数",
                children: _0x2ddfe2["jsx"](_0x549353, {
                  type: "number",
                  min: 0x1,
                  value: _0x165033["max_open_conns"],
                  onChange: (_0x3cb655) =>
                    _0x14dfce(
                      "max_open_conns",
                      Number(_0x3cb655["target"]["value"]),
                    ),
                }),
              }),
              _0x2ddfe2["jsx"](be, {
                label: "最大空闲连接",
                children: _0x2ddfe2["jsx"](_0x549353, {
                  type: "number",
                  min: 0x0,
                  value: _0x165033["max_idle_conns"],
                  onChange: (_0x149c92) =>
                    _0x14dfce(
                      "max_idle_conns",
                      Number(_0x149c92["target"]["value"]),
                    ),
                }),
              }),
            ],
          }),
          _0x2ddfe2["jsxs"]("div", {
            className: "flex\x20flex-wrap\x20gap-2",
            children: [
              _0x2ddfe2["jsxs"](_0x5185a8, {
                variant: "outline",
                onClick: () => _0x162b8f["mutate"](),
                disabled: _0x162b8f["isPending"],
                children: [
                  _0x162b8f["isPending"]
                    ? _0x2ddfe2["jsx"](_0x4e8a70, {
                        className: "mr-2\x20h-4\x20w-4\x20animate-spin",
                      })
                    : _0x2ddfe2["jsx"](_0x18a4f9, {
                        className: "mr-2\x20h-4\x20w-4",
                      }),
                  "测试连接",
                ],
              }),
              _0xc3b86["data"]?.["driver"] === "sqlite" &&
                _0x2ddfe2["jsx"](_0x5185a8, {
                  onClick: () => _0x103a7b(!0x0),
                  disabled:
                    _0x162b8f["isPending"] ||
                    _0x5de050["isPending"] ||
                    _0xc3b86["data"]["environment_override"],
                  children: "迁移到\x20PostgreSQL",
                }),
            ],
          }),
          _0x4953a1 &&
            _0x2ddfe2["jsxs"]("div", {
              className: "space-y-2\x20rounded-lg\x20border\x20p-4",
              children: [
                _0x2ddfe2["jsxs"]("div", {
                  className:
                    "flex\x20flex-wrap\x20items-center\x20justify-between\x20gap-2\x20text-sm",
                  children: [
                    _0x2ddfe2["jsx"]("span", {
                      className: "font-medium",
                      children: "正在迁移数据库",
                    }),
                    _0x2ddfe2["jsxs"]("span", {
                      className: "text-muted-foreground",
                      children: [
                        _0x28bc27?.["completed_tables"] ?? 0x0,
                        "\x20/",
                        "\x20",
                        _0x28bc27?.["total_tables"] ?? 0x0,
                        "\x20张表\x20·",
                        "\x20",
                        (_0x28bc27?.["rows"] ?? 0x0)["toLocaleString"](),
                        "\x20行",
                        (_0x28bc27?.["skipped"] ?? 0x0) > 0x0
                          ? "\x20·\x20跳过\x20" +
                            (_0x28bc27?.["skipped"] ?? 0x0)[
                              "toLocaleString"
                            ]() +
                            "\x20条孤儿数据"
                          : "",
                      ],
                    }),
                  ],
                }),
                _0x2ddfe2["jsx"](_0x4c3ba9, { value: _0xa96a5b }),
                _0x2ddfe2["jsx"]("p", {
                  className: "text-muted-foreground\x20text-sm\x20break-all",
                  children:
                    _0x28bc27?.["phase"] === "preparing"
                      ? "正在检查目标表结构和字段类型…"
                      : _0x28bc27?.["current_table"]
                        ? "正在复制：" + _0x28bc27["current_table"]
                        : "正在连接并准备迁移…",
                }),
              ],
            }),
          _0xc3b86["data"]?.["driver"] === "postgres" &&
            _0x2ddfe2["jsx"]("p", {
              className: "text-muted-foreground\x20text-sm",
              children:
                "当前已使用\x20PostgreSQL。为防止误覆盖，不支持从此页面迁移到另一个非空数据库。",
            }),
          _0xc3b86["data"]?.["environment_override"] &&
            _0x2ddfe2["jsx"]("p", {
              className: "text-sm\x20text-amber-600",
              children:
                "当前连接由\x20MMWX_DATABASE_*\x20环境变量管理。页面可以测试连接，但迁移切换已禁用，请先移除连接类环境变量并重启。",
            }),
        ],
      }),
      _0x2ddfe2["jsx"](_0x1e06f2, {
        open: _0xf55247,
        onOpenChange: _0x103a7b,
        title: "迁移并切换数据库？",
        desc: "迁移期间会短暂阻止写入；目标数据库必须为空。完成后将写入\x20data/database.json\x20并自动重启主控。请先备份当前数据目录。",
        confirmText: "开始迁移",
        cancelBtnText: "取消",
        isLoading: _0x5de050["isPending"],
        handleConfirm: () => _0x5de050["mutate"](),
      }),
    ],
  });
}
function be({ label: _0xa9edae, children: _0xc0a752 }) {
  return _0x2ddfe2["jsxs"]("div", {
    className: "space-y-2",
    children: [_0x2ddfe2["jsx"](_0x34df34, { children: _0xa9edae }), _0xc0a752],
  });
}
const qr = {
  enabled: !0x1,
  bot_token: "",
  admin_tg_ids: [],
  webapp_dev_preview: !0x1,
  running: !0x1,
};
function Er() {
  const _0xbc7e2f = _0x29acb4(),
    [_0x200bd5, _0xeced22] = _0x46e13c["useState"](null),
    [_0x1eb3d7, _0x3588c5] = _0x46e13c["useState"](null),
    _0x2b1c4d = _0x37d641({
      queryKey: ["tgbot-settings"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/tgbot"))["data"],
    }),
    _0x3a0168 = _0x200bd5 || _0x2b1c4d["data"] || qr,
    _0x1ec931 = _0x1eb3d7 ?? (_0x3a0168["admin_tg_ids"] ?? [])["join"](",\x20"),
    _0x17b93f = _0x144b3f({
      mutationFn: async (_0x2e19c2) => {
        const _0xb34983 = _0x1ec931["split"](",")
          ["map"]((_0x56a066) => Number(_0x56a066["trim"]()))
          ["filter"](
            (_0x559a58) =>
              Number["isSafeInteger"](_0x559a58) && _0x559a58 > 0x0,
          );
        return (
          await _0x495bb8["put"]("/api/admin/system-settings/tgbot", {
            ..._0x2e19c2,
            admin_tg_ids: _0xb34983,
          })
        )["data"];
      },
      onSuccess: (_0xc29cdd) => {
        (_0xeced22(_0xc29cdd),
          _0x3588c5((_0xc29cdd["admin_tg_ids"] ?? [])["join"](",\x20")),
          _0xbc7e2f["setQueryData"](["tgbot-settings"], _0xc29cdd),
          _0x54e43f["success"](
            _0xc29cdd["enabled"]
              ? "TGBot\x20配置已保存并重新启动"
              : "TGBot\x20已停用",
          ));
      },
      onError: (_0x1b345a) => {
        const _0x4f153f = _0x1b345a;
        (_0xeced22(null),
          _0x3588c5(null),
          _0xbc7e2f["invalidateQueries"]({ queryKey: ["tgbot-settings"] }),
          _0x54e43f["error"](
            _0x4f153f["response"]?.["data"]?.["error"] ||
              _0x4f153f["message"] ||
              "保存失败",
          ));
      },
    });
  return _0x2ddfe2["jsxs"](_0x2aeb40, {
    children: [
      _0x2ddfe2["jsx"](_0x1db8ce, {
        children: _0x2ddfe2["jsxs"]("div", {
          className: "flex\x20items-center\x20justify-between\x20gap-4",
          children: [
            _0x2ddfe2["jsxs"]("div", {
              children: [
                _0x2ddfe2["jsxs"](_0x30c50e, {
                  className: "flex\x20items-center\x20gap-2",
                  children: [
                    _0x2ddfe2["jsx"](_0x4fa2aa, { className: "h-5\x20w-5" }),
                    "TGBot",
                  ],
                }),
                _0x2ddfe2["jsx"](_0x54661d, {
                  children:
                    "Bot\x20与\x20Mini\x20App\x20已内置到主控，Mini\x20App\x20地址自动使用主控地址\x20+\x20/tg-app；保存后自动热重启。",
                }),
              ],
            }),
            _0x2ddfe2["jsxs"]("div", {
              className: "flex\x20items-center\x20gap-2",
              children: [
                _0x2ddfe2["jsx"]("span", {
                  className:
                    "h-2.5\x20w-2.5\x20rounded-full\x20" +
                    (_0x3a0168["running"]
                      ? "bg-green-500"
                      : "bg-muted-foreground/40"),
                }),
                _0x2ddfe2["jsx"]("span", {
                  className: "text-muted-foreground\x20text-xs",
                  children: _0x3a0168["running"] ? "运行中" : "未运行",
                }),
                _0x2ddfe2["jsx"](_0x543b01, {
                  disabled: _0x17b93f["isPending"] || _0x2b1c4d["isLoading"],
                  checked: _0x3a0168["enabled"],
                  onCheckedChange: (_0x5c66fc) => {
                    const _0xc0155d = { ..._0x3a0168, enabled: _0x5c66fc };
                    (_0xeced22(_0xc0155d), _0x17b93f["mutate"](_0xc0155d));
                  },
                }),
              ],
            }),
          ],
        }),
      }),
      _0x2ddfe2["jsxs"](_0x42cb32, {
        className: "space-y-5",
        children: [
          _0x2ddfe2["jsxs"]("div", {
            className: "space-y-2",
            children: [
              _0x2ddfe2["jsx"](_0x34df34, {
                htmlFor: "tgbot-token",
                children: "Bot\x20Token",
              }),
              _0x2ddfe2["jsx"](_0x549353, {
                id: "tgbot-token",
                type: "password",
                autoComplete: "new-password",
                value: _0x3a0168["bot_token"],
                onChange: (_0x1cb72e) =>
                  _0xeced22({
                    ..._0x3a0168,
                    bot_token: _0x1cb72e["target"]["value"],
                  }),
                placeholder: "从\x20@BotFather\x20获取",
              }),
              _0x2ddfe2["jsx"]("p", {
                className: "text-muted-foreground\x20text-xs",
                children:
                  "已保存的\x20Token\x20只显示掩码；保持掩码不变不会覆盖原\x20Token。",
              }),
            ],
          }),
          _0x2ddfe2["jsxs"]("div", {
            className: "space-y-2",
            children: [
              _0x2ddfe2["jsx"](_0x34df34, {
                htmlFor: "tgbot-admins",
                children: "管理员\x20Telegram\x20ID",
              }),
              _0x2ddfe2["jsx"](_0x549353, {
                id: "tgbot-admins",
                value: _0x1ec931,
                onChange: (_0x560ccf) =>
                  _0x3588c5(_0x560ccf["target"]["value"]),
                placeholder: "123456789,\x20987654321",
              }),
              _0x2ddfe2["jsx"]("p", {
                className: "text-muted-foreground\x20text-xs",
                children:
                  "多个\x20ID\x20用英文逗号分隔，可使用\x20/admin_*\x20命令及\x20Mini\x20App\x20管理功能。",
              }),
            ],
          }),
          _0x2ddfe2["jsxs"]("div", {
            className:
              "flex\x20items-start\x20justify-between\x20gap-4\x20rounded-lg\x20border\x20p-4",
            children: [
              _0x2ddfe2["jsxs"]("div", {
                children: [
                  _0x2ddfe2["jsx"]("div", {
                    className: "text-sm\x20font-medium",
                    children: "Mini\x20App\x20调试预览",
                  }),
                  _0x2ddfe2["jsx"]("p", {
                    className: "text-muted-foreground\x20mt-1\x20text-xs",
                    children:
                      "允许从\x20URL\x20读取\x20initData，仅限本地调试；生产环境开启会带来重放风险。",
                  }),
                ],
              }),
              _0x2ddfe2["jsx"](_0x543b01, {
                checked: _0x3a0168["webapp_dev_preview"],
                onCheckedChange: (_0xefd765) =>
                  _0xeced22({ ..._0x3a0168, webapp_dev_preview: _0xefd765 }),
              }),
            ],
          }),
          _0x2ddfe2["jsx"](_0x5185a8, {
            onClick: () => _0x17b93f["mutate"](_0x3a0168),
            disabled: _0x17b93f["isPending"] || _0x2b1c4d["isLoading"],
            children: _0x17b93f["isPending"] ? "保存并重启中…" : "保存配置",
          }),
        ],
      }),
    ],
  });
}
const Mr = [
    { key: "subscription", label: "订阅链接" },
    { key: "generator", label: "生成订阅" },
    { key: "templates", label: "模板管理" },
    { key: "subscribe-files", label: "订阅管理" },
    { key: "custom-rules", label: "覆写管理" },
    { key: "nodes", label: "节点管理" },
  ],
  Kr = {
    pages: [],
    quota_template: 0x0,
    quota_override: 0x0,
    quota_subscribe: 0x0,
  };
function Br() {
  const _0x2a0075 = _0x29acb4(),
    [_0x143ab8, _0x1547c4] = _0x46e13c["useState"](!0x1),
    [_0x3e6e71, _0x26e3f8] = _0x46e13c["useState"](Kr),
    { data: _0x8f43ea } = _0x37d641({
      queryKey: ["user-permissions-config"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/user-permissions"))[
          "data"
        ],
      enabled: _0x143ab8,
    });
  _0x46e13c["useEffect"](() => {
    _0x8f43ea?.["config"] &&
      _0x26e3f8({
        pages: _0x8f43ea["config"]["pages"] ?? [],
        quota_template: _0x8f43ea["config"]["quota_template"] ?? 0x0,
        quota_override: _0x8f43ea["config"]["quota_override"] ?? 0x0,
        quota_subscribe: _0x8f43ea["config"]["quota_subscribe"] ?? 0x0,
      });
  }, [_0x8f43ea]);
  const _0x3938b2 = _0x144b3f({
      mutationFn: async (_0x27eb86) => {
        const _0x22955b = _0x8f43ea?.["config"] ?? {};
        await _0x495bb8["put"]("/api/admin/system-settings/user-permissions", {
          ..._0x27eb86,
          routed_outbound_enabled: !!_0x22955b["routed_outbound_enabled"],
          quota_routed_outbound: Number(
            _0x22955b["quota_routed_outbound"] ?? 0x0,
          ),
        });
      },
      onSuccess: () => {
        (_0x54e43f["success"]("用户权限已保存"),
          _0x2a0075["invalidateQueries"]({
            queryKey: ["user-permissions-config"],
          }),
          _0x2a0075["invalidateQueries"]({ queryKey: ["user-permissions"] }),
          _0x1547c4(!0x1));
      },
      onError: () => {
        _0x54e43f["error"]("保存失败");
      },
    }),
    _0xb5ac4d = (_0x8d5def, _0x522abb) => {
      _0x26e3f8((_0x200b3b) => ({
        ..._0x200b3b,
        pages: _0x522abb
          ? [..._0x200b3b["pages"], _0x8d5def]
          : _0x200b3b["pages"]["filter"](
              (_0x3fd8ab) => _0x3fd8ab !== _0x8d5def,
            ),
      }));
    },
    _0x3de61c = (_0xb99873, _0x24fd27) => {
      const _0x31bdfc = Math["max"](0x0, parseInt(_0x24fd27, 0xa) || 0x0);
      _0x26e3f8((_0x322c27) => ({ ..._0x322c27, [_0xb99873]: _0x31bdfc }));
    };
  return _0x2ddfe2["jsxs"](_0x1cc136, {
    open: _0x143ab8,
    onOpenChange: _0x1547c4,
    children: [
      _0x2ddfe2["jsx"](_0x3f9d9d, {
        asChild: !0x0,
        children: _0x2ddfe2["jsx"](_0x5185a8, {
          variant: "outline",
          size: "icon",
          className: "h-7\x20w-7",
          title: "用户权限配置",
          children: _0x2ddfe2["jsx"](_0x195e80, {
            className: "h-3.5\x20w-3.5",
          }),
        }),
      }),
      _0x2ddfe2["jsxs"](_0x3e8fab, {
        className: "sm:max-w-md",
        children: [
          _0x2ddfe2["jsxs"](_0x333ff1, {
            children: [
              _0x2ddfe2["jsx"](_0x4b37d4, { children: "普通用户权限配置" }),
              _0x2ddfe2["jsx"](_0x1e6ae9, {
                children:
                  "统一配置普通用户可见的页面，以及可创建的资源数量上限（0\x20=\x20不限）。",
              }),
            ],
          }),
          _0x2ddfe2["jsxs"]("div", {
            className: "space-y-4\x20py-2",
            children: [
              _0x2ddfe2["jsxs"]("div", {
                className: "space-y-2",
                children: [
                  _0x2ddfe2["jsx"](_0x34df34, {
                    className: "text-sm\x20font-semibold",
                    children: "可见页面",
                  }),
                  _0x2ddfe2["jsx"]("div", {
                    className: "grid\x20grid-cols-2\x20gap-2",
                    children: Mr["map"]((_0x29ee75) =>
                      _0x2ddfe2["jsxs"](
                        "label",
                        {
                          className:
                            "flex\x20cursor-pointer\x20items-center\x20gap-2\x20rounded-md\x20border\x20p-2",
                          children: [
                            _0x2ddfe2["jsx"](_0x77832a, {
                              checked: _0x3e6e71["pages"]["includes"](
                                _0x29ee75["key"],
                              ),
                              onCheckedChange: (_0x239459) =>
                                _0xb5ac4d(_0x29ee75["key"], _0x239459 === !0x0),
                            }),
                            _0x2ddfe2["jsx"]("span", {
                              className: "text-sm",
                              children: _0x29ee75["label"],
                            }),
                          ],
                        },
                        _0x29ee75["key"],
                      ),
                    ),
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"]("div", {
                className: "space-y-3",
                children: [
                  _0x2ddfe2["jsx"](_0x34df34, {
                    className: "text-sm\x20font-semibold",
                    children: "数量限制（0\x20=\x20不限）",
                  }),
                  _0x2ddfe2["jsxs"]("div", {
                    className:
                      "flex\x20items-center\x20justify-between\x20gap-3",
                    children: [
                      _0x2ddfe2["jsx"]("span", {
                        className: "text-muted-foreground\x20text-sm",
                        children: "模板数量",
                      }),
                      _0x2ddfe2["jsx"](_0x549353, {
                        type: "number",
                        min: 0x0,
                        className: "w-28",
                        value: _0x3e6e71["quota_template"],
                        onChange: (_0x194128) =>
                          _0x3de61c(
                            "quota_template",
                            _0x194128["target"]["value"],
                          ),
                      }),
                    ],
                  }),
                  _0x2ddfe2["jsxs"]("div", {
                    className:
                      "flex\x20items-center\x20justify-between\x20gap-3",
                    children: [
                      _0x2ddfe2["jsx"]("span", {
                        className: "text-muted-foreground\x20text-sm",
                        children: "覆写规则数量",
                      }),
                      _0x2ddfe2["jsx"](_0x549353, {
                        type: "number",
                        min: 0x0,
                        className: "w-28",
                        value: _0x3e6e71["quota_override"],
                        onChange: (_0x3a9194) =>
                          _0x3de61c(
                            "quota_override",
                            _0x3a9194["target"]["value"],
                          ),
                      }),
                    ],
                  }),
                  _0x2ddfe2["jsxs"]("div", {
                    className:
                      "flex\x20items-center\x20justify-between\x20gap-3",
                    children: [
                      _0x2ddfe2["jsx"]("span", {
                        className: "text-muted-foreground\x20text-sm",
                        children: "订阅数量",
                      }),
                      _0x2ddfe2["jsx"](_0x549353, {
                        type: "number",
                        min: 0x0,
                        className: "w-28",
                        value: _0x3e6e71["quota_subscribe"],
                        onChange: (_0x5ce306) =>
                          _0x3de61c(
                            "quota_subscribe",
                            _0x5ce306["target"]["value"],
                          ),
                      }),
                    ],
                  }),
                ],
              }),
            ],
          }),
          _0x2ddfe2["jsxs"](_0x5881af, {
            children: [
              _0x2ddfe2["jsx"](_0x5185a8, {
                variant: "outline",
                onClick: () => _0x1547c4(!0x1),
                children: "取消",
              }),
              _0x2ddfe2["jsx"](_0x5185a8, {
                onClick: () => _0x3938b2["mutate"](_0x3e6e71),
                disabled: _0x3938b2["isPending"],
                children: _0x3938b2["isPending"] ? "保存中..." : "保存",
              }),
            ],
          }),
        ],
      }),
    ],
  });
}
function Dr(_0x3effe4) {
  const _0x15f5ef = _0x3effe4?.["trim"]() || "";
  if (_0x7c97c(_0x15f5ef)) return _0x15f5ef;
  const _0x594117 = _0x15f5ef["split"](/[·,\s]+/)[0x0]?.["toUpperCase"]();
  return /^[A-Z]{2}$/["test"](_0x594117 || "") ? _0x3f64f4(_0x594117) : "";
}
function Rr(_0x26cb9d) {
  if (!_0x26cb9d || _0x26cb9d < 0x1 || _0x26cb9d > 0x1f) return "";
  const _0x5e4b2a = new Date(),
    _0x4c6c89 = Date["UTC"](
      _0x5e4b2a["getUTCFullYear"](),
      _0x5e4b2a["getUTCMonth"](),
      _0x5e4b2a["getUTCDate"](),
    ),
    _0x4d3268 = (_0x5b29c2, _0x983380) =>
      new Date(
        Date["UTC"](
          _0x5b29c2,
          _0x983380,
          Math["min"](
            _0x26cb9d,
            new Date(Date["UTC"](_0x5b29c2, _0x983380 + 0x1, 0x0))[
              "getUTCDate"
            ](),
          ),
        ),
      );
  let _0x142e51 = _0x4d3268(
    _0x5e4b2a["getUTCFullYear"](),
    _0x5e4b2a["getUTCMonth"](),
  );
  return (
    _0x142e51["getTime"]() < _0x4c6c89 &&
      (_0x142e51 = _0x4d3268(
        _0x5e4b2a["getUTCFullYear"](),
        _0x5e4b2a["getUTCMonth"]() + 0x1,
      )),
    _0x142e51["toISOString"]()["slice"](0x0, 0xa)
  );
}
function Ir({ server: _0x391b30, saved: _0xc3ca28 }) {
  const [_0x13482f, _0x1a2b3b] = _0x46e13c["useState"](Dr(_0x391b30["region"])),
    [_0x4f0855, _0x3c3433] = _0x46e13c["useState"](!0x1),
    [_0x374bad, _0x8aa3a1] = _0x46e13c["useState"](
      _0x391b30["renewal_cycle"] || "month",
    ),
    [_0x5d22fc, _0x225f62] = _0x46e13c["useState"](
      _0x391b30["expires_at"]?.["slice"](0x0, 0xa) ||
        Rr(_0x391b30["traffic_reset_day"]),
    ),
    _0x1e34fe = _0x144b3f({
      mutationFn: () =>
        _0x495bb8["put"]("/api/admin/remote-servers/update", {
          id: _0x391b30["id"],
          name: _0x391b30["name"],
          traffic_limit: _0x391b30["traffic_limit"] || 0x0,
          traffic_reset_day: _0x391b30["traffic_reset_day"] || 0x0,
          region: _0x13482f,
          renewal_cycle: _0x374bad,
          expires_at: _0x5d22fc,
        }),
      onSuccess: () => {
        (_0x54e43f["success"](_0x391b30["name"] + "\x20探针信息已保存"),
          _0xc3ca28());
      },
      onError: _0x4deb0e,
    });
  return _0x2ddfe2["jsxs"]("div", {
    className: "grid\x20gap-2\x20border-t\x20p-3\x20first:border-t-0",
    children: [
      _0x2ddfe2["jsx"]("strong", {
        className: "text-sm",
        children: _0x391b30["name"],
      }),
      _0x2ddfe2["jsxs"]("div", {
        className: "grid\x20grid-cols-2\x20gap-2\x20md:grid-cols-3",
        children: [
          _0x2ddfe2["jsxs"]("div", {
            className:
              "border-input\x20bg-background\x20flex\x20h-9\x20items-center\x20gap-2\x20rounded-md\x20border\x20px-2\x20text-sm",
            children: [
              _0x2ddfe2["jsx"](_0x3e16bd, {
                currentFlag: _0x13482f,
                loading: _0x4f0855,
                onSelect: _0x1a2b3b,
                onAutoDetect: async () => {
                  if (!_0x391b30["ip_address"]) {
                    _0x54e43f["error"]("服务器\x20IP\x20为空，无法自动识别");
                    return;
                  }
                  _0x3c3433(!0x0);
                  try {
                    const _0x1233f0 = await _0xe15cc6(_0x391b30["ip_address"]);
                    _0x1a2b3b(_0x3f64f4(_0x1233f0["country_code"]));
                  } catch {
                    _0x54e43f["error"]("地区自动识别失败");
                  } finally {
                    _0x3c3433(!0x1);
                  }
                },
              }),
              _0x2ddfe2["jsx"]("span", {
                className: "text-muted-foreground",
                children: _0x13482f || "自动识别地区",
              }),
            ],
          }),
          _0x2ddfe2["jsx"](_0x549353, {
            type: "date",
            value: _0x5d22fc,
            onChange: (_0x4951ec) => _0x225f62(_0x4951ec["target"]["value"]),
          }),
          _0x2ddfe2["jsxs"](_0x3ba40b, {
            value: _0x374bad,
            onValueChange: _0x8aa3a1,
            children: [
              _0x2ddfe2["jsx"](_0x3d4a68, {
                children: _0x2ddfe2["jsx"](_0x1929c7, {}),
              }),
              _0x2ddfe2["jsxs"](_0x51ae8f, {
                children: [
                  _0x2ddfe2["jsx"](_0x1260f9, {
                    value: "month",
                    children: "每月",
                  }),
                  _0x2ddfe2["jsx"](_0x1260f9, {
                    value: "quarter",
                    children: "每季度",
                  }),
                  _0x2ddfe2["jsx"](_0x1260f9, {
                    value: "half_year",
                    children: "每半年",
                  }),
                  _0x2ddfe2["jsx"](_0x1260f9, {
                    value: "year",
                    children: "每年",
                  }),
                ],
              }),
            ],
          }),
        ],
      }),
      _0x2ddfe2["jsx"](_0x5185a8, {
        size: "sm",
        variant: "outline",
        disabled: _0x1e34fe["isPending"],
        onClick: () => _0x1e34fe["mutate"](),
        children: _0x1e34fe["isPending"] ? "保存中…" : "保存探针信息",
      }),
    ],
  });
}
const fn = [
  { value: "sub", labelKey: "tabs.subscription" },
  { value: "features", labelKey: "tabs.features" },
  { value: "push", labelKey: "tabs.notifications" },
  { value: "security", labelKey: "tabs.security" },
  { value: "probe", labelKey: "tabs.probe" },
  { value: "appearance", labelKey: "tabs.appearance" },
  { value: "announce", labelKey: "tabs.announcements" },
  { value: "tgbot", labelKey: "tabs.tgbot" },
  { value: "captcha", labelKey: "tabs.verification" },
  { value: "system", labelKey: "tabs.system" },
  { value: "database", labelKey: "tabs.database" },
];
function yi() {
  const { t: _0x1aa03d } = _0x2443c3("system"),
    _0x2f82ce = _0x29acb4(),
    { auth: _0x1bb57f } = _0x235baf(),
    [_0x4873e7, _0x1010c7] = _0x46e13c["useState"]("sub"),
    { data: _0x2d5c9d } = _0x26349a(),
    _0x497b59 = _0x24539f(_0x2d5c9d),
    [_0x137d5c, _0xec8154] = _0x46e13c["useState"](!0x1),
    [_0x53c069, _0xacab7f] = _0x46e13c["useState"]("node_name"),
    [_0x171e10, _0x3874f3] = _0x46e13c["useState"]("saved_only"),
    [_0x2b23a2, _0x5c52ba] = _0x46e13c["useState"](!0x0),
    [_0x2376de, _0x1802db] = _0x46e13c["useState"](0x0),
    [_0x49b722, _0x1df0ca] = _0x46e13c["useState"](!0x1),
    [_0x3cc118, _0x3b6a2a] =
      _0x46e13c["useState"]("剩余|流量|到期|订阅|时间|重置"),
    [_0x431d4a, _0x495896] = _0x46e13c["useState"](!0x1),
    [_0x18a19f, _0x281a92] = _0x46e13c["useState"](!0x1),
    [_0x1e8a26, _0x43ecb0] = _0x46e13c["useState"](!0x1),
    [_0xb9053e, _0xa30760] = _0x46e13c["useState"]("📅过期时间"),
    [_0x780f57, _0x18df6d] = _0x46e13c["useState"]("⌛剩余流量"),
    [_0x569631, _0x1c209a] = _0x46e13c["useState"](!0x1),
    [_0x1ff17e, _0x5f5553] = _0x46e13c["useState"](""),
    [_0x420042, _0x43f9e2] = _0x46e13c["useState"](""),
    [_0x5d21bf, _0x211f4e] = _0x46e13c["useState"](!0x1),
    [_0x5274a5, _0xd9dd4e] = _0x46e13c["useState"](!0x1),
    [_0x2af5a5, _0x1fc6ef] = _0x46e13c["useState"](""),
    [_0x5566fb, _0x63ad43] = _0x46e13c["useState"](0x5),
    [_0x54682b, _0x646247] = _0x46e13c["useState"](0xa),
    [_0x102305, _0xde3d23] = _0x46e13c["useState"](!0x1),
    [_0x3b6d17, _0x12a0cf] = _0x46e13c["useState"](!0x1),
    [_0x4378bf, _0x39db92] = _0x46e13c["useState"](!0x0),
    [_0x3f7c69, _0x5da0c4] = _0x46e13c["useState"](!0x1),
    [_0x406eb1, _0x2d0a4e] = _0x46e13c["useState"]([]),
    [_0x10f487, _0x946043] = _0x46e13c["useState"](!0x1),
    { data: _0x11f05e } = _0x37d641({
      queryKey: ["master-url"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/master-url"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    });
  _0x46e13c["useEffect"](() => {
    (_0x11f05e?.["master_url"] !== void 0x0 &&
      _0x5f5553(_0x11f05e["master_url"]),
      _0x11f05e?.["subscription_url"] !== void 0x0 &&
        _0x43f9e2(_0x11f05e["subscription_url"]),
      _0x11f05e?.["local_only"] !== void 0x0 &&
        _0x211f4e(_0x11f05e["local_only"]),
      _0x11f05e?.["https_recovery_enabled"] !== void 0x0 &&
        _0xd9dd4e(_0x11f05e["https_recovery_enabled"]),
      _0x11f05e?.["recovery_url_custom"] !== void 0x0 &&
        _0x1fc6ef(_0x11f05e["recovery_url_custom"]),
      _0x11f05e?.["recovery_failure_minutes"] &&
        _0x63ad43(_0x11f05e["recovery_failure_minutes"]),
      _0x11f05e?.["recovery_startup_grace_minutes"] &&
        _0x646247(_0x11f05e["recovery_startup_grace_minutes"]));
  }, [_0x11f05e]);
  const _0x24110e = _0x144b3f({
      mutationFn: async (_0x152cc9) =>
        (
          await _0x495bb8["post"](
            "/api/admin/system-settings/master-migration",
            {
              action: _0x152cc9,
              new_master_url: _0x1ff17e["trim"]()["replace"](/\/+$/, ""),
              change_domain: _0x3b6d17,
              move_host: _0x4378bf,
              force: _0x3f7c69,
            },
          )
        )["data"],
      onSuccess: (_0x32f15a) => {
        (_0x2d0a4e(_0x32f15a["agents"] || []),
          _0x946043(_0x32f15a["ready"]),
          _0x32f15a["committed"]
            ? (_0x2f82ce["invalidateQueries"]({ queryKey: ["master-url"] }),
              _0x2f82ce["invalidateQueries"]({
                queryKey: ["security-settings"],
              }),
              _0x54e43f["success"](
                _0x32f15a["turnstile_disabled"]
                  ? "主控迁移地址已下发并保存，Cloudflare\x20验证码已自动关闭"
                  : "主控迁移地址已下发并保存",
              ),
              _0xde3d23(!0x1),
              window["setTimeout"](() => {
                window["location"]["href"] = _0x1ff17e["trim"]()["replace"](
                  /\/+$/,
                  "",
                );
              }, 0x5dc))
            : _0x54e43f["success"](
                _0x32f15a["ready"]
                  ? "所有\x20Agent\x20已通过迁移检查"
                  : "检查完成，请处理风险项",
              ));
      },
      onError: (_0x41929d) => {
        const _0x560a18 = _0x41929d?.["response"]?.["data"];
        (Array["isArray"](_0x560a18?.["agents"]) &&
          (_0x2d0a4e(_0x560a18["agents"]), _0x946043(!!_0x560a18["ready"])),
          _0x4deb0e(_0x41929d));
      },
    }),
    _0x1dadac = _0x144b3f({
      mutationFn: async (_0x41da4e) => {
        await _0x495bb8["put"]("/api/admin/system-settings/master-url", {
          subscription_url: _0x41da4e,
        });
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["master-url"] }),
          _0x2f82ce["invalidateQueries"]({ queryKey: ["user-config"] }),
          _0x54e43f["success"](_0x1aa03d("masterUrl.subscriptionUpdated")));
      },
      onError: _0x4deb0e,
    }),
    _0xda88d3 = _0x144b3f({
      mutationFn: async (_0x197ede) => {
        await _0x495bb8["put"]("/api/admin/system-settings/master-url", {
          local_only: _0x197ede,
        });
      },
      onSuccess: (_0x1ef469, _0x137ab2) => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["master-url"] }),
          _0x54e43f["success"](
            _0x1aa03d(
              _0x137ab2
                ? "masterUrl.localOnlyEnabled"
                : "masterUrl.localOnlyDisabled",
            ),
          ));
      },
      onError: (_0xcfaf3d) => {
        (_0x211f4e(_0x11f05e?.["local_only"] ?? !0x1), _0x4deb0e(_0xcfaf3d));
      },
    }),
    _0x580940 = _0x144b3f({
      mutationFn: async () => {
        await _0x495bb8["put"]("/api/admin/system-settings/master-url", {
          https_recovery_enabled: _0x5274a5,
          recovery_url: _0x2af5a5["trim"]()["replace"](/\/+$/, ""),
          recovery_failure_minutes: _0x5566fb,
          recovery_startup_grace_minutes: _0x54682b,
        });
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["master-url"] }),
          _0x54e43f["success"]("HTTPS\x20故障自愈设置已保存"));
      },
      onError: _0x4deb0e,
    }),
    [_0x478a8c, _0x3b5e6d] = _0x46e13c["useState"](""),
    { data: _0x568d46 } = _0x37d641({
      queryKey: ["redeem-template"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/redeem-template"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    });
  _0x46e13c["useEffect"](() => {
    _0x568d46?.["redeem_template"] !== void 0x0 &&
      _0x3b5e6d(_0x568d46["redeem_template"]);
  }, [_0x568d46]);
  const _0x567a5b = _0x144b3f({
      mutationFn: async (_0x411374) => {
        await _0x495bb8["put"]("/api/admin/system-settings/redeem-template", {
          redeem_template: _0x411374,
        });
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["redeem-template"] }),
          _0x54e43f["success"](_0x1aa03d("redeemTemplate.updated")));
      },
      onError: _0x4deb0e,
    }),
    { data: _0xcc0e07, isLoading: _0x436e38 } = _0x37d641({
      queryKey: ["api-token"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/api-token"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x38366f = _0x144b3f({
      mutationFn: async () =>
        (
          await _0x495bb8["post"](
            "/api/admin/system-settings/api-token/regenerate",
          )
        )["data"],
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["api-token"] }),
          _0x54e43f["success"](_0x1aa03d("apiToken.regenerated")));
      },
      onError: _0x4deb0e,
    }),
    _0x290ac2 = () => {
      _0xcc0e07?.["token"] &&
        (navigator["clipboard"]["writeText"](_0xcc0e07["token"]),
        _0x54e43f["success"](_0x1aa03d("apiToken.copied")));
    },
    { data: _0x34bf94 } = _0x37d641({
      queryKey: ["short-link-enabled"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/short-link"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x2010e0 = _0x144b3f({
      mutationFn: async (_0xc5ba2a) => {
        await _0x495bb8["put"]("/api/admin/system-settings/short-link", {
          enable_short_link: _0xc5ba2a,
        });
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["short-link-enabled"] }),
          _0x2f82ce["invalidateQueries"]({ queryKey: ["user-subscriptions"] }),
          _0x54e43f["success"](_0x1aa03d("shortLink.updated")));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x34bf94?.["enable_short_link"] !== void 0x0 &&
      _0x521715(_0x34bf94["enable_short_link"]);
  }, [_0x34bf94]);
  const [_0x28c051, _0x521715] = _0x46e13c["useState"](!0x0),
    [_0x465b79, _0x71f3e] = _0x46e13c["useState"](!0x1),
    [_0x5e10d1, _0x2d46c8] = _0x46e13c["useState"](!0x0),
    [_0x1ec742, _0x44237a] = _0x46e13c["useState"](!0x0),
    [_0x40888d, _0x2f492d] = _0x46e13c["useState"](!0x1),
    [_0x2ac809, _0x5311a8] = _0x46e13c["useState"](0x2),
    [_0x4a293d, _0x59d309] = _0x46e13c["useState"](0x5),
    { data: _0x181019 } = _0x37d641({
      queryKey: ["override-scripts-enabled"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/override-scripts"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x407318 = _0x144b3f({
      mutationFn: async (_0x32645e) => {
        await _0x495bb8["put"]("/api/admin/system-settings/override-scripts", {
          enable_override_scripts: _0x32645e,
        });
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({
          queryKey: ["override-scripts-enabled"],
        }),
          _0x54e43f["success"](_0x1aa03d("overrideScripts.updated")));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x181019?.["enable_override_scripts"] !== void 0x0 &&
      _0x71f3e(_0x181019["enable_override_scripts"]);
  }, [_0x181019]);
  const [_0x47c4df, _0x2dd02f] = _0x46e13c["useState"](!0x0),
    { data: _0x53e6c4 } = _0x37d641({
      queryKey: ["update-cdn-enabled"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/update-cdn"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x431a9b = _0x144b3f({
      mutationFn: async (_0x2de835) => {
        await _0x495bb8["put"]("/api/admin/system-settings/update-cdn", {
          enabled: _0x2de835,
        });
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["update-cdn-enabled"] }),
          _0x54e43f["success"](_0x1aa03d("updateCDN.updated")));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x53e6c4?.["enabled"] !== void 0x0 && _0x2dd02f(_0x53e6c4["enabled"]);
  }, [_0x53e6c4]);
  const [_0x170746, _0x49d02f] = _0x46e13c["useState"]("yaml"),
    { data: _0x254602 } = _0x37d641({
      queryKey: ["subscription-output-format"],
      queryFn: async () =>
        (
          await _0x495bb8["get"](
            "/api/admin/system-settings/subscription-output-format",
          )
        )["data"],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x16cc3a = _0x144b3f({
      mutationFn: async (_0x9c0f39) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/subscription-output-format",
          { subscription_output_format: _0x9c0f39 },
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({
          queryKey: ["subscription-output-format"],
        }),
          _0x54e43f["success"](_0x1aa03d("subscriptionOutputFormat.updated")));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x254602?.["subscription_output_format"] &&
      _0x49d02f(_0x254602["subscription_output_format"]);
  }, [_0x254602]);
  const { data: _0x4eddf7 } = _0x37d641({
      queryKey: ["miaomiaowu-features-enabled"],
      queryFn: async () =>
        (
          await _0x495bb8["get"](
            "/api/admin/system-settings/miaomiaowu-features",
          )
        )["data"],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x146a37 = _0x144b3f({
      mutationFn: async (_0x5668a1) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/miaomiaowu-features",
          { enable_miaomiaowu_features: _0x5668a1 },
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({
          queryKey: ["miaomiaowu-features-enabled"],
        }),
          _0x54e43f["success"](_0x1aa03d("miaomiaowuFeatures.updated")));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x4eddf7?.["enable_miaomiaowu_features"] !== void 0x0 &&
      _0x44237a(_0x4eddf7["enable_miaomiaowu_features"]);
  }, [_0x4eddf7]);
  const { data: _0xa60d17 } = _0x37d641({
      queryKey: ["node-name-multiplier-prefix"],
      queryFn: async () =>
        (
          await _0x495bb8["get"](
            "/api/admin/system-settings/node-name-multiplier-prefix",
          )
        )["data"],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    [_0x256337, _0x5901ee] = _0x46e13c["useState"](!0x1),
    [_0x1fdf38, _0x5eddf2] = _0x46e13c["useState"]("「"),
    [_0x38ca4b, _0xb1b77b] = _0x46e13c["useState"]("」");
  _0x46e13c["useEffect"](() => {
    _0xa60d17 &&
      (_0x5901ee(!!_0xa60d17["enabled"]),
      _0xa60d17["left"] && _0x5eddf2(_0xa60d17["left"]),
      _0xa60d17["right"] && _0xb1b77b(_0xa60d17["right"]));
  }, [_0xa60d17]);
  const _0x34a342 = _0x144b3f({
      mutationFn: async (_0x455e5d) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/node-name-multiplier-prefix",
          _0x455e5d,
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({
          queryKey: ["node-name-multiplier-prefix"],
        }),
          _0x54e43f["success"]("设置已更新"));
      },
      onError: _0x4deb0e,
    }),
    { data: _0x1fd650 } = _0x37d641({
      queryKey: ["user-permissions-config"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/user-permissions"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    });
  _0x46e13c["useEffect"](() => {
    const _0x157fe9 = _0x1fd650?.["config"];
    if (_0x157fe9) {
      _0x2f492d(!!_0x157fe9["routed_outbound_enabled"]);
      const _0x529dcf = Number(_0x157fe9["quota_routed_outbound"]);
      _0x5311a8(_0x529dcf > 0x0 ? _0x529dcf : 0x2);
      const _0x428c82 = Number(_0x157fe9["daily_limit_routed_outbound"]);
      _0x59d309(_0x428c82 > 0x0 ? _0x428c82 : 0x5);
    }
  }, [_0x1fd650]);
  const _0x32b255 = _0x144b3f({
      mutationFn: async (_0x36bf9f) => {
        const _0x349fac = _0x1fd650?.["config"] ?? {};
        await _0x495bb8["put"]("/api/admin/system-settings/user-permissions", {
          pages: _0x349fac["pages"] ?? [],
          quota_template: Number(_0x349fac["quota_template"] ?? 0x0),
          quota_override: Number(_0x349fac["quota_override"] ?? 0x0),
          quota_subscribe: Number(_0x349fac["quota_subscribe"] ?? 0x0),
          routed_outbound_enabled: _0x36bf9f["enabled"],
          quota_routed_outbound: _0x36bf9f["enabled"]
            ? _0x36bf9f["quota"]
            : 0x0,
          daily_limit_routed_outbound: _0x36bf9f["enabled"]
            ? _0x36bf9f["daily"]
            : 0x0,
        });
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({
          queryKey: ["user-permissions-config"],
        }),
          _0x2f82ce["invalidateQueries"]({ queryKey: ["user-permissions"] }),
          _0x54e43f["success"]("已保存"));
      },
      onError: (_0x3cfa65) => {
        _0x54e43f["error"](
          _0x3cfa65?.["response"]?.["data"]?.["message"] || "保存失败",
        );
      },
    }),
    [_0x4a3bd8, _0x451e75] = _0x46e13c["useState"](!0x1),
    [_0x1f6fa6, _0x9f6337] = _0x46e13c["useState"](0xf),
    [_0x3f4852, _0x31e9f0] = _0x46e13c["useState"](!0x1),
    [_0x16e9a0, _0x492328] = _0x46e13c["useState"](!0x1),
    { data: _0x991a4e } = _0x37d641({
      queryKey: ["silent-mode"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/silent-mode"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x842682 = _0x144b3f({
      mutationFn: async (_0x3a4c6d) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/silent-mode",
          _0x3a4c6d,
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["silent-mode"] }),
          _0x54e43f["success"](_0x1aa03d("silentMode.updated")));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x991a4e &&
      (_0x451e75(_0x991a4e["silent_mode"]),
      _0x9f6337(_0x991a4e["silent_mode_timeout"]));
  }, [_0x991a4e]);
  const [_0x179ea5, _0x55e2c2] = _0x46e13c["useState"](!0x1),
    [_0x28af34, _0x195e6d] = _0x46e13c["useState"](!0x1),
    [_0x5289ba, _0x5aca07] = _0x46e13c["useState"](!0x1),
    [_0x39dd42, _0x4c73e] = _0x46e13c["useState"](""),
    [_0x307fbc, _0x593c3b] = _0x46e13c["useState"]("follow"),
    [_0xb66ae3, _0x1e1227] = _0x46e13c["useState"](""),
    [_0x385ae1, _0x34add0] = _0x46e13c["useState"](""),
    [_0x2f3e37, _0x320f7a] = _0x46e13c["useState"](!0x1),
    [_0x4832dd, _0x200bc7] = _0x46e13c["useState"](!0x1),
    [_0x3846db, _0x586e39] = _0x46e13c["useState"](!0x1),
    [_0x57745b, _0x4ecb0d] = _0x46e13c["useState"](""),
    [_0x2f7373, _0x696237] = _0x46e13c["useState"](!0x0),
    [_0x132fa8, _0x2af46d] = _0x46e13c["useState"]([]),
    [_0x580edb, _0x32438b] = _0x46e13c["useState"](0x7),
    [_0x123118, _0x1f3edf] = _0x46e13c["useState"](!0x1),
    [_0x3810c9, _0x510199] = _0x46e13c["useState"](0x1e),
    [_0x40c023, _0x5e7f11] = _0x46e13c["useState"](!0x1),
    [_0x237d87, _0x3c1046] = _0x46e13c["useState"](!0x1),
    [_0x334fa2, _0x54bd7b] = _0x46e13c["useState"](!0x1),
    [_0x36e07d, _0x2169f6] = _0x46e13c["useState"](!0x1),
    [_0x201bbd, _0x566d88] = _0x46e13c["useState"](!0x0),
    [_0x32c971, _0x745469] = _0x46e13c["useState"](!0x0),
    [_0x48c9e8, _0x298897] = _0x46e13c["useState"](!0x1),
    [_0xf265e3, _0x3addbf] = _0x46e13c["useState"](!0x1),
    [_0xc58657, _0x19ab74] = _0x46e13c["useState"](!0x0),
    [_0x12d674, _0x1f17de] = _0x46e13c["useState"](!0x0),
    [_0x3cdab9, _0x98befa] = _0x46e13c["useState"](!0x0),
    [_0x4894c7, _0xe71a4f] = _0x46e13c["useState"](!0x0),
    [_0x4c2853, _0x21e7c8] = _0x46e13c["useState"](!0x0),
    [_0x3fc691, _0xa4bc46] = _0x46e13c["useState"](!0x0),
    [_0x52c280, _0x206944] = _0x46e13c["useState"](!0x0),
    [_0x584361, _0x509596] = _0x46e13c["useState"](!0x1),
    [_0x2c452d, _0x397aa6] = _0x46e13c["useState"](!0x1),
    [_0x43391c, _0x34987e] = _0x46e13c["useState"](!0x1),
    [_0x25dc28, _0x1cab7e] = _0x46e13c["useState"]([]),
    [_0x14d133, _0x39e60b] = _0x46e13c["useState"]({}),
    { data: _0x41f4f1 } = _0x37d641({
      queryKey: ["probe-disguise-settings"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/probe-disguise"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x45b042 = _0x144b3f({
      mutationFn: async (_0x269cff) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/probe-disguise",
          _0x269cff,
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({
          queryKey: ["probe-disguise-settings"],
        }),
          _0x54e43f["success"]("伪装配置已保存"));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x41f4f1 &&
      (_0x195e6d(_0x41f4f1["internal_enabled"]),
      _0x5aca07(_0x41f4f1["external_enabled"]),
      _0x55e2c2(_0x41f4f1["internal_enabled"] || _0x41f4f1["external_enabled"]),
      _0x4c73e(_0x41f4f1["title"] || ""),
      _0x593c3b(_0x41f4f1["theme"] || "follow"),
      _0x1e1227(_0x41f4f1["logo"] || ""),
      _0x34add0(_0x41f4f1["icon"] || ""),
      _0x320f7a(!!_0x41f4f1["block_login"]),
      _0x200bc7(!!_0x41f4f1["external_access_only"]),
      _0x586e39(!!_0x41f4f1["external_token_configured"]),
      _0x696237(_0x41f4f1["show_name"]),
      _0x2af46d(_0x41f4f1["server_ids"] || []),
      _0x5e7f11(_0x41f4f1["metric_cpu"]),
      _0x3c1046(_0x41f4f1["metric_mem"]),
      _0x54bd7b(_0x41f4f1["metric_disk"]),
      _0x2169f6(_0x41f4f1["metric_ping"]),
      _0x566d88(_0x41f4f1["metric_traffic"]),
      _0x745469(_0x41f4f1["metric_speed"]),
      _0x298897(!!_0x41f4f1["show_expiry"]),
      _0x3addbf(!!_0x41f4f1["show_globe"]),
      _0x19ab74(_0x41f4f1["show_forward"] !== !0x1),
      _0x1f17de(_0x41f4f1["show_daily_trend"] !== !0x1),
      _0x98befa(_0x41f4f1["show_traffic_hotspots"] !== !0x1),
      _0xe71a4f(_0x41f4f1["show_traffic_7d"] !== !0x1),
      _0x21e7c8(_0x41f4f1["show_resource_heatmap"] !== !0x1),
      _0xa4bc46(_0x41f4f1["show_traffic_quota"] !== !0x1),
      _0x206944(_0x41f4f1["show_renewal_timeline"] !== !0x1),
      _0x509596(!!_0x41f4f1["show_health_score"]),
      _0x397aa6(!!_0x41f4f1["show_return_route"]),
      _0x34987e(!!_0x41f4f1["show_external_license"]),
      _0x1cab7e(_0x41f4f1["ping_targets"] || []),
      _0x39e60b(_0x41f4f1["ping_targets_override"] || {}),
      _0x32438b(_0x41f4f1["forward_metrics_retention_days"] || 0x7),
      _0x510199(_0x41f4f1["forward_daily_retention_days"] || 0x1e),
      _0x1f3edf(!!_0x41f4f1["metrics_persist_enabled"]));
  }, [_0x41f4f1]);
  const _0x6360b0 = (_0x16692c) => {
      const _0x4dc829 = _0x16692c["internal_enabled"] ?? _0x28af34,
        _0x1a052d = _0x16692c["external_enabled"] ?? _0x5289ba,
        _0x2ae073 = {
          enabled: _0x4dc829 || _0x1a052d,
          internal_enabled: _0x4dc829,
          external_enabled: _0x1a052d,
          title: _0x16692c["title"] ?? _0x39dd42,
          theme: _0x16692c["theme"] ?? _0x307fbc,
          logo: _0x16692c["logo"] ?? _0xb66ae3,
          icon: _0x16692c["icon"] ?? _0x385ae1,
          block_login: _0x16692c["block_login"] ?? _0x2f3e37,
          external_access_only: _0x16692c["external_access_only"] ?? _0x4832dd,
          server_ids: _0x16692c["server_ids"] ?? _0x132fa8,
          show_name: _0x16692c["show_name"] ?? _0x2f7373,
          metric_cpu: _0x16692c["metric_cpu"] ?? _0x40c023,
          metric_mem: _0x16692c["metric_mem"] ?? _0x237d87,
          metric_disk: _0x16692c["metric_disk"] ?? _0x334fa2,
          metric_ping: _0x16692c["metric_ping"] ?? _0x36e07d,
          metric_traffic: _0x16692c["metric_traffic"] ?? _0x201bbd,
          metric_speed: _0x16692c["metric_speed"] ?? _0x32c971,
          show_expiry: _0x16692c["show_expiry"] ?? _0x48c9e8,
          show_globe: _0x16692c["show_globe"] ?? _0xf265e3,
          show_forward: _0x16692c["show_forward"] ?? _0xc58657,
          show_daily_trend: _0x16692c["show_daily_trend"] ?? _0x12d674,
          show_traffic_hotspots:
            _0x16692c["show_traffic_hotspots"] ?? _0x3cdab9,
          show_traffic_7d: _0x16692c["show_traffic_7d"] ?? _0x4894c7,
          show_resource_heatmap:
            _0x16692c["show_resource_heatmap"] ?? _0x4c2853,
          show_traffic_quota: _0x16692c["show_traffic_quota"] ?? _0x3fc691,
          show_renewal_timeline:
            _0x16692c["show_renewal_timeline"] ?? _0x52c280,
          show_health_score: _0x16692c["show_health_score"] ?? _0x584361,
          show_return_route: _0x16692c["show_return_route"] ?? _0x2c452d,
          show_external_license:
            _0x16692c["show_external_license"] ?? _0x43391c,
          ping_targets: _0x16692c["ping_targets"] ?? _0x25dc28,
          ping_targets_override:
            _0x16692c["ping_targets_override"] ?? _0x14d133,
          forward_metrics_retention_days:
            _0x16692c["forward_metrics_retention_days"] ?? _0x580edb,
          forward_daily_retention_days:
            _0x16692c["forward_daily_retention_days"] ?? _0x3810c9,
          metrics_persist_enabled:
            _0x16692c["metrics_persist_enabled"] ?? _0x123118,
        };
      (_0x16692c["external_access_token"] &&
        (_0x2ae073["external_access_token"] =
          _0x16692c["external_access_token"]),
        _0x45b042["mutate"](_0x2ae073));
    },
    [_0x3fd9b7, _0x35426f] = _0x46e13c["useState"]("flat"),
    { data: _0x5048a0 } = _0x37d641({
      queryKey: ["default-theme-settings"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/default-theme"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x47362d = _0x144b3f({
      mutationFn: async (_0x22f752) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/default-theme",
          _0x22f752,
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({
          queryKey: ["default-theme-settings"],
        }),
          _0x54e43f["success"]("默认主题已保存"));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x5048a0 && _0x35426f(_0x5048a0["default_theme"]);
  }, [_0x5048a0]);
  const [_0x1f5853, _0x46b92e] = _0x46e13c["useState"](""),
    { data: _0x510a11 } = _0x37d641({
      queryKey: ["login-wallpaper-settings"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/login-wallpaper"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0x3e6952 = _0x144b3f({
      mutationFn: async (_0x229a4e) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/login-wallpaper",
          _0x229a4e,
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({
          queryKey: ["login-wallpaper-settings"],
        }),
          _0x54e43f["success"]("登录页壁纸已保存"));
      },
      onError: _0x4deb0e,
    });
  _0x46e13c["useEffect"](() => {
    _0x510a11 && _0x46b92e(_0x510a11["login_wallpaper"] || "");
  }, [_0x510a11]);
  const { data: _0x3ded3d } = _0x37d641({
      queryKey: ["remote-servers"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/remote-servers"))["data"],
      enabled: !!_0x1bb57f["accessToken"] && _0x179ea5,
      staleTime: 0x3c * 0x3e8,
    }),
    { data: _0x12640b } = Cr(
      !!_0x1bb57f["accessToken"] && _0x179ea5 && _0x36e07d,
    ),
    { data: _0x5b58be } = _0x37d641({
      queryKey: ["require-encryption"],
      queryFn: async () =>
        (
          await _0x495bb8["get"](
            "/api/admin/system-settings/require-encryption",
          )
        )["data"],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    }),
    _0xae5e96 = _0x144b3f({
      mutationFn: async (_0x4e3a3d) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/require-encryption",
          _0x4e3a3d,
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["require-encryption"] }),
          _0x54e43f["success"](_0x1aa03d("encryption.updated")));
      },
      onError: _0x4deb0e,
    }),
    [_0x2a5d0a, _0x44faac] = _0x46e13c["useState"](0x3),
    [_0x245fb9, _0x6bd1a7] = _0x46e13c["useState"](0x5),
    [_0x265cb2, _0x14a567] = _0x46e13c["useState"](0x78),
    [_0x1f7866, _0x270780] = _0x46e13c["useState"](0x1e),
    { data: _0x210ced } = _0x37d641({
      queryKey: ["system-intervals"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/system-settings/intervals"))[
          "data"
        ],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    });
  _0x46e13c["useEffect"](() => {
    _0x210ced &&
      (_0x44faac(_0x210ced["speed_collect_interval"]),
      _0x6bd1a7(_0x210ced["report_interval"]),
      _0x14a567(_0x210ced["traffic_check_interval"]),
      _0x270780(_0x210ced["heartbeat_interval"]));
  }, [_0x210ced]);
  const _0xcd2c2c = _0x144b3f({
      mutationFn: async (_0x3d64f8) => {
        await _0x495bb8["put"](
          "/api/admin/system-settings/intervals",
          _0x3d64f8,
        );
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["system-intervals"] }),
          _0x54e43f["success"](_0x1aa03d("intervals.updated")));
      },
      onError: _0x4deb0e,
    }),
    [_0x10f532, _0x19f16a] = _0x46e13c["useState"](!0x1),
    [_0x4710bf, _0x3be089] = _0x46e13c["useState"](""),
    [_0x416f64, _0x341556] = _0x46e13c["useState"](!0x1),
    _0x3fe589 = _0xd1c38e(),
    [_0x26396c, _0x7c81bc] = _0x46e13c["useState"]({
      notify_enabled: !0x1,
      telegram_bot_token: "",
      telegram_chat_id: "",
      notify_login: !0x1,
      notify_subscribe_fetch: !0x1,
      notify_daily_traffic: !0x1,
      notify_server_offline: !0x1,
      notify_server_online: !0x1,
      notify_traffic_threshold: !0x1,
      notify_daily_traffic_time: "08:00",
      notify_traffic_threshold_percent: 0x50,
      notify_traffic_threshold_80: !0x1,
      notify_over_limit: !0x1,
      notify_package_expiring: !0x1,
      notify_package_expiring_days: 0x3,
      notify_package_expired: !0x1,
      notify_user_registered: !0x1,
      notify_telegram_bound: !0x1,
      notify_cert_result: !0x1,
      notify_server_renewal: !0x1,
      notify_agent_long_offline: !0x1,
      notify_agent_long_offline_minutes: 0x1e,
      notify_device_limit_exceeded: !0x1,
      notify_ip_ban: !0x1,
      notify_server_tolerance_seconds: 0x78,
      notify_probe_quality: !0x1,
      probe_jitter_threshold_ms: 0x50,
      probe_loss_threshold_pct: 0x14,
      probe_window_minutes: 0x5,
      probe_min_samples: 0x5,
      probe_trigger_consecutive: 0x2,
      probe_recover_consecutive: 0x2,
      probe_cooldown_minutes: 0x1e,
    }),
    [_0x446de5, _0x1939c3] = _0x46e13c["useState"](""),
    [_0x2854a5, _0x5c956a] = _0x46e13c["useState"](""),
    [_0x4bdc60, _0x14d4c3] = _0x46e13c["useState"](""),
    [_0x680825, _0x594d48] = _0x46e13c["useState"]([]),
    _0x12883a = _0x46e13c["useRef"](!0x1),
    { data: _0x4e4821 } = _0x37d641({
      queryKey: ["notify-config"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/notify-config"))["data"],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    });
  _0x46e13c["useEffect"](() => {
    if (_0x4e4821) {
      const {
        notify_daily_traffic_template: _0x1a3866,
        notify_daily_traffic_template_default: _0x1ae0d9,
        notify_daily_traffic_placeholders: _0x310a61,
        ..._0x171bbb
      } = _0x4e4821;
      (_0x14d4c3(_0x1ae0d9 ?? ""),
        _0x594d48(_0x310a61 ?? []),
        _0x12883a["current"] ||
          ((_0x12883a["current"] = !0x0),
          _0x5c956a(_0x1a3866 || _0x1ae0d9 || "")),
        _0x7c81bc(_0x171bbb),
        _0x1939c3(_0x171bbb["telegram_bot_token"]));
    }
  }, [_0x4e4821]);
  const _0x61c4fc = _0x144b3f({
      mutationFn: async (_0xc5b7d0) => {
        await _0x495bb8["put"]("/api/admin/notify-config", _0xc5b7d0);
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["notify-config"] }),
          _0x54e43f["success"](_0x1aa03d("telegram.configUpdated")));
      },
      onError: _0x4deb0e,
    }),
    _0xeeec1 = _0x144b3f({
      mutationFn: async (_0x4d50d4) => {
        await _0x495bb8["put"]("/api/admin/notify-config", {
          ..._0x26396c,
          notify_daily_traffic_template: _0x4d50d4,
        });
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["notify-config"] }),
          _0x54e43f["success"]("推送文案已保存"));
      },
      onError: _0x4deb0e,
    }),
    _0x440fff = _0x144b3f({
      mutationFn: async (_0x4da366) =>
        (
          await _0x495bb8["post"]("/api/admin/notify-config/preview", {
            template: _0x4da366,
          })
        )["data"],
      onError: _0x4deb0e,
    }),
    _0x46c7e5 = _0x144b3f({
      mutationFn: async () => {
        await _0x495bb8["post"]("/api/admin/notify-config/test");
      },
      onSuccess: () => {
        _0x54e43f["success"](_0x1aa03d("telegram.testSent"));
      },
      onError: _0x4deb0e,
    }),
    _0x24abe3 = (_0x4d90cc) => {
      const _0x5d69c7 = { ..._0x26396c, ..._0x4d90cc };
      (_0x7c81bc(_0x5d69c7), _0x61c4fc["mutate"](_0x5d69c7));
    },
    [_0x585db0, _0x28c5b7] = _0x46e13c["useState"]({
      login_rate_max_attempts: 0x5,
      login_rate_window_minutes: 0x3c,
      login_rate_lock_minutes: 0x3c,
      brute_force_enabled: !0x0,
      brute_force_max_failures: 0x5,
      brute_force_window_minutes: 0x5a0,
      brute_force_block_minutes: 0x5a0,
      sub_rate_enabled: !0x0,
      sub_rate_limit: 0x3c,
      sub_rate_window_minutes: 0x1,
      block_unknown_subscription_ua: !0x1,
      skip_local_ip: !0x0,
      turnstile_site_key: "",
      turnstile_secret_key: "",
    }),
    { data: _0x39d6ec } = _0x37d641({
      queryKey: ["security-settings"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/security-settings"))["data"],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    });
  _0x46e13c["useEffect"](() => {
    _0x39d6ec && _0x28c5b7(_0x39d6ec);
  }, [_0x39d6ec]);
  const _0x5f2acc = _0x144b3f({
      mutationFn: async (_0x5f1cd6) => {
        await _0x495bb8["put"]("/api/admin/security-settings", _0x5f1cd6);
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["security-settings"] }),
          _0x54e43f["success"](_0x1aa03d("autoSaved"), {
            id: "security-config-autosaved",
          }));
      },
      onError: _0x4deb0e,
    }),
    _0x1e3e85 = (_0x44e308) => {
      const _0x33a9d1 = { ..._0x585db0, ..._0x44e308 };
      (_0x28c5b7(_0x33a9d1), _0x5f2acc["mutate"](_0x33a9d1));
    },
    _0x25971f = _0x144b3f({
      mutationFn: async (_0x34d293) => {
        await _0x495bb8["put"]("/api/admin/security-settings", _0x34d293);
      },
      onSuccess: () => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["security-settings"] }),
          _0x54e43f["success"](_0x1aa03d("turnstile.saved")));
      },
      onError: _0x4deb0e,
    }),
    { data: _0x486d96, isLoading: _0xa03ea3 } = _0x37d641({
      queryKey: ["user-config"],
      queryFn: async () => (await _0x495bb8["get"]("/api/user/config"))["data"],
      enabled: !!_0x1bb57f["accessToken"],
      staleTime: 0x12c * 0x3e8,
    });
  _0x46e13c["useEffect"](() => {
    _0x486d96 &&
      (_0xec8154(_0x486d96["force_sync_external"]),
      _0xacab7f(_0x486d96["match_rule"]),
      _0x3874f3(_0x486d96["sync_scope"] || "saved_only"),
      _0x5c52ba(_0x486d96["keep_node_name"] !== !0x1),
      _0x1802db(_0x486d96["cache_expire_minutes"]),
      _0x1df0ca(_0x486d96["sync_traffic"]),
      _0x3b6a2a(
        _0x486d96["node_name_filter"] || "剩余|流量|到期|订阅|时间|重置",
      ),
      _0x495896(!!_0x486d96["append_sub_info"]),
      _0x2d46c8(_0x486d96["use_new_template_system"] !== !0x1),
      _0x19f16a(_0x486d96["enable_proxy_provider"] || !0x1),
      _0x3be089(_0x486d96["proxy_groups_source_url"] || ""),
      _0x341556(_0x486d96["client_compatibility_mode"] || !0x1),
      _0x281a92(!!_0x486d96["enable_sub_info_nodes"]),
      _0x43ecb0(!!_0x486d96["sub_info_v2ray_only"]),
      _0xa30760(_0x486d96["sub_info_expire_prefix"] || "📅过期时间"),
      _0x18df6d(_0x486d96["sub_info_traffic_prefix"] || "⌛剩余流量"));
  }, [_0x486d96]);
  const _0x3002e9 = _0x144b3f({
      mutationFn: async (_0x2af40c) => {
        await _0x495bb8["put"]("/api/user/config", _0x2af40c);
      },
      onSuccess: (_0x452dab, _0x48b236) => {
        (_0x2f82ce["invalidateQueries"]({ queryKey: ["user-config"] }),
          _0x48b236["enable_short_link"] !== _0x28c051 &&
            _0x2f82ce["invalidateQueries"]({
              queryKey: ["user-subscriptions"],
            }),
          _0xec8154(_0x48b236["force_sync_external"]),
          _0xacab7f(_0x48b236["match_rule"]),
          _0x3874f3(_0x48b236["sync_scope"]),
          _0x5c52ba(_0x48b236["keep_node_name"]),
          _0x1802db(_0x48b236["cache_expire_minutes"]),
          _0x1df0ca(_0x48b236["sync_traffic"]),
          _0x3b6a2a(_0x48b236["node_name_filter"]),
          _0x495896(_0x48b236["append_sub_info"]),
          _0x521715(_0x48b236["enable_short_link"]),
          _0x2d46c8(_0x48b236["use_new_template_system"]),
          _0x19f16a(_0x48b236["enable_proxy_provider"]),
          _0x3be089(_0x48b236["proxy_groups_source_url"] || ""),
          _0x341556(_0x48b236["client_compatibility_mode"]),
          _0x281a92(_0x48b236["enable_sub_info_nodes"]),
          _0x43ecb0(_0x48b236["sub_info_v2ray_only"]),
          _0xa30760(_0x48b236["sub_info_expire_prefix"]),
          _0x18df6d(_0x48b236["sub_info_traffic_prefix"]),
          _0x54e43f["success"](_0x1aa03d("configUpdated")));
      },
      onError: (_0x559b31) => {
        (_0x4deb0e(_0x559b31),
          _0x54e43f["error"](_0x1aa03d("configUpdateFailed")));
      },
    }),
    _0x5637c8 = (_0x5c32ee) => {
      _0x3002e9["mutate"]({
        force_sync_external: _0x137d5c,
        match_rule: _0x53c069,
        sync_scope: _0x171e10,
        keep_node_name: _0x2b23a2,
        cache_expire_minutes: _0x2376de,
        sync_traffic: _0x49b722,
        node_name_filter: _0x3cc118,
        append_sub_info: _0x431d4a,
        enable_short_link: _0x28c051,
        use_new_template_system: _0x5e10d1,
        enable_proxy_provider: _0x10f532,
        proxy_groups_source_url: _0x4710bf,
        client_compatibility_mode: _0x416f64,
        enable_sub_info_nodes: _0x18a19f,
        sub_info_v2ray_only: _0x1e8a26,
        sub_info_expire_prefix: _0xb9053e,
        sub_info_traffic_prefix: _0x780f57,
        ..._0x5c32ee,
      });
    };
  return _0x2ddfe2["jsxs"]("div", {
    className: "bg-background\x20min-h-svh",
    children: [
      _0x2ddfe2["jsx"](_0x5d8e9a, {}),
      _0x2ddfe2["jsxs"]("main", {
        className:
          "mx-auto\x20w-full\x20max-w-4xl\x20px-4\x20py-8\x20pt-24\x20sm:px-6",
        children: [
          _0x2ddfe2["jsxs"]("section", {
            className: "space-y-2",
            children: [
              _0x2ddfe2["jsx"]("h1", {
                className: "text-3xl\x20font-semibold\x20tracking-tight",
                children: _0x1aa03d("title"),
              }),
              _0x2ddfe2["jsx"]("p", {
                className: "text-muted-foreground",
                children: _0x1aa03d("description"),
              }),
            ],
          }),
          _0x2ddfe2["jsxs"](_0x2b16ce, {
            value: _0x4873e7,
            onValueChange: _0x1010c7,
            className: "mt-8",
            children: [
              _0x2ddfe2["jsx"]("div", {
                className: "mb-4\x20sm:hidden",
                children: _0x2ddfe2["jsxs"](_0x3ba40b, {
                  value: _0x4873e7,
                  onValueChange: _0x1010c7,
                  children: [
                    _0x2ddfe2["jsx"](_0x3d4a68, {
                      className: "w-full",
                      children: _0x2ddfe2["jsx"](_0x1929c7, {}),
                    }),
                    _0x2ddfe2["jsx"](_0x51ae8f, {
                      children: fn["map"]((_0x50761f) =>
                        _0x2ddfe2["jsx"](
                          _0x1260f9,
                          {
                            value: _0x50761f["value"],
                            children: _0x1aa03d(_0x50761f["labelKey"]),
                          },
                          _0x50761f["value"],
                        ),
                      ),
                    }),
                  ],
                }),
              }),
              _0x2ddfe2["jsx"](_0x383cce, {
                className:
                  "mb-6\x20hidden\x20h-auto\x20w-full\x20flex-wrap\x20sm:flex",
                children: fn["map"]((_0x123439) =>
                  _0x2ddfe2["jsx"](
                    _0x3102ee,
                    {
                      value: _0x123439["value"],
                      children: _0x1aa03d(_0x123439["labelKey"]),
                    },
                    _0x123439["value"],
                  ),
                ),
              }),
              _0x2ddfe2["jsx"](_0x550af3, {
                value: "sub",
                className: "space-y-6",
                children: _0x2ddfe2["jsxs"](_0x2aeb40, {
                  children: [
                    _0x2ddfe2["jsxs"](_0x1db8ce, {
                      className: "pb-4",
                      children: [
                        _0x2ddfe2["jsx"](_0x30c50e, {
                          children: _0x1aa03d("sync.title"),
                        }),
                        _0x2ddfe2["jsx"](_0x54661d, {
                          children: _0x1aa03d("sync.description"),
                        }),
                      ],
                    }),
                    _0x2ddfe2["jsxs"](_0x42cb32, {
                      className: "space-y-4",
                      children: [
                        _0x2ddfe2["jsxs"]("div", {
                          className: "flex\x20items-center\x20justify-between",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className: "flex\x20items-center\x20gap-2",
                              children: [
                                _0x2ddfe2["jsx"](_0x34df34, {
                                  htmlFor: "sync-traffic",
                                  className: "cursor-pointer",
                                  children: _0x1aa03d("sync.syncTraffic"),
                                }),
                                _0x2ddfe2["jsxs"](_0x30391b, {
                                  children: [
                                    _0x2ddfe2["jsx"](_0x2b1216, {
                                      asChild: !0x0,
                                      children: _0x2ddfe2["jsx"](_0x1e002a, {
                                        className:
                                          "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                      }),
                                    }),
                                    _0x2ddfe2["jsx"](_0x42f438, {
                                      side: "right",
                                      className: "max-w-xs",
                                      children: _0x2ddfe2["jsx"]("p", {
                                        children: _0x1aa03d(
                                          "sync.syncTrafficHint",
                                        ),
                                      }),
                                    }),
                                  ],
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsx"](_0x543b01, {
                              id: "sync-traffic",
                              checked: _0x49b722,
                              onCheckedChange: (_0x25e833) =>
                                _0x5637c8({ sync_traffic: _0x25e833 }),
                              disabled: _0x3002e9["isPending"],
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className: "space-y-2\x20border-t\x20pt-3",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className: "flex\x20items-center\x20gap-2",
                              children: [
                                _0x2ddfe2["jsx"](_0x34df34, {
                                  htmlFor: "node-name-filter",
                                  children: _0x1aa03d("sync.nodeNameFilter"),
                                }),
                                _0x2ddfe2["jsxs"](_0x30391b, {
                                  children: [
                                    _0x2ddfe2["jsx"](_0x2b1216, {
                                      asChild: !0x0,
                                      children: _0x2ddfe2["jsx"](_0x1e002a, {
                                        className:
                                          "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                      }),
                                    }),
                                    _0x2ddfe2["jsx"](_0x42f438, {
                                      side: "right",
                                      className: "max-w-xs",
                                      children: _0x2ddfe2["jsx"]("p", {
                                        children: _0x1aa03d(
                                          "sync.nodeNameFilterHint",
                                        ),
                                      }),
                                    }),
                                  ],
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsx"](_0x549353, {
                              id: "node-name-filter",
                              value: _0x3cc118,
                              onChange: (_0xdebf34) =>
                                _0x3b6a2a(_0xdebf34["target"]["value"]),
                              onBlur: () =>
                                _0x5637c8({ node_name_filter: _0x3cc118 }),
                              disabled: _0x3002e9["isPending"],
                              placeholder: "剩余|流量|到期|订阅|时间|重置",
                            }),
                            _0x2ddfe2["jsx"]("p", {
                              className: "text-muted-foreground\x20text-xs",
                              children: _0x1aa03d("sync.nodeNameFilterDesc"),
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className:
                            "flex\x20items-center\x20justify-between\x20border-t\x20pt-3",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className: "flex\x20items-center\x20gap-2",
                              children: [
                                _0x2ddfe2["jsx"](_0x34df34, {
                                  htmlFor: "append-sub-info",
                                  className: "cursor-pointer",
                                  children: "节点名称追加订阅信息",
                                }),
                                _0x2ddfe2["jsxs"](_0x30391b, {
                                  children: [
                                    _0x2ddfe2["jsx"](_0x2b1216, {
                                      asChild: !0x0,
                                      children: _0x2ddfe2["jsx"](_0x1e002a, {
                                        className:
                                          "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                      }),
                                    }),
                                    _0x2ddfe2["jsx"](_0x42f438, {
                                      side: "right",
                                      className: "max-w-xs",
                                      children: _0x2ddfe2["jsx"]("p", {
                                        children:
                                          "开启后,同步外部订阅时在节点名称后追加剩余流量和剩余天数,例如:节点名\x20398.22GB📊\x2026Days⏳",
                                      }),
                                    }),
                                  ],
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsx"](_0x543b01, {
                              id: "append-sub-info",
                              checked: _0x431d4a,
                              onCheckedChange: (_0x2f0d0d) =>
                                _0x5637c8({ append_sub_info: _0x2f0d0d }),
                              disabled: _0x3002e9["isPending"],
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className:
                            "flex\x20items-center\x20justify-between\x20border-t\x20pt-3",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className: "flex\x20items-center\x20gap-2",
                              children: [
                                _0x2ddfe2["jsx"](_0x34df34, {
                                  htmlFor: "force-sync-external",
                                  className: "cursor-pointer",
                                  children: _0x1aa03d("sync.forceSyncExternal"),
                                }),
                                _0x2ddfe2["jsxs"](_0x30391b, {
                                  children: [
                                    _0x2ddfe2["jsx"](_0x2b1216, {
                                      asChild: !0x0,
                                      children: _0x2ddfe2["jsx"](_0x1e002a, {
                                        className:
                                          "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                      }),
                                    }),
                                    _0x2ddfe2["jsx"](_0x42f438, {
                                      side: "right",
                                      className: "max-w-xs",
                                      children: _0x2ddfe2["jsx"]("p", {
                                        children: _0x1aa03d(
                                          "sync.forceSyncExternalHint",
                                        ),
                                      }),
                                    }),
                                  ],
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsx"](_0x543b01, {
                              id: "force-sync-external",
                              checked: _0x137d5c,
                              onCheckedChange: (_0x3d6b07) =>
                                _0x5637c8({ force_sync_external: _0x3d6b07 }),
                              disabled: _0x3002e9["isPending"],
                            }),
                          ],
                        }),
                        _0x137d5c &&
                          _0x2ddfe2["jsxs"]("div", {
                            className:
                              "bg-muted/30\x20-mx-6\x20space-y-4\x20rounded-b-lg\x20border-t\x20px-6\x20py-4\x20pt-3",
                            children: [
                              _0x2ddfe2["jsxs"]("div", {
                                className: "space-y-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    children: _0x1aa03d("sync.matchRule"),
                                  }),
                                  _0x2ddfe2["jsxs"](_0x2db305, {
                                    value: _0x53c069,
                                    onValueChange: (_0x4b64a9) => {
                                      (_0xacab7f(_0x4b64a9),
                                        _0x5637c8({ match_rule: _0x4b64a9 }));
                                    },
                                    disabled: _0x3002e9["isPending"],
                                    className: "flex\x20flex-wrap\x20gap-4",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "flex\x20items-center\x20space-x-2",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x356e54, {
                                            value: "node_name",
                                            id: "match-node-name",
                                          }),
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "match-node-name",
                                            className:
                                              "cursor-pointer\x20font-normal",
                                            children: _0x1aa03d(
                                              "sync.matchRuleNodeName",
                                            ),
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "flex\x20items-center\x20space-x-2",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x356e54, {
                                            value: "server_port",
                                            id: "match-server-port",
                                          }),
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "match-server-port",
                                            className:
                                              "cursor-pointer\x20font-normal",
                                            children: _0x1aa03d(
                                              "sync.matchRuleServerPort",
                                            ),
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "flex\x20items-center\x20space-x-2",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x356e54, {
                                            value: "type_server_port",
                                            id: "match-type-server-port",
                                          }),
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "match-type-server-port",
                                            className:
                                              "cursor-pointer\x20font-normal",
                                            children: _0x1aa03d(
                                              "sync.matchRuleTypeServerPort",
                                            ),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "border-border/50\x20space-y-2\x20border-t\x20pt-3",
                                children: [
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    children: _0x1aa03d("sync.syncScope"),
                                  }),
                                  _0x2ddfe2["jsxs"](_0x2db305, {
                                    value: _0x171e10,
                                    onValueChange: (_0x255a5c) => {
                                      (_0x3874f3(_0x255a5c),
                                        _0x5637c8({ sync_scope: _0x255a5c }));
                                    },
                                    disabled: _0x3002e9["isPending"],
                                    className: "flex\x20flex-wrap\x20gap-4",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "flex\x20items-center\x20space-x-2",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x356e54, {
                                            value: "saved_only",
                                            id: "sync-saved-only",
                                          }),
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "sync-saved-only",
                                            className:
                                              "cursor-pointer\x20font-normal",
                                            children: _0x1aa03d(
                                              "sync.syncScopeSavedOnly",
                                            ),
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "flex\x20items-center\x20space-x-2",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x356e54, {
                                            value: "all",
                                            id: "sync-all",
                                          }),
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "sync-all",
                                            className:
                                              "cursor-pointer\x20font-normal",
                                            children:
                                              _0x1aa03d("sync.syncScopeAll"),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "border-border/50\x20flex\x20items-center\x20justify-between\x20border-t\x20pt-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "keep-node-name",
                                        className: "cursor-pointer",
                                        children:
                                          _0x1aa03d("sync.keepNodeName"),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "right",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "sync.keepNodeNameHint",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x543b01, {
                                    id: "keep-node-name",
                                    checked: _0x2b23a2,
                                    onCheckedChange: (_0x169eec) => {
                                      (_0x5c52ba(_0x169eec),
                                        _0x5637c8({
                                          keep_node_name: _0x169eec,
                                        }));
                                    },
                                    disabled: _0x3002e9["isPending"],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "border-border/50\x20space-y-2\x20border-t\x20pt-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "cache-expire-minutes",
                                        children: _0x1aa03d(
                                          "sync.cacheExpireMinutes",
                                        ),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "right",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "sync.cacheExpireMinutesHint",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x549353, {
                                    id: "cache-expire-minutes",
                                    type: "number",
                                    min: "0",
                                    value: _0x2376de,
                                    onChange: (_0x5b67f1) =>
                                      _0x1802db(
                                        parseInt(
                                          _0x5b67f1["target"]["value"],
                                        ) || 0x0,
                                      ),
                                    onBlur: () =>
                                      _0x5637c8({
                                        cache_expire_minutes: _0x2376de,
                                      }),
                                    disabled: _0x3002e9["isPending"],
                                    placeholder: "0",
                                    className: "w-32",
                                  }),
                                  _0x2ddfe2["jsx"]("p", {
                                    className: "text-destructive\x20text-xs",
                                    children: _0x1aa03d(
                                      "sync.cacheExpireWarning",
                                    ),
                                  }),
                                ],
                              }),
                            ],
                          }),
                      ],
                    }),
                  ],
                }),
              }),
              _0x2ddfe2["jsxs"](_0x550af3, {
                value: "features",
                className: "space-y-6",
                children: [
                  _0x2ddfe2["jsx"](Vr, {}),
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsx"](_0x1db8ce, {
                        className: "pb-4",
                        children: _0x2ddfe2["jsxs"]("div", {
                          className:
                            "flex\x20items-center\x20justify-between\x20gap-4",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className: "space-y-1",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "flex\x20items-center\x20gap-2",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x30c50e, {
                                      children: "订阅信息节点",
                                    }),
                                    _0x2ddfe2["jsxs"](_0x30391b, {
                                      children: [
                                        _0x2ddfe2["jsx"](_0x2b1216, {
                                          asChild: !0x0,
                                          children: _0x2ddfe2["jsx"](
                                            _0x1e002a,
                                            {
                                              className:
                                                "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                            },
                                          ),
                                        }),
                                        _0x2ddfe2["jsx"](_0x42f438, {
                                          side: "right",
                                          className: "max-w-xs",
                                          children: _0x2ddfe2["jsx"]("p", {
                                            children:
                                              "开启后，在\x20Clash\x20订阅节点列表顶部添加过期时间和剩余流量信息节点。",
                                          }),
                                        }),
                                      ],
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsx"](_0x54661d, {
                                  children:
                                    "自定义订阅中显示的过期时间与剩余流量节点名称。",
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsx"](_0x543b01, {
                              id: "enable-sub-info-nodes",
                              checked: _0x18a19f,
                              onCheckedChange: (_0x5e349e) =>
                                _0x5637c8({ enable_sub_info_nodes: _0x5e349e }),
                              disabled: _0x3002e9["isPending"],
                              "aria-label": "启用订阅信息节点",
                            }),
                          ],
                        }),
                      }),
                      _0x18a19f &&
                        _0x2ddfe2["jsxs"](_0x42cb32, {
                          className:
                            "grid\x20gap-4\x20border-t\x20pt-4\x20sm:grid-cols-2",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className:
                                "flex\x20items-center\x20justify-between\x20gap-4\x20sm:col-span-2",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "space-y-1",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      htmlFor: "sub-info-v2ray-only",
                                      children:
                                        "仅对\x20V2Ray\x20系列客户端生效",
                                    }),
                                    _0x2ddfe2["jsx"]("p", {
                                      className:
                                        "text-muted-foreground\x20text-xs",
                                      children:
                                        "开启后，只有\x20V2Ray、V2RayN、V2RayNG\x20等客户端请求订阅时才输出这两个信息节点。",
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsx"](_0x543b01, {
                                  id: "sub-info-v2ray-only",
                                  checked: _0x1e8a26,
                                  onCheckedChange: (_0x55e284) =>
                                    _0x5637c8({
                                      sub_info_v2ray_only: _0x55e284,
                                    }),
                                  disabled: _0x3002e9["isPending"],
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x2ddfe2["jsx"](_0x34df34, {
                                  htmlFor: "sub-info-expire-prefix",
                                  children: "过期时间前缀",
                                }),
                                _0x2ddfe2["jsx"](_0x549353, {
                                  id: "sub-info-expire-prefix",
                                  value: _0xb9053e,
                                  onChange: (_0xfce65b) =>
                                    _0xa30760(_0xfce65b["target"]["value"]),
                                  onBlur: () =>
                                    _0x5637c8({
                                      sub_info_expire_prefix: _0xb9053e,
                                    }),
                                  disabled: _0x3002e9["isPending"],
                                  placeholder: "📅过期时间",
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x2ddfe2["jsx"](_0x34df34, {
                                  htmlFor: "sub-info-traffic-prefix",
                                  children: "剩余流量前缀",
                                }),
                                _0x2ddfe2["jsx"](_0x549353, {
                                  id: "sub-info-traffic-prefix",
                                  value: _0x780f57,
                                  onChange: (_0x19f6c7) =>
                                    _0x18df6d(_0x19f6c7["target"]["value"]),
                                  onBlur: () =>
                                    _0x5637c8({
                                      sub_info_traffic_prefix: _0x780f57,
                                    }),
                                  disabled: _0x3002e9["isPending"],
                                  placeholder: "⌛剩余流量",
                                }),
                              ],
                            }),
                          ],
                        }),
                    ],
                  }),
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsxs"](_0x1db8ce, {
                        className: "pb-4",
                        children: [
                          _0x2ddfe2["jsx"](_0x30c50e, {
                            children: _0x1aa03d("title"),
                          }),
                          _0x2ddfe2["jsx"](_0x54661d, {
                            children: _0x1aa03d("description"),
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsxs"](_0x42cb32, {
                        children: [
                          _0x2ddfe2["jsxs"]("div", {
                            className:
                              "grid\x20grid-cols-1\x20gap-4\x20sm:grid-cols-2",
                            children: [
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "short-link-toggle",
                                        className: "cursor-pointer",
                                        children: _0x1aa03d(
                                          "shortLink.enableLabel",
                                        ),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "top",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "shortLink.description",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x543b01, {
                                    id: "short-link-toggle",
                                    checked: _0x28c051,
                                    onCheckedChange: (_0x340a65) => {
                                      (_0x521715(_0x340a65),
                                        _0x2010e0["mutate"](_0x340a65));
                                    },
                                    disabled: _0x2010e0["isPending"],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "override-scripts-toggle",
                                        className: "cursor-pointer",
                                        children: _0x1aa03d(
                                          "overrideScripts.enableLabel",
                                        ),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "top",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "overrideScripts.description",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x543b01, {
                                    id: "override-scripts-toggle",
                                    checked: _0x465b79,
                                    onCheckedChange: (_0x5ddc66) => {
                                      (_0x71f3e(_0x5ddc66),
                                        _0x407318["mutate"](_0x5ddc66));
                                    },
                                    disabled: _0x407318["isPending"],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "update-cdn-toggle",
                                        className: "cursor-pointer",
                                        children: _0x1aa03d(
                                          "updateCDN.enableLabel",
                                        ),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "top",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "updateCDN.description",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x543b01, {
                                    id: "update-cdn-toggle",
                                    checked: _0x47c4df,
                                    onCheckedChange: (_0x37bc26) => {
                                      (_0x2dd02f(_0x37bc26),
                                        _0x431a9b["mutate"](_0x37bc26));
                                    },
                                    disabled: _0x431a9b["isPending"],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        className: "cursor-default",
                                        children: _0x1aa03d(
                                          "subscriptionOutputFormat.label",
                                        ),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "top",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "subscriptionOutputFormat.description",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"]("div", {
                                    className: "flex\x20gap-1",
                                    children: ["yaml", "json"]["map"](
                                      (_0x473a9e) =>
                                        _0x2ddfe2["jsx"](
                                          "button",
                                          {
                                            type: "button",
                                            onClick: () => {
                                              (_0x49d02f(_0x473a9e),
                                                _0x16cc3a["mutate"](_0x473a9e));
                                            },
                                            disabled: _0x16cc3a["isPending"],
                                            className:
                                              "rounded-md\x20border\x20px-3\x20py-1\x20text-xs\x20transition-colors\x20" +
                                              (_0x170746 === _0x473a9e
                                                ? "bg-primary\x20text-primary-foreground\x20border-primary"
                                                : "bg-background\x20hover:bg-muted\x20border-border"),
                                            children:
                                              _0x473a9e["toUpperCase"](),
                                          },
                                          _0x473a9e,
                                        ),
                                    ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](Br, {}),
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "mmw-features-toggle",
                                        className: "cursor-pointer",
                                        children: _0x1aa03d(
                                          "miaomiaowuFeatures.enableLabel",
                                        ),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "top",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "miaomiaowuFeatures.description",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x543b01, {
                                    id: "mmw-features-toggle",
                                    checked: _0x1ec742,
                                    onCheckedChange: (_0x3b1992) => {
                                      (_0x44237a(_0x3b1992),
                                        _0x146a37["mutate"](_0x3b1992));
                                    },
                                    disabled: _0x146a37["isPending"],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        className: "cursor-default",
                                        children: "从旧版迁移",
                                      }),
                                      _0x2ddfe2["jsx"]("span", {
                                        className:
                                          "text-muted-foreground\x20text-xs",
                                        children:
                                          "一次性把旧版数据导入到当前实例",
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x5185a8, {
                                    variant: "outline",
                                    size: "sm",
                                    asChild: !0x0,
                                    children: _0x2ddfe2["jsx"](_0x51d33c, {
                                      to: "/migrate-from-mmw",
                                      children: "打开迁移向导",
                                    }),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "user-routed-outbound-toggle",
                                        className: "cursor-pointer",
                                        children: "允许用户创建路由出站",
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "top",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children:
                                                "开启后,普通用户可在自己可见的节点上创建路由出站(routed_owner=user),不依赖套餐分配。",
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "flex\x20flex-col\x20items-end\x20gap-0.5",
                                        children: [
                                          _0x2ddfe2["jsx"]("span", {
                                            className:
                                              "text-muted-foreground\x20text-[10px]",
                                            children: "数量",
                                          }),
                                          _0x2ddfe2["jsx"](_0x549353, {
                                            type: "number",
                                            min: 0x1,
                                            className: "h-7\x20w-16\x20text-xs",
                                            value: _0x2ac809,
                                            onChange: (_0x1a81f0) => {
                                              const _0xc23f1b = Math["max"](
                                                0x1,
                                                parseInt(
                                                  _0x1a81f0["target"]["value"],
                                                  0xa,
                                                ) || 0x0,
                                              );
                                              _0x5311a8(_0xc23f1b);
                                            },
                                            onBlur: () => {
                                              _0x40888d &&
                                                _0x32b255["mutate"]({
                                                  enabled: !0x0,
                                                  quota: _0x2ac809,
                                                  daily: _0x4a293d,
                                                });
                                            },
                                            disabled: !_0x40888d,
                                            placeholder: "2",
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "flex\x20flex-col\x20items-end\x20gap-0.5",
                                        children: [
                                          _0x2ddfe2["jsx"]("span", {
                                            className:
                                              "text-muted-foreground\x20text-[10px]",
                                            children: "每日次数",
                                          }),
                                          _0x2ddfe2["jsx"](_0x549353, {
                                            type: "number",
                                            min: 0x1,
                                            className: "h-7\x20w-16\x20text-xs",
                                            value: _0x4a293d,
                                            onChange: (_0x11c24d) => {
                                              const _0xb60516 = Math["max"](
                                                0x1,
                                                parseInt(
                                                  _0x11c24d["target"]["value"],
                                                  0xa,
                                                ) || 0x0,
                                              );
                                              _0x59d309(_0xb60516);
                                            },
                                            onBlur: () => {
                                              _0x40888d &&
                                                _0x32b255["mutate"]({
                                                  enabled: !0x0,
                                                  quota: _0x2ac809,
                                                  daily: _0x4a293d,
                                                });
                                            },
                                            disabled: !_0x40888d,
                                            placeholder: "5",
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsx"](_0x543b01, {
                                        id: "user-routed-outbound-toggle",
                                        checked: _0x40888d,
                                        onCheckedChange: (_0x290de8) => {
                                          (_0x2f492d(_0x290de8),
                                            _0x32b255["mutate"]({
                                              enabled: _0x290de8,
                                              quota: _0x2ac809,
                                              daily: _0x4a293d,
                                            }));
                                        },
                                        disabled: _0x32b255["isPending"],
                                      }),
                                    ],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20border-orange-500/30\x20bg-orange-500/5\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "silent-mode-toggle",
                                        className: "cursor-pointer",
                                        children: _0x1aa03d(
                                          "silentMode.enableLabel",
                                        ),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "top",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "silentMode.description",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x543b01, {
                                    id: "silent-mode-toggle",
                                    checked: _0x4a3bd8,
                                    onCheckedChange: (_0x5bff9f) => {
                                      if (_0x5bff9f) {
                                        _0x31e9f0(!0x0);
                                        return;
                                      }
                                      (_0x451e75(!0x1),
                                        _0x842682["mutate"]({
                                          silent_mode: !0x1,
                                          silent_mode_timeout: _0x1f6fa6,
                                        }));
                                    },
                                    disabled: _0x842682["isPending"],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "flex\x20items-center\x20gap-2",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "require-encryption-toggle",
                                        className: "cursor-pointer",
                                        children: _0x1aa03d(
                                          "encryption.enableLabel",
                                        ),
                                      }),
                                      _0x2ddfe2["jsxs"](_0x30391b, {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x2b1216, {
                                            asChild: !0x0,
                                            children: _0x2ddfe2["jsx"](
                                              _0x1e002a,
                                              {
                                                className:
                                                  "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                              },
                                            ),
                                          }),
                                          _0x2ddfe2["jsx"](_0x42f438, {
                                            side: "top",
                                            className: "max-w-xs",
                                            children: _0x2ddfe2["jsx"]("p", {
                                              children: _0x1aa03d(
                                                "encryption.description",
                                              ),
                                            }),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x543b01, {
                                    id: "require-encryption-toggle",
                                    checked:
                                      _0x5b58be?.["require_encryption"] ?? !0x1,
                                    onCheckedChange: (_0x1a483c) => {
                                      _0xae5e96["mutate"]({
                                        require_encryption: _0x1a483c,
                                      });
                                    },
                                    disabled: _0xae5e96["isPending"],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsx"]("div", {
                                className:
                                  "rounded-lg\x20border\x20p-3\x20sm:col-span-2",
                                children: _0x2ddfe2["jsxs"]("div", {
                                  className:
                                    "flex\x20flex-wrap\x20items-center\x20justify-between\x20gap-3",
                                  children: [
                                    _0x2ddfe2["jsxs"]("div", {
                                      className:
                                        "flex\x20min-w-0\x20items-center\x20gap-2",
                                      children: [
                                        _0x2ddfe2["jsx"](_0x34df34, {
                                          htmlFor: "node-mult-prefix-toggle",
                                          className:
                                            "cursor-pointer\x20whitespace-nowrap",
                                          children: "节点名称追加倍率",
                                        }),
                                        _0x2ddfe2["jsxs"](_0x30391b, {
                                          children: [
                                            _0x2ddfe2["jsx"](_0x2b1216, {
                                              asChild: !0x0,
                                              children: _0x2ddfe2["jsx"](
                                                _0x1e002a,
                                                {
                                                  className:
                                                    "text-muted-foreground\x20h-4\x20w-4\x20shrink-0\x20cursor-help",
                                                },
                                              ),
                                            }),
                                            _0x2ddfe2["jsx"](_0x42f438, {
                                              side: "top",
                                              className: "max-w-xs",
                                              children: _0x2ddfe2["jsxs"]("p", {
                                                children: [
                                                  "开启后,订阅生成时套餐内倍率\x20≠\x201\x20的节点会在名称前加上",
                                                  "\x20",
                                                  _0x2ddfe2["jsxs"]("code", {
                                                    children: [
                                                      _0x1fdf38,
                                                      "<倍率>",
                                                      _0x38ca4b,
                                                    ],
                                                  }),
                                                  "。例如倍率\x202\x20→",
                                                  "\x20",
                                                  _0x2ddfe2["jsx"]("code", {
                                                    children:
                                                      _0x1fdf38 +
                                                      "2" +
                                                      _0x38ca4b +
                                                      "原节点名",
                                                  }),
                                                  "。默认关闭。",
                                                ],
                                              }),
                                            }),
                                          ],
                                        }),
                                      ],
                                    }),
                                    _0x2ddfe2["jsxs"]("div", {
                                      className:
                                        "flex\x20flex-wrap\x20items-center\x20justify-end\x20gap-2",
                                      children: [
                                        _0x256337 &&
                                          _0x2ddfe2["jsxs"](
                                            _0x2ddfe2["Fragment"],
                                            {
                                              children: [
                                                _0x2ddfe2["jsx"](_0x549353, {
                                                  id: "node-mult-prefix-left",
                                                  "aria-label": "左分隔符",
                                                  value: _0x1fdf38,
                                                  maxLength: 0x4,
                                                  onChange: (_0x51cf22) =>
                                                    _0x5eddf2(
                                                      _0x51cf22["target"][
                                                        "value"
                                                      ],
                                                    ),
                                                  onBlur: () =>
                                                    _0x34a342["mutate"]({
                                                      enabled: _0x256337,
                                                      left: _0x1fdf38,
                                                      right: _0x38ca4b,
                                                    }),
                                                  className:
                                                    "h-8\x20w-14\x20text-center",
                                                }),
                                                _0x2ddfe2["jsx"]("span", {
                                                  className:
                                                    "text-muted-foreground\x20text-xs\x20tabular-nums\x20select-none",
                                                  children: "2",
                                                }),
                                                _0x2ddfe2["jsx"](_0x549353, {
                                                  id: "node-mult-prefix-right",
                                                  "aria-label": "右分隔符",
                                                  value: _0x38ca4b,
                                                  maxLength: 0x4,
                                                  onChange: (_0x5182f2) =>
                                                    _0xb1b77b(
                                                      _0x5182f2["target"][
                                                        "value"
                                                      ],
                                                    ),
                                                  onBlur: () =>
                                                    _0x34a342["mutate"]({
                                                      enabled: _0x256337,
                                                      left: _0x1fdf38,
                                                      right: _0x38ca4b,
                                                    }),
                                                  className:
                                                    "h-8\x20w-14\x20text-center",
                                                }),
                                                _0x2ddfe2["jsx"]("span", {
                                                  className:
                                                    "text-muted-foreground\x20text-xs\x20whitespace-nowrap\x20select-none",
                                                  children: "节点名",
                                                }),
                                                _0x2ddfe2["jsxs"]("span", {
                                                  className:
                                                    "text-muted-foreground\x20text-xs\x20whitespace-nowrap",
                                                  children: [
                                                    "预览:",
                                                    _0x2ddfe2["jsxs"]("span", {
                                                      className: "font-mono",
                                                      children: [
                                                        _0x1fdf38,
                                                        "2",
                                                        _0x38ca4b,
                                                        "🇯🇵\x20日本节点",
                                                      ],
                                                    }),
                                                  ],
                                                }),
                                              ],
                                            },
                                          ),
                                        _0x2ddfe2["jsx"](_0x543b01, {
                                          id: "node-mult-prefix-toggle",
                                          checked: _0x256337,
                                          onCheckedChange: (_0xb60c77) => {
                                            (_0x5901ee(_0xb60c77),
                                              _0x34a342["mutate"]({
                                                enabled: _0xb60c77,
                                                left: _0x1fdf38,
                                                right: _0x38ca4b,
                                              }));
                                          },
                                          disabled: _0x34a342["isPending"],
                                        }),
                                      ],
                                    }),
                                  ],
                                }),
                              }),
                            ],
                          }),
                          _0x4a3bd8 &&
                            _0x2ddfe2["jsxs"]("div", {
                              className:
                                "mt-4\x20space-y-2\x20rounded-lg\x20border\x20border-orange-500/30\x20bg-orange-500/5\x20p-3",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "flex\x20items-center\x20gap-2",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      htmlFor: "silent-mode-timeout",
                                      children: _0x1aa03d("silentMode.timeout"),
                                    }),
                                    _0x2ddfe2["jsxs"](_0x30391b, {
                                      children: [
                                        _0x2ddfe2["jsx"](_0x2b1216, {
                                          asChild: !0x0,
                                          children: _0x2ddfe2["jsx"](
                                            _0x1e002a,
                                            {
                                              className:
                                                "text-muted-foreground\x20h-4\x20w-4\x20cursor-help",
                                            },
                                          ),
                                        }),
                                        _0x2ddfe2["jsx"](_0x42f438, {
                                          side: "top",
                                          className: "max-w-xs",
                                          children: _0x2ddfe2["jsx"]("p", {
                                            children: _0x1aa03d(
                                              "silentMode.hint",
                                              { timeout: _0x1f6fa6 },
                                            ),
                                          }),
                                        }),
                                      ],
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsx"](_0x549353, {
                                  id: "silent-mode-timeout",
                                  type: "number",
                                  min: 0x1,
                                  max: 0x5a0,
                                  value: _0x1f6fa6,
                                  onChange: (_0x1b1ad7) =>
                                    _0x9f6337(
                                      parseInt(_0x1b1ad7["target"]["value"]) ||
                                        0xf,
                                    ),
                                  onBlur: () =>
                                    _0x842682["mutate"]({
                                      silent_mode: _0x4a3bd8,
                                      silent_mode_timeout: _0x1f6fa6,
                                    }),
                                  disabled: _0x842682["isPending"],
                                  className: "max-w-32",
                                }),
                              ],
                            }),
                          _0x5b58be?.["require_encryption"] &&
                            _0x2ddfe2["jsx"]("p", {
                              className:
                                "mt-1\x20text-xs\x20text-amber-600\x20dark:text-amber-400",
                              children: _0x1aa03d("encryption.warning"),
                            }),
                        ],
                      }),
                    ],
                  }),
                ],
              }),
              _0x2ddfe2["jsxs"](_0x550af3, {
                value: "push",
                className: "space-y-6",
                children: [
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsxs"](_0x1db8ce, {
                        className: "pb-4",
                        children: [
                          _0x2ddfe2["jsxs"](_0x30c50e, {
                            className: "flex\x20items-center\x20gap-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x14801b, {
                                className: "h-5\x20w-5",
                              }),
                              _0x1aa03d("telegram.enableLabel"),
                            ],
                          }),
                          _0x2ddfe2["jsx"](_0x54661d, {
                            children: _0x1aa03d("telegram.description"),
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsxs"](_0x42cb32, {
                        className: "space-y-4",
                        children: [
                          _0x2ddfe2["jsxs"]("div", {
                            className:
                              "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
                            children: [
                              _0x2ddfe2["jsx"](_0x34df34, {
                                htmlFor: "notify-enabled",
                                className: "cursor-pointer",
                                children: _0x1aa03d("telegram.enableLabel"),
                              }),
                              _0x2ddfe2["jsx"](_0x543b01, {
                                id: "notify-enabled",
                                checked: _0x26396c["notify_enabled"],
                                onCheckedChange: (_0x4afb19) =>
                                  _0x24abe3({ notify_enabled: _0x4afb19 }),
                                disabled: _0x61c4fc["isPending"],
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsxs"]("div", {
                            className: "space-y-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x34df34, {
                                htmlFor: "bot-token",
                                children: _0x1aa03d("telegram.botToken"),
                              }),
                              _0x2ddfe2["jsx"](_0x549353, {
                                id: "bot-token",
                                value: _0x446de5,
                                onChange: (_0x5d10f3) =>
                                  _0x1939c3(_0x5d10f3["target"]["value"]),
                                onBlur: () => {
                                  _0x446de5 !==
                                    _0x26396c["telegram_bot_token"] &&
                                    _0x24abe3({
                                      telegram_bot_token: _0x446de5,
                                    });
                                },
                                placeholder: _0x1aa03d(
                                  "telegram.botTokenPlaceholder",
                                ),
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsxs"]("div", {
                            className: "space-y-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x34df34, {
                                htmlFor: "chat-id",
                                children: _0x1aa03d("telegram.chatId"),
                              }),
                              _0x2ddfe2["jsx"](_0x549353, {
                                id: "chat-id",
                                value: _0x26396c["telegram_chat_id"],
                                onChange: (_0x365f06) =>
                                  _0x7c81bc({
                                    ..._0x26396c,
                                    telegram_chat_id:
                                      _0x365f06["target"]["value"],
                                  }),
                                onBlur: () => {
                                  _0x26396c["telegram_chat_id"] !==
                                    _0x4e4821?.["telegram_chat_id"] &&
                                    _0x24abe3({
                                      telegram_chat_id:
                                        _0x26396c["telegram_chat_id"],
                                    });
                                },
                                placeholder: _0x1aa03d(
                                  "telegram.chatIdPlaceholder",
                                ),
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsx"](_0x5185a8, {
                            variant: "outline",
                            size: "sm",
                            className: "w-full",
                            onClick: () => _0x46c7e5["mutate"](),
                            disabled:
                              _0x46c7e5["isPending"] ||
                              !_0x26396c["telegram_bot_token"] ||
                              !_0x26396c["telegram_chat_id"],
                            children: _0x46c7e5["isPending"]
                              ? "..."
                              : _0x1aa03d("telegram.sendTest"),
                          }),
                          _0x2ddfe2["jsxs"]("div", {
                            className: "space-y-2\x20border-t\x20pt-3",
                            children: [
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-login",
                                    checked: _0x26396c["notify_login"],
                                    onCheckedChange: (_0x1f59cc) =>
                                      _0x24abe3({
                                        notify_login: _0x1f59cc === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-login",
                                    className: "cursor-pointer\x20text-sm",
                                    children: _0x1aa03d(
                                      "telegram.events.login",
                                    ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-subscribe",
                                    checked:
                                      _0x26396c["notify_subscribe_fetch"],
                                    onCheckedChange: (_0x5f0f22) =>
                                      _0x24abe3({
                                        notify_subscribe_fetch:
                                          _0x5f0f22 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-subscribe",
                                    className: "cursor-pointer\x20text-sm",
                                    children: _0x1aa03d(
                                      "telegram.events.subscribe",
                                    ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-online",
                                    checked: _0x26396c["notify_server_online"],
                                    onCheckedChange: (_0xd0c022) =>
                                      _0x24abe3({
                                        notify_server_online:
                                          _0xd0c022 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-online",
                                    className: "cursor-pointer\x20text-sm",
                                    children: _0x1aa03d(
                                      "telegram.events.serverOnline",
                                    ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-offline",
                                    checked: _0x26396c["notify_server_offline"],
                                    onCheckedChange: (_0x34ba26) =>
                                      _0x24abe3({
                                        notify_server_offline:
                                          _0x34ba26 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-offline",
                                    className: "cursor-pointer\x20text-sm",
                                    children: _0x1aa03d(
                                      "telegram.events.serverOffline",
                                    ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20flex-wrap\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-server-tolerance",
                                    className:
                                      "cursor-pointer\x20text-sm\x20whitespace-nowrap",
                                    children: _0x1aa03d(
                                      "telegram.events.serverToleranceSeconds",
                                      { defaultValue: "上下线容忍阈值(秒)" },
                                    ),
                                  }),
                                  _0x2ddfe2["jsx"](_0x549353, {
                                    id: "notify-server-tolerance",
                                    type: "number",
                                    min: 0x0,
                                    className: "h-8\x20w-24",
                                    value:
                                      _0x26396c[
                                        "notify_server_tolerance_seconds"
                                      ],
                                    onChange: (_0x1293c7) =>
                                      _0x7c81bc({
                                        ..._0x26396c,
                                        notify_server_tolerance_seconds: Math[
                                          "max"
                                        ](
                                          0x0,
                                          parseInt(
                                            _0x1293c7["target"]["value"],
                                          ) || 0x0,
                                        ),
                                      }),
                                    onBlur: () =>
                                      _0x24abe3({
                                        notify_server_tolerance_seconds:
                                          _0x26396c[
                                            "notify_server_tolerance_seconds"
                                          ],
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"]("span", {
                                    className:
                                      "text-muted-foreground\x20text-xs",
                                    children: _0x1aa03d(
                                      "telegram.events.serverToleranceHint",
                                      {
                                        defaultValue:
                                          "离线满该秒数才发下线通知;阈值内又上线不发(压抖动+主控重启误报)。0=关闭",
                                      },
                                    ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "space-y-3\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className:
                                      "flex\x20items-start\x20justify-between\x20gap-3",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        children: [
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "notify-probe-quality",
                                            className:
                                              "cursor-pointer\x20text-sm\x20font-medium",
                                            children: "探针延迟波动与丢包告警",
                                          }),
                                          _0x2ddfe2["jsx"]("p", {
                                            className:
                                              "text-muted-foreground\x20mt-1\x20text-xs",
                                            children:
                                              "检测探针页面已选择的服务器和目标；连续异常后通知，恢复正常后发送恢复通知。",
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsx"](_0x543b01, {
                                        id: "notify-probe-quality",
                                        checked:
                                          _0x26396c["notify_probe_quality"],
                                        onCheckedChange: (_0x53cee5) =>
                                          _0x24abe3({
                                            notify_probe_quality: _0x53cee5,
                                          }),
                                        disabled: _0x61c4fc["isPending"],
                                      }),
                                    ],
                                  }),
                                  _0x26396c["notify_probe_quality"] &&
                                    _0x2ddfe2["jsx"]("div", {
                                      className:
                                        "grid\x20grid-cols-2\x20gap-3\x20sm:grid-cols-4",
                                      children: [
                                        [
                                          "波动阈值(ms)",
                                          "probe_jitter_threshold_ms",
                                          0x1,
                                          0x1388,
                                        ],
                                        [
                                          "丢包阈值(%)",
                                          "probe_loss_threshold_pct",
                                          0x1,
                                          0x64,
                                        ],
                                        [
                                          "统计窗口(分钟)",
                                          "probe_window_minutes",
                                          0x2,
                                          0x3c,
                                        ],
                                        [
                                          "最少样本",
                                          "probe_min_samples",
                                          0x2,
                                          0x3e8,
                                        ],
                                        [
                                          "异常确认次数",
                                          "probe_trigger_consecutive",
                                          0x1,
                                          0xa,
                                        ],
                                        [
                                          "恢复确认次数",
                                          "probe_recover_consecutive",
                                          0x1,
                                          0xa,
                                        ],
                                        [
                                          "重复提醒间隔(分钟)",
                                          "probe_cooldown_minutes",
                                          0x1,
                                          0x5a0,
                                        ],
                                      ]["map"](
                                        ([
                                          _0x245626,
                                          _0xde747b,
                                          _0x5c5623,
                                          _0x11c3bd,
                                        ]) =>
                                          _0x2ddfe2["jsxs"](
                                            "div",
                                            {
                                              className: "space-y-1",
                                              children: [
                                                _0x2ddfe2["jsx"](_0x34df34, {
                                                  className: "text-xs",
                                                  children: _0x245626,
                                                }),
                                                _0x2ddfe2["jsx"](_0x549353, {
                                                  type: "number",
                                                  min: Number(_0x5c5623),
                                                  max: Number(_0x11c3bd),
                                                  className: "h-8",
                                                  value: _0x26396c[_0xde747b],
                                                  onChange: (_0x3a9fe1) =>
                                                    _0x7c81bc({
                                                      ..._0x26396c,
                                                      [_0xde747b]: Number(
                                                        _0x3a9fe1["target"][
                                                          "value"
                                                        ],
                                                      ),
                                                    }),
                                                  onBlur: () =>
                                                    _0x24abe3({
                                                      [_0xde747b]:
                                                        _0x26396c[_0xde747b],
                                                    }),
                                                }),
                                              ],
                                            },
                                            String(_0xde747b),
                                          ),
                                      ),
                                    }),
                                  _0x2ddfe2["jsx"]("p", {
                                    className:
                                      "text-muted-foreground\x20text-xs",
                                    children:
                                      "延迟波动按最近窗口的\x20P95\x20−\x20P50\x20计算；服务器离线或样本不足不会按丢包告警。",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-daily",
                                    checked: _0x26396c["notify_daily_traffic"],
                                    onCheckedChange: (_0x1738d9) =>
                                      _0x24abe3({
                                        notify_daily_traffic:
                                          _0x1738d9 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-daily",
                                    className: "cursor-pointer\x20text-sm",
                                    children: _0x1aa03d(
                                      "telegram.events.dailyTraffic",
                                    ),
                                  }),
                                  _0x26396c["notify_daily_traffic"] &&
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      type: "time",
                                      value:
                                        _0x26396c["notify_daily_traffic_time"],
                                      onChange: (_0x3886) =>
                                        _0x7c81bc({
                                          ..._0x26396c,
                                          notify_daily_traffic_time:
                                            _0x3886["target"]["value"],
                                        }),
                                      onBlur: () =>
                                        _0x24abe3({
                                          notify_daily_traffic_time:
                                            _0x26396c[
                                              "notify_daily_traffic_time"
                                            ],
                                        }),
                                      className: "h-7\x20w-24\x20text-xs",
                                    }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-threshold",
                                    checked:
                                      _0x26396c["notify_traffic_threshold"],
                                    onCheckedChange: (_0x2e474b) =>
                                      _0x24abe3({
                                        notify_traffic_threshold:
                                          _0x2e474b === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-threshold",
                                    className: "cursor-pointer\x20text-sm",
                                    children: _0x1aa03d(
                                      "telegram.events.trafficThreshold",
                                    ),
                                  }),
                                  _0x26396c["notify_traffic_threshold"] &&
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      type: "number",
                                      min: 0x1,
                                      max: 0x64,
                                      value:
                                        _0x26396c[
                                          "notify_traffic_threshold_percent"
                                        ],
                                      onChange: (_0x517c46) =>
                                        _0x7c81bc({
                                          ..._0x26396c,
                                          notify_traffic_threshold_percent:
                                            Number(
                                              _0x517c46["target"]["value"],
                                            ),
                                        }),
                                      onBlur: () =>
                                        _0x24abe3({
                                          notify_traffic_threshold_percent:
                                            _0x26396c[
                                              "notify_traffic_threshold_percent"
                                            ],
                                        }),
                                      className: "h-7\x20w-16\x20text-xs",
                                    }),
                                ],
                              }),
                              _0x2ddfe2["jsx"]("div", {
                                className:
                                  "text-muted-foreground\x20mt-2\x20border-t\x20pt-2\x20text-xs",
                                children: "用户生命周期\x20/\x20运营事件",
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-threshold-80",
                                    checked:
                                      _0x26396c["notify_traffic_threshold_80"],
                                    onCheckedChange: (_0x583d35) =>
                                      _0x24abe3({
                                        notify_traffic_threshold_80:
                                          _0x583d35 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-threshold-80",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "用户流量达\x2080%\x20预警",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-over-limit",
                                    checked: _0x26396c["notify_over_limit"],
                                    onCheckedChange: (_0x30b26f) =>
                                      _0x24abe3({
                                        notify_over_limit: _0x30b26f === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-over-limit",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "用户流量超\x20100%(已踢)",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-package-expiring",
                                    checked:
                                      _0x26396c["notify_package_expiring"],
                                    onCheckedChange: (_0x15c435) =>
                                      _0x24abe3({
                                        notify_package_expiring:
                                          _0x15c435 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-package-expiring",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "套餐即将到期(N\x20天前)",
                                  }),
                                  _0x26396c["notify_package_expiring"] &&
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      type: "number",
                                      min: 0x1,
                                      max: 0x1e,
                                      value:
                                        _0x26396c[
                                          "notify_package_expiring_days"
                                        ],
                                      onChange: (_0x2602b4) =>
                                        _0x7c81bc({
                                          ..._0x26396c,
                                          notify_package_expiring_days: Number(
                                            _0x2602b4["target"]["value"],
                                          ),
                                        }),
                                      onBlur: () =>
                                        _0x24abe3({
                                          notify_package_expiring_days:
                                            _0x26396c[
                                              "notify_package_expiring_days"
                                            ],
                                        }),
                                      className: "h-7\x20w-16\x20text-xs",
                                    }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-package-expired",
                                    checked:
                                      _0x26396c["notify_package_expired"],
                                    onCheckedChange: (_0x587543) =>
                                      _0x24abe3({
                                        notify_package_expired:
                                          _0x587543 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-package-expired",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "套餐已到期",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-user-registered",
                                    checked:
                                      _0x26396c["notify_user_registered"],
                                    onCheckedChange: (_0x2334c5) =>
                                      _0x24abe3({
                                        notify_user_registered:
                                          _0x2334c5 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-user-registered",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "新用户注册",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-telegram-bound",
                                    checked: _0x26396c["notify_telegram_bound"],
                                    onCheckedChange: (_0x4bd340) =>
                                      _0x24abe3({
                                        notify_telegram_bound:
                                          _0x4bd340 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-telegram-bound",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "首次绑定\x20Telegram",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-cert-result",
                                    checked: _0x26396c["notify_cert_result"],
                                    onCheckedChange: (_0x4f33c0) =>
                                      _0x24abe3({
                                        notify_cert_result: _0x4f33c0 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-cert-result",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "证书申请成败",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-server-renewal",
                                    checked: _0x26396c["notify_server_renewal"],
                                    onCheckedChange: (_0x516f76) =>
                                      _0x24abe3({
                                        notify_server_renewal:
                                          _0x516f76 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-server-renewal",
                                    className: "cursor-pointer\x20text-sm",
                                    children:
                                      "服务器续费(到期前\x207/3\x20天\x20+\x20续费成功)",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-agent-long-offline",
                                    checked:
                                      _0x26396c["notify_agent_long_offline"],
                                    onCheckedChange: (_0x229e2e) =>
                                      _0x24abe3({
                                        notify_agent_long_offline:
                                          _0x229e2e === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-agent-long-offline",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "Agent\x20长期离线(N\x20分钟)",
                                  }),
                                  _0x26396c["notify_agent_long_offline"] &&
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      type: "number",
                                      min: 0x5,
                                      max: 0x5a0,
                                      value:
                                        _0x26396c[
                                          "notify_agent_long_offline_minutes"
                                        ],
                                      onChange: (_0x129a9f) =>
                                        _0x7c81bc({
                                          ..._0x26396c,
                                          notify_agent_long_offline_minutes:
                                            Number(
                                              _0x129a9f["target"]["value"],
                                            ),
                                        }),
                                      onBlur: () =>
                                        _0x24abe3({
                                          notify_agent_long_offline_minutes:
                                            _0x26396c[
                                              "notify_agent_long_offline_minutes"
                                            ],
                                        }),
                                      className: "h-7\x20w-16\x20text-xs",
                                    }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-device-limit",
                                    checked:
                                      _0x26396c["notify_device_limit_exceeded"],
                                    onCheckedChange: (_0x2d169e) =>
                                      _0x24abe3({
                                        notify_device_limit_exceeded:
                                          _0x2d169e === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-device-limit",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "连接数超限(agent\x20上报)",
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "flex\x20items-center\x20gap-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x77832a, {
                                    id: "notify-ip-ban",
                                    checked: _0x26396c["notify_ip_ban"],
                                    onCheckedChange: (_0x301530) =>
                                      _0x24abe3({
                                        notify_ip_ban: _0x301530 === !0x0,
                                      }),
                                  }),
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "notify-ip-ban",
                                    className: "cursor-pointer\x20text-sm",
                                    children: "IP\x20封禁(暴力防护触发)",
                                  }),
                                ],
                              }),
                            ],
                          }),
                        ],
                      }),
                    ],
                  }),
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsxs"](_0x1db8ce, {
                        className: "pb-4",
                        children: [
                          _0x2ddfe2["jsxs"](_0x30c50e, {
                            className: "flex\x20items-center\x20gap-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x6a84da, {
                                className: "h-5\x20w-5",
                              }),
                              "每日推送文案",
                            ],
                          }),
                          _0x2ddfe2["jsx"](_0x54661d, {
                            children:
                              "自定义每日流量及昨日节点/用户增量的推送文案。数据用占位符表示,发送时替换成真实数据。留空则使用默认文案。",
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsxs"](_0x42cb32, {
                        className: "space-y-4",
                        children: [
                          _0x2ddfe2["jsxs"]("div", {
                            className: "space-y-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x34df34, {
                                htmlFor: "daily-template",
                                children: "文案模板",
                              }),
                              _0x2ddfe2["jsx"](_0x5d9440, {
                                id: "daily-template",
                                className:
                                  "min-h-[220px]\x20font-mono\x20text-xs",
                                value: _0x2854a5,
                                onChange: (_0x276934) =>
                                  _0x5c956a(_0x276934["target"]["value"]),
                                spellCheck: !0x1,
                              }),
                              _0x2ddfe2["jsx"]("p", {
                                className: "text-muted-foreground\x20text-xs",
                                children:
                                  "文案按\x20Telegram\x20Markdown\x20渲染:*粗体*、`等宽`。列表为空时,占位符所在的整行会被移除,但它上面的段落标题会保留。",
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsxs"]("div", {
                            className:
                              "space-y-2\x20rounded-lg\x20border\x20p-3",
                            children: [
                              _0x2ddfe2["jsx"]("div", {
                                className: "text-sm\x20font-medium",
                                children: "可用占位符",
                              }),
                              _0x2ddfe2["jsx"]("div", {
                                className: "space-y-1.5",
                                children: _0x680825["map"]((_0x1d4a69) =>
                                  _0x2ddfe2["jsxs"](
                                    "div",
                                    {
                                      className:
                                        "flex\x20flex-col\x20gap-0.5\x20sm:flex-row\x20sm:items-baseline\x20sm:gap-2",
                                      children: [
                                        _0x2ddfe2["jsx"]("button", {
                                          type: "button",
                                          className:
                                            "bg-muted\x20hover:bg-muted-foreground/20\x20shrink-0\x20self-start\x20rounded\x20px-1.5\x20py-0.5\x20font-mono\x20text-xs",
                                          title: "点击插入到文案末尾",
                                          onClick: () =>
                                            _0x5c956a(
                                              (_0x5ecde8) =>
                                                _0x5ecde8 + _0x1d4a69["name"],
                                            ),
                                          children: _0x1d4a69["name"],
                                        }),
                                        _0x2ddfe2["jsx"]("span", {
                                          className:
                                            "text-muted-foreground\x20text-xs",
                                          children: _0x1d4a69["desc"],
                                        }),
                                      ],
                                    },
                                    _0x1d4a69["name"],
                                  ),
                                ),
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsxs"]("div", {
                            className:
                              "flex\x20flex-wrap\x20items-center\x20gap-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x5185a8, {
                                onClick: () => _0xeeec1["mutate"](_0x2854a5),
                                disabled: _0xeeec1["isPending"],
                                children: _0xeeec1["isPending"]
                                  ? "保存中..."
                                  : "保存文案",
                              }),
                              _0x2ddfe2["jsx"](_0x5185a8, {
                                variant: "outline",
                                onClick: () => _0x440fff["mutate"](_0x2854a5),
                                disabled: _0x440fff["isPending"],
                                children: _0x440fff["isPending"]
                                  ? "渲染中..."
                                  : "预览",
                              }),
                              _0x2ddfe2["jsx"](_0x5185a8, {
                                variant: "outline",
                                onClick: () => _0x5c956a(_0x4bdc60),
                                disabled: !_0x4bdc60 || _0x2854a5 === _0x4bdc60,
                                children: "恢复默认",
                              }),
                            ],
                          }),
                          _0x440fff["data"] &&
                            _0x2ddfe2["jsxs"]("div", {
                              className: "space-y-2",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "flex\x20items-center\x20gap-2",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      children: "预览",
                                    }),
                                    _0x2ddfe2["jsx"]("span", {
                                      className:
                                        "text-muted-foreground\x20text-xs",
                                      children: _0x440fff["data"]["sample"]
                                        ? "(暂无真实流量数据,用示例数据渲染)"
                                        : "(用当前真实数据渲染)",
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsx"]("pre", {
                                  className:
                                    "bg-muted/40\x20overflow-x-auto\x20rounded-lg\x20border\x20p-3\x20text-xs\x20whitespace-pre-wrap",
                                  children: _0x440fff["data"]["preview"],
                                }),
                              ],
                            }),
                        ],
                      }),
                    ],
                  }),
                ],
              }),
              _0x2ddfe2["jsx"](_0x550af3, {
                value: "security",
                className: "space-y-6",
                children: _0x2ddfe2["jsxs"](_0x2aeb40, {
                  children: [
                    _0x2ddfe2["jsxs"](_0x1db8ce, {
                      className: "pb-4",
                      children: [
                        _0x2ddfe2["jsx"](_0x30c50e, {
                          children: "自定义安全阈值",
                        }),
                        _0x2ddfe2["jsx"](_0x54661d, {
                          children:
                            "调整登录限流、订阅暴力防护、订阅频率限制的阈值。修改后立即生效,无需重启主控。",
                        }),
                      ],
                    }),
                    _0x2ddfe2["jsxs"](_0x42cb32, {
                      className: "space-y-6",
                      children: [
                        _0x2ddfe2["jsxs"]("div", {
                          className:
                            "flex\x20items-start\x20justify-between\x20gap-4",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className: "flex-1",
                              children: [
                                _0x2ddfe2["jsx"]("div", {
                                  className: "text-sm\x20font-semibold",
                                  children: "不封禁本地\x20IP",
                                }),
                                _0x2ddfe2["jsx"]("p", {
                                  className:
                                    "text-muted-foreground\x20mt-1\x20text-xs",
                                  children:
                                    "反代/Docker\x20场景下,若上游未传\x20X-Forwarded-For,主控可能将所有用户视作同一本机\x20IP\x20—\x20一次封禁会让所有人连不上。\x20开启后,loopback\x20/\x20内网\x20/\x20私有网段\x20IP\x20跳过封禁与频率限制(登录账户维度仍生效)。",
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsx"](_0x543b01, {
                              checked: _0x585db0["skip_local_ip"],
                              onCheckedChange: (_0x953ee7) =>
                                _0x1e3e85({ skip_local_ip: _0x953ee7 }),
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className: "space-y-3\x20border-t\x20pt-4",
                          children: [
                            _0x2ddfe2["jsx"]("div", {
                              className: "text-sm\x20font-semibold",
                              children: "登录限流",
                            }),
                            _0x2ddfe2["jsx"]("p", {
                              className: "text-muted-foreground\x20text-xs",
                              children:
                                "失败次数达上限后,IP/账户在锁定时长内禁止登录。",
                            }),
                            _0x2ddfe2["jsxs"]("div", {
                              className:
                                "grid\x20grid-cols-1\x20gap-3\x20sm:grid-cols-3",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "space-y-1",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      htmlFor: "login-rate-max-attempts",
                                      className: "text-xs",
                                      children: "最大尝试次数",
                                    }),
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      id: "login-rate-max-attempts",
                                      type: "number",
                                      min: 0x1,
                                      value:
                                        _0x585db0["login_rate_max_attempts"],
                                      onChange: (_0x241c34) =>
                                        _0x28c5b7({
                                          ..._0x585db0,
                                          login_rate_max_attempts: Number(
                                            _0x241c34["target"]["value"],
                                          ),
                                        }),
                                      onBlur: () =>
                                        _0x1e3e85({
                                          login_rate_max_attempts:
                                            _0x585db0[
                                              "login_rate_max_attempts"
                                            ],
                                        }),
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "space-y-1",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      htmlFor: "login-rate-window",
                                      className: "text-xs",
                                      children: "统计窗口(分钟)",
                                    }),
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      id: "login-rate-window",
                                      type: "number",
                                      min: 0x1,
                                      value:
                                        _0x585db0["login_rate_window_minutes"],
                                      onChange: (_0x2f947b) =>
                                        _0x28c5b7({
                                          ..._0x585db0,
                                          login_rate_window_minutes: Number(
                                            _0x2f947b["target"]["value"],
                                          ),
                                        }),
                                      onBlur: () =>
                                        _0x1e3e85({
                                          login_rate_window_minutes:
                                            _0x585db0[
                                              "login_rate_window_minutes"
                                            ],
                                        }),
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "space-y-1",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      htmlFor: "login-rate-lock",
                                      className: "text-xs",
                                      children: "锁定时长(分钟)",
                                    }),
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      id: "login-rate-lock",
                                      type: "number",
                                      min: 0x1,
                                      value:
                                        _0x585db0["login_rate_lock_minutes"],
                                      onChange: (_0x5ebd93) =>
                                        _0x28c5b7({
                                          ..._0x585db0,
                                          login_rate_lock_minutes: Number(
                                            _0x5ebd93["target"]["value"],
                                          ),
                                        }),
                                      onBlur: () =>
                                        _0x1e3e85({
                                          login_rate_lock_minutes:
                                            _0x585db0[
                                              "login_rate_lock_minutes"
                                            ],
                                        }),
                                    }),
                                  ],
                                }),
                              ],
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className: "space-y-3\x20border-t\x20pt-4",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className:
                                "flex\x20items-center\x20justify-between",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  children: [
                                    _0x2ddfe2["jsx"]("div", {
                                      className: "text-sm\x20font-semibold",
                                      children: "订阅暴力防护",
                                    }),
                                    _0x2ddfe2["jsx"]("p", {
                                      className:
                                        "text-muted-foreground\x20mt-1\x20text-xs",
                                      children:
                                        "短链接\x20/x/\x20与临时订阅\x20/t/\x20访问失败次数达上限\x20→\x20封禁\x20IP,并发\x20Telegram\x20通知。",
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsx"](_0x543b01, {
                                  checked: _0x585db0["brute_force_enabled"],
                                  onCheckedChange: (_0x435f77) =>
                                    _0x1e3e85({
                                      brute_force_enabled: _0x435f77,
                                    }),
                                }),
                              ],
                            }),
                            _0x585db0["brute_force_enabled"] &&
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "grid\x20grid-cols-1\x20gap-3\x20sm:grid-cols-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "brute-force-max-failures",
                                        className: "text-xs",
                                        children: "最大失败次数",
                                      }),
                                      _0x2ddfe2["jsx"](_0x549353, {
                                        id: "brute-force-max-failures",
                                        type: "number",
                                        min: 0x1,
                                        value:
                                          _0x585db0["brute_force_max_failures"],
                                        onChange: (_0x2c2a67) =>
                                          _0x28c5b7({
                                            ..._0x585db0,
                                            brute_force_max_failures: Number(
                                              _0x2c2a67["target"]["value"],
                                            ),
                                          }),
                                        onBlur: () =>
                                          _0x1e3e85({
                                            brute_force_max_failures:
                                              _0x585db0[
                                                "brute_force_max_failures"
                                              ],
                                          }),
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "brute-force-window",
                                        className: "text-xs",
                                        children: "统计窗口(分钟)",
                                      }),
                                      _0x2ddfe2["jsx"](_0x549353, {
                                        id: "brute-force-window",
                                        type: "number",
                                        min: 0x1,
                                        value:
                                          _0x585db0[
                                            "brute_force_window_minutes"
                                          ],
                                        onChange: (_0x1b951d) =>
                                          _0x28c5b7({
                                            ..._0x585db0,
                                            brute_force_window_minutes: Number(
                                              _0x1b951d["target"]["value"],
                                            ),
                                          }),
                                        onBlur: () =>
                                          _0x1e3e85({
                                            brute_force_window_minutes:
                                              _0x585db0[
                                                "brute_force_window_minutes"
                                              ],
                                          }),
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "brute-force-block",
                                        className: "text-xs",
                                        children: "封禁时长(分钟)",
                                      }),
                                      _0x2ddfe2["jsx"](_0x549353, {
                                        id: "brute-force-block",
                                        type: "number",
                                        min: 0x1,
                                        value:
                                          _0x585db0[
                                            "brute_force_block_minutes"
                                          ],
                                        onChange: (_0x18264d) =>
                                          _0x28c5b7({
                                            ..._0x585db0,
                                            brute_force_block_minutes: Number(
                                              _0x18264d["target"]["value"],
                                            ),
                                          }),
                                        onBlur: () =>
                                          _0x1e3e85({
                                            brute_force_block_minutes:
                                              _0x585db0[
                                                "brute_force_block_minutes"
                                              ],
                                          }),
                                      }),
                                    ],
                                  }),
                                ],
                              }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className:
                            "flex\x20items-start\x20justify-between\x20gap-4\x20border-t\x20pt-4",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className: "flex-1",
                              children: [
                                _0x2ddfe2["jsx"]("div", {
                                  className: "text-sm\x20font-semibold",
                                  children: "禁止浏览器访问订阅",
                                }),
                                _0x2ddfe2["jsx"]("p", {
                                  className:
                                    "text-muted-foreground\x20mt-1\x20text-xs",
                                  children:
                                    "开启后仅允许系统可识别的代理客户端获取订阅，浏览器、聊天软件链接预览、空\x20UA\x20与未知\x20UA\x20将返回\x20403。",
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsx"](_0x543b01, {
                              checked:
                                _0x585db0["block_unknown_subscription_ua"],
                              onCheckedChange: (_0x7048a0) =>
                                _0x1e3e85({
                                  block_unknown_subscription_ua: _0x7048a0,
                                }),
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className: "space-y-3\x20border-t\x20pt-4",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className:
                                "flex\x20items-center\x20justify-between",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  children: [
                                    _0x2ddfe2["jsx"]("div", {
                                      className: "text-sm\x20font-semibold",
                                      children: "订阅频率限制",
                                    }),
                                    _0x2ddfe2["jsx"]("p", {
                                      className:
                                        "text-muted-foreground\x20mt-1\x20text-xs",
                                      children:
                                        "同一\x20IP\x20在窗口内最多访问订阅的次数,超出返回\x20429,防机器抓取。",
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsx"](_0x543b01, {
                                  checked: _0x585db0["sub_rate_enabled"],
                                  onCheckedChange: (_0xa9346e) =>
                                    _0x1e3e85({ sub_rate_enabled: _0xa9346e }),
                                }),
                              ],
                            }),
                            _0x585db0["sub_rate_enabled"] &&
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "grid\x20grid-cols-1\x20gap-3\x20sm:grid-cols-2",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "sub-rate-limit",
                                        className: "text-xs",
                                        children: "最大请求数",
                                      }),
                                      _0x2ddfe2["jsx"](_0x549353, {
                                        id: "sub-rate-limit",
                                        type: "number",
                                        min: 0x1,
                                        value: _0x585db0["sub_rate_limit"],
                                        onChange: (_0x20d408) =>
                                          _0x28c5b7({
                                            ..._0x585db0,
                                            sub_rate_limit: Number(
                                              _0x20d408["target"]["value"],
                                            ),
                                          }),
                                        onBlur: () =>
                                          _0x1e3e85({
                                            sub_rate_limit:
                                              _0x585db0["sub_rate_limit"],
                                          }),
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "sub-rate-window",
                                        className: "text-xs",
                                        children: "统计窗口(分钟)",
                                      }),
                                      _0x2ddfe2["jsx"](_0x549353, {
                                        id: "sub-rate-window",
                                        type: "number",
                                        min: 0x1,
                                        value:
                                          _0x585db0["sub_rate_window_minutes"],
                                        onChange: (_0x4b386b) =>
                                          _0x28c5b7({
                                            ..._0x585db0,
                                            sub_rate_window_minutes: Number(
                                              _0x4b386b["target"]["value"],
                                            ),
                                          }),
                                        onBlur: () =>
                                          _0x1e3e85({
                                            sub_rate_window_minutes:
                                              _0x585db0[
                                                "sub_rate_window_minutes"
                                              ],
                                          }),
                                      }),
                                    ],
                                  }),
                                ],
                              }),
                          ],
                        }),
                      ],
                    }),
                  ],
                }),
              }),
              _0x2ddfe2["jsx"](_0x550af3, {
                value: "probe",
                className: "space-y-6",
                children: _0x2ddfe2["jsxs"](_0x2aeb40, {
                  children: [
                    _0x2ddfe2["jsxs"](_0x1db8ce, {
                      className: "pb-4",
                      children: [
                        _0x2ddfe2["jsx"](_0x30c50e, { children: "伪装探针" }),
                        _0x2ddfe2["jsx"](_0x54661d, {
                          children:
                            "把面板首页伪装成公开的服务器监控页,右上角保留隐蔽登录入口",
                        }),
                      ],
                    }),
                    _0x2ddfe2["jsxs"](_0x42cb32, {
                      className: "space-y-4",
                      children: [
                        _0x2ddfe2["jsxs"]("div", {
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className:
                                "flex\x20items-start\x20justify-between\x20gap-4",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "flex-1",
                                  children: [
                                    _0x2ddfe2["jsx"]("div", {
                                      className: "text-sm\x20font-semibold",
                                      children: "伪装成探针",
                                    }),
                                    _0x2ddfe2["jsx"]("p", {
                                      className:
                                        "text-muted-foreground\x20mt-1\x20text-xs",
                                      children:
                                        "开启后,未登录访客访问域名首页会看到一个只读的公开服务器状态页(仅显示网速/流量/重置日期/在线状态,无任何编辑按钮),右上角有隐蔽登录入口。用于把订阅管理面板伪装成无害的监控站。已登录的管理员不受影响,仍直接进后台。",
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsxs"]("div", {
                                  className:
                                    "flex\x20shrink-0\x20flex-col\x20gap-3",
                                  children: [
                                    _0x2ddfe2["jsxs"]("label", {
                                      className:
                                        "flex\x20items-center\x20justify-between\x20gap-3\x20text-xs",
                                      children: [
                                        _0x2ddfe2["jsx"]("span", {
                                          children: "内置探针",
                                        }),
                                        _0x2ddfe2["jsx"](_0x543b01, {
                                          checked: _0x28af34,
                                          onCheckedChange: (_0x27ecb3) => {
                                            (_0x195e6d(_0x27ecb3),
                                              _0x55e2c2(_0x27ecb3 || _0x5289ba),
                                              _0x6360b0({
                                                internal_enabled: _0x27ecb3,
                                              }));
                                          },
                                        }),
                                      ],
                                    }),
                                    _0x2ddfe2["jsxs"]("label", {
                                      className:
                                        "flex\x20items-center\x20justify-between\x20gap-3\x20text-xs",
                                      children: [
                                        _0x2ddfe2["jsx"]("span", {
                                          children: "外置探针",
                                        }),
                                        _0x2ddfe2["jsx"](_0x543b01, {
                                          checked: _0x5289ba,
                                          onCheckedChange: (_0xf664c5) => {
                                            (_0x5aca07(_0xf664c5),
                                              _0x55e2c2(_0xf664c5 || _0x28af34),
                                              _0x6360b0({
                                                external_enabled: _0xf664c5,
                                              }));
                                          },
                                        }),
                                      ],
                                    }),
                                  ],
                                }),
                              ],
                            }),
                            _0x179ea5 &&
                              _0x2ddfe2["jsxs"]("div", {
                                className: "mt-4\x20space-y-4",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1.5",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "probe-title",
                                        className: "text-xs",
                                        children:
                                          "页面标题(浏览器标签与页头都用它)",
                                      }),
                                      _0x2ddfe2["jsx"](_0x549353, {
                                        id: "probe-title",
                                        placeholder:
                                          "如:XX\x20网络监控\x20/\x20服务器状态",
                                        value: _0x39dd42,
                                        onChange: (_0x51d619) =>
                                          _0x4c73e(
                                            _0x51d619["target"]["value"],
                                          ),
                                        onBlur: () =>
                                          _0x6360b0({ title: _0x39dd42 }),
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1.5",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "probe-theme",
                                        className: "text-xs",
                                        children: "探针主题",
                                      }),
                                      _0x2ddfe2["jsxs"](_0x3ba40b, {
                                        value: [
                                          "follow",
                                          "flat",
                                        ]["includes"](_0x307fbc)
                                          ? _0x307fbc
                                          : "custom",
                                        onValueChange: (_0x5502e6) => {
                                          if (_0x5502e6 === "custom") {
                                            _0x593c3b("");
                                            return;
                                          }
                                          (_0x593c3b(_0x5502e6),
                                            _0x6360b0({ theme: _0x5502e6 }));
                                        },
                                        children: [
                                          _0x2ddfe2["jsx"](_0x3d4a68, {
                                            id: "probe-theme",
                                            children: _0x2ddfe2["jsx"](
                                              _0x1929c7,
                                              {},
                                            ),
                                          }),
                                          _0x2ddfe2["jsxs"](_0x51ae8f, {
                                            children: [
                                              _0x2ddfe2["jsx"](_0x1260f9, {
                                                value: "follow",
                                                children: "跟随系统默认主题",
                                              }),
                                              _0x2ddfe2["jsx"](_0x1260f9, {
                                                value: "flat",
                                                children: "MEO 简约",
                                              }),
                                              _0x2ddfe2["jsx"](_0x1260f9, {
                                                value: "custom",
                                                children: "自定义主题名称",
                                              }),
                                            ],
                                          }),
                                        ],
                                      }),
                                      ![
                                        "follow",
                                        "flat",
                                      ]["includes"](_0x307fbc) &&
                                        _0x2ddfe2["jsx"](_0x549353, {
                                          id: "probe-custom-theme",
                                          value: _0x307fbc,
                                          maxLength: 0x40,
                                          placeholder: "例如\x20sakura-dark",
                                          onChange: (_0x3ccd5d) =>
                                            _0x593c3b(
                                              _0x3ccd5d["target"]["value"],
                                            ),
                                          onBlur: () => {
                                            const _0x4ebc5f =
                                              _0x307fbc["trim"]();
                                            _0x4ebc5f &&
                                              (_0x593c3b(_0x4ebc5f),
                                              _0x6360b0({ theme: _0x4ebc5f }));
                                          },
                                          onKeyDown: (_0x2ba796) => {
                                            _0x2ba796["key"] === "Enter" &&
                                              _0x2ba796["currentTarget"][
                                                "blur"
                                              ]();
                                          },
                                        }),
                                      _0x2ddfe2["jsx"]("p", {
                                        className:
                                          "text-muted-foreground\x20text-xs",
                                        children:
                                          "自定义名称会由主控接口下发给外置探针，对应\x20theme-名称\x20CSS\x20类；内置探针对未知主题使用默认样式。",
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1.5",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "probe-logo",
                                        className: "text-xs",
                                        children:
                                          "页头\x20Logo(可选,显示在标题左侧)",
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className: "flex\x20gap-2",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x549353, {
                                            id: "probe-logo",
                                            placeholder:
                                              "图片链接,如\x20https://example.com/logo.png",
                                            value: _0xb66ae3["startsWith"](
                                              "data:",
                                            )
                                              ? "(已上传的图片)"
                                              : _0xb66ae3,
                                            disabled:
                                              _0xb66ae3["startsWith"]("data:"),
                                            onChange: (_0x248c35) =>
                                              _0x1e1227(
                                                _0x248c35["target"]["value"],
                                              ),
                                            onBlur: () =>
                                              _0x6360b0({ logo: _0xb66ae3 }),
                                          }),
                                          _0x2ddfe2["jsx"]("input", {
                                            id: "probe-logo-file",
                                            type: "file",
                                            accept: "image/*",
                                            className: "hidden",
                                            onChange: (_0x269e83) => {
                                              const _0x131234 =
                                                _0x269e83["target"][
                                                  "files"
                                                ]?.[0x0];
                                              if (
                                                ((_0x269e83["target"]["value"] =
                                                  ""),
                                                !_0x131234)
                                              )
                                                return;
                                              if (
                                                _0x131234["size"] >
                                                0x5a * 0x400
                                              ) {
                                                _0x54e43f["error"](
                                                  "图片过大(建议\x2090KB\x20以内),请压缩或改用图片链接",
                                                );
                                                return;
                                              }
                                              const _0x486866 =
                                                new FileReader();
                                              ((_0x486866["onload"] = () => {
                                                const _0x58f685 = String(
                                                  _0x486866["result"] || "",
                                                );
                                                (_0x1e1227(_0x58f685),
                                                  _0x6360b0({
                                                    logo: _0x58f685,
                                                  }));
                                              }),
                                                _0x486866["readAsDataURL"](
                                                  _0x131234,
                                                ));
                                            },
                                          }),
                                          _0x2ddfe2["jsx"](_0x5185a8, {
                                            variant: "outline",
                                            onClick: () =>
                                              document["getElementById"](
                                                "probe-logo-file",
                                              )?.["click"](),
                                            children: "上传",
                                          }),
                                          _0xb66ae3 &&
                                            _0x2ddfe2["jsx"](_0x5185a8, {
                                              variant: "outline",
                                              onClick: () => {
                                                (_0x1e1227(""),
                                                  _0x6360b0({ logo: "" }));
                                              },
                                              children: "清除",
                                            }),
                                        ],
                                      }),
                                      _0xb66ae3 &&
                                        _0x2ddfe2["jsxs"]("div", {
                                          className:
                                            "flex\x20items-center\x20gap-2\x20pt-1",
                                          children: [
                                            _0x2ddfe2["jsx"]("img", {
                                              src: _0xb66ae3,
                                              alt: "logo\x20预览",
                                              className:
                                                "h-8\x20w-auto\x20max-w-[160px]\x20object-contain",
                                            }),
                                            _0x2ddfe2["jsx"]("span", {
                                              className:
                                                "text-muted-foreground\x20text-xs",
                                              children: "预览",
                                            }),
                                          ],
                                        }),
                                      _0x2ddfe2["jsx"]("p", {
                                        className:
                                          "text-muted-foreground\x20text-xs",
                                        children:
                                          "推荐填图片链接;上传的图片会内嵌进配置,探针页每\x205\x20秒轮询一次,过大会持续占带宽",
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1.5",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "probe-icon",
                                        className: "text-xs",
                                        children:
                                          "浏览器标签页图标(内置/外置探针共用)",
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className: "flex\x20gap-2",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x549353, {
                                            id: "probe-icon",
                                            placeholder:
                                              "图片链接,如\x20https://example.com/favicon.ico",
                                            value: _0x385ae1["startsWith"](
                                              "data:",
                                            )
                                              ? "(已上传的图片)"
                                              : _0x385ae1,
                                            disabled:
                                              _0x385ae1["startsWith"]("data:"),
                                            onChange: (_0x258084) =>
                                              _0x34add0(
                                                _0x258084["target"]["value"],
                                              ),
                                            onBlur: () =>
                                              _0x6360b0({ icon: _0x385ae1 }),
                                          }),
                                          _0x2ddfe2["jsx"]("input", {
                                            id: "probe-icon-file",
                                            type: "file",
                                            accept: "image/*,.ico",
                                            className: "hidden",
                                            onChange: (_0x53f800) => {
                                              const _0x3b4236 =
                                                _0x53f800["target"][
                                                  "files"
                                                ]?.[0x0];
                                              if (
                                                ((_0x53f800["target"]["value"] =
                                                  ""),
                                                !_0x3b4236)
                                              )
                                                return;
                                              if (
                                                _0x3b4236["size"] >
                                                0x5a * 0x400
                                              ) {
                                                _0x54e43f["error"](
                                                  "图片过大(建议\x2090KB\x20以内),请压缩或改用图片链接",
                                                );
                                                return;
                                              }
                                              const _0x6f0661 =
                                                new FileReader();
                                              ((_0x6f0661["onload"] = () => {
                                                const _0x5c3884 = String(
                                                  _0x6f0661["result"] || "",
                                                );
                                                (_0x34add0(_0x5c3884),
                                                  _0x6360b0({
                                                    icon: _0x5c3884,
                                                  }));
                                              }),
                                                _0x6f0661["readAsDataURL"](
                                                  _0x3b4236,
                                                ));
                                            },
                                          }),
                                          _0x2ddfe2["jsx"](_0x5185a8, {
                                            variant: "outline",
                                            onClick: () =>
                                              document["getElementById"](
                                                "probe-icon-file",
                                              )?.["click"](),
                                            children: "上传",
                                          }),
                                          _0x385ae1 &&
                                            _0x2ddfe2["jsx"](_0x5185a8, {
                                              variant: "outline",
                                              onClick: () => {
                                                (_0x34add0(""),
                                                  _0x6360b0({ icon: "" }));
                                              },
                                              children: "清除",
                                            }),
                                        ],
                                      }),
                                      _0x385ae1 &&
                                        _0x2ddfe2["jsxs"]("div", {
                                          className:
                                            "flex\x20items-center\x20gap-2\x20pt-1",
                                          children: [
                                            _0x2ddfe2["jsx"]("img", {
                                              src: _0x385ae1,
                                              alt: "浏览器图标预览",
                                              className:
                                                "h-8\x20w-8\x20rounded\x20object-contain",
                                            }),
                                            _0x2ddfe2["jsx"]("span", {
                                              className:
                                                "text-muted-foreground\x20text-xs",
                                              children: "预览",
                                            }),
                                          ],
                                        }),
                                      _0x2ddfe2["jsx"]("p", {
                                        className:
                                          "text-muted-foreground\x20text-xs",
                                        children:
                                          "仅用于浏览器标签页，不会替代页头\x20Logo",
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className:
                                      "flex\x20items-center\x20justify-between\x20gap-4",
                                    children: [
                                      _0x2ddfe2["jsx"]("div", {
                                        className: "text-sm",
                                        children: "显示服务器名称",
                                      }),
                                      _0x2ddfe2["jsx"](_0x543b01, {
                                        checked: _0x2f7373,
                                        onCheckedChange: (_0x359391) => {
                                          (_0x696237(_0x359391),
                                            _0x6360b0({
                                              show_name: _0x359391,
                                            }));
                                        },
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className:
                                      "flex\x20items-center\x20justify-between\x20gap-4",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        children: [
                                          _0x2ddfe2["jsx"]("div", {
                                            className: "text-sm",
                                            children: "禁止访问原登录页",
                                          }),
                                          _0x2ddfe2["jsx"]("p", {
                                            className:
                                              "text-muted-foreground\x20mt-1\x20text-xs",
                                            children:
                                              "开启后,未登录访客访问\x20/login\x20会被弹回探针页,站点对外只剩一个监控页。\x20管理员仍可通过探针页右上角的隐蔽入口登录。",
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsx"](_0x543b01, {
                                        className: "shrink-0",
                                        checked: _0x2f3e37,
                                        onCheckedChange: (_0x41fe3d) => {
                                          (_0x320f7a(_0x41fe3d),
                                            _0x6360b0({
                                              block_login: _0x41fe3d,
                                            }));
                                        },
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className:
                                      "probe-security-panel\x20space-y-3\x20rounded-md\x20border\x20border-amber-300/70\x20bg-amber-50/50\x20p-3\x20dark:bg-amber-950/10",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "flex\x20items-start\x20justify-between\x20gap-4",
                                        children: [
                                          _0x2ddfe2["jsxs"]("div", {
                                            children: [
                                              _0x2ddfe2["jsx"]("div", {
                                                className:
                                                  "text-sm\x20font-medium",
                                                children: "保护探针数据接口",
                                              }),
                                              _0x2ddfe2["jsx"]("p", {
                                                className:
                                                  "text-muted-foreground\x20mt-1\x20text-xs",
                                                children:
                                                  "开启后，已启用的内置探针仅允许同源请求，已启用的外置探针必须携带专用密钥；跨站调用、无来源扫描和错误密钥请求返回\x20404。两个探针的启用状态互不影响。密钥只应保存到\x20Cloudflare\x20Worker\x20Secret，不能写进前端代码。",
                                              }),
                                            ],
                                          }),
                                          _0x2ddfe2["jsx"](_0x543b01, {
                                            className: "shrink-0",
                                            checked: _0x4832dd,
                                            onCheckedChange: (_0x588b0c) => {
                                              (_0x200bc7(_0x588b0c),
                                                _0x6360b0({
                                                  external_access_only:
                                                    _0x588b0c,
                                                }));
                                            },
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className: "flex\x20flex-wrap\x20gap-2",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x5185a8, {
                                            type: "button",
                                            variant: "outline",
                                            onClick: () => {
                                              const _0x5aef73 = new Uint8Array(
                                                0x20,
                                              );
                                              crypto["getRandomValues"](
                                                _0x5aef73,
                                              );
                                              const _0x5a66cd = btoa(
                                                String["fromCharCode"](
                                                  ..._0x5aef73,
                                                ),
                                              )
                                                ["replace"](/\+/g, "-")
                                                ["replace"](/\//g, "_")
                                                ["replace"](/=/g, "");
                                              (_0x4ecb0d(_0x5a66cd),
                                                _0x586e39(!0x0),
                                                _0x6360b0({
                                                  external_access_token:
                                                    _0x5a66cd,
                                                }));
                                            },
                                            children: _0x3846db
                                              ? "轮换访问密钥"
                                              : "生成访问密钥",
                                          }),
                                          _0x57745b &&
                                            _0x2ddfe2["jsx"](_0x5185a8, {
                                              type: "button",
                                              variant: "outline",
                                              onClick: async () => {
                                                try {
                                                  (await navigator["clipboard"][
                                                    "writeText"
                                                  ](_0x57745b),
                                                    _0x54e43f["success"](
                                                      "访问密钥已复制",
                                                    ));
                                                } catch {
                                                  _0x54e43f["error"](
                                                    "复制失败，请手动复制",
                                                  );
                                                }
                                              },
                                              children: "复制密钥",
                                            }),
                                        ],
                                      }),
                                      _0x57745b
                                        ? _0x2ddfe2["jsxs"]("div", {
                                            className: "space-y-1",
                                            children: [
                                              _0x2ddfe2["jsx"](_0x549353, {
                                                readOnly: !0x0,
                                                value: _0x57745b,
                                                className:
                                                  "font-mono\x20text-xs",
                                              }),
                                              _0x2ddfe2["jsx"]("p", {
                                                className:
                                                  "text-xs\x20font-medium\x20text-red-600\x20dark:text-red-400",
                                                children:
                                                  "请立即复制并执行\x20wrangler\x20secret\x20put\x20PROBE_TOKEN；刷新页面后无法再次查看此密钥。",
                                              }),
                                            ],
                                          })
                                        : _0x2ddfe2["jsx"]("p", {
                                            className:
                                              "text-muted-foreground\x20text-xs",
                                            children: _0x3846db
                                              ? "访问密钥已配置，主控只保存\x20SHA-256\x20哈希。"
                                              : "仅使用内置探针时无需生成密钥；独立探针需要生成并配置密钥。",
                                          }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-2",
                                    children: [
                                      _0x2ddfe2["jsx"]("div", {
                                        className:
                                          "text-muted-foreground\x20text-xs",
                                        children:
                                          "展示的服务器(不选\x20=\x20探针页不展示任何服务器)",
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className:
                                          "max-h-48\x20divide-y\x20overflow-y-auto\x20rounded-md\x20border",
                                        children: [
                                          (_0x3ded3d?.["servers"] || [])["map"](
                                            (_0xf2adc4) =>
                                              _0x2ddfe2["jsxs"](
                                                "label",
                                                {
                                                  className:
                                                    "hover:bg-accent/40\x20flex\x20cursor-pointer\x20items-center\x20gap-2\x20px-3\x20py-2\x20text-sm",
                                                  children: [
                                                    _0x2ddfe2["jsx"](
                                                      _0x77832a,
                                                      {
                                                        checked: _0x132fa8[
                                                          "includes"
                                                        ](_0xf2adc4["id"]),
                                                        onCheckedChange: (
                                                          _0x57faa3,
                                                        ) => {
                                                          const _0x2a4472 =
                                                            _0x57faa3 === !0x0
                                                              ? [
                                                                  ..._0x132fa8,
                                                                  _0xf2adc4[
                                                                    "id"
                                                                  ],
                                                                ]
                                                              : _0x132fa8[
                                                                  "filter"
                                                                ](
                                                                  (_0x3fa71e) =>
                                                                    _0x3fa71e !==
                                                                    _0xf2adc4[
                                                                      "id"
                                                                    ],
                                                                );
                                                          (_0x2af46d(_0x2a4472),
                                                            _0x6360b0({
                                                              server_ids:
                                                                _0x2a4472,
                                                            }));
                                                        },
                                                      },
                                                    ),
                                                    _0x2ddfe2["jsx"]("span", {
                                                      className: "truncate",
                                                      children:
                                                        _0xf2adc4["name"],
                                                    }),
                                                  ],
                                                },
                                                _0xf2adc4["id"],
                                              ),
                                          ),
                                          (_0x3ded3d?.["servers"] || [])[
                                            "length"
                                          ] === 0x0 &&
                                            _0x2ddfe2["jsx"]("div", {
                                              className:
                                                "text-muted-foreground\x20px-3\x20py-4\x20text-center\x20text-xs",
                                              children: "暂无服务器",
                                            }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x132fa8["length"] > 0x0 &&
                                    _0x2ddfe2["jsxs"]("div", {
                                      className: "space-y-2",
                                      children: [
                                        _0x2ddfe2["jsx"]("div", {
                                          className:
                                            "text-muted-foreground\x20text-xs",
                                          children:
                                            "服务器地域与续费信息（可手动维护）",
                                        }),
                                        _0x2ddfe2["jsx"]("div", {
                                          className: "rounded-md\x20border",
                                          children: (_0x3ded3d?.["servers"] ||
                                            [])
                                            ["filter"]((_0x23934e) =>
                                              _0x132fa8["includes"](
                                                _0x23934e["id"],
                                              ),
                                            )
                                            ["map"]((_0x2b018a) =>
                                              _0x2ddfe2["jsx"](
                                                Ir,
                                                {
                                                  server: _0x2b018a,
                                                  saved: () =>
                                                    _0x2f82ce[
                                                      "invalidateQueries"
                                                    ]({
                                                      queryKey: [
                                                        "remote-servers",
                                                      ],
                                                    }),
                                                },
                                                _0x2b018a["id"],
                                              ),
                                            ),
                                        }),
                                      ],
                                    }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-2\x20border-t\x20pt-3",
                                    children: [
                                      _0x2ddfe2["jsx"]("div", {
                                        className:
                                          "text-muted-foreground\x20text-xs",
                                        children:
                                          "探针数据采集(CPU/内存/硬盘/Ping\x20需\x20agent\x20采集上报;流量/网速为展示开关,数据主控已有,关闭则伪装页不显示该项)",
                                      }),
                                      [
                                        [
                                          "metric_cpu",
                                          "CPU\x20负载",
                                          _0x40c023,
                                          _0x5e7f11,
                                        ],
                                        [
                                          "metric_mem",
                                          "内存",
                                          _0x237d87,
                                          _0x3c1046,
                                        ],
                                        [
                                          "metric_disk",
                                          "硬盘",
                                          _0x334fa2,
                                          _0x54bd7b,
                                        ],
                                        [
                                          "metric_ping",
                                          "Ping\x20探测(各省市三网延迟)",
                                          _0x36e07d,
                                          _0x2169f6,
                                        ],
                                        [
                                          "metric_traffic",
                                          "流量信息",
                                          _0x201bbd,
                                          _0x566d88,
                                        ],
                                        [
                                          "metric_speed",
                                          "网速信息",
                                          _0x32c971,
                                          _0x745469,
                                        ],
                                        [
                                          "show_expiry",
                                          "服务器到期时间",
                                          _0x48c9e8,
                                          _0x298897,
                                        ],
                                        [
                                          "show_globe",
                                          "3D\x20地球地区分布",
                                          _0xf265e3,
                                          _0x3addbf,
                                        ],
                                        [
                                          "show_forward",
                                          "转发链拓扑与延迟",
                                          _0xc58657,
                                          _0x19ab74,
                                        ],
                                      ]["map"](
                                        ([
                                          _0x3fb503,
                                          _0x5181ce,
                                          _0x1a6b58,
                                          _0x52bd8a,
                                        ]) =>
                                          _0x2ddfe2["jsxs"](
                                            "label",
                                            {
                                              className:
                                                "flex\x20cursor-pointer\x20items-center\x20gap-2\x20text-sm",
                                              children: [
                                                _0x2ddfe2["jsx"](_0x77832a, {
                                                  checked: _0x1a6b58,
                                                  onCheckedChange: (
                                                    _0x5ca8ef,
                                                  ) => {
                                                    const _0x422e81 =
                                                      _0x5ca8ef === !0x0;
                                                    (_0x52bd8a(_0x422e81),
                                                      _0x6360b0({
                                                        [_0x3fb503]: _0x422e81,
                                                      }));
                                                  },
                                                }),
                                                _0x2ddfe2["jsx"]("span", {
                                                  children: _0x5181ce,
                                                }),
                                              ],
                                            },
                                            _0x3fb503,
                                          ),
                                      ),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-2\x20border-t\x20pt-3",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        children: [
                                          _0x2ddfe2["jsx"]("div", {
                                            className:
                                              "text-sm\x20font-semibold",
                                            children: "探针数据模块",
                                          }),
                                          _0x2ddfe2["jsx"]("p", {
                                            className:
                                              "text-muted-foreground\x20mt-1\x20text-xs",
                                            children:
                                              "分别控制首页的数据模块；关闭后公开接口仍只返回其他已启用模块需要的数据。",
                                          }),
                                        ],
                                      }),
                                      [
                                        [
                                          "show_daily_trend",
                                          "首页每日流量趋势",
                                          _0x12d674,
                                          _0x1f17de,
                                        ],
                                        [
                                          "show_traffic_hotspots",
                                          "首页实时流量热点",
                                          _0x3cdab9,
                                          _0x98befa,
                                        ],
                                        [
                                          "show_traffic_7d",
                                          "近\x207\x20日上下行流量堆叠图",
                                          _0x4894c7,
                                          _0xe71a4f,
                                        ],
                                        [
                                          "show_resource_heatmap",
                                          "CPU\x20/\x20内存\x20/\x20硬盘压力热力图",
                                          _0x4c2853,
                                          _0x21e7c8,
                                        ],
                                        [
                                          "show_traffic_quota",
                                          "流量额度使用率排行榜",
                                          _0x3fc691,
                                          _0xa4bc46,
                                        ],
                                        [
                                          "show_renewal_timeline",
                                          "服务器续费与到期时间轴",
                                          _0x52c280,
                                          _0x206944,
                                        ],
                                        [
                                          "show_health_score",
                                          "服务器卡片健康评分",
                                          _0x584361,
                                          _0x509596,
                                        ],
                                      ]["map"](
                                        ([
                                          _0x304e0e,
                                          _0x492de0,
                                          _0x492c47,
                                          _0x52aa2f,
                                        ]) =>
                                          _0x2ddfe2["jsxs"](
                                            "div",
                                            {
                                              className:
                                                "flex\x20items-center\x20justify-between\x20gap-4\x20rounded-md\x20border\x20p-3",
                                              children: [
                                                _0x2ddfe2["jsx"]("span", {
                                                  className: "text-sm",
                                                  children: _0x492de0,
                                                }),
                                                _0x2ddfe2["jsx"](_0x543b01, {
                                                  className: "shrink-0",
                                                  checked: _0x492c47,
                                                  onCheckedChange: (
                                                    _0x2aedda,
                                                  ) => {
                                                    (_0x52aa2f(_0x2aedda),
                                                      _0x6360b0({
                                                        [_0x304e0e]: _0x2aedda,
                                                      }));
                                                  },
                                                }),
                                              ],
                                            },
                                            _0x304e0e,
                                          ),
                                      ),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className:
                                      "flex\x20items-start\x20justify-between\x20gap-4\x20rounded-md\x20border\x20p-3",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        children: [
                                          _0x2ddfe2["jsx"]("div", {
                                            className: "text-sm\x20font-medium",
                                            children: "三网回程测试与铭牌",
                                          }),
                                          _0x2ddfe2["jsx"]("p", {
                                            className:
                                              "text-muted-foreground\x20mt-1\x20text-xs",
                                            children:
                                              "开启后，探针\x20HTTP\x20与\x20WebSocket\x20接口返回三网回程数据，并在服务器卡片显示铭牌；主控每天\x2004:20\x20自动测试一次。关闭后不再自动测试，也不会向探针接口返回该数据，日志管理仍可手动运行。",
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsx"](_0x543b01, {
                                        className: "shrink-0",
                                        checked: _0x2c452d,
                                        onCheckedChange: (_0x57681e) => {
                                          (_0x397aa6(_0x57681e),
                                            _0x6360b0({
                                              show_return_route: _0x57681e,
                                            }));
                                        },
                                      }),
                                    ],
                                  }),
                                  _0x36e07d &&
                                    _0x2ddfe2["jsxs"](_0x2ddfe2["Fragment"], {
                                      children: [
                                        _0x2ddfe2["jsxs"]("div", {
                                          className: "space-y-1",
                                          children: [
                                            _0x2ddfe2["jsx"]("div", {
                                              className:
                                                "text-xs\x20font-medium",
                                              children: "全局\x20Ping\x20目标",
                                            }),
                                            _0x2ddfe2["jsx"]("div", {
                                              className:
                                                "text-muted-foreground\x20text-xs",
                                              children:
                                                "未单独指定的服务器都用这批目标",
                                            }),
                                          ],
                                        }),
                                        _0x2ddfe2["jsx"](wn, {
                                          regions: _0x12640b?.["regions"],
                                          source: _0x12640b?.["source"],
                                          selected: _0x25dc28,
                                          onChange: (_0x1c9ddf) => {
                                            (_0x1cab7e(_0x1c9ddf),
                                              _0x6360b0({
                                                ping_targets: _0x1c9ddf,
                                              }));
                                          },
                                        }),
                                        _0x2ddfe2["jsxs"]("div", {
                                          className:
                                            "space-y-2\x20border-t\x20pt-3",
                                          children: [
                                            _0x2ddfe2["jsx"]("div", {
                                              className:
                                                "text-muted-foreground\x20text-xs",
                                              children:
                                                "按服务器单独指定\x20Ping\x20目标(不设置则跟随上方全局)",
                                            }),
                                            (_0x3ded3d?.["servers"] || [])
                                              ["filter"]((_0x590481) =>
                                                _0x132fa8["includes"](
                                                  _0x590481["id"],
                                                ),
                                              )
                                              ["map"]((_0x548900) =>
                                                _0x2ddfe2["jsx"](
                                                  zr,
                                                  {
                                                    server: _0x548900,
                                                    regions:
                                                      _0x12640b?.["regions"],
                                                    source:
                                                      _0x12640b?.["source"],
                                                    globalTargets: _0x25dc28,
                                                    override:
                                                      _0x14d133[
                                                        String(_0x548900["id"])
                                                      ],
                                                    onChange: (_0x4e9666) => {
                                                      const _0x5722f0 = {
                                                        ..._0x14d133,
                                                      };
                                                      (_0x4e9666 === void 0x0
                                                        ? delete _0x5722f0[
                                                            String(
                                                              _0x548900["id"],
                                                            )
                                                          ]
                                                        : (_0x5722f0[
                                                            String(
                                                              _0x548900["id"],
                                                            )
                                                          ] = _0x4e9666),
                                                        _0x39e60b(_0x5722f0),
                                                        _0x6360b0({
                                                          ping_targets_override:
                                                            _0x5722f0,
                                                        }));
                                                    },
                                                  },
                                                  _0x548900["id"],
                                                ),
                                              ),
                                            _0x132fa8["length"] === 0x0 &&
                                              _0x2ddfe2["jsx"]("div", {
                                                className:
                                                  "text-muted-foreground\x20text-xs",
                                                children:
                                                  "先在上方勾选要展示的服务器",
                                              }),
                                          ],
                                        }),
                                      ],
                                    }),
                                ],
                              }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className: "space-y-3\x20rounded-lg\x20border\x20p-3",
                          children: [
                            _0x2ddfe2["jsxs"]("div", {
                              className:
                                "flex\x20items-start\x20justify-between\x20gap-4\x20border-b\x20pb-3",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "space-y-1",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      htmlFor: "probe-metrics-persist",
                                      children: "保存历史数据",
                                    }),
                                    _0x2ddfe2["jsx"]("p", {
                                      className:
                                        "text-muted-foreground\x20text-xs",
                                      children: _0x41f4f1?.[
                                        "metrics_persist_supported"
                                      ]
                                        ? "关闭后不再写入历史表,延迟趋势图将没有数据;当前状态、健康判定与入口\x20DNS\x20编排不受影响。"
                                        : "当前数据库为\x20SQLite,不支持保存历史数据\x20——\x20分钟级时序会让\x20WAL\x20文件持续膨胀到几十\x20GB。切换到\x20PostgreSQL\x20后可开启。",
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsx"](_0x543b01, {
                                  id: "probe-metrics-persist",
                                  checked: _0x123118,
                                  disabled:
                                    !_0x41f4f1?.["metrics_persist_supported"],
                                  onCheckedChange: (_0x2b9c68) => {
                                    (_0x1f3edf(_0x2b9c68),
                                      _0x6360b0({
                                        metrics_persist_enabled: _0x2b9c68,
                                      }));
                                  },
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsxs"]("div", {
                              className: "space-y-1",
                              children: [
                                _0x2ddfe2["jsx"]("div", {
                                  className: "text-sm\x20font-semibold",
                                  children: "转发链数据保留",
                                }),
                                _0x2ddfe2["jsx"]("p", {
                                  className: "text-muted-foreground\x20text-xs",
                                  children:
                                    "转发链的延迟采样点与每日流量分开设置保留期,超期数据每天自动清理。\x20修改后下次清理即按新值执行,无需重启。",
                                }),
                              ],
                            }),
                            _0x2ddfe2["jsxs"]("div", {
                              className: "grid\x20gap-3\x20sm:grid-cols-2",
                              children: [
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "space-y-1.5",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      htmlFor: "forward-metrics-retention",
                                      children: "延迟采样点(天)",
                                    }),
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      id: "forward-metrics-retention",
                                      type: "number",
                                      min: 0x1,
                                      max: 0x7,
                                      value: _0x580edb,
                                      onChange: (_0x9c442) =>
                                        _0x32438b(
                                          Number(_0x9c442["target"]["value"]),
                                        ),
                                      onBlur: () =>
                                        _0x6360b0({
                                          forward_metrics_retention_days:
                                            _0x580edb,
                                        }),
                                    }),
                                    _0x2ddfe2["jsx"]("p", {
                                      className:
                                        "text-muted-foreground\x20text-xs",
                                      children:
                                        "1–7\x20天,默认\x207。分钟级采样,是延迟趋势图的数据来源;体量最大,更久的趋势请看每日流量。",
                                    }),
                                  ],
                                }),
                                _0x2ddfe2["jsxs"]("div", {
                                  className: "space-y-1.5",
                                  children: [
                                    _0x2ddfe2["jsx"](_0x34df34, {
                                      htmlFor: "forward-daily-retention",
                                      children: "每日流量(天)",
                                    }),
                                    _0x2ddfe2["jsx"](_0x549353, {
                                      id: "forward-daily-retention",
                                      type: "number",
                                      min: 0x1,
                                      max: 0x1e,
                                      value: _0x3810c9,
                                      onChange: (_0x22379f) =>
                                        _0x510199(
                                          Number(_0x22379f["target"]["value"]),
                                        ),
                                      onBlur: () =>
                                        _0x6360b0({
                                          forward_daily_retention_days:
                                            _0x3810c9,
                                        }),
                                    }),
                                    _0x2ddfe2["jsx"]("p", {
                                      className:
                                        "text-muted-foreground\x20text-xs",
                                      children:
                                        "1–30\x20天,默认\x2030。每天每链每服务器一行,占用极小。",
                                    }),
                                  ],
                                }),
                              ],
                            }),
                          ],
                        }),
                      ],
                    }),
                  ],
                }),
              }),
              _0x2ddfe2["jsxs"](_0x550af3, {
                value: "appearance",
                className: "space-y-6",
                children: [
                  _0x2ddfe2["jsx"](Pr, {}),
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsxs"](_0x1db8ce, {
                        className: "pb-4",
                        children: [
                          _0x2ddfe2["jsx"](_0x30c50e, { children: "默认主题" }),
                          _0x2ddfe2["jsx"](_0x54661d, {
                            children:
                              "全站统一使用 MEO 简约主题，仅保留浅色、深色与跟随系统模式。",
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsxs"](_0x42cb32, {
                        children: [
                          _0x2ddfe2["jsx"]("div", {
                            className: "flex\x20gap-2",
                            children: [
                              { value: "flat", label: "MEO 简约" },
                            ]["map"]((_0x2e220e) =>
                              _0x2ddfe2["jsx"](
                                _0x5185a8,
                                {
                                  type: "button",
                                  variant:
                                    _0x3fd9b7 === _0x2e220e["value"]
                                      ? "default"
                                      : "outline",
                                  disabled:
                                    _0x47362d["isPending"],
                                  onClick: () => {
                                    (_0x35426f(_0x2e220e["value"]),
                                    _0x47362d["mutate"]({
                                      default_theme: _0x2e220e["value"],
                                    }));
                                  },
                                  children: _0x2e220e["label"],
                                },
                                _0x2e220e["value"],
                              ),
                            ),
                          }),
                        ],
                      }),
                    ],
                  }),
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsxs"](_0x1db8ce, {
                        className: "pb-4",
                        children: [
                          _0x2ddfe2["jsx"](_0x30c50e, {
                            children: "登录页壁纸",
                          }),
                          _0x2ddfe2["jsxs"](_0x54661d, {
                            children: [
                              "自定义登录页背景图（填写图片 URL 或站内路径）。留空时使用 MEO 默认背景。",
                              _0x2ddfe2["jsx"]("br", {}),
                              "想用自己的图片:把图片放进主控的",
                              "\x20",
                              _0x2ddfe2["jsx"]("code", {
                                className: "bg-muted\x20rounded\x20px-1",
                                children: "data/public/",
                              }),
                              "\x20",
                              "目录(Docker\x20下就是挂载的\x20data\x20卷),再填",
                              "\x20",
                              _0x2ddfe2["jsx"]("code", {
                                className: "bg-muted\x20rounded\x20px-1",
                                children: "/public/图片名.jpg",
                              }),
                              "。注意",
                              "\x20",
                              _0x2ddfe2["jsx"]("code", {
                                className: "bg-muted\x20rounded\x20px-1",
                                children: "/images/xxx",
                              }),
                              "\x20",
                              "是编译进程序的只读资源,放不进自己的图。",
                            ],
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsxs"](_0x42cb32, {
                        className: "space-y-3",
                        children: [
                          _0x2ddfe2["jsx"](_0x549353, {
                            value: _0x1f5853,
                            onChange: (_0x2d1c76) =>
                              _0x46b92e(_0x2d1c76["target"]["value"]),
                            placeholder:
                              "https://example.com/wallpaper.webp(留空用内置网格背景)",
                          }),
                          _0x2ddfe2["jsxs"]("div", {
                            className: "flex\x20gap-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x5185a8, {
                                disabled: _0x3e6952["isPending"],
                                onClick: () =>
                                  _0x3e6952["mutate"]({
                                    login_wallpaper: _0x1f5853["trim"](),
                                  }),
                                children: "保存",
                              }),
                              _0x1f5853 &&
                                _0x2ddfe2["jsx"](_0x5185a8, {
                                  variant: "outline",
                                  disabled: _0x3e6952["isPending"],
                                  onClick: () => {
                                    (_0x46b92e(""),
                                      _0x3e6952["mutate"]({
                                        login_wallpaper: "",
                                      }));
                                  },
                                  children: "清除",
                                }),
                            ],
                          }),
                          _0x1f5853 &&
                            _0x2ddfe2["jsx"]("img", {
                              src: _0x1f5853,
                              alt: "",
                              className:
                                "max-h-40\x20w-full\x20rounded-md\x20border\x20object-cover",
                            }),
                        ],
                      }),
                    ],
                  }),
                ],
              }),
              _0x2ddfe2["jsx"](_0x550af3, {
                value: "announce",
                className: "space-y-6",
                children: _0x2ddfe2["jsx"](Sr, {}),
              }),
              _0x2ddfe2["jsx"](_0x550af3, {
                value: "tgbot",
                className: "space-y-6",
                children: _0x2ddfe2["jsx"](Er, {}),
              }),
              _0x2ddfe2["jsx"](_0x550af3, {
                value: "captcha",
                className: "space-y-6",
                children: _0x2ddfe2["jsxs"](_0x2aeb40, {
                  children: [
                    _0x2ddfe2["jsxs"](_0x1db8ce, {
                      className: "pb-4",
                      children: [
                        _0x2ddfe2["jsxs"](_0x30c50e, {
                          className: "flex\x20items-center\x20gap-2",
                          children: [
                            _0x2ddfe2["jsx"](_0x1618b9, {
                              className: "h-5\x20w-5",
                            }),
                            _0x1aa03d("turnstile.title"),
                          ],
                        }),
                        _0x2ddfe2["jsx"](_0x54661d, {
                          children: _0x1aa03d("turnstile.description"),
                        }),
                      ],
                    }),
                    _0x2ddfe2["jsxs"](_0x42cb32, {
                      className: "space-y-4",
                      children: [
                        _0x2ddfe2["jsxs"]("div", {
                          className: "space-y-1",
                          children: [
                            _0x2ddfe2["jsx"](_0x34df34, {
                              htmlFor: "turnstile-site-key",
                              className: "text-sm",
                              children: _0x1aa03d("turnstile.siteKey"),
                            }),
                            _0x2ddfe2["jsx"](_0x549353, {
                              id: "turnstile-site-key",
                              name: "mmwx-turnstile-site-key",
                              type: "text",
                              autoComplete: "off",
                              "data-1p-ignore": !0x0,
                              "data-lpignore": "true",
                              "data-form-type": "other",
                              placeholder: _0x1aa03d(
                                "turnstile.siteKeyPlaceholder",
                              ),
                              value: _0x585db0["turnstile_site_key"],
                              onChange: (_0x2eba13) =>
                                _0x28c5b7({
                                  ..._0x585db0,
                                  turnstile_site_key:
                                    _0x2eba13["target"]["value"],
                                }),
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsxs"]("div", {
                          className: "space-y-1",
                          children: [
                            _0x2ddfe2["jsx"](_0x34df34, {
                              htmlFor: "turnstile-secret-key",
                              className: "text-sm",
                              children: _0x1aa03d("turnstile.secretKey"),
                            }),
                            _0x2ddfe2["jsx"](_0x549353, {
                              id: "turnstile-secret-key",
                              name: "mmwx-turnstile-secret-key",
                              type: "password",
                              autoComplete: "off",
                              "data-1p-ignore": !0x0,
                              "data-lpignore": "true",
                              "data-form-type": "other",
                              placeholder: _0x585db0["turnstile_secret_key"]
                                ? _0x1aa03d("turnstile.secretKeyConfigured")
                                : _0x1aa03d("turnstile.secretKeyPlaceholder"),
                              value: _0x585db0["turnstile_secret_key"],
                              onChange: (_0x295b36) =>
                                _0x28c5b7({
                                  ..._0x585db0,
                                  turnstile_secret_key:
                                    _0x295b36["target"]["value"],
                                }),
                            }),
                          ],
                        }),
                        _0x2ddfe2["jsx"]("div", {
                          children: _0x2ddfe2["jsxs"](_0x5185a8, {
                            size: "sm",
                            onClick: () => _0x25971f["mutate"](_0x585db0),
                            disabled: _0x25971f["isPending"],
                            children: [
                              _0x25971f["isPending"]
                                ? _0x2ddfe2["jsx"](_0x4e8a70, {
                                    className:
                                      "mr-1\x20h-4\x20w-4\x20animate-spin",
                                  })
                                : null,
                              _0x1aa03d("actions.save", { ns: "common" }),
                            ],
                          }),
                        }),
                        _0x2ddfe2["jsx"](Gr, {
                          siteKey: _0x585db0["turnstile_site_key"],
                          secretKeyConfigured:
                            !!_0x585db0["turnstile_secret_key"],
                        }),
                        _0x2ddfe2["jsxs"]("p", {
                          className: "text-muted-foreground\x20text-xs",
                          children: [
                            _0x1aa03d("turnstile.hint"),
                            "\x20",
                            _0x2ddfe2["jsx"]("a", {
                              href: "https://github.com/aoomee/MEO",
                              target: "_blank",
                              rel: "noopener\x20noreferrer",
                              className: "text-primary\x20hover:underline",
                              children: _0x1aa03d("turnstile.docLinkLabel"),
                            }),
                          ],
                        }),
                      ],
                    }),
                  ],
                }),
              }),
              _0x2ddfe2["jsxs"](_0x550af3, {
                value: "system",
                className: "space-y-6",
                children: [
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsxs"](_0x1db8ce, {
                        className: "pb-4",
                        children: [
                          _0x2ddfe2["jsx"](_0x30c50e, {
                            children: _0x1aa03d("proxyGroups.title"),
                          }),
                          _0x2ddfe2["jsx"](_0x54661d, {
                            children: _0x1aa03d("proxyGroups.description"),
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsxs"](_0x42cb32, {
                        className: "space-y-4",
                        children: [
                          _0x2ddfe2["jsxs"]("div", {
                            className: "space-y-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x34df34, {
                                htmlFor: "proxy-groups-source-url",
                                children: _0x1aa03d("proxyGroups.sourceUrl"),
                              }),
                              _0x2ddfe2["jsx"](_0x549353, {
                                id: "proxy-groups-source-url",
                                value: _0x4710bf,
                                placeholder:
                                  "https://raw.githubusercontent.com/iluobei/miaomiaowu/refs/heads/main/proxy_groups/proxy-groups-lite.json",
                                disabled: _0x3002e9["isPending"],
                                onChange: (_0x10c412) =>
                                  _0x3be089(_0x10c412["target"]["value"]),
                                onBlur: () => {
                                  const _0x3984eb = _0x4710bf["trim"]();
                                  (_0x3be089(_0x3984eb),
                                    _0x3984eb !==
                                      (_0x486d96?.["proxy_groups_source_url"] ||
                                        "") &&
                                      _0x5637c8({
                                        proxy_groups_source_url: _0x3984eb,
                                      }));
                                },
                              }),
                              _0x2ddfe2["jsx"]("p", {
                                className: "text-muted-foreground\x20text-xs",
                                children: _0x1aa03d("proxyGroups.sourceHint"),
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsxs"](_0x5185a8, {
                            onClick: () => {
                              const _0x45d143 = _0x4710bf["trim"]() || void 0x0;
                              _0x3fe589["mutate"](_0x45d143, {
                                onError: _0x4deb0e,
                              });
                            },
                            disabled: _0x3fe589["isPending"],
                            children: [
                              _0x2ddfe2["jsx"](_0x18a4f9, {
                                className: _0x7633f5(
                                  "mr-2\x20h-4\x20w-4",
                                  _0x3fe589["isPending"] && "animate-spin",
                                ),
                              }),
                              _0x3fe589["isPending"]
                                ? _0x1aa03d("proxyGroups.syncing")
                                : _0x1aa03d("proxyGroups.sync"),
                            ],
                          }),
                          _0x3fe589["isSuccess"]
                            ? _0x2ddfe2["jsx"]("p", {
                                className:
                                  "text-sm\x20text-green-600\x20dark:text-green-400",
                                children: _0x1aa03d("proxyGroups.synced"),
                              })
                            : null,
                        ],
                      }),
                    ],
                  }),
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsxs"](_0x1db8ce, {
                        className: "pb-4",
                        children: [
                          _0x2ddfe2["jsxs"](_0x30c50e, {
                            className: "flex\x20items-center\x20gap-2",
                            children: [
                              _0x2ddfe2["jsx"](_0x18ddc9, {
                                className: "h-5\x20w-5",
                              }),
                              _0x1aa03d("intervals.title"),
                            ],
                          }),
                          _0x2ddfe2["jsx"](_0x54661d, {
                            children: _0x1aa03d("intervals.description"),
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsxs"](_0x42cb32, {
                        className: "space-y-4",
                        children: [
                          _0x2ddfe2["jsxs"]("div", {
                            className: "grid\x20grid-cols-2\x20gap-4",
                            children: [
                              _0x2ddfe2["jsxs"]("div", {
                                className: "space-y-1",
                                children: [
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "speed-interval",
                                    className: "text-sm",
                                    children: _0x1aa03d(
                                      "intervals.speedCollect",
                                    ),
                                  }),
                                  _0x2ddfe2["jsx"](_0x549353, {
                                    id: "speed-interval",
                                    type: "number",
                                    min: 0x1,
                                    value: _0x2a5d0a,
                                    onChange: (_0x350843) =>
                                      _0x44faac(
                                        Number(_0x350843["target"]["value"]),
                                      ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "space-y-1",
                                children: [
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "report-interval",
                                    className: "text-sm",
                                    children: _0x1aa03d(
                                      "intervals.reportInterval",
                                    ),
                                  }),
                                  _0x2ddfe2["jsx"](_0x549353, {
                                    id: "report-interval",
                                    type: "number",
                                    min: 0x1,
                                    max: 0x3c,
                                    value: _0x245fb9,
                                    onChange: (_0x54d6c7) =>
                                      _0x6bd1a7(
                                        Number(_0x54d6c7["target"]["value"]),
                                      ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "space-y-1",
                                children: [
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "check-interval",
                                    className: "text-sm",
                                    children: _0x1aa03d(
                                      "intervals.trafficCheck",
                                    ),
                                  }),
                                  _0x2ddfe2["jsx"](_0x549353, {
                                    id: "check-interval",
                                    type: "number",
                                    min: 0xa,
                                    value: _0x265cb2,
                                    onChange: (_0x195124) =>
                                      _0x14a567(
                                        Number(_0x195124["target"]["value"]),
                                      ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "space-y-1",
                                children: [
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "heartbeat-interval",
                                    className: "text-sm",
                                    children: _0x1aa03d("intervals.heartbeat"),
                                  }),
                                  _0x2ddfe2["jsx"](_0x549353, {
                                    id: "heartbeat-interval",
                                    type: "number",
                                    min: 0x5,
                                    value: _0x1f7866,
                                    onChange: (_0x260f9d) =>
                                      _0x270780(
                                        Number(_0x260f9d["target"]["value"]),
                                      ),
                                  }),
                                ],
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsx"](_0x5185a8, {
                            size: "sm",
                            onClick: () =>
                              _0xcd2c2c["mutate"]({
                                speed_collect_interval: _0x2a5d0a,
                                report_interval: _0x245fb9,
                                traffic_check_interval: _0x265cb2,
                                heartbeat_interval: _0x1f7866,
                              }),
                            disabled: _0xcd2c2c["isPending"],
                            children: _0x1aa03d("actions.save", {
                              ns: "common",
                            }),
                          }),
                          _0x2ddfe2["jsxs"]("div", {
                            className: "space-y-2\x20border-t\x20pt-4",
                            children: [
                              _0x2ddfe2["jsx"](_0x34df34, {
                                htmlFor: "master-url",
                                children: _0x1aa03d("masterUrl.title"),
                              }),
                              _0x2ddfe2["jsx"](_0x549353, {
                                id: "master-url",
                                placeholder: _0x1aa03d("masterUrl.placeholder"),
                                value: _0x1ff17e,
                                onChange: (_0x4404da) =>
                                  _0x5f5553(_0x4404da["target"]["value"]),
                                onBlur: () =>
                                  _0x5f5553(
                                    _0x1ff17e["trim"]()["replace"](/\/+$/, ""),
                                  ),
                              }),
                              _0x2ddfe2["jsx"]("p", {
                                className: "text-muted-foreground\x20text-xs",
                                children: _0x1aa03d("masterUrl.hint"),
                              }),
                              _0x2ddfe2["jsx"](_0x5185a8, {
                                size: "sm",
                                variant: "outline",
                                disabled:
                                  !_0x1ff17e["trim"]() ||
                                  _0x1ff17e["trim"]()["replace"](/\/+$/, "") ===
                                    (_0x11f05e?.["master_url"] || ""),
                                onClick: () => {
                                  (_0x2d0a4e([]),
                                    _0x946043(!0x1),
                                    _0x5da0c4(!0x1),
                                    _0xde3d23(!0x0));
                                },
                                children: "检查并迁移主控地址",
                              }),
                              _0x2ddfe2["jsx"](_0x1cc136, {
                                open: _0x102305,
                                onOpenChange: _0xde3d23,
                                children: _0x2ddfe2["jsxs"](_0x3e8fab, {
                                  className:
                                    "max-h-[85vh]\x20max-w-2xl\x20overflow-y-auto",
                                  children: [
                                    _0x2ddfe2["jsxs"](_0x333ff1, {
                                      children: [
                                        _0x2ddfe2["jsx"](_0x4b37d4, {
                                          children: "迁移主控地址",
                                        }),
                                        _0x2ddfe2["jsx"](_0x1e6ae9, {
                                          children:
                                            "系统会先让在线\x20Agent\x20从其所在服务器探测新地址，确认后才写入配置。请先在新主控恢复相同数据库和密钥，并在新主控把“主服务器地址”设置为这里的新地址。",
                                        }),
                                      ],
                                    }),
                                    _0x2ddfe2["jsxs"]("div", {
                                      className: "space-y-4",
                                      children: [
                                        _0x2ddfe2["jsxs"]("div", {
                                          className:
                                            "rounded-md\x20border\x20p-3\x20text-sm\x20break-all",
                                          children: [
                                            _0x2ddfe2["jsx"]("div", {
                                              className:
                                                "text-muted-foreground",
                                              children: "当前地址",
                                            }),
                                            _0x2ddfe2["jsx"]("div", {
                                              children:
                                                _0x11f05e?.["master_url"] ||
                                                "-",
                                            }),
                                            _0x2ddfe2["jsx"]("div", {
                                              className:
                                                "text-muted-foreground\x20mt-2",
                                              children: "新地址",
                                            }),
                                            _0x2ddfe2["jsx"]("div", {
                                              children: _0x1ff17e || "-",
                                            }),
                                          ],
                                        }),
                                        _0x2ddfe2["jsxs"]("div", {
                                          className:
                                            "flex\x20items-start\x20justify-between\x20gap-4\x20rounded-md\x20border\x20p-3",
                                          children: [
                                            _0x2ddfe2["jsxs"]("div", {
                                              children: [
                                                _0x2ddfe2["jsx"](_0x34df34, {
                                                  children: "是否更换域名",
                                                }),
                                                _0x2ddfe2["jsx"]("p", {
                                                  className:
                                                    "text-muted-foreground\x20text-xs",
                                                  children:
                                                    "关闭时新旧地址必须使用相同域名，只允许更换协议或端口。开启并确认迁移后，会自动关闭\x20Cloudflare\x20验证码，避免新域名无法登录；迁移后请重新配置并测试验证码。",
                                                }),
                                              ],
                                            }),
                                            _0x2ddfe2["jsx"](_0x543b01, {
                                              checked: _0x3b6d17,
                                              onCheckedChange: (_0x26ad17) => {
                                                (_0x12a0cf(_0x26ad17),
                                                  _0x2d0a4e([]),
                                                  _0x946043(!0x1));
                                              },
                                            }),
                                          ],
                                        }),
                                        _0x2ddfe2["jsxs"]("div", {
                                          className:
                                            "flex\x20items-start\x20justify-between\x20gap-4\x20rounded-md\x20border\x20p-3",
                                          children: [
                                            _0x2ddfe2["jsxs"]("div", {
                                              children: [
                                                _0x2ddfe2["jsx"](_0x34df34, {
                                                  children:
                                                    "是否迁移到另一台服务器",
                                                }),
                                                _0x2ddfe2["jsx"]("p", {
                                                  className:
                                                    "text-muted-foreground\x20text-xs",
                                                  children:
                                                    "关闭时保留本机回环和\x20Docker\x20内网\x20Agent\x20地址；开启时这些\x20Agent\x20也必须通过新地址检查。",
                                                }),
                                              ],
                                            }),
                                            _0x2ddfe2["jsx"](_0x543b01, {
                                              checked: _0x4378bf,
                                              onCheckedChange: (_0x2c79bf) => {
                                                (_0x39db92(_0x2c79bf),
                                                  _0x2d0a4e([]),
                                                  _0x946043(!0x1));
                                              },
                                            }),
                                          ],
                                        }),
                                        _0x406eb1["length"] > 0x0 &&
                                          _0x2ddfe2["jsx"]("div", {
                                            className: "space-y-2",
                                            children: _0x406eb1["map"](
                                              (_0x3c76fa) =>
                                                _0x2ddfe2["jsxs"](
                                                  "div",
                                                  {
                                                    className:
                                                      "flex\x20items-center\x20justify-between\x20gap-3\x20rounded-md\x20border\x20p-3\x20text-sm",
                                                    children: [
                                                      _0x2ddfe2["jsxs"]("div", {
                                                        children: [
                                                          _0x2ddfe2["jsx"](
                                                            "div",
                                                            {
                                                              className:
                                                                "font-medium",
                                                              children:
                                                                _0x3c76fa[
                                                                  "name"
                                                                ],
                                                            },
                                                          ),
                                                          _0x2ddfe2["jsx"](
                                                            "div",
                                                            {
                                                              className:
                                                                "text-muted-foreground\x20text-xs",
                                                              children:
                                                                _0x3c76fa[
                                                                  "message"
                                                                ] ||
                                                                (_0x3c76fa[
                                                                  "latency_ms"
                                                                ] != null
                                                                  ? _0x3c76fa[
                                                                      "latency_ms"
                                                                    ] + "\x20ms"
                                                                  : ""),
                                                            },
                                                          ),
                                                        ],
                                                      }),
                                                      _0x2ddfe2["jsx"]("span", {
                                                        className: _0x7633f5(
                                                          "font-medium",
                                                          _0x3c76fa[
                                                            "status"
                                                          ] === "ready" ||
                                                            _0x3c76fa[
                                                              "status"
                                                            ] === "preserved"
                                                            ? "text-green-600"
                                                            : _0x3c76fa[
                                                                  "status"
                                                                ] === "skipped"
                                                              ? "text-muted-foreground"
                                                              : "text-destructive",
                                                        ),
                                                        children:
                                                          _0x3c76fa["status"],
                                                      }),
                                                    ],
                                                  },
                                                  _0x3c76fa["server_id"],
                                                ),
                                            ),
                                          }),
                                        !_0x10f487 &&
                                          _0x406eb1["some"](
                                            (_0x15f04c) =>
                                              _0x15f04c["status"] ===
                                                "offline" ||
                                              _0x15f04c["status"] === "failed",
                                          ) &&
                                          _0x2ddfe2["jsxs"]("label", {
                                            className:
                                              "border-destructive/40\x20flex\x20items-start\x20gap-2\x20rounded-md\x20border\x20p-3\x20text-sm",
                                            children: [
                                              _0x2ddfe2["jsx"](_0x77832a, {
                                                checked: _0x3f7c69,
                                                onCheckedChange: (_0x14e7c4) =>
                                                  _0x5da0c4(_0x14e7c4 === !0x0),
                                              }),
                                              _0x2ddfe2["jsx"]("span", {
                                                children:
                                                  "我确认未通过检查的\x20Agent\x20可能失联，仍然迁移。离线\x20Agent\x20无法通过新主控自动获知新域名。",
                                              }),
                                            ],
                                          }),
                                      ],
                                    }),
                                    _0x2ddfe2["jsxs"](_0x5881af, {
                                      children: [
                                        _0x2ddfe2["jsx"](_0x5185a8, {
                                          variant: "outline",
                                          onClick: () =>
                                            _0x24110e["mutate"]("preview"),
                                          disabled: _0x24110e["isPending"],
                                          children: "检查\x20Agent",
                                        }),
                                        _0x2ddfe2["jsx"](_0x5185a8, {
                                          onClick: () =>
                                            _0x24110e["mutate"]("commit"),
                                          disabled:
                                            _0x24110e["isPending"] ||
                                            (!_0x10f487 && !_0x3f7c69),
                                          children: "确认迁移",
                                        }),
                                      ],
                                    }),
                                  ],
                                }),
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className: "space-y-2\x20pt-2",
                                children: [
                                  _0x2ddfe2["jsx"](_0x34df34, {
                                    htmlFor: "subscription-url",
                                    children: _0x1aa03d(
                                      "masterUrl.subscriptionTitle",
                                    ),
                                  }),
                                  _0x2ddfe2["jsx"](_0x549353, {
                                    id: "subscription-url",
                                    placeholder: _0x1aa03d(
                                      "masterUrl.subscriptionPlaceholder",
                                    ),
                                    value: _0x420042,
                                    onChange: (_0x3182cf) =>
                                      _0x43f9e2(_0x3182cf["target"]["value"]),
                                    onBlur: () => {
                                      const _0x601e7e = _0x420042["trim"]()[
                                        "replace"
                                      ](/\/+$/, "");
                                      (_0x43f9e2(_0x601e7e),
                                        _0x601e7e !==
                                          (_0x11f05e?.["subscription_url"] ||
                                            "") &&
                                          _0x1dadac["mutate"](_0x601e7e));
                                    },
                                    disabled: _0x1dadac["isPending"],
                                  }),
                                  _0x2ddfe2["jsx"]("p", {
                                    className:
                                      "text-muted-foreground\x20text-xs",
                                    children: _0x1aa03d(
                                      "masterUrl.subscriptionHint",
                                    ),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "flex\x20items-start\x20justify-between\x20gap-4\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "master-local-only",
                                        children: _0x1aa03d(
                                          "masterUrl.localOnly",
                                        ),
                                      }),
                                      _0x2ddfe2["jsx"]("p", {
                                        className:
                                          "text-destructive\x20text-xs",
                                        children: _0x11f05e?.["is_docker"]
                                          ? "Docker\x20环境不支持此功能，请使用端口映射或宿主机防火墙限制公网访问。"
                                          : _0x1aa03d(
                                              "masterUrl.localOnlyHint",
                                            ),
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsx"](_0x543b01, {
                                    id: "master-local-only",
                                    checked: _0x5d21bf,
                                    onCheckedChange: (_0x43382d) => {
                                      (_0x211f4e(_0x43382d),
                                        _0xda88d3["mutate"](_0x43382d));
                                    },
                                    disabled:
                                      _0xda88d3["isPending"] ||
                                      _0x11f05e?.["is_docker"],
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsxs"]("div", {
                                className:
                                  "space-y-3\x20rounded-lg\x20border\x20p-3",
                                children: [
                                  _0x2ddfe2["jsxs"]("div", {
                                    className:
                                      "flex\x20items-start\x20justify-between\x20gap-4",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        className: "space-y-1",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "https-recovery",
                                            children: "HTTPS\x20故障自愈",
                                          }),
                                          _0x2ddfe2["jsx"]("p", {
                                            className:
                                              "text-muted-foreground\x20text-xs",
                                            children:
                                              "仅当曾在线的\x20Agent\x20全部超时且公网\x20HTTPS\x20连续探测失败时，恢复\x20HTTP\x20公网监听；首次部署不会触发。",
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsx"](_0x543b01, {
                                        id: "https-recovery",
                                        checked: _0x5274a5,
                                        onCheckedChange: _0xd9dd4e,
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className: "space-y-1",
                                    children: [
                                      _0x2ddfe2["jsx"](_0x34df34, {
                                        htmlFor: "recovery-url",
                                        children: "HTTP\x20恢复地址",
                                      }),
                                      _0x2ddfe2["jsx"](_0x549353, {
                                        id: "recovery-url",
                                        placeholder:
                                          _0x11f05e?.["recovery_url"] ||
                                          "自动使用\x20http://主控域名:主控端口",
                                        value: _0x2af5a5,
                                        onChange: (_0xc136e) =>
                                          _0x1fc6ef(
                                            _0xc136e["target"]["value"],
                                          ),
                                      }),
                                      _0x2ddfe2["jsxs"]("p", {
                                        className:
                                          "text-muted-foreground\x20text-xs",
                                        children: [
                                          "留空自动使用：",
                                          _0x2ddfe2["jsx"]("span", {
                                            className: "font-mono",
                                            children:
                                              _0x11f05e?.["recovery_url"] ||
                                              "等待主控地址配置",
                                          }),
                                          "。仅在\x20CDN、NAT\x20或内外网地址不同时手动覆盖。",
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x2ddfe2["jsxs"]("div", {
                                    className:
                                      "grid\x20gap-3\x20sm:grid-cols-2",
                                    children: [
                                      _0x2ddfe2["jsxs"]("div", {
                                        className: "space-y-1",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "recovery-failure",
                                            children: "连续失败（分钟）",
                                          }),
                                          _0x2ddfe2["jsx"](_0x549353, {
                                            id: "recovery-failure",
                                            type: "number",
                                            min: 0x1,
                                            max: 0x3c,
                                            value: _0x5566fb,
                                            onChange: (_0xabc938) =>
                                              _0x63ad43(
                                                Number(
                                                  _0xabc938["target"]["value"],
                                                ),
                                              ),
                                          }),
                                        ],
                                      }),
                                      _0x2ddfe2["jsxs"]("div", {
                                        className: "space-y-1",
                                        children: [
                                          _0x2ddfe2["jsx"](_0x34df34, {
                                            htmlFor: "recovery-grace",
                                            children: "启动保护（分钟）",
                                          }),
                                          _0x2ddfe2["jsx"](_0x549353, {
                                            id: "recovery-grace",
                                            type: "number",
                                            min: 0x1,
                                            max: 0x78,
                                            value: _0x54682b,
                                            onChange: (_0x5a7db7) =>
                                              _0x646247(
                                                Number(
                                                  _0x5a7db7["target"]["value"],
                                                ),
                                              ),
                                          }),
                                        ],
                                      }),
                                    ],
                                  }),
                                  _0x11f05e?.["recovery_pending"] &&
                                    _0x2ddfe2["jsxs"]("p", {
                                      className: "text-destructive\x20text-xs",
                                      children: [
                                        "恢复待完成：",
                                        _0x11f05e["recovery_reason"],
                                      ],
                                    }),
                                  _0x2ddfe2["jsx"](_0x5185a8, {
                                    size: "sm",
                                    onClick: () => _0x580940["mutate"](),
                                    disabled: _0x580940["isPending"],
                                    children: "保存自愈设置",
                                  }),
                                ],
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsxs"]("div", {
                            className: "space-y-2\x20border-t\x20pt-4",
                            children: [
                              _0x2ddfe2["jsx"](_0x34df34, {
                                htmlFor: "redeem-template",
                                children: _0x1aa03d("redeemTemplate.title"),
                              }),
                              _0x2ddfe2["jsx"](_0x5d9440, {
                                id: "redeem-template",
                                rows: 0x8,
                                placeholder: _0x1aa03d(
                                  "redeemTemplate.placeholder",
                                ),
                                value: _0x478a8c,
                                onChange: (_0xe6ead0) =>
                                  _0x3b5e6d(_0xe6ead0["target"]["value"]),
                                onBlur: () => {
                                  _0x478a8c !==
                                    (_0x568d46?.["redeem_template"] || "") &&
                                    _0x567a5b["mutate"](_0x478a8c);
                                },
                                disabled: _0x567a5b["isPending"],
                              }),
                              _0x2ddfe2["jsx"]("p", {
                                className: "text-muted-foreground\x20text-xs",
                                children: _0x1aa03d("redeemTemplate.hint"),
                              }),
                            ],
                          }),
                        ],
                      }),
                    ],
                  }),
                  _0x2ddfe2["jsxs"](_0x2aeb40, {
                    children: [
                      _0x2ddfe2["jsxs"](_0x1db8ce, {
                        className: "pb-4",
                        children: [
                          _0x2ddfe2["jsx"](_0x30c50e, {
                            children: _0x1aa03d("apiToken.title"),
                          }),
                          _0x2ddfe2["jsx"](_0x54661d, {
                            children: _0x1aa03d("apiToken.description"),
                          }),
                        ],
                      }),
                      _0x2ddfe2["jsxs"](_0x42cb32, {
                        className: "space-y-3",
                        children: [
                          _0x2ddfe2["jsxs"]("div", {
                            className: "flex\x20items-center\x20gap-2",
                            children: [
                              _0x2ddfe2["jsxs"]("div", {
                                className: "relative\x20flex-1",
                                children: [
                                  _0x2ddfe2["jsx"](_0x549353, {
                                    type: _0x569631 ? "text" : "password",
                                    value: _0x436e38
                                      ? "..."
                                      : _0xcc0e07?.["token"] || "",
                                    readOnly: !0x0,
                                    className: "pr-10\x20font-mono\x20text-sm",
                                  }),
                                  _0x2ddfe2["jsx"](_0x5185a8, {
                                    type: "button",
                                    variant: "ghost",
                                    size: "sm",
                                    className:
                                      "absolute\x20top-0\x20right-0\x20h-full\x20px-3\x20hover:bg-transparent",
                                    onClick: () => _0x1c209a(!_0x569631),
                                    children: _0x569631
                                      ? _0x2ddfe2["jsx"](_0x2d8607, {
                                          className:
                                            "text-muted-foreground\x20h-4\x20w-4",
                                        })
                                      : _0x2ddfe2["jsx"](_0x1814ad, {
                                          className:
                                            "text-muted-foreground\x20h-4\x20w-4",
                                        }),
                                  }),
                                ],
                              }),
                              _0x2ddfe2["jsx"](_0x5185a8, {
                                type: "button",
                                variant: "outline",
                                size: "icon",
                                onClick: _0x290ac2,
                                disabled: _0x436e38 || !_0xcc0e07?.["token"],
                                children: _0x2ddfe2["jsx"](_0x3061be, {
                                  className: "h-4\x20w-4",
                                }),
                              }),
                              _0x2ddfe2["jsx"](_0x5185a8, {
                                type: "button",
                                variant: "outline",
                                size: "icon",
                                onClick: () => _0x492328(!0x0),
                                disabled: _0x436e38 || _0x38366f["isPending"],
                                children: _0x2ddfe2["jsx"](_0x18a4f9, {
                                  className:
                                    "h-4\x20w-4\x20" +
                                    (_0x38366f["isPending"]
                                      ? "animate-spin"
                                      : ""),
                                }),
                              }),
                            ],
                          }),
                          _0x2ddfe2["jsx"]("p", {
                            className:
                              "text-muted-foreground\x20text-xs\x20whitespace-pre-line",
                            children: _0x1aa03d("apiToken.warning"),
                          }),
                        ],
                      }),
                    ],
                  }),
                ],
              }),
              _0x2ddfe2["jsx"](_0x550af3, {
                value: "database",
                className: "space-y-6",
                children: _0x2ddfe2["jsx"](Fr, {}),
              }),
            ],
          }),
        ],
      }),
      _0x2ddfe2["jsx"](_0x1e06f2, {
        open: _0x3f4852,
        onOpenChange: _0x31e9f0,
        title: _0x1aa03d("silentMode.enableConfirm"),
        desc: _0x2ddfe2["jsx"]("div", {
          className: "text-muted-foreground",
          children: _0x1aa03d("silentMode.description"),
        }),
        cancelBtnText: _0x1aa03d("actions.cancel", { ns: "common" }),
        confirmText: _0x1aa03d("silentMode.enableLabel"),
        destructive: !0x0,
        handleConfirm: () => {
          (_0x31e9f0(!0x1),
            _0x451e75(!0x0),
            _0x842682["mutate"]({
              silent_mode: !0x0,
              silent_mode_timeout: _0x1f6fa6,
            }));
        },
      }),
      _0x2ddfe2["jsx"](_0x1e06f2, {
        open: _0x16e9a0,
        onOpenChange: _0x492328,
        title: _0x1aa03d("apiToken.regenerateConfirm"),
        desc: _0x2ddfe2["jsx"]("div", {
          className: "text-muted-foreground\x20whitespace-pre-line",
          children: _0x1aa03d("apiToken.warning"),
        }),
        cancelBtnText: _0x1aa03d("actions.cancel", { ns: "common" }),
        confirmText: _0x1aa03d("apiToken.regenerate", {
          defaultValue: "重新生成",
        }),
        destructive: !0x0,
        isLoading: _0x38366f["isPending"],
        handleConfirm: () => {
          (_0x492328(!0x1), _0x38366f["mutate"]());
        },
      }),
    ],
  });
}
function Gr({ siteKey: _0x2daff0, secretKeyConfigured: _0x4ac7dd }) {
  const { t: _0x4ccf9c } = _0x2443c3("system"),
    {
      containerRef: _0x5a0ac8,
      token: _0x30f33c,
      reset: _0x309152,
    } = _0x472a2b(_0x2daff0 || void 0x0),
    [_0x432334, _0x45da84] = _0x46e13c["useState"]("idle"),
    [_0x139e0f, _0x5529c8] = _0x46e13c["useState"](""),
    [_0x4a5fa3, _0x5b51b0] = _0x46e13c["useState"]([]),
    [_0x590ffc, _0x588c0d] = _0x46e13c["useState"](""),
    _0x3e4b5f = () => {
      _0x30f33c &&
        (_0x45da84("verifying"),
        _0x5b51b0([]),
        _0x588c0d(""),
        _0x495bb8["post"]("/api/admin/security-settings/turnstile/test", {
          token: _0x30f33c,
        })
          ["then"]((_0x2e8ee1) => {
            const _0x2cc45b = _0x2e8ee1["data"];
            _0x2cc45b["success"]
              ? (_0x45da84("success"), _0x5529c8(_0x2cc45b["hostname"] || ""))
              : (_0x45da84("failed"),
                _0x5b51b0(_0x2cc45b["error_codes"] || []));
          })
          ["catch"]((_0x1c7ae5) => {
            (_0x45da84("failed"),
              _0x588c0d(
                _0x1c7ae5?.["response"]?.["data"]?.["error"] ||
                  _0x1c7ae5?.["message"] ||
                  "request\x20failed",
              ));
          }));
    };
  if (!_0x2daff0) return null;
  if (!_0x4ac7dd)
    return _0x2ddfe2["jsx"]("div", {
      className:
        "rounded-md\x20border\x20border-amber-500/30\x20bg-amber-500/5\x20p-3\x20text-xs\x20text-amber-700\x20dark:text-amber-400",
      children: _0x4ccf9c("turnstile.testRequiresSecret"),
    });
  const _0x401bee = () => {
    (_0x309152(),
      _0x45da84("idle"),
      _0x5b51b0([]),
      _0x588c0d(""),
      _0x5529c8(""));
  };
  return _0x2ddfe2["jsxs"]("div", {
    className: "bg-muted/30\x20space-y-2\x20rounded-md\x20border\x20p-3",
    children: [
      _0x2ddfe2["jsxs"]("div", {
        className: "flex\x20items-center\x20gap-2\x20text-sm\x20font-medium",
        children: [
          _0x4ccf9c("turnstile.testTitle"),
          _0x432334 === "idle" &&
            _0x2ddfe2["jsxs"]("span", {
              className: "text-muted-foreground\x20text-xs",
              children: ["·\x20", _0x4ccf9c("turnstile.testWaitingUser")],
            }),
          _0x432334 === "verifying" &&
            _0x2ddfe2["jsx"](_0x4e8a70, {
              className: "h-3\x20w-3\x20animate-spin",
            }),
          _0x432334 === "success" &&
            _0x2ddfe2["jsx"](_0x157104, {
              className: "h-4\x20w-4\x20text-emerald-600",
            }),
          _0x432334 === "failed" &&
            _0x2ddfe2["jsx"](_0x5749cc, {
              className: "h-4\x20w-4\x20text-red-600",
            }),
        ],
      }),
      _0x2ddfe2["jsx"]("div", { ref: _0x5a0ac8 }),
      _0x432334 === "success" &&
        _0x2ddfe2["jsxs"]("p", {
          className: "text-xs\x20text-emerald-700\x20dark:text-emerald-400",
          children: [
            "✓\x20",
            _0x4ccf9c("turnstile.testSuccess", { hostname: _0x139e0f || "—" }),
          ],
        }),
      _0x432334 === "failed" &&
        _0x2ddfe2["jsxs"]("div", {
          className:
            "space-y-1\x20text-xs\x20text-red-700\x20dark:text-red-400",
          children: [
            _0x2ddfe2["jsxs"]("p", {
              children: ["✗\x20", _0x4ccf9c("turnstile.testFailed")],
            }),
            _0x4a5fa3["length"] > 0x0 &&
              _0x2ddfe2["jsx"]("p", {
                className: "font-mono",
                children: _0x4a5fa3["join"](",\x20"),
              }),
            _0x590ffc &&
              _0x4a5fa3["length"] === 0x0 &&
              _0x2ddfe2["jsx"]("p", { children: _0x590ffc }),
            _0x4a5fa3["includes"]("invalid-input-secret") &&
              _0x2ddfe2["jsx"]("p", {
                className: "text-muted-foreground",
                children: _0x4ccf9c("turnstile.errInvalidSecret"),
              }),
            _0x4a5fa3["includes"]("missing-input-secret") &&
              _0x2ddfe2["jsx"]("p", {
                className: "text-muted-foreground",
                children: _0x4ccf9c("turnstile.errMissingSecret"),
              }),
            _0x4a5fa3["includes"]("timeout-or-duplicate") &&
              _0x2ddfe2["jsx"]("p", {
                className: "text-muted-foreground",
                children: _0x4ccf9c("turnstile.errTimeoutOrDup"),
              }),
            _0x4a5fa3["includes"]("invalid-input-response") &&
              _0x2ddfe2["jsx"]("p", {
                className: "text-muted-foreground",
                children: _0x4ccf9c("turnstile.errInvalidResponse"),
              }),
          ],
        }),
      _0x432334 === "idle" || _0x432334 === "verifying"
        ? _0x2ddfe2["jsxs"](_0x5185a8, {
            size: "sm",
            onClick: _0x3e4b5f,
            disabled: !_0x30f33c || _0x432334 === "verifying",
            children: [
              _0x432334 === "verifying"
                ? _0x2ddfe2["jsx"](_0x4e8a70, {
                    className: "mr-1\x20h-3\x20w-3\x20animate-spin",
                  })
                : null,
              _0x4ccf9c("turnstile.testButton"),
            ],
          })
        : _0x2ddfe2["jsx"](_0x5185a8, {
            size: "sm",
            variant: "outline",
            onClick: _0x401bee,
            children: _0x4ccf9c("turnstile.testRetry"),
          }),
    ],
  });
}
function zr({
  server: _0x3c94bc,
  regions: _0xcb8827,
  source: _0x145f34,
  globalTargets: _0x36ff70,
  override: _0x191924,
  onChange: _0x4a9625,
}) {
  const [_0x532dd6, _0x3126c8] = _0x46e13c["useState"](!0x1),
    _0x320f56 = _0x191924 !== void 0x0;
  return _0x2ddfe2["jsxs"]("div", {
    className: "rounded-md\x20border",
    children: [
      _0x2ddfe2["jsxs"]("button", {
        type: "button",
        className:
          "hover:bg-accent/40\x20flex\x20w-full\x20items-center\x20gap-2\x20px-3\x20py-2\x20text-sm",
        onClick: () => _0x3126c8((_0x12b834) => !_0x12b834),
        children: [
          _0x2ddfe2["jsx"]("span", {
            className: "text-muted-foreground",
            children: _0x532dd6 ? "▾" : "▸",
          }),
          _0x2ddfe2["jsx"]("span", {
            className: "truncate",
            children: _0x3c94bc["name"],
          }),
          _0x2ddfe2["jsx"]("span", {
            className:
              "text-muted-foreground\x20ml-auto\x20shrink-0\x20text-xs",
            children: _0x320f56
              ? "单独指定(" + _0x191924["length"] + ")"
              : "跟随全局",
          }),
        ],
      }),
      _0x532dd6 &&
        _0x2ddfe2["jsxs"]("div", {
          className: "space-y-3\x20border-t\x20px-3\x20py-3",
          children: [
            _0x2ddfe2["jsxs"](_0x2db305, {
              value: _0x320f56 ? "custom" : "global",
              onValueChange: (_0x32900e) =>
                _0x4a9625(
                  _0x32900e === "custom" ? (_0x191924 ?? []) : void 0x0,
                ),
              children: [
                _0x2ddfe2["jsxs"]("label", {
                  className:
                    "flex\x20cursor-pointer\x20items-center\x20gap-2\x20text-sm",
                  children: [
                    _0x2ddfe2["jsx"](_0x356e54, {
                      value: "global",
                      id: "ping-global-" + _0x3c94bc["id"],
                    }),
                    _0x2ddfe2["jsx"]("span", { children: "跟随全局设置" }),
                  ],
                }),
                !_0x320f56 &&
                  _0x2ddfe2["jsx"]("div", {
                    className: "text-muted-foreground\x20pl-6\x20text-xs",
                    children:
                      _0x36ff70["length"] > 0x0
                        ? "当前全局:" +
                          _0x36ff70["map"]((_0x462663) => _0x462663["label"])[
                            "join"
                          ]("、")
                        : "全局未选任何目标\x20→\x20此服务器不做\x20ping\x20探测",
                  }),
                _0x2ddfe2["jsxs"]("label", {
                  className:
                    "flex\x20cursor-pointer\x20items-center\x20gap-2\x20text-sm",
                  children: [
                    _0x2ddfe2["jsx"](_0x356e54, {
                      value: "custom",
                      id: "ping-custom-" + _0x3c94bc["id"],
                    }),
                    _0x2ddfe2["jsx"]("span", {
                      children: "为此服务器单独指定",
                    }),
                  ],
                }),
              ],
            }),
            _0x320f56 &&
              _0x2ddfe2["jsx"](wn, {
                regions: _0xcb8827,
                source: _0x145f34,
                selected: _0x191924,
                max: Nn,
                onChange: (_0x17239e) => _0x4a9625(_0x17239e),
                emptyHint: _0x2ddfe2["jsx"]("div", {
                  className: "text-xs\x20text-yellow-600",
                  children:
                    "⚠\x20一个目标都不选\x20=\x20此服务器不做\x20ping\x20探测(想沿用全局请选上面的「跟随全局设置」)",
                }),
              }),
          ],
        }),
    ],
  });
}
function Vr() {
  const { t: _0x55d653 } = _0x2443c3("system"),
    _0x46712b = _0x29acb4(),
    [_0x59d229, _0x226636] = _0x46e13c["useState"](new Set()),
    [_0x5bd0f3, _0x1a731b] = _0x46e13c["useState"](!0x1),
    { data: _0x158ca3, isLoading: _0x5cfbea } = _0x37d641({
      queryKey: ["reality-domain-share"],
      queryFn: async () =>
        (await _0x495bb8["get"]("/api/admin/remote/reality-domains/share"))[
          "data"
        ],
    }),
    _0x216c15 = () =>
      _0x46712b["invalidateQueries"]({ queryKey: ["reality-domain-share"] }),
    _0x624694 = _0x144b3f({
      mutationFn: async (_0x3bd59d) =>
        (
          await _0x495bb8["post"](
            "/api/admin/remote/reality-domains/share/toggle",
            _0x3bd59d,
          )
        )["data"],
      onSuccess: (_0x5d4df8, _0x5326f2) => {
        (_0x1a731b(!0x1),
          _0x216c15(),
          _0x5326f2["enabled"]
            ? _0x54e43f["success"](
                _0x55d653("realityShare.shared", {
                  count: _0x5d4df8?.["accepted"]?.["length"] ?? 0x0,
                }),
              )
            : _0x54e43f["success"](_0x55d653("realityShare.disabled")));
      },
      onError: (_0x1721d3) =>
        _0x54e43f["error"](
          _0x1721d3?.["response"]?.["data"]?.["error"] ||
            _0x55d653("realityShare.toggleFailed"),
        ),
    }),
    _0x15175e = _0x144b3f({
      mutationFn: (_0x34447f) =>
        _0x495bb8["post"]("/api/admin/remote/reality-domains/share/withdraw", {
          domain: _0x34447f,
        }),
      onSuccess: () => {
        (_0x216c15(),
          _0x54e43f["success"](_0x55d653("realityShare.withdrawn")));
      },
      onError: (_0x1e4cbd) =>
        _0x54e43f["error"](
          _0x1e4cbd?.["response"]?.["data"]?.["error"] ||
            _0x55d653("realityShare.withdrawFailed"),
        ),
    }),
    _0x5ba136 = _0x158ca3?.["pending"] ?? [],
    _0x11f95c = _0x158ca3?.["shared"] ?? [],
    _0x1a89f4 = _0x158ca3?.["licensed"] ?? !0x1,
    _0x152cdb = () => {
      (_0x226636(new Set(_0x5ba136)), _0x1a731b(!0x0));
    },
    _0x3017e1 = (_0x483bec) =>
      _0x226636((_0x1e8193) => {
        const _0x54b976 = new Set(_0x1e8193);
        return (
          _0x54b976["has"](_0x483bec)
            ? _0x54b976["delete"](_0x483bec)
            : _0x54b976["add"](_0x483bec),
          _0x54b976
        );
      });
  return _0x2ddfe2["jsxs"](_0x2aeb40, {
    children: [
      _0x2ddfe2["jsxs"](_0x1db8ce, {
        className: "pb-4",
        children: [
          _0x2ddfe2["jsx"](_0x30c50e, {
            children: _0x55d653("realityShare.title"),
          }),
          _0x2ddfe2["jsx"](_0x54661d, {
            children: _0x55d653("realityShare.description"),
          }),
        ],
      }),
      _0x2ddfe2["jsxs"](_0x42cb32, {
        className: "space-y-4",
        children: [
          _0x2ddfe2["jsxs"]("div", {
            className:
              "flex\x20items-center\x20justify-between\x20rounded-lg\x20border\x20p-3",
            children: [
              _0x2ddfe2["jsxs"]("div", {
                className: "space-y-1",
                children: [
                  _0x2ddfe2["jsx"](_0x34df34, {
                    htmlFor: "reality-share-toggle",
                    className: "cursor-pointer",
                    children: _0x55d653("realityShare.enableLabel"),
                  }),
                  _0x2ddfe2["jsx"]("p", {
                    className: "text-muted-foreground\x20text-xs",
                    children: _0x1a89f4
                      ? _0x55d653("realityShare.poolSize", {
                          count: _0x158ca3?.["pool_size"] ?? 0x0,
                        })
                      : _0x55d653("realityShare.proRequired"),
                  }),
                ],
              }),
              _0x2ddfe2["jsx"](_0x543b01, {
                id: "reality-share-toggle",
                checked: _0x158ca3?.["enabled"] ?? !0x1,
                disabled: !_0x1a89f4 || _0x5cfbea || _0x624694["isPending"],
                onCheckedChange: (_0x18c505) => {
                  _0x18c505
                    ? _0x152cdb()
                    : _0x624694["mutate"]({ enabled: !0x1, domains: [] });
                },
              }),
            ],
          }),
          _0x11f95c["length"] > 0x0 &&
            _0x2ddfe2["jsxs"]("div", {
              className: "space-y-2",
              children: [
                _0x2ddfe2["jsx"]("p", {
                  className: "text-sm\x20font-medium",
                  children: _0x55d653("realityShare.sharedList", {
                    count: _0x11f95c["length"],
                  }),
                }),
                _0x2ddfe2["jsx"]("div", {
                  className: "max-h-48\x20space-y-2\x20overflow-auto",
                  children: _0x11f95c["map"]((_0x353263) =>
                    _0x2ddfe2["jsxs"](
                      "div",
                      {
                        className:
                          "flex\x20items-center\x20justify-between\x20gap-2\x20rounded-md\x20border\x20p-2",
                        children: [
                          _0x2ddfe2["jsx"]("span", {
                            className:
                              "min-w-0\x20flex-1\x20truncate\x20text-sm",
                            children: _0x353263,
                          }),
                          _0x2ddfe2["jsx"](_0x5185a8, {
                            type: "button",
                            variant: "ghost",
                            size: "sm",
                            disabled: _0x15175e["isPending"],
                            onClick: () => _0x15175e["mutate"](_0x353263),
                            children: _0x55d653("realityShare.withdraw"),
                          }),
                        ],
                      },
                      _0x353263,
                    ),
                  ),
                }),
              ],
            }),
        ],
      }),
      _0x2ddfe2["jsx"](_0x1cc136, {
        open: _0x5bd0f3,
        onOpenChange: _0x1a731b,
        children: _0x2ddfe2["jsxs"](_0x3e8fab, {
          className: "max-w-lg",
          children: [
            _0x2ddfe2["jsxs"](_0x333ff1, {
              children: [
                _0x2ddfe2["jsx"](_0x4b37d4, {
                  children: _0x55d653("realityShare.previewTitle"),
                }),
                _0x2ddfe2["jsx"](_0x1e6ae9, {
                  children: _0x55d653("realityShare.previewDesc"),
                }),
              ],
            }),
            _0x5ba136["length"] === 0x0
              ? _0x2ddfe2["jsx"]("p", {
                  className: "text-muted-foreground\x20text-sm",
                  children: _0x55d653("realityShare.nothingToShare"),
                })
              : _0x2ddfe2["jsx"]("div", {
                  className: "max-h-[45vh]\x20space-y-2\x20overflow-auto",
                  children: _0x5ba136["map"]((_0x19c372) =>
                    _0x2ddfe2["jsxs"](
                      "label",
                      {
                        className:
                          "flex\x20cursor-pointer\x20items-center\x20gap-3\x20rounded-md\x20border\x20p-2",
                        children: [
                          _0x2ddfe2["jsx"](_0x77832a, {
                            checked: _0x59d229["has"](_0x19c372),
                            onCheckedChange: () => _0x3017e1(_0x19c372),
                          }),
                          _0x2ddfe2["jsx"]("span", {
                            className:
                              "min-w-0\x20flex-1\x20truncate\x20text-sm",
                            children: _0x19c372,
                          }),
                        ],
                      },
                      _0x19c372,
                    ),
                  ),
                }),
            _0x2ddfe2["jsxs"](_0x5881af, {
              children: [
                _0x2ddfe2["jsx"](_0x5185a8, {
                  variant: "outline",
                  onClick: () => _0x1a731b(!0x1),
                  children: _0x55d653("realityShare.cancel"),
                }),
                _0x2ddfe2["jsx"](_0x5185a8, {
                  disabled: _0x59d229["size"] === 0x0 || _0x624694["isPending"],
                  onClick: () =>
                    _0x624694["mutate"]({
                      enabled: !0x0,
                      domains: Array["from"](_0x59d229),
                    }),
                  children: _0x55d653("realityShare.confirmShare", {
                    count: _0x59d229["size"],
                  }),
                }),
              ],
            }),
          ],
        }),
      }),
    ],
  });
}
export { yi as component };
