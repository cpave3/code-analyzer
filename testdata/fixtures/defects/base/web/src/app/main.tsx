import React from "react";
import { title, parse, Counter } from "../lib/format";

export function App(props: { name: string }) {
  const c = new Counter();
  return <h1 data-n={c.inc() + parse("1")}>{title(props.name)}</h1>;
}
