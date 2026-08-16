import { getJSON, postJSON, subscribeEvents } from "./api.js";
import { t } from "./i18n.js";
import { applyTemplate, toggle } from "./plan.js";

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
  app.innerHTML = "";
  const hero = el("div", "hero");
  hero.append(el("h1", "", t(state.lang, "welcome.title")), el("p", "subtitle", t(state.lang, "welcome.subtitle")));
  app.append(hero, renderTemplates(), renderRegion(), renderChecklist(), renderFooter(), renderLog(), renderSettings(), renderReport());
}

function renderTemplates() {
  const wrap = el("div");
  wrap.append(el("div", "section-label", t(state.lang, "templates.title")));
  const grid = el("div", "templates");
  for (const tmpl of state.templates) {
    const card = el("div", "template-card");
    card.append(
      el("b", "", tmpl[`name_${state.lang}`]),
      el("p", "", tmpl[`desc_${state.lang}`]),
      el("span", "count", `${t(state.lang, "summary.selected", { count: tmpl.count })}`)
    );
    card.addEventListener("click", () => {
      state.selected = applyTemplate(state.selected, tmpl.tool_ids);
      loadPlanStatus().then(render);
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
    btn.addEventListener("click", async () => {
      state.region = value;
      await postJSON("/api/settings", { language: state.lang, region: state.region });
      render();
    });
    seg.append(btn);
  }
  wrap.append(seg, el("span", "region-hint", t(state.lang, "region.hint")));
  return wrap;
}

function renderChecklist() {
  const wrap = el("div");
  wrap.append(el("div", "section-label", t(state.lang, "checklist.title")));
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
      row.append(el("span", "box", "✓"), el("span", "name", item[`name_${state.lang}`]), el("span", "method", item.method));
      if (item.status === "installed") {
        row.append(el("span", "badge", t(state.lang, "tool.installed", { version: item.version || "" })));
      }
      row.addEventListener("click", () => {
        state.selected = toggle(state.selected, item.id);
        loadPlanStatus().then(render);
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
  footer.append(el("span", "summary", t(state.lang, "summary.selected", { count: state.selected.length })));
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
  log.append(el("div", "", `[${time}] ${line}`));
  log.scrollTop = log.scrollHeight;
}

async function runFlow() {
  const cta = document.querySelector(".cta");
  cta.disabled = true;
  clearError();
  try {
    const plan = await postJSON("/api/plan", { tool_ids: state.selected });
    const run = await postJSON("/api/run", { tool_ids: plan.items.map((item) => item.tool_id) });
    const unsubscribe = subscribeEvents(run.run_id, async (event) => {
      const text = t(state.lang, event.message_key, event.params || {});
      const mark = event.type === "success" ? "✓" : event.type === "failed" ? "✗" : event.type === "skipped" ? "⏭" : "→";
      if (event.type === "run_done") {
        appendLog(`${mark} ${text} · ${event.status}`);
        unsubscribe();
        await loadReport();
        cta.disabled = false;
        return;
      }
      appendLog(`${mark} ${text}`);
    });
  } catch (err) {
    showError(t(state.lang, "error.network", { message: err.message }));
    cta.disabled = false;
  }
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

init();
