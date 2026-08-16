export const messages = {
  zh: {
    "app.name": "Envly",
    "welcome.title": "欢迎使用 Envly",
    "welcome.subtitle": "选一个模板，或自己勾选，然后一键配置你的开发环境。",
    "step.template": "1 选模板",
    "step.checklist": "2 确认清单",
    "step.run": "3 执行配置",
    "templates.title": "选择模板",
    "checklist.title": "确认清单 · 点击可增删",
    "region.label": "网络区域",
    "region.auto": "自动检测",
    "region.cn": "中国大陆",
    "region.global": "海外",
    "region.hint": "海外网络将跳过镜像源配置",
    "summary.selected": "已选 {count} 项",
    "run.button": "一键配置",
    "log.title": "执行日志（实时）",
    "tool.start": "正在配置 {tool}",
    "tool.progress": "{tool} 进行中 {percent}%",
    "tool.done": "{tool} 配置完成",
    "run.done": "执行结束",
  },
  en: {
    "app.name": "Envly",
    "welcome.title": "Welcome to Envly",
    "welcome.subtitle": "Pick a template or build your own checklist, then set up your environment in one click.",
    "step.template": "1 Template",
    "step.checklist": "2 Checklist",
    "step.run": "3 Run",
    "templates.title": "Choose a template",
    "checklist.title": "Checklist · click to toggle",
    "region.label": "Network region",
    "region.auto": "Auto detect",
    "region.cn": "Mainland China",
    "region.global": "Overseas",
    "region.hint": "Mirror sources are skipped for overseas networks",
    "summary.selected": "{count} selected",
    "run.button": "Set up",
    "log.title": "Live log",
    "tool.start": "Configuring {tool}",
    "tool.progress": "{tool} in progress {percent}%",
    "tool.done": "{tool} ready",
    "run.done": "Finished",
  },
};

export function t(lang, key, params = {}) {
  let text = (messages[lang] && messages[lang][key]) || messages.zh[key] || key;
  for (const [k, v] of Object.entries(params)) {
    text = text.replaceAll(`{${k}}`, String(v));
  }
  return text;
}
