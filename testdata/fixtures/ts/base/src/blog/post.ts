import { slugify } from "../text/slug";

export interface Post {
  title: string;
  body: string;
}

export function postPath(post: Post): string {
  return "/posts/" + slugify(post.title);
}
