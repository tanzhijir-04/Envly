const messages = {
  zh: {
    "nav.github": "GitHub",
    "hero.eyebrow": "v1.0.0 · Windows · macOS · Linux",
    "hero.title": "一键配置开发环境",
    "hero.lede": "选模板、勾清单、一键执行——实时日志与执行报告都在一个 Liquid Glass 单页向导里完成，海外网络自动跳过镜像源。",
    "hero.download": "下载 Envly",
    "hero.quickstart": "快速开始",
    "features.title": "功能特性",
    "features.wizard.title": "单页向导",
    "features.wizard.desc": "模板 → 勾选 → 一键配置 → 实时日志，全流程一个页面。",
    "features.install.title": "真实安装",
    "features.install.desc": "winget / npm / pip / rustup / download，已安装自动跳过。",
    "features.env.title": "环境配置",
    "features.env.desc": "npm/pip 镜像、代理、PowerShell Profile、PATH，可审计、可一键还原。",
    "features.region.title": "区域自适应",
    "features.region.desc": "自动检测网络区域，海外不强制切换镜像源。",
    "features.report.title": "执行报告",
    "features.report.desc": "安装记录与环境变更一目了然。",
    "features.desktop.title": "桌面体验",
    "features.desktop.desc": "无边框窗口、自绘控件、圆角与 Liquid Glass 动效。",
    "download.title": "下载",
    "download.win.desc": "MSI 安装包 · 免安装 exe · 引擎二进制",
    "download.mac.desc": "Envly.dmg + envly-engine-darwin",
    "download.linux.desc": "envly.deb · envly.AppImage + envly-engine-linux",
    "download.note": "macOS / Linux 为「Pake 壳 + 独立引擎」双进程结构，先运行引擎再启动应用。",
    "stack.title": "技术栈",
    "cta.title": "开始使用 Envly",
    "cta.desc": "下载对应平台的安装包，或从源码跑起来。",
    "cta.button": "前往 Releases",
    "footer.github": "GitHub",
    "footer.readme": "README",
    "footer.releases": "Releases",
  },
  en: {
    "nav.github": "GitHub",
    "hero.eyebrow": "v1.0.0 · Windows · macOS · Linux",
    "hero.title": "Set up your dev environment in one click",
    "hero.lede": "Pick a template, tweak the checklist, run — live logs and reports in a single Liquid Glass wizard. Mirrors are skipped automatically outside mainland China.",
    "hero.download": "Download Envly",
    "hero.quickstart": "Quick start",
    "features.title": "Features",
    "features.wizard.title": "Single-page wizard",
    "features.wizard.desc": "Template → checklist → one-click run → live log, all on one page.",
    "features.install.title": "Real installs",
    "features.install.desc": "winget / npm / pip / rustup / download; installed items are skipped.",
    "features.env.title": "Environment config",
    "features.env.desc": "npm/pip mirrors, proxy, PowerShell profile, PATH — auditable and restorable.",
    "features.region.title": "Region-aware",
    "features.region.desc": "Auto-detects your network region; mirrors are skipped overseas.",
    "features.report.title": "Execution report",
    "features.report.desc": "Install records and environment changes at a glance.",
    "features.desktop.title": "Desktop experience",
    "features.desktop.desc": "Frameless window, custom controls, rounded corners, Liquid Glass motion.",
    "download.title": "Download",
    "download.win.desc": "MSI installer · portable exe · engine binary",
    "download.mac.desc": "Envly.dmg + envly-engine-darwin",
    "download.linux.desc": "envly.deb · envly.AppImage + envly-engine-linux",
    "download.note": "macOS / Linux use a Pake shell + standalone engine. Start the engine before launching the app.",
    "stack.title": "Tech stack",
    "cta.title": "Get started with Envly",
    "cta.desc": "Download the installer for your platform, or run from source.",
    "cta.button": "Go to Releases",
    "footer.github": "GitHub",
    "footer.readme": "README",
    "footer.releases": "Releases",
  },
};

function applyLang(lang) {
  document.documentElement.lang = lang;
  const dict = messages[lang];
  document.querySelectorAll("[data-i18n]").forEach((el) => {
    const key = el.dataset.i18n;
    if (dict[key]) el.textContent = dict[key];
  });
  document.title = lang === "zh" ? "Envly — 一键开发环境配置" : "Envly — One-click dev environment setup";
}

const seg = document.getElementById("lang-seg");
const saved = localStorage.getItem("envly-lang") || "zh";
applyLang(saved);

seg.addEventListener("click", (e) => {
  const btn = e.target.closest("button");
  if (!btn) return;
  seg.querySelectorAll("button").forEach((b) => b.classList.remove("active"));
  btn.classList.add("active");
  const lang = btn.dataset.lang;
  localStorage.setItem("envly-lang", lang);
  applyLang(lang);
});

const nav = document.getElementById("nav");
addEventListener("scroll", () => {
  nav.classList.toggle("scrolled", scrollY > 8);
}, { passive: true });
