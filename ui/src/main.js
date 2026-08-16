import { getJSON, postJSON, subscribeEvents } from "./api.js";
import { t } from "./i18n.js";
import { applyTemplate, toggle } from "./plan.js";
import { initMotion, prefersReducedMotion } from "./motion.js";

const app = document.getElementById("app");
const state = {
  lang: "zh",
  region: "auto",
  tools: [],
  groups: [],
  templates: [],
  selected: [],
  planItems: [],
  report: null,
  activeTemplate: null,
};

function el(tag, className, text) {
  const node = document.createElement(tag);
  if (className) node.className = className;
  if (text != null) node.textContent = text;
  return node;
}

function groupItems() {
  const selected = new Set(state.selected);
  const statusMap = new Map((state.planItems || []).map((item) => [item.tool_id, item]));
  return state.groups.map((group) => ({
    ...group,
    items: state.tools
      .filter((tool) => tool.group_id === group.id)
      .map((tool) => {
        const plan = statusMap.get(tool.id);
        return { ...tool, checked: selected.has(tool.id), status: plan ? plan.status : "pending", version: plan ? plan.version : "" };
      }),
  }));
}

function render() {
  const scrollY = window.scrollY;
  app.innerHTML = "";
  app.append(renderHero(), renderTemplates(), renderRegion(), renderChecklist(), renderFooter(), renderLog(), renderSettings(), renderReport());
  window.scrollTo(0, scrollY);
  if (state.running) showCancelButton();
  initMotion();
}

function refreshTemplates() {
  document.querySelectorAll(".template-card").forEach((card) => {
    card.classList.toggle("on", card.dataset.id === state.activeTemplate);
  });
}

function refreshRows() {
  const selected = new Set(state.selected);
  document.querySelectorAll(".row").forEach((row) => {
    row.classList.toggle("on", selected.has(row.dataset.id));
  });
}

function refreshCounts() {
  const selected = new Set(state.selected);
  document.querySelectorAll(".group").forEach((section) => {
    const rows = section.querySelectorAll(".row");
    const checked = [...rows].filter((r) => selected.has(r.dataset.id)).length;
    const cnt = section.querySelector(".cnt");
    if (cnt) cnt.textContent = `${t(state.lang, "summary.selected", { count: checked })} / ${rows.length}`;
  });
  const summary = document.querySelector(".summary");
  if (summary) summary.textContent = `${t(state.lang, "summary.selected", { count: state.selected.length })} · ${t(state.lang, "summary.auto.skip")}`;
}

function renderHero() {
  const hero = el("div", "hero");
  hero.append(el("h1", "", t(state.lang, "welcome.title")), el("p", "subtitle", t(state.lang, "welcome.subtitle")));
  return hero;
}

function renderTemplates() {
  const wrap = el("div");
  wrap.append(el("div", "section-label", t(state.lang, "templates.title")));
  const grid = el("div", "templates");
  for (const tmpl of state.templates) {
    const card = el("div", "template-card" + (state.activeTemplate === tmpl.id ? " on" : ""));
    card.dataset.id = tmpl.id;
    card.append(
      el("b", "", tmpl[`name_${state.lang}`]),
      el("p", "", tmpl[`desc_${state.lang}`]),
      el("span", "count", `${tmpl.count} ${t(state.lang, "template.unit")}`)
    );
    card.addEventListener("click", () => {
      state.activeTemplate = tmpl.id;
      state.selected = applyTemplate(state.selected, tmpl.tool_ids);
      refreshTemplates();
      refreshRows();
      refreshCounts();
    });
    grid.append(card);
  }
  wrap.append(grid);
  return wrap;
}

function renderRegion() {
  const wrap = el("div", "region");
  wrap.append(el("span", "", t(state.lang, "region.label")));
  const seg = el("div", "seg");
  for (const value of ["auto", "cn", "global"]) {
    const btn = el("button", value === state.region ? "on" : "", t(state.lang, `region.${value}`));
    btn.dataset.value = value;
    btn.addEventListener("click", () => {
      state.region = value;
      document.querySelectorAll(".seg button").forEach((b) => b.classList.toggle("on", b.dataset.value === value));
      postJSON("/api/settings", { language: state.lang, region: state.region }).catch((err) => {
        showError(t(state.lang, "error.network", { message: err.message }));
      });
    });
    seg.append(btn);
  }
  wrap.append(seg, el("span", "region-hint", t(state.lang, "region.hint")));
  return wrap;
}

