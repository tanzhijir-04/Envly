import { describe, expect, it } from "vitest";
import { prefersReducedMotion } from "./motion.js";

describe("prefersReducedMotion", () => {
  it("returns true when media query matches", () => {
    const fake = () => ({ matches: true });
    expect(prefersReducedMotion(fake)).toBe(true);
  });

  it("returns false when media query does not match", () => {
    const fake = () => ({ matches: false });
    expect(prefersReducedMotion(fake)).toBe(false);
  });

  it("returns false when no matcher is available", () => {
    expect(prefersReducedMotion(null)).toBe(false);
  });
});
