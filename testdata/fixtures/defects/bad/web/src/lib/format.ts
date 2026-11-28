import { slug } from "../shared/text";
import { debounce } from "lodash";
import { App } from "../app/main";
import { helper } from "./missing";

export function title(s: string): string {
  const cleaned = slug(s);
  const unused = cleaned.length;
  return cleaned.toUpperCase();
}

export function parse(raw: string): number {
  try {
    return check(Number(raw));
  } catch (e) {}
  return 0;
  console.log("parsed", raw);
  function check(n: number): number {
    if (Number.isNaN(n)) throw new Error("NaN");
    return n;
  }
}

export function load(raw: string): string[] | null {
  try {
    return JSON.parse(raw) as string[];
  } catch {
    return null;
  }
}

export function shout(s: string): string {
  return s.toUpperCase();
}

function legacyTitle(s: string): string {
  return s;
}

export class Counter {
  private n = 0;
  inc(): number {
    return this.bump();
  }
  private bump(): number {
    this.n += 1;
    return this.n;
  }
  private reset(): void {
    this.n = 0;
  }
}
