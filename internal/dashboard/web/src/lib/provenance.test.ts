import type { Commit } from "../core/api";
import { keyColors, stepKey } from "./provenance";

const c = (prov: Commit["provenance"]): Commit => ({ index: 0, sha: "a", parent: "", subject: "", author: "", time: "", merge: false, provenance: prov });

describe("stepKey", () => {
  it("keys a step by the chosen fields", () => {
    expect(stepKey(c([{ commit: "a", model: ["opus"], harness: ["claude-code"], source: "trailer" }]), ["model", "harness"])).toBe("model=opus, harness=claude-code");
  });
  it("marks a merge of differently attributed commits as mixed, ignoring the merge commit itself", () => {
    const k = stepKey(
      c([
        { commit: "m", source: "none", merge: true },
        { commit: "a", model: ["opus"], source: "trailer" },
        { commit: "b", model: ["sonnet"], source: "trailer" },
      ]),
      ["model"],
    );
    expect(k).toBe("(mixed)");
  });
  it("reads custom fields from meta", () => {
    expect(stepKey(c([{ commit: "a", meta: { workflow: "tdd" }, source: "note" }]), ["workflow"])).toBe("workflow=tdd");
  });
});

describe("keyColors", () => {
  it("keeps greys for mixed and unattributed keys and colours the rest in order", () => {
    const m = keyColors(["model=opus", "(mixed)", "model=(none)", "model=sonnet", "model=opus"]);
    expect(m.get("model=opus")).not.toBe(m.get("model=sonnet"));
    expect(m.get("(mixed)")).toBe("#8b8f98");
    expect(m.get("model=(none)")).toBe("#a7abb3");
    expect(m.size).toBe(4);
  });
});
