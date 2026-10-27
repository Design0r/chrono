import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useMemo, useRef, useState } from "react";
import { chrono } from "../api/chrono/client";
import { timestampQueries, userQueries } from "../api/queries";
import {
  durationFromTimestamps,
  formatCounter,
  secondsToCounter,
  yearOfRange,
} from "../lib/timestamp-utils";
import type { User } from "../types/auth";
import type { WorkTime } from "../types/response";
import { CreateTimestampModal } from "./timestamps/CreateTimestamp";
import { Timer } from "./timestamps/Timer";
import { TimestampTable } from "./timestamps/TimestampTable";
import { TimestampTableByWeek } from "./timestamps/TimestampTableByWeek";
import { useCalendarEventsByDay } from "../lib/useCalendarEventsByDay";
import { LoadingSpinner } from "./LoadingSpinner";
import { useToast } from "./Toast";

export function Timestamps({ user }: { user: User }) {
  const { addToast, addErrorToast } = useToast();
  const [runningTimer, setRunningTimer] = useState<number>(0);
  const queryClient = useQueryClient();

  const latestTimestampQ = useQuery(timestampQueries.latest());
  const timestampsQ = useQuery(timestampQueries.today());

  const timestamps = useMemo(
    () => (timestampsQ.isError ? [] : (timestampsQ.data ?? [])),
    [timestampsQ.data, timestampsQ.isError],
  );
  const eventsByDay = useCalendarEventsByDay(timestamps, user.id);

  // Single source of truth: der neueste Timestamp ohne end_time läuft.
  // Direkt aus den Query-Daten abgeleitet, damit ein Refetch den
  // angezeigten Zustand nie aus dem Tritt bringen kann.
  const currTimer = useMemo(() => {
    const latest = latestTimestampQ.data;
    if (!latest || latest.end_time !== null) return null;
    // Unlesbare start_time würde eine sinnlose Laufzeit anzeigen.
    if (Number.isNaN(Date.parse(latest.start_time))) return null;
    return latest;
  }, [latestTimestampQ.data]);

  const paused = currTimer === null;

  // Bei pausiertem Timer ungenutzt: die Anzeige steht dann fest auf 0.
  const startTime = currTimer ? Date.parse(currTimer.start_time) : 0;

  // Timer, die wir selbst gestartet haben, sollen keinen "Resuming"-Toast auslösen.
  const startedTimerId = useRef<number | null>(null);
  const announcedTimerId = useRef<number | null>(null);

  const startMut = useMutation({
    mutationFn: () => chrono.timestamps.start(),
    onError: (e) => addErrorToast(e),
    onSuccess: async (data) => {
      startedTimerId.current = data.id;
      queryClient.setQueryData(timestampQueries.latest().queryKey, data);
      await queryClient.invalidateQueries({ queryKey: timestampQueries.all() });
    },
  });

  const stopMut = useMutation({
    mutationFn: (id: number) => chrono.timestamps.stop(id),
    onError: (e) => addErrorToast(e),
    onSuccess: async (data) => {
      queryClient.setQueryData(timestampQueries.latest().queryKey, data);
      await queryClient.invalidateQueries({ queryKey: timestampQueries.all() });
    },
  });

  useEffect(() => {
    if (!currTimer) {
      announcedTimerId.current = null;
      return;
    }
    if (announcedTimerId.current === currTimer.id) return;
    announcedTimerId.current = currTimer.id;
  }, [currTimer, addToast]);

  useEffect(() => {
    if (timestampsQ.isError) addErrorToast(timestampsQ.error);
  }, [timestampsQ.isError, timestampsQ.error, addErrorToast]);

  // Bei pausiertem Timer zählt nur die abgeschlossene Zeit, unabhängig davon,
  // was der letzte Tick noch gemeldet hat.
  const totalTime = secondsToCounter(
    durationFromTimestamps(timestamps) + (paused ? 0 : runningTimer),
  );

  return (
    <div className="mx-auto flex lg:flex-row flex-col w-full gap-4 justify-center ">
      <div className="flex items-center justify-center flex-1 rounded-2xl bg-base-200/60 border border-base-300/80 p-6 lg:p-8">
        <div className="flex flex-col items-center justify-center gap-4">
          <Timer
            paused={paused}
            startUnix={startTime}
            onUpdate={(seconds: number) => setRunningTimer(paused ? 0 : seconds)}
          />
          <div className="flex mt-1 gap-3 justify-center items-center">
            <button
              disabled={!paused || startMut.isPending}
              className="btn btn-lg btn-success w-22 rounded-full shadow-md hover:shadow-lg transition-shadow disabled:opacity-40 disabled:shadow-none icon-filled"
              onClick={() => startMut.mutate()}
              title="Timer starten"
            >
              <span className="text-xl icon-filled scale-145">play_arrow</span>
            </button>
            <button
              disabled={paused || stopMut.isPending}
              className="btn btn-circle btn-lg btn-error shadow-md hover:shadow-lg transition-shadow disabled:opacity-40 disabled:shadow-none icon-outlined"
              onClick={() => currTimer && stopMut.mutate(currTimer.id)}
              title="Timer stoppen"
            >
              <span className="text-xl icon-filled scale-125">stop</span>
            </button>
          </div>
          <div className="flex whitespace-nowrap items-center gap-2 mt-0 text-base-content/70">
            <span className={`font-semibold ${paused ? "text-success" : "text-accent/80"}`}>
              {(() => {
                const f = formatCounter(totalTime);
                return (
                  <>
                    <span>{f.hours}</span>
                    <span
                      className={
                        paused ? "text-success/90 animate-pulse" : "text-accent/80 animate-pulse"
                      }
                    >
                      :
                    </span>
                    <span>{f.minutes} h</span>
                  </>
                );
              })()}
              <span className={`pl-2.5 ${paused ? "text-success/50" : "text-accent/60"}`}>
                Today
              </span>
            </span>
          </div>
        </div>
      </div>
      <div className="rounded-xl overflow-hidden w-full h-full">
        <TimestampTable timestamps={timestamps} user={user} eventsByDay={eventsByDay} />
      </div>
    </div>
  );
}

