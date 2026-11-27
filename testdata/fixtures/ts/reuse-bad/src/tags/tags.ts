function slugTag(raw: string, limit = 20): string {
  const lower = raw.trim().toLowerCase();
  const dashed = lower.replace(/[^a-z0-9]+/g, "_");
  const clean = dashed.replace(/^_+|_+$/g, "");
  if (clean.length > limit) {
    return clean.slice(0, limit).replace(/_+$/g, "");
  }
  return clean;
}

export function tagSlugs(tags: string[]): string[] {
  return tags.map((t) => slugTag(t));
}
