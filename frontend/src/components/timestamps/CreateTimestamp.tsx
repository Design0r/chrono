import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { ChronoClient } from "../../api/chrono/client";
import {
  datetimeLocalToIso,
  isoToDatetimeLocal,
} from "../../lib/timestamp-utils";
import type { User } from "../../types/auth";
import { ErrorPage } from "../ErrorPage";
import { LoadingSpinnerPage } from "../LoadingSpinner";
import { useToast } from "../Toast";

export function CreateTimestampModal({ onClose }: { onClose: () => void }) {
  const queryClient = useQueryClient();
  const [startDate, setStartDate] = useState(
    isoToDatetimeLocal(new Date().toISOString()),
  );
  const [endDate, setEndDate] = useState<string>(
    isoToDatetimeLocal(new Date().toISOString()),
  );

  const [selectedUserId, setSelectedUserId] = useState<number>(0);

  const chrono = new ChronoClient();
  const { addToast, addErrorToast } = useToast();

  const usersQ = useQuery({
    queryKey: ["users"],
    queryFn: () => chrono.users.getUsers({ includeInactiveUsers: false }),
    staleTime: 1000 * 60 * 30, // 30min
    gcTime: 1000 * 60 * 60 * 1, // 1h
    retry: false,
  });

  const mutation = useMutation({
    mutationKey: ["timestamps", startDate, endDate, selectedUserId],
    mutationFn: ({
      userId,
      start,
      end,
    }: {
      userId: number;
      start: string;
      end: string;
    }) => {
      return chrono.timestamps.create({
        id: userId,
        user_id: userId,
        start_time: datetimeLocalToIso(start),
        end_time: datetimeLocalToIso(end),
      });
    },
    onError: (e) => addErrorToast(e),
    onSuccess: () => {
      addToast("Created Timestamp", "success");
      queryClient.invalidateQueries({ queryKey: ["timestamps"] });
      onClose();
    },
    retry: false,
  });

  const queries = [usersQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending) return <LoadingSpinnerPage />;
  if (firstError) return <ErrorPage error={firstError} />;

  const users = usersQ.data! as User[];

  return (
    <div className="fixed inset-0 z-50 flex text-white items-center justify-center p-4">
      <button
        aria-label="Close modal"
        onClick={onClose}
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
      />

      <div
        role="dialog"
        aria-modal="true"
        className="relative w-full max-w-lg rounded-2xl bg-base-100 shadow-2xl ring-1 ring-black/10"
      >
        <div className="flex items-center justify-between px-5 py-4 border-b border-black/10">
          <h2 className="text-base font-semibold">Create Timestamp</h2>

          <button
            type="button"
            onClick={onClose}
            className="inline-flex h-9 w-9 items-center justify-center rounded-full hover:bg-slate-100 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-400"
          >
            <span className="icon-outlined text-[20px] leading-none">
              close
            </span>
          </button>
        </div>

        <div className="px-5 py-4 space-y-4 flex flex-col">
          <label className="space-y-2">
            <span className="text-sm font-medium">User</span>
            <select
              defaultValue={""}
              onChange={(e) => setSelectedUserId(Number(e.target.value))}
              className="w-full cursor-pointer bg-base-100 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-xl pl-4 text-base border rounded-md"
            >
              <option value="">-</option>
              {users.map((u) => (
                <option key={u.id} value={u.id}>
                  {u.username}
                </option>
              ))}
            </select>
          </label>

          <div className="grid gap-4 sm:grid-cols-2">
            <label className="space-y-2">
              <span className="text-sm font-medium">Start</span>
              <input
                type="datetime-local"
                className="input"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
              />
            </label>

            <label className="space-y-2">
              <span className="text-sm font-medium">End</span>
              <input
                type="datetime-local"
                className="input"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
              />
            </label>
          </div>
        </div>

        <div className="flex items-center justify-end gap-2 px-5 py-4 border-t border-black/10">
          <button
            type="button"
            onClick={onClose}
            className="btn btn-soft btn-error"
          >
            Cancel
          </button>
          <button
            type="button"
            className="btn btn-soft btn-success"
            onClick={() => {
              mutation.mutate({
                userId: selectedUserId,
                start: startDate,
                end: endDate,
              });
              onClose();
            }}
            disabled={mutation.isPending}
          >
            {mutation.isPending ? "Creating..." : "Create"}
          </button>
        </div>
      </div>
    </div>
  );
}
