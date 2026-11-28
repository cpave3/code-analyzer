import { describe, expect, it } from "vitest";
import { parse, title } from "./format";

describe("format", () => {
  it("titles", () => {
    expect(title("a b")).toBe("A-B");
  });
  it("parses async", async () => {
    await expect(Promise.resolve(parse("2"))).resolves.toBe(2);
  });
});
