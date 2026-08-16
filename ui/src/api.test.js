import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { getJSON, postJSON } from "./api.js";

describe("api client", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("getJSON returns parsed body", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: true, json: async () => ({ ok: true }) }));
    await expect(getJSON("/api/health")).resolves.toEqual({ ok: true });
  });

  it("getJSON throws on non-ok", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: false, status: 500 }));
    await expect(getJSON("/api/health")).rejects.toThrow("HTTP 500");
  });

  it("postJSON sends JSON body", async () => {
    const mock = vi.fn().mockResolvedValue({ ok: true, json: async () => ({ run_id: "abc" }) });
    vi.stubGlobal("fetch", mock);
    await postJSON("/api/run", { tool_ids: ["nodejs"] });
    expect(mock).toHaveBeenCalledWith("/api/run", expect.objectContaining({ method: "POST" }));
  });
});
