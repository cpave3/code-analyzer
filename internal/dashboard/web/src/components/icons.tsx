// Inline stroke icons (no icon font, no network).
const base = { width: 16, height: 16, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: 2, strokeLinecap: "round" as const, strokeLinejoin: "round" as const, "aria-hidden": true };

export const Logo = () => (
  <svg {...base} width={22} height={22} stroke="var(--accent)">
    <path d="M4 19V9M10 19V5M16 19v-7M22 19H2" />
  </svg>
);
export const Moon = () => (
  <svg {...base}>
    <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z" />
  </svg>
);
export const Sun = () => (
  <svg {...base}>
    <circle cx="12" cy="12" r="4" />
    <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
  </svg>
);
export const Info = () => (
  <svg {...base}>
    <circle cx="12" cy="12" r="9" />
    <path d="M12 11v6M12 7.5v.5" />
  </svg>
);
export const Alert = () => (
  <svg {...base} width={18} height={18} stroke="var(--bad)">
    <circle cx="12" cy="12" r="9" />
    <path d="M12 7v6M12 16.5v.5" />
  </svg>
);
export const Check = () => (
  <svg {...base} strokeWidth={2.5}>
    <path d="M5 12l5 5 9-10" />
  </svg>
);
export const Prev = () => (
  <svg {...base}>
    <path d="M15 18l-6-6 6-6" />
  </svg>
);
export const Next = () => (
  <svg {...base}>
    <path d="M9 18l6-6-6-6" />
  </svg>
);
export const Play = () => (
  <svg {...base}>
    <path d="M7 5l12 7-12 7z" />
  </svg>
);
export const Pause = () => (
  <svg {...base}>
    <path d="M8 5v14M16 5v14" />
  </svg>
);
export const Close = () => (
  <svg {...base}>
    <path d="M6 6l12 12M18 6L6 18" />
  </svg>
);
export const Chevron = () => (
  <svg {...base} width={14} height={14}>
    <path d="M6 9l6 6 6-6" />
  </svg>
);
export const Empty = () => (
  <svg {...base} width={40} height={40} strokeWidth={1.6} stroke="var(--dim)">
    <path d="M3 7h18M3 12h18M3 17h10" />
  </svg>
);
