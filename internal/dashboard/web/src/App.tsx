import { lazy, Suspense, useState } from "react";
import { Route, Routes } from "react-router";
import { About } from "./components/About";
import { Header } from "./components/Header";
import { PlaybackProvider } from "./components/Playback";
import { ErrorBanner, Skeleton } from "./components/states";
import { useMeta, useMetrics } from "./data";
import { Overview } from "./views/Overview";

// Charts and the graph renderer load with the views that use them.
const Modules = lazy(() => import("./views/Modules"));
const Changes = lazy(() => import("./views/Changes"));
const GraphView = lazy(() => import("./views/Graph"));

export function App() {
  const meta = useMeta();
  const metrics = useMetrics();
  const [about, setAbout] = useState(false);
  if (meta.isPending) return <LoadingPage />;
  if (meta.isError)
    return (
      <main className="page">
        <ErrorBanner title="The dashboard could not load this repository's history" error={meta.error} onRetry={() => meta.refetch()} />
      </main>
    );
  return (
    <PlaybackProvider>
      <Header onAbout={() => setAbout(true)} />
      {metrics.isError && metrics.data && (
        <div style={{ maxWidth: 1400, margin: "12px auto 0", padding: "0 24px" }}>
          <ErrorBanner title="Can't reach codeeval serve" onRetry={() => metrics.refetch()}>
            <div style={{ fontSize: 13 }}>Showing the results loaded at {new Date(metrics.dataUpdatedAt).toLocaleTimeString()}. Start the server again, then retry.</div>
          </ErrorBanner>
        </div>
      )}
      <Suspense fallback={<LoadingPage bare />}>
        <Routes>
          <Route path="/" element={<Overview />} />
          <Route path="/modules/*" element={<Modules />} />
          <Route path="/changes" element={<Changes />} />
          <Route path="/graph" element={<GraphView />} />
          <Route path="*" element={<Overview />} />
        </Routes>
      </Suspense>
      {about && <About onClose={() => setAbout(false)} />}
    </PlaybackProvider>
  );
}

export function LoadingPage({ bare }: { bare?: boolean }) {
  return (
    <main className="page" aria-busy="true">
      {!bare && <Skeleton height={44} />}
      <div className="kpis">
        {[0, 1, 2, 3, 4].map((i) => (
          <Skeleton key={i} height={88} />
        ))}
      </div>
      <div className="row">
        <div className="main-col">
          <Skeleton height={420} />
        </div>
        <div className="side-col">
          <Skeleton height={200} />
          <Skeleton height={200} />
        </div>
      </div>
    </main>
  );
}
