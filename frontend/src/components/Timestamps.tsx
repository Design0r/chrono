import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ChronoClient } from "../api/chrono/client";
import { useEffect, useMemo, useState } from "react";
import type { Timestamp } from "../types/response";
import { useToast } from "./Toast";

export function Timestamps() {
  const chrono = useMemo(() => new ChronoClient(), []);

  const [timestamps, setTimestemps] = useState<Timestamp[]>([]);
  const [paused, setPaused] = useState<boolean>(true);
  const { addToast, addErrorToast } = useToast();
  const [currTimer, setCurrTimer] = useState<Timestamp | null>(null);
  const [startTime, setStartTime] = useState<number>(0);

  const queryClient = useQueryClient();

  const latestTimestampQ = useQuery({
    queryKey: ["timestamps", "latest"],
    queryFn: () => chrono.timestamps.getLatest(),
    staleTime: 1000 * 60 * 10, // 10min
    gcTime: 1000 * 60 * 20, // 20min
    retry: false,
  });

  const timestampsQ = useQuery({
    queryKey: ["timestamps"],
    queryFn: () => chrono.timestamps.getForToday(),
    staleTime: 1000 * 60 * 10, // 10min
    gcTime: 1000 * 60 * 20, // 20min
    retry: false,
  });

  const startMut = useMutation({
    mutationKey: ["timestamps", "start"],
    mutationFn: () => chrono.timestamps.start(),
    onError: (e) => addErrorToast(e),
    onSuccess: (data) => {
      setPaused(false);
      setCurrTimer(data);
      queryClient.invalidateQueries({ queryKey: ["timestamps"] });
      addToast("Started Timer");
    },
    retry: false,
  });

  const stopMut = useMutation({
    mutationKey: ["timestamps", "stop"],
    mutationFn: (id: number) => chrono.timestamps.stop(id),
    onError: (e) => addErrorToast(e),
    onSuccess: () => {
      setPaused(true);
      setCurrTimer(null);
      queryClient.invalidateQueries({ queryKey: ["timestamps"] });
      addToast("Stopped Timer");
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
    if (!hasEnded) {
      addToast("Resuming latest unfinished Timer", "info");
      setCurrTimer(latest);
      setStartTime(Date.parse(latest.start_time));
      setPaused(false);
    } else setPaused(true);
  }, [latestTimestampQ.data, latestTimestampQ.isError]);

  useEffect(() => {
    if (timestampsQ.isError) addErrorToast(timestampsQ.error);
  }, [timestampsQ.isError]);

  return (
    <div>
      <Timer paused={paused} startUnix={startTime} />
      <button className="btn" onClick={() => startMut.mutate()}>
        Start
      </button>
      <button
        className="btn"
        onClick={() => {
          if (!currTimer) return;
          stopMut.mutate(currTimer.id);
        }}
      >
        Stop
      </button>
      <div>
        {timestamps.map((t, i) => (
          <div key={i}>
            {t.start_time} - {t.end_time}
          </div>
        ))}
      </div>
    </div>
  );
}

type TimeCounter = {
  hours: number;
  minutes: number;
  seconds: number;
};

function secondsToCounter(totalSeconds: number): TimeCounter {
  const seconds = Math.max(0, Math.floor(totalSeconds));
  const hours = Math.floor(seconds / 60 / 60);
  const minutes = Math.floor(seconds / 60) % 60;
  const s = seconds % 60;
  return { hours, minutes, seconds: s };
}

function Timer({ startUnix, paused }: { startUnix: number; paused: boolean }) {
  const [timer, setTimer] = useState<TimeCounter>(() => secondsToCounter(0));

  useEffect(() => {
    if (!startUnix) {
      setTimer(secondsToCounter(0));
      return;
    }

    const tick = () => {
      const elapsedSeconds = (Date.now() - startUnix) / 1000;
      setTimer(secondsToCounter(elapsedSeconds));
    };

    tick();

    if (paused) return; // no interval while paused

    const interval = setInterval(tick, 1000);
    return () => clearInterval(interval);
  }, [startUnix, paused]);

  return (
    <div>
      {timer.hours}h {timer.minutes}m {timer.seconds}s
    </div>
  );
}
