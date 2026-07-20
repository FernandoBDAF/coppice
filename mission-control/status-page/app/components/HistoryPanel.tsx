"use client";

// Collapsible run history from /api/runs. A row reopens that action's modal
// (record, plus live stream if it is still running). Refetches when an action
// launches or completes (runsVersion) while the panel is open.

import { useCallback, useEffect, useState } from "react";
import { getJSON } from "../lib/api";
import { actionStateClass, fmtClock, fmtDuration } from "../lib/format";
import type { ActionRecord } from "../lib/types";
import { useCockpit } from "../lib/store";

export function HistoryPanel() {
  const { openAction, runsVersion } = useCockpit();
  const [open, setOpen] = useState(false);
  const [runs, setRuns] = useState<ActionRecord[]>([]);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    try {
      const data = await getJSON<ActionRecord[]>("/api/runs?limit=20");
      setRuns(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load runs");
    }
  }, []);

  useEffect(() => {
    if (open) void load();
  }, [open, load]);

  // Refetch when actions change, but only while the panel is showing.
  useEffect(() => {
    if (open && runsVersion > 0) void load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [runsVersion]);

  return (
    <section className="section history">
      <button
        className="history-toggle"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
      >
        <span className="section-label">{open ? "▾" : "▸"} run history</span>
      </button>

      {open && (
        <div className="history-body">
          <div className="history-controls">
            <button className="btn-ghost" onClick={() => void load()}>
              ↻ refresh
            </button>
          </div>
          {error ? (
            <div className="notice">
              <strong>runs unavailable</strong> — {error}
            </div>
          ) : runs.length === 0 ? (
            <div className="notice">no runs recorded yet</div>
          ) : (
            <div className="table-wrap">
              <table className="runs-table">
                <thead>
                  <tr>
                    <th>time</th>
                    <th>system/target/verb</th>
                    <th>command</th>
                    <th>state</th>
                    <th>dur</th>
                  </tr>
                </thead>
                <tbody>
                  {runs.map((r) => (
                    <tr
                      key={r.id}
                      className="run-row"
                      onClick={() => openAction({ kind: "open", record: r })}
                    >
                      <td>{fmtClock(r.started_at)}</td>
                      <td>
                        {r.request.system}/{r.request.target}/{r.request.verb}
                      </td>
                      <td className="run-cmd">
                        <code>{r.command}</code>
                      </td>
                      <td>
                        <span className={`badge ${actionStateClass(r.state)}`}>
                          {r.state.toUpperCase()}
                          {r.exit_code !== undefined ? ` ${r.exit_code}` : ""}
                        </span>
                      </td>
                      <td>{fmtDuration(r.started_at, r.ended_at)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}
    </section>
  );
}
