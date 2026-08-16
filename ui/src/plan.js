export function toggle(selected, id) {
  return selected.includes(id) ? selected.filter((x) => x !== id) : [...selected, id];
}

export function applyTemplate(selected, templateToolIDs) {
  return [...new Set([...selected, ...templateToolIDs])];
}