function renderChecklist() {
  const wrap = el("div");
  const labelRow = el("div", "section-label-row");
  labelRow.append(el("div", "section-label", t(state.lang, "checklist.title")));
  const actions = el("div", "label-actions");
  const selectAll = el("button", "mini-btn", t(state.lang, "checklist.select_all"));
  selectAll.addEventListener("click", () => {
    state.selected = state.tools.map((tool) => tool.id);
    refreshRows();
    refreshCounts();
  });
  const clear = el("button", "mini-btn", t(state.lang, "checklist.clear"));
  clear.addEventListener("click", () => {
    state.selected = [];
    refreshRows();
    refreshCounts();
  });
  actions.append(selectAll, clear);
  labelRow.append(actions);
  wrap.append(labelRow);
  const panel = el("div", "panel");
  for (const group of groupItems()) {
    const section = el("div", "group");
    const head = el("div", "group-head");
    head.append(
      el("b", "", group[`name_${state.lang}`]),
      el("span", "cnt", `${t(state.lang, "summary.selected", { count: group.items.filter((i) => i.checked).length })} / ${group.items.length}`)
    );
    section.append(head);
    for (const item of group.items) {
      const row = el("div", item.checked ? "row on" : "row");
      row.dataset.id = item.id;
      row.append(el("span", "box", "✓"), el("span", "name", item[`name_${state.lang}`]), el("span", "ver", item.version || ""));
      if (item.status === "installed") {
        row.append(el("span", "badge", t(state.lang, "tool.installed", { version: item.version || "" })));
      }
      row.append(el("span", "method", item.method));
      row.addEventListener("click", () => {
        state.selected = toggle(state.selected, item.id);
        row.classList.toggle("on", state.selected.includes(item.id));
        refreshCounts();
      });
      section.append(row);
    }
    panel.append(section);
  }
  wrap.append(panel);
  return wrap;
}

function renderFooter() {
  const footer = el("div", "footer");
  footer.append(
    el("span", "summary", `${t(state.lang, "summary.selected", { count: state.selected.length })} · ${t(state.lang, "summary.auto.skip")}`)
  );
  const cta = el("button", "cta", t(state.lang, "run.button"));
  cta.addEventListener("click", runFlow);
  footer.append(cta);
  return footer;
}

function renderLog() {
  const wrap = el("div");
  wrap.append(el("div", "section-label", t(state.lang, "log.title")));
  wrap.append(el("div", "log", ""));
  return wrap;
}

function renderSettings() {
  const wrap = el("div");
  wrap.append(el("div", "section-label", t(state.lang, "settings.title")));
  const panel = el("div", "panel");
  const row = el("div", "row");
  const btn = el("button", "cta", t(state.lang, "settings.restore"));
  btn.addEventListener("click", async () => {
    try {
      await postJSON("/api/settings/restore-env", {});
      showError(t(state.lang, "settings.restore.done"));
      setTimeout(clearError, 3000);
      await loadReport();
      render();
    } catch (err) {
      showError(t(state.lang, "settings.restore.fail", { message: err.message }));
    }
  });
  row.append(btn);
  panel.append(row);
  wrap.append(panel);
  return wrap;
}

function renderReport() {
  const wrap = el("div");
  wrap.append(el("div", "section-label", t(state.lang, "report.title")));
  const panel = el("div", "panel");
  const data = state.report || {};
  const head = el("div", "group-head");
  head.append(el("b", "", t(state.lang, "report.status", { status: data.status || "—" })));
  panel.append(head);
  const records = (data.records || []).map((r) => `${r.name} ${r.version || ""} · ${r.method}`).join("\n") || t(state.lang, "report.empty");
  const envOps = (data.env_ops || []).map((op) => `${op.key}: ${op.before || "—"} → ${op.after || "—"}`).join("\n") || t(state.lang, "report.empty");
  const recordsSection = el("div", "group");
  recordsSection.append(el("pre", "report-pre", `${t(state.lang, "report.records")}\n${records}`));
  const envSection = el("div", "group");
  envSection.append(el("pre", "report-pre", `${t(state.lang, "report.envops")}\n${envOps}`));
  panel.append(recordsSection, envSection);
  wrap.append(panel);
  return wrap;
}

async function loadReport() {
  try {
    state.report = await getJSON("/api/report");
  } catch {
    state.report = null;
  }
}

