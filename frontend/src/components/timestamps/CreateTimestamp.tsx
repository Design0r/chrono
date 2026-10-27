import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { useState } from "react";
import { chrono } from "../../api/chrono/client";
import { timestampQueries, userQueries } from "../../api/queries";
import { datetimeLocalToIso, isoToDatetimeLocal } from "../../lib/timestamp-utils";
import { useToast } from "../Toast";

export function CreateTimestampModal({ onClose }: { onClose: () => void }) {
  const queryClient = useQueryClient();
  const [isoDate] = useState(() => new Date().toISOString());
  const [startDate, setStartDate] = useState(isoToDatetimeLocal(isoDate));
  const [endDate, setEndDate] = useState<string>(isoToDatetimeLocal(isoDate));

  const [selectedUserId, setSelectedUserId] = useState<number>(0);

  const { addToast, addErrorToast } = useToast();

  const { data: users } = useSuspenseQuery(userQueries.list());

  const mutation = useMutation({
    mutationFn: ({ userId, start, end }: { userId: number; start: string; end: string }) => {
      return chrono.timestamps.create({
        id: userId,
        user_id: userId,
        start_time: datetimeLocalToIso(start),
        end_time: datetimeLocalToIso(end),
      });
    },
    onError: (e) => addErrorToast(e),
    onSuccess: async () => {
      addToast("Created Timestamp", "success");
      await queryClient.invalidateQueries({ queryKey: timestampQueries.all() });
      onClose();
    },
  });

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
            <span className="icon-outlined text-[20px] leading-none">close</span>
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
          <button type="button" onClick={onClose} className="btn btn-soft btn-error">
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
