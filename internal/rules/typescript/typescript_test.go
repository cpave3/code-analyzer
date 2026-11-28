package typescript

import (
	"strings"
	"testing"

	lts "github.com/cameron/code-evaluator/internal/lang/typescript"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
	"github.com/cameron/code-evaluator/internal/rules/rulestest"
)

type tc struct {
	name  string
	rule  string
	files map[string]string
	want  []string // "severity:substring" per expected finding, in order
}

func check(t *testing.T, cases []tc) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := rulestest.Run(t, lts.New(), Register, c.rule, c.files)
			if len(got) != len(c.want) {
				t.Fatalf("got %d findings, want %d:\n%s", len(got), len(c.want), strings.Join(rulestest.Summary(got), "\n"))
			}
			for i, w := range c.want {
				sev, sub, _ := strings.Cut(w, ":")
				if string(got[i].Severity) != sev || !strings.Contains(got[i].Message, sub) {
					t.Errorf("finding %d = %s %q, want %s containing %q", i, got[i].Severity, got[i].Message, sev, sub)
				}
			}
		})
	}
}

func TestDeadCode(t *testing.T) {
	check(t, []tc{
		{"imports", catalog.UnusedImport, map[string]string{
			"src/a.tsx": `import React from "react";
import type { Props } from "./types";
import * as ns from "./ns";
import def, { a, b as bee, c } from "./x";
import "./side-effect.css";
export function A(p: Props) {
  return <div>{a}{ns.v}</div>;
}
export { c };
`,
			"src/b.ts": `import React from "react";
export const x = 1;
`}, []string{"gate:imported def from", "gate:imported bee from", "gate:imported React from \"react\""}},
		{"variables", catalog.UnusedVariable, map[string]string{
			"src/a.ts": `import { g } from "./g";
const moduleUnused = 1;
const moduleUsed = 2;
export const exported = 3;
const helper = () => 1;
export function f(): number {
  const used = g();
  let unused = 0;
  const _ignored = 1;
  const [d] = [1];
  const inner = () => used;
  const deadInner = () => 0;
  for (const k of [1]) {}
  return inner() + moduleUsed;
}
`,
			"src/script.ts": "const globalScript = 1;\n",
		}, []string{"gate:variable moduleUnused", "gate:variable unused is", "gate:variable deadInner"}},
		{"unreachable", catalog.UnreachableCode, map[string]string{"src/a.ts": `export function f(x: number): number {
  if (x) {
    return 1;
    x++;
  }
  switch (x) {
    case 1:
      throw new Error("x");
      x--;
    default:
      break;
  }
  return g();
  function g() { return 2; }
  type T = number;
  interface I {}
}
`}, []string{"gate:`x++;` follows `return 1;`", "gate:`x--;` follows `throw new Error(\"x\");`"}},
		{"private functions", catalog.UnusedPrivateFunction, map[string]string{
			"src/a.ts": `export function pub() { return used(); }
function used() { return 1; }
function dead() { return 2; }
const deadArrow = () => 3;
const wrapped = memo(() => 4);
function reExported() {}
export { reExported };
export class C {
  private kept() { return 1; }
  private dropped() { return 2; }
  #hidden() { return 3; }
  run() { return this.kept(); }
}
`,
			"src/script.ts": "function globalHelper() {}\n",
		}, []string{"gate:function dead is", "gate:function deadArrow", "gate:function wrapped", "gate:private method C.dropped", "gate:private method C.#hidden"}},
		{"exports", catalog.UnusedExport, map[string]string{
			"src/lib.ts":   "export function used() {}\nexport function dead() {}\nexport function internal() {}\nexport const deadArrow = () => 1;\nexport default function main() {}\nconst x = internal;\n",
			"src/ns.ts":    "export function anything() {}\n",
			"src/entry.ts": "export function neverImported() {}\n",
			"src/app.ts":   "import { used } from './lib';\nimport * as ns from './ns';\nused(); ns;\n",
		}, []string{"advisory:exported function dead", "advisory:exported function deadArrow"}},
	})
}

