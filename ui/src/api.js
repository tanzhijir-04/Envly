export async function getJSON(path) {
  const res = await fetch(path);
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: ${path}`);
  }
  return res.json();
}

export async function postJSON(path, body) {
  const res = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: ${path}`);
  }
  return res.json();
}

export function subscribeEvents(runID, onEvent) {
  const es = new EventSource(`/api/events?run_id=${encodeURIComponent(runID)}`);
  es.onmessage = (event) => onEvent(JSON.parse(event.data));
  return () => es.close();
}
