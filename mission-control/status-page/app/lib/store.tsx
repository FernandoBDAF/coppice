"use client";

// Cockpit-wide shared state: the auth token, the single active action modal,
// the toast stack, and a "runs changed" tick so the history panel can refresh
// after an action launches or completes. Kept in one provider so any component
// can fire an action or a toast without prop drilling.

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { getToken as readToken, setToken as persistToken } from "./api";
import type { ActionRecord, ActionRequest, SystemDef } from "./types";
import { ActionModal } from "../components/ActionModal";
import { ToastStack, type Toast } from "../components/ToastStack";

// The modal either launches a fresh action (POST) or reopens a known record.
export type ActionModalMode =
  | { kind: "launch"; request: ActionRequest; system: SystemDef; confirm: boolean }
  | { kind: "open"; record: ActionRecord };

interface CockpitCtx {
  token: string;
  saveToken: (t: string) => void;
  addToast: (t: Omit<Toast, "id">) => void;
  openAction: (mode: ActionModalMode) => void;
  openRunning: (system: string, target: string) => Promise<void>;
  bumpRuns: () => void;
  runsVersion: number;
}

const Ctx = createContext<CockpitCtx | null>(null);

export function useCockpit(): CockpitCtx {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error("useCockpit used outside CockpitProvider");
  return ctx;
}

let nextToastId = 1;

export function CockpitProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string>(() => readToken());
  const [toasts, setToasts] = useState<Toast[]>([]);
  const [mode, setMode] = useState<ActionModalMode | null>(null);
  const [runsVersion, setRunsVersion] = useState(0);
  const fetchRunsRef = useRef<
    ((s: string, t: string) => Promise<ActionRecord | null>) | null
  >(null);

  const saveToken = useCallback((t: string) => {
    persistToken(t);
    setTokenState(t);
  }, []);

  const dismissToast = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const addToast = useCallback((t: Omit<Toast, "id">) => {
    const id = nextToastId++;
    setToasts((prev) => [...prev, { ...t, id }]);
    if (!t.action) {
      window.setTimeout(() => dismissToast(id), 6000);
    }
  }, [dismissToast]);

  const openAction = useCallback((m: ActionModalMode) => setMode(m), []);
  const bumpRuns = useCallback(() => setRunsVersion((v) => v + 1), []);

  // Best-effort: find the running action for a (system,target) and reopen it.
  // 409 responses do not carry the running id, so we look it up in the runs log.
  const openRunning = useCallback(
    async (system: string, target: string) => {
      const finder = fetchRunsRef.current;
      if (!finder) return;
      try {
        const rec = await finder(system, target);
        if (rec) setMode({ kind: "open", record: rec });
        else
          addToast({
            tone: "warn",
            message: "Couldn't locate the running action in the runs log.",
          });
      } catch {
        addToast({ tone: "warn", message: "Failed to load the runs log." });
      }
    },
    [addToast]
  );

  // Registered lazily to avoid importing the runs fetch here directly.
  fetchRunsRef.current = async (system, target) => {
    const { getJSON } = await import("./api");
    const runs = await getJSON<ActionRecord[]>("/api/runs?limit=20");
    return (
      runs.find(
        (r) =>
          r.request.system === system &&
          r.request.target === target &&
          (r.state === "running" || r.state === "pending")
      ) ?? null
    );
  };

  const value = useMemo<CockpitCtx>(
    () => ({
      token,
      saveToken,
      addToast,
      openAction,
      openRunning,
      bumpRuns,
      runsVersion,
    }),
    [token, saveToken, addToast, openAction, openRunning, bumpRuns, runsVersion]
  );

  return (
    <Ctx.Provider value={value}>
      {children}
      {mode && (
        <ActionModal
          mode={mode}
          onClose={() => setMode(null)}
          onOpenRecord={(record) => setMode({ kind: "open", record })}
        />
      )}
      <ToastStack toasts={toasts} onDismiss={dismissToast} />
    </Ctx.Provider>
  );
}