export function TeamTimestamps({
  startDate,
  endDate,
}: {
  startDate?: string;
  endDate?: string;
  user: User;
}) {
  const [createTimestampModal, setCreateTimestampModal] = useState(false);

  const workedQ = useQuery(timestampQueries.workedSecondsPerUser({ startDate, endDate }));
  const usersQ = useQuery(userQueries.list());
  const allWorktimesQ = useQuery(timestampQueries.workHoursForAllUsers(yearOfRange(startDate)));

  const queries = [workedQ, usersQ, allWorktimesQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending || firstError) return <></>;

  const worked = workedQ.data!;
  const worktimes = allWorktimesQ.data!;
  const members = usersQ.data!.filter((u) => u.id in worked);

  return (
    <>
      <hr />
      <div className="flex gap-4 items-center">
        <h2>Team Timestamps</h2>
        <button
          className="btn btn-soft btn-primary icon-outlined"
          onClick={() => setCreateTimestampModal(true)}
        >
          add
        </button>
        {createTimestampModal && (
          <CreateTimestampModal onClose={() => setCreateTimestampModal(false)} />
        )}
      </div>
      <div className="w-full">
        {members.map((user) => (
          <TeamMemberTimestamps
            key={user.id}
            user={user}
            range={{ startDate, endDate }}
            workedSeconds={worked[user.id]}
            worktime={worktimes[user.id]}
          />
        ))}
      </div>
    </>
  );
}

function TeamMemberTimestamps({
  user,
  range,
  workedSeconds,
  worktime,
}: {
  user: User;
  range: { startDate?: string; endDate?: string };
  workedSeconds: number;
  worktime: WorkTime;
}) {
  const [open, setOpen] = useState(false);
  const timestampsQ = useQuery({ ...timestampQueries.forUser(user.id, range), enabled: open });

  const fWorked = formatCounter(secondsToCounter(workedSeconds));
  const fExpected = formatCounter(secondsToCounter(worktime.expected * 3600));
  const overtime = (worktime.worked - worktime.expected) * 3600;
  const overtimeLabel = overtime < 0 ? "Missing Time" : "Overtime";
  const fOvertime = formatCounter(secondsToCounter(Math.abs(overtime)));

  return (
    <div className="my-4">
      <details
        className="collapse border-base-300 border collapse-arrow"
        onToggle={(e) => setOpen(e.currentTarget.open)}
      >
        <summary className="collapse-title bg-base-300/50 focus-within:bg-info/25 focus:text-white hover:bg-info/25 flex flex-wrap items-center gap-x-8 gap-y-1">
          <span className="font-semibold min-w-32">{user.username}</span>
          <span className="flex flex-wrap gap-x-6 gap-y-1 text-sm font-normal">
            <SummaryValue label="Worked" value={fWorked} highlight />
            <SummaryValue label="Expected" value={fExpected} />
            <SummaryValue label={overtimeLabel} value={fOvertime} />
          </span>
        </summary>
        {open && (
          <div className="collapse-content px-0 pt-6 pb-0 bg-black/20 flex flex-col text-sm">
            {timestampsQ.isPending ? (
              <div className="flex justify-center py-8">
                <LoadingSpinner size="md" />
              </div>
            ) : timestampsQ.isError ? (
              <p className="px-4 pb-6 text-error">{timestampsQ.error.message}</p>
            ) : (
              <TimestampTableByWeek timestamps={timestampsQ.data} user={user} />
            )}
          </div>
        )}
      </details>
    </div>
  );
}

function SummaryValue({
  label,
  value,
  highlight = false,
}: {
  label: string;
  value: { hours: string; minutes: string };
  highlight?: boolean;
}) {
  return (
    <span className="whitespace-nowrap">
      <span className="text-accent/70 pr-1.5">{label}</span>
      <span className={highlight ? "text-info" : "text-info/90"}>
        {value.hours}:{value.minutes} h
      </span>
    </span>
  );
}
