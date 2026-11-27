export function greet(name: string): string {
  const n = name.trim();
  return n ? `hello ${n}` : "hello";
}
