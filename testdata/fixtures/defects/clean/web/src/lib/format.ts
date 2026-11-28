import dayjs from "dayjs";
import { slug } from "../shared/text";

export function title(s: string): string {
  const cleaned = slug(s);
  return cleaned.toUpperCase();
}

export function stamp(d: Date): string {
  return dayjs(d).format("YYYY-MM-DD");
}

export function parse(raw: string): number {
  try {
    return check(Number(raw));
  } catch (e) {
    throw new Error(`parse ${raw}: ${e}`);
  }
  function check(n: number): number {
    if (Number.isNaN(n)) throw new Error("NaN");
    return n;
  }
}

export function load(raw: string): string[] {
  try {
    return JSON.parse(raw) as string[];
  } catch (e) {
    throw new Error(`load: ${e}`);
  }
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
}
