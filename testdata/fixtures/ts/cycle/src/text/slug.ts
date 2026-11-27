import type { Post } from "../blog/post";

export function slugify(input: string, maxLength = 40): string {
  const lowered = input.trim().toLowerCase();
  const replaced = lowered.replace(/[^a-z0-9]+/g, "-");
  const trimmed = replaced.replace(/^-+|-+$/g, "");
  if (trimmed.length > maxLength) {
    return trimmed.slice(0, maxLength).replace(/-+$/g, "");
  }
  return trimmed;
}

export function postSlug(post: Post): string {
  return slugify(post.title);
}
