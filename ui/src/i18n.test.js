import { describe, expect, it } from "vitest";
import { t } from "./i18n.js";

describe("i18n", () => {
  it("translates zh and en keys", () => {
    expect(t("zh", "run.button")).toContain("一键配置");
    expect(t("en", "run.button")).toContain("Set up");
  });

  it("falls back to zh for unknown language", () => {
    expect(t("fr", "welcome.title")).toContain("Envly");
  });

  it("interpolates params", () => {
    expect(t("zh", "summary.selected", { count: 5 })).toContain("5");
  });
});
