import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useMemo, useRef, useState } from "react";
import { ChronoClient } from "../api/chrono/client";
import {
  durationFromTimestamps,
  formatCounter,
  secondsToCounter,
} from "../lib/timestamp-utils";
import type { User } from "../types/auth";
import type { Timestamp } from "../types/response";
import { CreateTimestampModal } from "./timestamps/CreateTimestamp";
import { Timer } from "./timestamps/Timer";
import { TimestampTable } from "./timestamps/TimestampTable";
import { TimestampTableByWeek } from "./timestamps/TimestampTableByWeek";
import { useToast } from "./Toast";

export function Timestamps({ user }: { user: User }) {
  const chrono = useMemo(() => new ChronoClient(), []);
  const { addToast, addErrorToast } = useToast();
  const [runningTimer, setRunningTimer] = useState<number>(0);
  const queryClient = useQueryClient();

  const latestTimestampQ = useQuery({
    queryKey: ["timestamps", "latest"],
    queryFn: () => chrono.timestamps.getLatest(),
    staleTime: 1000 * 30,
    gcTime: 1000 * 60 * 20,
    refetchOnWindowFocus: true,
    refetchOnReconnect: true,
    retry: false,
  });

  const timestampsQ = useQuery({
    queryKey: ["timestamps"],
    queryFn: () => chrono.timestamps.getForToday(),
    staleTime: 1000 * 60 * 10,
    gcTime: 1000 * 60 * 20,
    retry: false,
  });

  const timestamps = useMemo(
    () => (timestampsQ.isError ? [] : (timestampsQ.data ?? [])),
    [timestampsQ.data, timestampsQ.isError],
  );

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
    mutationKey: ["timestamps", "start"],
    mutationFn: () => chrono.timestamps.start(),
    onError: (e) => addErrorToast(e),
    onSuccess: (data) => {
      startedTimerId.current = data.id;
      queryClient.setQueryData(["timestamps", "latest"], data);
      queryClient.invalidateQueries({ queryKey: ["timestamps"] });
    },
    retry: false,
  });

  const stopMut = useMutation({
    mutationKey: ["timestamps", "stop"],
    mutationFn: (id: number) => chrono.timestamps.stop(id),
    onError: (e) => addErrorToast(e),
    onSuccess: (data) => {
      queryClient.setQueryData(["timestamps", "latest"], data);
      queryClient.invalidateQueries({ queryKey: ["timestamps"] });
    },
    retry: false,
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
            onUpdate={(seconds: number) =>
              setRunningTimer(paused ? 0 : seconds)
            }
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
            <span
              className={`font-semibold ${paused ? "text-success" : "text-accent/80"}`}
            >
              {(() => {
                const f = formatCounter(totalTime);
                return (
                  <>
                    <span>{f.hours}</span>
                    <span
                      className={
                        paused
                          ? "text-success/90 animate-pulse"
                          : "text-accent/80 animate-pulse"
                      }
                    >
                      :
                    </span>
                    <span>{f.minutes} h</span>
                  </>
                );
              })()}
              <span
                className={`pl-2.5 ${paused ? "text-success/50" : "text-accent/60"}`}
              >
                Today
              </span>
            </span>
          </div>
        </div>
      </div>
      <div className="rounded-xl overflow-hidden w-full h-full">
        <TimestampTable timestamps={timestamps} user={user} />
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
  const chrono = useMemo(() => new ChronoClient(), []);
  const currYear = startDate
    ? new Date(startDate).getFullYear()
    : new Date().getFullYear();

  const [createTimestampModal, setCreateTimestampModal] = useState(false);

  const allTimestampsQ = useQuery({
    queryKey: ["timestamps", "all", startDate, endDate],
    queryFn: () => chrono.timestamps.getAll(startDate, endDate),
    staleTime: 1000 * 60 * 1,
    gcTime: 1000 * 60 * 30,
    retry: false,
  });

  const usersQ = useQuery({
    queryKey: ["users"],
    queryFn: () => chrono.users.getUsers({ includeInactiveUsers: false }),
    staleTime: 1000 * 60 * 1,
    gcTime: 1000 * 60 * 30,
    retry: false,
  });

  const allWorktimesQ = useQuery({
    queryKey: ["worktimes", "all", startDate, endDate],
    queryFn: () => chrono.timestamps.getWorkHoursForAllUsers(currYear),
    staleTime: 1000 * 60 * 1,
    gcTime: 1000 * 60 * 30,
    retry: false,
  });

  const queries = [allTimestampsQ, usersQ, allWorktimesQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  const usersMap = useMemo(() => {
    const users = (usersQ.data ?? []) as User[];
    return users.reduce(
      (map, u) => {
        map[u.id] = u;
        return map;
      },
      {} as Record<number, User>,
    );
  }, [usersQ.data]);

  const timestampsMap = useMemo(() => {
    const timestamps = (allTimestampsQ.data ?? []) as Timestamp[];
    return timestamps.reduce(
      (map, ts) => {
        (map[ts.user_id] ??= []).push(ts);
        return map;
      },
      {} as Record<number, Timestamp[]>,
    );
  }, [allTimestampsQ.data]);

  if (anyPending || firstError) return <></>;

  const worktimes = allWorktimesQ.data!;

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
          <CreateTimestampModal
            onClose={() => setCreateTimestampModal(false)}
          />
        )}
      </div>
      <div className="w-full">
        {Object.entries(timestampsMap).map(([k, v]) => {
          const user = usersMap[Number(k)];
          if (!user) return <div key={0} />;

          const counter = secondsToCounter(durationFromTimestamps(v));
          const worktime = worktimes[user.id];
          const expectedCounter = secondsToCounter(worktime.expected * 3600);
          let overtime = (worktime.worked - worktime.expected) * 3600;
          let overtimeLabel = "Overtime";
          if (overtime < 0) {
            overtime *= -1;
            overtimeLabel = "Missing Time";
          }
          const overtimeCounter = secondsToCounter(overtime);
          const fCounter = formatCounter(counter);
          const fExpected = formatCounter(expectedCounter);
          const fOvertime = formatCounter(overtimeCounter);

          return (
            <div key={user.id} className="my-4">
              <details className="collapse border-base-300 border collapse-arrow">
                <summary className="collapse-title bg-base-300/50 focus-within:bg-info/25 focus:text-white hover:bg-info/25 font-semibold">
                  {user.username}
                </summary>
                <div className="collapse-content px-0 pt-6 pb-0 bg-black/20 flex flex-col text-sm">
                  <h3 className="text-base font-semibold mb-3.5 px-4">
                    Overview
                  </h3>
                  <table className="bg-base-300/50 rounded-none table mb-12">
                    <thead>
                      <tr className="text-accent/80 *:w-1/3 *:font-normal">
                        <th>Worked</th>
                        <th>Expected</th>
                        <th>{overtimeLabel}</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="hover:bg-base-300 text-info/70 bg-base-200/40">
                        <td className="*:text-info">
                          <span>{fCounter.hours}</span>
                          <span>:</span>
                          <span>{fCounter.minutes}</span> h
                        </td>
                        <td className="*:text-info/90">
                          <span>{fExpected.hours}</span>
                          <span>:</span>
                          <span>{fExpected.minutes}</span> h
                        </td>
                        <td className="*:text-info/90">
                          <span>{fOvertime.hours}</span>
                          <span>:</span>
                          <span>{fOvertime.minutes}</span> h
                        </td>
                      </tr>
                    </tbody>
                  </table>
                  <TimestampTableByWeek timestamps={v} user={user} />
                </div>
              </details>
            </div>
          );
        })}
      </div>
    </>
  );
}
