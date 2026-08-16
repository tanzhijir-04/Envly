import { describe, expect, it } from "vitest";
import { summarize } from "./status.js";

describe("summarize", () => {
  it("counts success, failed, skipped", () => {
    const events = [
      { type: "success" },
      { type: "success" },
      { type: "failed" },
      { type: "skipped" },
    ];
    expect(summarize(events)).toEqual({ success: 2, failed: 1, skipped: 1 });
  });

  it("ignores run_done", () => {
    expect(summarize([{ type: "run_done" }])).toEqual({ success: 0, failed: 0, skipped: 0 });
  });

  it("counts run_done via params", () => {
    expect(summarize([{ type: "run_done", status: "success" }])).toEqual({ success: 0, failed: 0, skipped: 0 });
  });
});
