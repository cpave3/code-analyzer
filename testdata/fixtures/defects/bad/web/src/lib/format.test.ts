import { describe, expect, it } from "vitest";
import { parse, title } from "./format";

describe("format", () => {
  it("titles", () => {
    const t = title("a b");
    expect(t).toBe(t);
  });
  it("parses async", async () => {
    expect(Promise.resolve(parse("2"))).resolves.toBe(2);
  });
});
