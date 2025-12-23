import { useMutation, useQuery } from "@tanstack/react-query";
import { ChronoClient } from "../api/chrono/client";
import { useEffect, useState } from "react";
import type { Timestamp } from "../types/response";
import { useToast } from "./Toast";

export function Timestamps() {
  const chrono = new ChronoClient();

  const [timestamps, setTimestemps] = useState<Timestamp[]>([]);
  const [timer, setTimer] = useState<boolean>(false);
  const { addToast, addErrorToast } = useToast();

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
    // onSuccess: () => setNotifications([]),
    retry: false,
  });

  const stopMut = useMutation({
    mutationKey: ["timestamps", "stop"],
    mutationFn: () => chrono.timestamps.start(),
    onError: (e) => addErrorToast(e),
    // onSuccess: () => setNotifications([]),
    retry: false,
  });

  useEffect(() => {
    if (timestampsQ.isError) return;
    setTimestemps(timestampsQ.data || []);
  }, [timestampsQ.data, timestampsQ.isError]);

  useEffect(() => {
    if (timestampsQ.isError) addErrorToast(timestampsQ.error);
  }, [timestampsQ.isError]);

  return (
    <div>
      {timer && <div>Running...</div>}
      <button className="btn" onClick={() => startMut.mutate()}>
        Start
      </button>
      <button className="btn" onClick={() => stopMut.mutate()}>
        Stop
      </button>
      <div>
        {timestamps.map((t, i) => (
          <div key={i}>
            {t.start_date} - {t.end_date}
          </div>
        ))}
      </div>
    </div>
  );
}
