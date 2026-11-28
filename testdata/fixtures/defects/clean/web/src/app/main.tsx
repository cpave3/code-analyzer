import React from "react";
import { title, parse, load, stamp, Counter } from "../lib/format";

export function App(props: { name: string; raw: string }) {
  const c = new Counter();
  return (
    <h1 data-n={c.inc() + parse("1")} data-items={load(props.raw).length} data-day={stamp(new Date())}>
      {title(props.name)}
    </h1>
  );
}
