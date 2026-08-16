import { describe, expect, it } from "vitest";
import { applyTemplate, toggle } from "./plan.js";

describe("plan helpers", () => {
  it("toggle adds and removes ids", () => {
    expect(toggle([], "nodejs")).toEqual(["nodejs"]);
    expect(toggle(["nodejs"], "nodejs")).toEqual([]);
  });

  it("applyTemplate merges and dedupes", () => {
    expect(applyTemplate(["nodejs"], ["git", "nodejs"])).toEqual(["nodejs", "git"]);
  });
});
