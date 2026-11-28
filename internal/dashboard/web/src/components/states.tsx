import type { ReactNode } from "react";
import { HttpError } from "../core/api";
import { Alert, Empty } from "./icons";

export function Skeleton({ height, width = "100%" }: { height: number; width?: number | string }) {
  return <div className="skeleton" style={{ height, width }} aria-hidden="true" />;
}

export function errorText(e: unknown): string {
  if (e instanceof HttpError) return e.message;
  if (e instanceof TypeError) return "Can't reach codeeval serve. Start it again, then retry.";
  return e instanceof Error ? e.message : String(e);
}

export function ErrorBanner({ title, error, onRetry, children }: { title: string; error?: unknown; onRetry?: () => void; children?: ReactNode }) {
  return (
    <div className="banner" role="alert">
      <Alert />
      <div style={{ flex: 1, display: "flex", flexDirection: "column", gap: 4 }}>
        <div className="t">{title}</div>
        {error != null && <div style={{ fontSize: 13 }}>{errorText(error)}</div>}
        {children}
      </div>
      {onRetry && (
        <button className="btn small" onClick={onRetry}>
          Retry
        </button>
      )}
    </div>
  );
}

export function EmptyState({ title, children }: { title: string; children?: ReactNode }) {
  return (
    <div className="card empty">
      <div>
        <Empty />
        <h2>{title}</h2>
        {children}
      </div>
    </div>
  );
}
