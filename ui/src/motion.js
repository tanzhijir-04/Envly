// 轻量动效层：遵循 liquid-glass 规范（非线性缓动、可中断、尊重 reduced-motion）。

export function prefersReducedMotion(media) {
  const matcher = media || (typeof window !== "undefined" && window.matchMedia ? window.matchMedia.bind(window) : null);
  if (!matcher) return false;
  return matcher("(prefers-reduced-motion: reduce)").matches;
}

export function setupSegmentedMotion(seg) {
  if (!window.gsap) return;
  let thumb = seg.querySelector(".seg-thumb");
  if (!thumb) {
    thumb = document.createElement("span");
    thumb.className = "seg-thumb";
    seg.appendChild(thumb);
  }
  const active = seg.querySelector("button.on");
  const move = (btn, animate) => {
    const pos = { left: btn.offsetLeft, width: btn.offsetWidth };
    if (animate) {
      window.gsap.to(thumb, { left: pos.left, width: pos.width, duration: 0.3, ease: "power3.out", overwrite: "auto" });
    } else {
      window.gsap.set(thumb, { left: pos.left, width: pos.width });
    }
  };
  if (active) move(active, false);
  seg.querySelectorAll("button").forEach((btn) => {
    btn.addEventListener("click", () => move(btn, true));
  });
}

export function initMotion() {
  if (!window.gsap || prefersReducedMotion()) return;
  document.querySelectorAll(".seg").forEach(setupSegmentedMotion);
}
