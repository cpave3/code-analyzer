import { slugify } from "../text/slug";

export function tagSlugs(tags: string[]): string[] {
  return tags.map((t) => slugify(t, 20));
}
