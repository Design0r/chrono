import { useMutation, useQueryClient } from "@tanstack/react-query";
import { chrono } from "../../api/chrono/client";
import { timestampQueries } from "../../api/queries";
import { useAppForm } from "../../lib/form";
import { datetimeLocalToIso, isoToDatetimeLocal, rangeError } from "../../lib/timestamp-utils";
import type { Timestamp } from "../../types/response";
import { useToast } from "../Toast";

export function EditTimestampModal({
  timestamp,
  onClose,
}: {
  timestamp: Timestamp;
  onClose: () => void;
}) {
  const queryClient = useQueryClient();
  const { addToast, addErrorToast } = useToast();

  const mutation = useMutation({
    mutationFn: ({ start, end }: { start: string; end: string }) =>
      chrono.timestamps.update({
        id: timestamp.id,
        user_id: timestamp.user_id,
        start_time: datetimeLocalToIso(start),
        end_time: end ? datetimeLocalToIso(end) : null,
      }),
    onError: (e) => addErrorToast(e),
    onSuccess: async () => {
      addToast("Updated Timestamp", "success");
      await queryClient.invalidateQueries({ queryKey: timestampQueries.all() });
      onClose();
    },
  });

  const form = useAppForm({
    defaultValues: {
      start: isoToDatetimeLocal(timestamp.start_time),
      end: timestamp.end_time ? isoToDatetimeLocal(timestamp.end_time) : "",
    },
    validators: {
      onChange: ({ value }) => rangeError(value.start, value.end),
    },
    onSubmit: ({ value }) => mutation.mutateAsync(value).catch(() => {}),
  });

  return (
    <div className="fixed inset-0 z-50 flex text-white items-center justify-center p-4">
      <button
        aria-label="Close modal"
        onClick={onClose}
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
      />

      <form.AppForm>
        <form.Form
          role="dialog"
          aria-modal="true"
          className="relative w-full max-w-lg rounded-2xl bg-base-100 shadow-2xl ring-1 ring-black/10"
        >
          <div className="flex items-center justify-between px-5 py-4 border-b border-black/10">
            <h2 className="text-base font-semibold">Edit Timestamp</h2>

            <button
              type="button"
              onClick={onClose}
              className="inline-flex h-9 w-9 items-center justify-center rounded-full hover:bg-slate-100 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-400"
            >
              <span className="icon-outlined text-[20px] leading-none">close</span>
            </button>
          </div>

          <div className="px-5 py-4 space-y-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <label className="space-y-2">
                <span className="text-sm font-medium">Start</span>
                <form.AppField name="start">
                  {(field) => <field.TextInput type="datetime-local" className="input" required />}
                </form.AppField>
              </label>

              <label className="space-y-2">
                <span className="text-sm font-medium">End</span>
                <form.AppField name="end">
                  {(field) => <field.TextInput type="datetime-local" className="input" />}
                </form.AppField>
              </label>
            </div>
            <form.FormErrors className="text-sm text-error" />
          </div>

          <div className="flex items-center justify-end gap-2 px-5 py-4 border-t border-black/10">
            <button type="button" onClick={onClose} className="btn btn-soft btn-error">
              Cancel
            </button>
            <form.SubmitButton className="btn btn-soft btn-success" pending="Saving...">
              Save
            </form.SubmitButton>
          </div>
        </form.Form>
      </form.AppForm>
    </div>
  );
}
