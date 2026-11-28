// The folder tree behind the Modules treemap: modules (Go packages,
// TypeScript or C# directories) placed in their folders, with chains of
// single-child folders merged ("internal/lang" rather than two levels).
import type { GraphModule } from "../core/api";
import { moduleLabel } from "./format";

export interface TreeNode {
  name: string; // shown label (may be a merged chain "a/b")
  path: string; // repository path ("" = root)
  module?: GraphModule; // set on a module tile
  children: TreeNode[];
  lines: number;
  debt: number;
  density: number; // debt per 1,000 source lines
}

export type SizeBy = "lines" | "debt";

const OWN = "(own files)";

export function buildTree(modules: GraphModule[]): TreeNode {
  const root: TreeNode = { name: "", path: "", children: [], lines: 0, debt: 0, density: 0 };
  const folder = (parts: string[]): TreeNode => {
    let n = root;
    for (let i = 0; i < parts.length; i++) {
      const path = parts.slice(0, i + 1).join("/");
      let c = n.children.find((x) => x.path === path && !x.module);
      if (!c) n.children.push((c = { name: parts[i], path, children: [], lines: 0, debt: 0, density: 0 }));
      n = c;
    }
    return n;
  };
  for (const m of modules) {
    const label = moduleLabel(m.id);
    const parts = label === "." || label === "" ? [] : label.split("/");
    const f = folder(parts);
    // A module is a tile inside its own folder; the tile is the folder's
    // "(own files)" when the folder also has subfolders, else the folder.
    f.children.push({ name: OWN, path: parts.join("/"), module: m, children: [], lines: m.lines, debt: m.debt.total, density: m.density });
  }
  total(root);
  return collapse(root);
}

function total(n: TreeNode) {
  if (n.module) return;
  n.lines = 0;
  n.debt = 0;
  for (const c of n.children) {
    total(c);
    n.lines += c.lines;
    n.debt += c.debt;
  }
  n.density = n.lines ? (n.debt * 1000) / n.lines : 0;
}

// collapse turns a folder whose only child is its own module into that
// module's tile, and merges single-child folder chains.
function collapse(n: TreeNode): TreeNode {
  n.children = n.children.map(collapse);
  if (n.module) return n;
  if (n.children.length === 1 && n.path !== "") {
    const c = n.children[0];
    if (c.module && c.path === n.path) return { ...c, name: n.name };
    return { ...c, name: n.name + "/" + c.name };
  }
  return n;
}

// findNode returns the folder or module tile at path, or null.
export function findNode(root: TreeNode, path: string): TreeNode | null {
  if (root.path === path && !root.module) return root;
  const here = root.children.find((c) => !c.module && c.path === path) ?? root.children.find((c) => c.module && c.path === path);
  if (here) return here;
  for (const c of root.children.filter((x) => !x.module && within(path, x.path))) {
    const r = findNode(c, path);
    if (r) return r;
  }
  return null;
}

const within = (path: string, dir: string) => dir === "" || path.startsWith(dir + "/");

// parentPath is the folder above path ("" for top-level entries).
export function parentPath(path: string): string {
  const i = path.lastIndexOf("/");
  return i < 0 ? "" : path.slice(0, i);
}

export function sizeOf(n: TreeNode, by: SizeBy): number {
  return Math.max(1, by === "debt" ? n.debt : n.lines);
}

// Debt-per-kLOC bands of the sequential ramp; band 0 is lowest.
export const DENSITY_STEPS = [30, 60, 105, 150];
export function densityBand(v: number): number {
  let b = 0;
  while (b < DENSITY_STEPS.length && v >= DENSITY_STEPS[b]) b++;
  return b;
}
