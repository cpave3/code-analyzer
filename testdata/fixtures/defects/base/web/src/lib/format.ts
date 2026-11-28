import { slug } from "../shared/text";

export function title(s: string): string {
  const cleaned = slug(s);
  return cleaned.toUpperCase();
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
