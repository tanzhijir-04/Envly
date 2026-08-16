export function summarize(events) {
  const result = { success: 0, failed: 0, skipped: 0 };
  for (const event of events) {
    if (result[event.type] !== undefined) {
      result[event.type] += 1;
    }
  }
  return result;
}
