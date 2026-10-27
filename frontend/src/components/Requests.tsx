import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { chrono } from "../api/chrono/client";
import { requestQueries } from "../api/queries";
import { useAppForm } from "../lib/form";
import { invalidateEventData } from "../lib/invalidate";
import type { BatchRequest, State } from "../types/response";
import { capitalize } from "../utils/string";
import { useToast } from "./Toast";

export function RequestRow({ request }: { request: BatchRequest }) {
  const reqStartDate = new Date(request.start_date);
  const reqEndDate = new Date(request.end_date);
  const { addToast, addErrorToast } = useToast();

  const queryClient = useQueryClient();

  const [dialog, setDialog] = useState(false);

  const mutation = useMutation({
    mutationFn: ({ state, reason }: { state: State; reason: string }) =>
      chrono.requests.patchRequest({
        userId: request.request.user_id,
        state: state,
        reason: reason,
        start_date: request.start_date,
        end_date: request.end_date,
      }),
    onSuccess: (_, vars) => {
      addToast(
        `${capitalize(vars.state)} request ${request.request.name} from ${request.request.username}`,
        "success",
      );
      queryClient.setQueryData(requestQueries.list().queryKey, (prev) =>
        prev?.filter((r) => r.request.request_id !== request.request.request_id),
      );
      return invalidateEventData(queryClient);
    },
    onError: (error) => addErrorToast(error),
  });

  return (
    <>
      {dialog && (
        <DeclineModal
          onClose={() => setDialog(false)}
          onDecline={(data: { state: State; reason: string }) => mutation.mutateAsync(data)}
          message={request.request.message}
        />
      )}

      <tr className="hover:bg-primary/10 pt-8 border-b border-primary/25 text-base-content/80 hover:text-primary animate-color">
        <td>{request.request.request_id}</td>
        <td>{request.request.username}</td>
        <td>{request.request.name}</td>
        <td>
          {request.event_count}{" "}
          <span className="opacity-40 px-1">{request.event_count > 1 ? "days" : "day"}</span>
        </td>
        <td>
          {reqStartDate.getDate()}.{reqStartDate.getMonth() + 1}.{reqStartDate.getFullYear()}
        </td>
        <td>
          {reqEndDate.getDate()}.{reqEndDate.getMonth() + 1}.{reqEndDate.getFullYear()}
        </td>
        <td className="align-middle">
          {request.conflicts && (
            <div className="flex flex-wrap space-x-2 space-y-2">
              {request.conflicts.map((c) => (
                <span
                  key={c.id}
                  className="badge badge-soft badge-info badge-sm cursor-default hover:text-white/60 hover:bg-info/20 animate-color"
                >
                  {c.username}
                </span>
              ))}
            </div>
          )}
        </td>
        <td className="flex justify-center align-middle gap-2">
          <button
            onClick={() => mutation.mutate({ state: "accepted", reason: "" })}
            className="btn btn-soft btn-sm btn-primary icon-outlined animate-color"
          >
            check
          </button>
          <button
            onClick={() => setDialog(true)}
            className="btn btn-soft btn-sm btn-error icon-outlined animate-color"
          >
            close
          </button>
        </td>
      </tr>
    </>
  );
}

export function DeclineModal({
  message,
  onClose,
  onDecline,
}: {
  message: string | null;
  onClose: () => void;
  onDecline: ({ state, reason }: { state: State; reason: string }) => Promise<unknown>;
}) {
  const form = useAppForm({
    defaultValues: { reason: "" },
    onSubmit: ({ value }) => onDecline({ state: "declined", reason: value.reason }).catch(() => {}),
  });

  return (
    <div className="fixed z-10 inset-0 bg-opacity-50 backdrop-blur-md flex items-center justify-center">
      <dialog className="modal modal-open">
        <div className="modal-box">
          <button
            type="button"
            onClick={onClose}
            className="absolute cursor-pointer right-3 top-3 icon-outlined items-end text-xl justify-end"
          >
            close
          </button>
          <h1 className="text-xl font-bold">Reject Request</h1>
          {message && <h4>{message}</h4>}
          <form.AppForm>
            <form.Form className="modal-backdrop">
              <label htmlFor="reason">
                <h5>Reject Reason</h5>
              </label>
              <form.AppField name="reason">
                {(field) => (
                  <field.TextArea className="textarea w-full text-white textarea-bordered" />
                )}
              </form.AppField>
              <br />
              <form.SubmitButton className="btn btn-error text-xl animate-color">
                Decline
              </form.SubmitButton>
            </form.Form>
          </form.AppForm>
        </div>
      </dialog>
    </div>
  );
}
