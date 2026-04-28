import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";
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
  const [timestamps, setTimestemps] = useState<Timestamp[]>([]);
  const [paused, setPaused] = useState<boolean>(true);
  const { addToast, addErrorToast } = useToast();
  const [currTimer, setCurrTimer] = useState<Timestamp | null>(null);
  const [startTime, setStartTime] = useState<number>(Date.now());
  const [runningTimer, setRunningTimer] = useState<number>(0);
  const queryClient = useQueryClient();

  const latestTimestampQ = useQuery({
    queryKey: ["timestamps", "latest"],
    queryFn: () => chrono.timestamps.getLatest(),
    staleTime: 1000 * 60 * 10,
    gcTime: 1000 * 60 * 20,
    retry: false,
  });

  const timestampsQ = useQuery({
    queryKey: ["timestamps"],
    queryFn: () => chrono.timestamps.getForToday(),
    staleTime: 1000 * 60 * 10,
    gcTime: 1000 * 60 * 20,
    retry: false,
  });

  const startMut = useMutation({
    mutationKey: ["timestamps", "start"],
    mutationFn: () => chrono.timestamps.start(),
    onError: (e) => addErrorToast(e),
    onSuccess: (data) => {
      setCurrTimer(data);
      setPaused(false);
      setStartTime(Date.now());
      addToast("Started Timer", "success");
    },
    retry: false,
  });

  const stopMut = useMutation({
    mutationKey: ["timestamps", "stop"],
    mutationFn: (id: number) => chrono.timestamps.stop(id),
    onError: (e) => addErrorToast(e),
    onSuccess: () => {
      setCurrTimer(null);
      setStartTime(Date.now());
      setPaused(true);
      queryClient.invalidateQueries({ queryKey: ["timestamps"] });
      addToast("Stopped Timer", "success");
    },
    retry: false,
  });

  useEffect(() => {
    if (timestampsQ.isError) return;
    setTimestemps(timestampsQ.data || []);
  }, [timestampsQ.data, timestampsQ.isError]);

  useEffect(() => {
    if (latestTimestampQ.isError) return;
    const latest = latestTimestampQ.data;
    if (!latest) return;
    const hasEnded = latest.end_time !== null;
    if (!hasEnded && latest.id !== currTimer?.id) {
      addToast("Resuming latest unfinished Timer", "info");
      setCurrTimer(latest);
      setStartTime(Date.parse(latest.start_time));
      setPaused(false);
    } else setPaused(true);
  }, [latestTimestampQ.data, latestTimestampQ.isError]);

  useEffect(() => {
    if (timestampsQ.isError) addErrorToast(timestampsQ.error);
  }, [timestampsQ.isError]);

  const totalTime = secondsToCounter(
    durationFromTimestamps(timestamps) + runningTimer,
  );

  return (
    <div className="mx-auto flex lg:flex-row flex-col w-full gap-4 justify-center ">
      <div className="flex items-center justify-center flex-1 rounded-2xl bg-base-200/60 border border-base-300/80 p-6 lg:p-8">
        <div className="flex flex-col items-center justify-center gap-4">
          <Timer
            paused={paused}
            startUnix={startTime}
            onUpdate={(seconds: number) => setRunningTimer(seconds)}
          />
          <div className="flex mt-1 gap-3 justify-center items-center">
            <button
              disabled={!paused}
              className="btn btn-lg btn-success w-22 rounded-full shadow-md hover:shadow-lg transition-shadow disabled:opacity-40 disabled:shadow-none icon-filled"
              onClick={() => startMut.mutate()}
              title="Timer starten"
            >
              <span className="text-xl icon-filled scale-145">play_arrow</span>
            </button>
            <button
              disabled={paused}
              className="btn btn-circle btn-lg btn-error shadow-md hover:shadow-lg transition-shadow disabled:opacity-40 disabled:shadow-none icon-outlined"
              onClick={() => currTimer && stopMut.mutate(currTimer.id)}
              title="Timer stoppen"
            >
              <span className="text-xl icon-filled scale-125">stop</span>
            </button>
          </div>
          <div className="flex whitespace-nowrap items-center gap-2 mt-0 text-base-content/70">
            <span
              className={`font-mono font-semibold ${paused ? "text-success" : "text-accent/80"}`}
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
  user: _currUser,
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

// Re-exports for routes and other components
export {
  datetimeLocalToIso,
  durationFromTimestamps,
  formatCounter,
  groupTimestampsByDay,
  groupTimestampsByWeek,
  isoToDateLocal,
  isoToDatetimeLocal,
  secondsToCounter,
} from "../lib/timestamp-utils";
export type { DayGroup, TimeCounter, WeekGroup } from "../lib/timestamp-utils";
export { TimestampTableByWeek } from "./timestamps/TimestampTableByWeek";