function appendLog(line) {
  const log = document.querySelector(".log");
  if (!log) return;
  const time = new Date().toLocaleTimeString();
  const row = el("div", "", `[${time}] ${line}`);
  log.append(row);
  if (window.gsap && !prefersReducedMotion()) {
    window.gsap.fromTo(row, { opacity: 0, y: 4 }, { opacity: 1, y: 0, duration: 0.25, ease: "power3.out" });
  }
  log.scrollTop = log.scrollHeight;
}

async function runFlow() {
  const cta = document.querySelector(".cta");
  cta.disabled = true;
  state.running = true;
  clearError();
  showCancelButton();
  try {
    const plan = await postJSON("/api/plan", { tool_ids: state.selected });
    const run = await postJSON("/api/run", { tool_ids: plan.items.map((item) => item.tool_id) });
    const unsubscribe = subscribeEvents(run.run_id, async (event) => {
      const text = t(state.lang, event.message_key, event.params || {});
      const mark = event.type === "success" ? "✓" : event.type === "failed" ? "✗" : event.type === "skipped" ? "⏭" : "→";
      if (event.type === "run_done") {
        appendLog(`${mark} ${text} · ${event.status}`);
        unsubscribe();
        state.running = false;
        hideCancelButton();
        await loadReport();
        cta.disabled = false;
        return;
      }
      appendLog(`${mark} ${text}`);
    });
  } catch (err) {
    showError(t(state.lang, "error.network", { message: err.message }));
    state.running = false;
    hideCancelButton();
    cta.disabled = false;
  }
}

function showCancelButton() {
  const footer = document.querySelector(".footer");
  if (!footer || footer.querySelector(".btn-cancel")) return;
  const btn = el("button", "btn-cancel", t(state.lang, "run.cancel"));
  btn.addEventListener("click", async () => {
    btn.disabled = true;
    try {
      await postJSON("/api/cancel", {});
    } catch (err) {
      showError(t(state.lang, "error.network", { message: err.message }));
      btn.disabled = false;
    }
  });
  footer.insertBefore(btn, footer.querySelector(".cta"));
}

function hideCancelButton() {
  const btn = document.querySelector(".btn-cancel");
  if (btn) btn.remove();
}

function showError(message) {
  let banner = document.querySelector(".banner");
  if (!banner) {
    banner = el("div", "banner");
    const footer = document.querySelector(".footer");
    if (footer) {
      footer.prepend(banner);
    } else {
      banner.style.position = "fixed";
      banner.style.top = "60px";
      banner.style.left = "22px";
      banner.style.right = "22px";
      banner.style.zIndex = "99";
      document.body.appendChild(banner);
    }
  }
  banner.textContent = message;
}

function clearError() {
  const banner = document.querySelector(".banner");
  if (banner) banner.remove();
}

async function loadPlanStatus() {
  if (state.selected.length === 0) {
    state.planItems = [];
    return;
  }
  const plan = await postJSON("/api/plan", { tool_ids: state.selected });
  state.planItems = plan.items;
}

document.querySelectorAll("[data-lang]").forEach((btn) => {
  btn.addEventListener("click", async () => {
    state.lang = btn.dataset.lang;
    document.documentElement.lang = state.lang;
    document.querySelectorAll("[data-lang]").forEach((b) => b.classList.toggle("on", b.dataset.lang === state.lang));
    await postJSON("/api/settings", { language: state.lang, region: state.region });
    render();
  });
});

async function init() {
  try {
    const catalog = await getJSON("/api/catalog");
    state.tools = catalog.tools;
    state.groups = catalog.groups;
    state.templates = await getJSON("/api/templates");
    state.planItems = [];
    const settings = await getJSON("/api/settings");
    state.lang = settings.language || "zh";
    state.region = settings.region || "auto";
    document.documentElement.lang = state.lang;
    document.querySelectorAll("[data-lang]").forEach((b) => b.classList.toggle("on", b.dataset.lang === state.lang));
    await loadPlanStatus();
    await loadReport();
    render();
  } catch (err) {
    console.error("Envly init failed:", err);
    showError(t(state.lang, "error.network", { message: err.message }));
  }
}

window.addEventListener("error", (event) => {
  console.error("Envly window error:", event.error || event.message);
  showError(String(event.error || event.message));
});

document.querySelectorAll(".win-btn").forEach((btn) => {
  btn.addEventListener("click", async () => {
    try {
      await postJSON("/api/window/action", { action: btn.dataset.action });
    } catch (err) {
      showError(t(state.lang, "error.network", { message: err.message }));
    }
  });
});

init();
