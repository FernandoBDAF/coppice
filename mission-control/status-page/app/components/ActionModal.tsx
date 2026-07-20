"use client";

// The action modal is the teaching surface and the control plane in one:
// it shows the request, the RESOLVED make command, live streamed output, and
// the terminal state. It drives three entry points — launch (POST), a
// destructive-confirm gate, and reopening a known record (streaming only if
// that record is still running).

import { useCallback, useEffect, useRef, useState } from "react";
import { ApiError, getJSON, postJSON, streamUrl } from "../lib/api";
import { actionStateClass, previewCommand, reportSummary } from "../lib/format";
import type {
  ActionRecord,
  ActionState,
  ExperimentReport,
} from "../lib/types";
import { useCockpit, type ActionModalMode } from "../lib/store";

const MAX_LINES = 2000;

type Phase = "confirm" | "posting" | "streaming" | "done" | "error";

export function ActionModal({
  mode,
  onClose,
  onOpenRecord,
}: {
  mode: ActionModalMode;
  onClose: () => void;
  onOpenRecord: (record: ActionRecord) => void;
}) {
  const { addToast, openRunning, bumpRuns } = useCockpit();

  const initialPhase: Phase =
    mode.kind === "launch"
      ? mode.confirm
        ? "confirm"
        : "posting"
      : mode.record.state === "running" || mode.record.state === "pending"
      ? "streaming"
      : "done";

  const [phase, setPhase] = useState<Phase>(initialPhase);
  const [actionId, setActionId] = useState<string | null>(
    mode.kind === "open" ? mode.record.id : null
  );
  const [command, setCommand] = useState<string>(
    mode.kind === "open"
      ? mode.record.command
      : previewCommand(
          mode.system,
          mode.request.target,
          mode.request.verb,
          mode.request.params
        ) ?? ""
  );
  const [resolved, setResolved] = useState<boolean>(mode.kind === "open");
  const [lines, setLines] = useState<string[]>([]);
  const [finalState, setFinalState] = useState<ActionState | null>(
    mode.kind === "open" &&
      mode.record.state !== "running" &&
      mode.record.state !== "pending"
      ? mode.record.state
      : null
  );
  const [exitCode, setExitCode] = useState<number | undefined>(
    mode.kind === "open" ? mode.record.exit_code : undefined
  );
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [streamError, setStreamError] = useState<string | null>(null);
  // Scored-experiment report: present up-front when reopening a completed
  // record; for a fresh launch it is fetched once the run ends (the SSE end
  // event carries only state + exit_code).
  const [report, setReport] = useState<ExperimentReport | null>(
    mode.kind === "open" ? mode.record.report ?? null : null
  );

  const esRef = useRef<EventSource | null>(null);
  const outRef = useRef<HTMLPreElement | null>(null);
  const request = mode.kind === "launch" ? mode.request : mode.record.request;

  const appendLine = useCallback((line: string) => {
    setLines((prev) => {
      const next = prev.length >= MAX_LINES ? prev.slice(1) : prev.slice();
      next.push(line);
      return next;
    });
  }, []);

  // Best-effort: pull the just-completed record to surface its scored-run
  // report (assertion breakdown). The engine attaches the report before the
  // SSE end event fires, so a direct GET here is race-free. Absence is
  // normal — non-experiment verbs and older runs carry no report — so
  // failures stay silent.
  const fetchReport = useCallback(async (id: string) => {
    try {
      const rec = await getJSON<ActionRecord>(
        `/api/actions/${encodeURIComponent(id)}`
      );
      if (rec.report) setReport(rec.report);
    } catch {
      /* report is optional; leave it unset */
    }
  }, []);

  const startStream = useCallback(
    (id: string) => {
      const es = new EventSource(streamUrl(id));
      esRef.current = es;
      es.addEventListener("line", (e) => appendLine((e as MessageEvent).data));
      es.addEventListener("end", (e) => {
        try {
          const payload = JSON.parse((e as MessageEvent).data) as {
            state: ActionState;
            exit_code?: number;
          };
          setFinalState(payload.state);
          setExitCode(payload.exit_code);
        } catch {
          setFinalState("succeeded");
        }
        setPhase("done");
        es.close();
        esRef.current = null;
        bumpRuns();
        void fetchReport(id);
      });
      es.onerror = () => {
        // EventSource auto-reconnects; close to avoid a storm and surface a
        // non-fatal note. If the action already ended this is harmless.
        if (esRef.current) {
          setStreamError("stream disconnected");
          es.close();
          esRef.current = null;
          setPhase((p) => (p === "streaming" ? "done" : p));
        }
      };
    },
    [appendLine, bumpRuns, fetchReport]
  );

  const doPost = useCallback(async () => {
    if (mode.kind !== "launch") return;
    setPhase("posting");
    setErrorMsg(null);
    try {
      const res = await postJSON<{ id: string; command: string }>(
        "/api/actions",
        mode.request
      );
      setActionId(res.id);
      setCommand(res.command);
      setResolved(true);
      setPhase("streaming");
      bumpRuns();
      startStream(res.id);
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        addToast({
          tone: "warn",
          message: `An action is already running for ${mode.request.system}/${mode.request.target}.`,
          action: {
            label: "View running",
            onClick: () =>
              openRunning(mode.request.system, mode.request.target),
          },
        });
        onClose();
        return;
      }
      setErrorMsg(err instanceof Error ? err.message : String(err));
      setPhase("error");
    }
  }, [mode, addToast, openRunning, onClose, bumpRuns, startStream]);

  // Kick off on mount: auto-post (non-destructive launch) or stream (reopen of
  // a still-running record). Destructive launches wait on the confirm gate.
  useEffect(() => {
    if (mode.kind === "launch" && !mode.confirm) {
      void doPost();
    } else if (mode.kind === "open" && initialPhase === "streaming") {
      startStream(mode.record.id);
    }
    return () => {
      esRef.current?.close();
      esRef.current = null;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Auto-scroll the output to the newest line.
  useEffect(() => {
    if (outRef.current) outRef.current.scrollTop = outRef.current.scrollHeight;
  }, [lines]);

  // Escape closes.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  const title = `${request.verb} · ${request.system} · ${request.target}`;
  const destructive = request.verb === "down";

  const paramEntries = Object.entries(request.params ?? {}).filter(
    ([, v]) => v !== undefined && v !== ""
  );

  return (
    <div className="overlay" onMouseDown={onClose}>
      <div
        className="modal"
        role="dialog"
        aria-modal="true"
        aria-label={title}
        onMouseDown={(e) => e.stopPropagation()}
      >
        <div className="modal-head">
          <span className="modal-title">
            <span className={destructive ? "verb-destructive" : "verb"}>
              {request.verb}
            </span>{" "}
            {request.system}
            <span className="prompt-sep"> :: </span>
            {request.target}
          </span>
          <button className="modal-close" aria-label="Close" onClick={onClose}>
            ×
          </button>
        </div>

        <div className="modal-body">
          <div className="req-grid">
            <span className="req-key">system</span>
            <span>{request.system}</span>
            <span className="req-key">target</span>
            <span>{request.target}</span>
            <span className="req-key">verb</span>
            <span>{request.verb}</span>
            {paramEntries.map(([k, v]) => (
              <span key={k} style={{ display: "contents" }}>
                <span className="req-key">{k}</span>
                <span>{String(v)}</span>
              </span>
            ))}
          </div>

          <div className="cmd-line">
            <span className="cmd-arrow">this runs →</span>
            <code className={`cmd ${resolved ? "resolved" : "preview"}`}>
              {command || "(command resolved by controld)"}
            </code>
            {!resolved && command && (
              <span className="cmd-note">preview</span>
            )}
          </div>

          {report && (
            <div
              className={`report-panel ${report.passed ? "ok" : "bad"}`}
              role="group"
              aria-label="assertion results"
            >
              <div className="report-head">
                <span className={`badge ${report.passed ? "ok" : "bad"}`}>
                  {report.passed ? "PASS" : "FAIL"}
                </span>
                <span className="report-summary">{reportSummary(report)}</span>
              </div>
              {report.assertions.length > 0 && (
                <ul className="report-list">
                  {report.assertions.map((a, i) => (
                    <li
                      key={i}
                      className={`report-item ${a.passed ? "ok" : "bad"}`}
                    >
                      <span className="report-marker" aria-hidden="true">
                        {a.passed ? "✓" : "✗"}
                      </span>
                      <span className="report-name">{a.name}</span>
                      {a.detail && (
                        <pre className="report-detail">{a.detail}</pre>
                      )}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          )}

          {phase === "confirm" && (
            <div className="confirm-box">
              <p>
                <strong>Destructive action.</strong> This will tear down{" "}
                <code>{request.system}</code> on <code>{request.target}</code>.
              </p>
              <div className="confirm-actions">
                <button className="btn-danger" onClick={() => void doPost()}>
                  Confirm — run it
                </button>
                <button className="btn-ghost" onClick={onClose}>
                  Cancel
                </button>
              </div>
            </div>
          )}

          {phase === "posting" && (
            <div className="notice">submitting action to controld…</div>
          )}

          {phase === "error" && errorMsg && (
            <div className="error-box">
              <strong>action rejected</strong>
              <span>{errorMsg}</span>
            </div>
          )}

          {(phase === "streaming" || phase === "done") && (
            <>
              <div className="output-label">
                output
                {actionId && <span className="output-id"> · {actionId}</span>}
              </div>
              <pre className="output" ref={outRef}>
                {lines.length === 0
                  ? finalState
                    ? "(no retained output for this completed action)"
                    : "waiting for output…"
                  : lines.join("\n")}
              </pre>
              {streamError && (
                <div className="stream-note">{streamError}</div>
              )}
            </>
          )}
        </div>

        <div className="modal-foot">
          {finalState ? (
            <span className={`badge ${actionStateClass(finalState)}`}>
              {finalState.toUpperCase()}
              {exitCode !== undefined ? ` · exit ${exitCode}` : ""}
            </span>
          ) : phase === "streaming" ? (
            <span className="badge warn">RUNNING</span>
          ) : (
            <span className="prompt-meta">idle</span>
          )}
          <button className="btn-ghost" style={{ marginLeft: "auto" }} onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