func TestLostErrors(t *testing.T) {
	check(t, []tc{
		{"empty handlers", catalog.EmptyHandler, map[string]string{"src/a.ts": `export async function f(p: Promise<number>) {
  try { await p; } catch (e) {}
  try { await p; } catch {
    // ignored: best effort
  }
  try { await p; } catch (e) { console.error(e); }
  p.catch(() => {});
  p.then((v) => v, function () {});
  p.catch((e) => log(e));
}
function log(e: unknown) { return e; }
`}, []string{"gate:empty catch block", "advisory:holds only a comment", "gate:.catch() with an empty", "gate:.then() with an empty"}},
		{"error to success", catalog.ErrorToSuccess, map[string]string{"src/a.ts": `export function f(raw: string) {
  try { return JSON.parse(raw); } catch { return null; }
}
export function g(raw: string) {
  try { return JSON.parse(raw); } catch (e) { return []; }
}
export function h(raw: string) {
  try { return JSON.parse(raw); } catch (e) { return { error: String(e) }; }
}
export const p = fetch("x").catch(() => undefined);
export const q = fetch("x").catch((e) => ({ e }));
`}, []string{"advisory:`return null;`", "advisory:`return [];`", "advisory:.catch(() => undefined)"}},
	})
}

func TestDependencies(t *testing.T) {
	files := map[string]string{
		"package.json":           `{"name":"root","devDependencies":{"typescript":"5"}}`,
		"web/package.json":       `{"name":"web","dependencies":{"react":"18","left-pad":"1","tailwind-thing":"1","@types/node":"20"},"devDependencies":{"@types/estree":"1"},"scripts":{"x":"cli-only-dep run"}}`,
		"web/tailwind.config.js": `module.exports = { plugins: [require("tailwind-thing")] };`,
		"web/src/a.ts": `import { useState } from "react";
import fs from "node:fs";
import path from "path";
import type { Node } from "estree";
import typescript from "typescript";
import missingPkg from "missing-pkg";
import { x } from "./x";
import { y } from "./y.js";
import data from "./data.json";
import icon from "./icon.svg?url";
import gone from "./gone";
import virt from "virtual:pwa";
import amb from "ambient-mod";
import { shared } from "@acme/shared";
import { nope } from "@acme/shared/nope";
`,
		"web/src/x.ts":             "export const x = 1;\n",
		"web/src/y.js":             "export const y = 1;\n",
		"web/src/data.json":        "{}",
		"web/src/icon.svg":         "<svg/>",
		"web/src/env.d.ts":         "declare module \"ambient-mod\";\n",
		"pkgs/shared/package.json": `{"name":"@acme/shared"}`,
		"pkgs/shared/src/index.ts": "export const shared = 1;\n",
	}
	check(t, []tc{
		{"unresolved", catalog.UnresolvedImport, files, []string{"gate:relative import \"./gone\"", "advisory:workspace import \"@acme/shared/nope\""}},
		{"undeclared", catalog.UndeclaredDependency, files, []string{"gate:package \"missing-pkg\""}},
		{"unused", catalog.UnusedDependency, files, []string{"advisory:dependency \"left-pad\""}},
	})
}

func TestTestOracle(t *testing.T) {
	check(t, []tc{
		{"tautologies", catalog.TautologicalAssertion, map[string]string{"src/a.test.ts": `import assert from "node:assert";
it("x", () => {
  const a = 1;
  expect(a).toBe(a);
  expect(1).toEqual(1);
  expect(a).not.toBe(a);
  expect(f()).toBe(f());
  expect(true).toBeTruthy();
  expect(a).toBeTruthy();
  assert.deepEqual(a, a);
  assert(true);
  expect(a.b).to.deep.equal(a.b);
});
`, "src/a.ts": "expect(1).toBe(1);\n"},
			[]string{"gate:compares a with itself", "gate:compares two literals", "gate:asserts a constant", "gate:compares a with itself", "gate:asserts a constant", "gate:compares a.b with itself"}},
		{"unawaited", catalog.LostAsyncAssertion, map[string]string{"src/a.test.ts": `test("p", async () => {
  expect(load()).resolves.toBe(1);
  expect(load()).rejects.toThrow("x");
  await expect(load()).resolves.toBe(1);
  return expect(load()).resolves.toBe(1);
});
test("q", () => expect(load()).resolves.toBe(1));
`}, []string{"gate:`expect(load()).resolves.toBe(1)` is neither", "gate:`expect(load()).rejects.toThrow(\"x\")`"}},
	})
}
