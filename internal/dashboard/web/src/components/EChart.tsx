// A thin React wrapper over Apache ECharts (tree-shaken core): it creates
// the chart once, applies option changes and resizes with its box.
import { BarChart, LineChart, TreemapChart } from "echarts/charts";
import { BrushComponent, GridComponent, MarkLineComponent, MarkPointComponent, TooltipComponent } from "echarts/components";
import * as echarts from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { useEffect, useRef } from "react";

echarts.use([BrushComponent, BarChart, LineChart, TreemapChart, GridComponent, TooltipComponent, MarkLineComponent, MarkPointComponent, CanvasRenderer]);

export type EChartsOption = echarts.EChartsCoreOption;
export type Handlers = Record<string, (p: unknown) => void>;

export function EChart({ option, height, onEvents, label, after }: { option: EChartsOption; height: number | string; onEvents?: Handlers; label: string; after?: (c: echarts.ECharts) => void }) {
  const el = useRef<HTMLDivElement>(null);
  const chart = useRef<echarts.ECharts | null>(null);
  const handlers = useRef<Handlers | undefined>(onEvents);
  handlers.current = onEvents;
  useEffect(() => {
    const c = echarts.init(el.current!, undefined, { renderer: "canvas" });
    chart.current = c;
    const ro = new ResizeObserver(() => c.resize());
    ro.observe(el.current!);
    for (const name of ["click", "dblclick", "mouseover", "brushEnd"]) c.on(name, (p: unknown) => handlers.current?.[name]?.(p));
    c.getZr().on("click", (e: { target?: unknown }) => {
      if (!e.target) handlers.current?.["blank"]?.(e);
    });
    return () => {
      ro.disconnect();
      c.dispose();
      chart.current = null;
    };
  }, []);
  useEffect(() => {
    if (!chart.current) return;
    chart.current.setOption(option, { notMerge: true });
    after?.(chart.current);
  }, [option, after]);
  return <div ref={el} role="img" aria-label={label} style={{ width: "100%", height }} />;
}
