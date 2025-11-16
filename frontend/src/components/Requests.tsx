import { useMutation } from "@tanstack/react-query";
import type { BatchRequest, State } from "../types/response";
import { ChronoClient } from "../api/chrono/client";
import { useState } from "react";
import { useToast } from "./Toast";
import { capitalize } from "../utils/string";

export function RequestRow({ request }: { request: BatchRequest }) {
  const reqStartDate = new Date(request.start_date);
  const reqEndDate = new Date(request.end_date);
  const { addToast, addErrorToast } = useToast();

  const chrono = new ChronoClient();

  const [visible, setVisible] = useState(true);
  const [dialog, setDialog] = useState(false);

  const mutation = useMutation({
    mutationKey: ["requests", request.request.request_id],
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
      setVisible(false);
    },
    onError: (error) => addErrorToast(error),
  });

  if (!visible) return <></>;
  return (
    <>
      {dialog && (
        <DeclineModal
          onClose={() => setDialog(false)}
          onDecline={(data: { state: State; reason: string }) =>
            mutation.mutate(data)
          }
          message={request.request.message}
        />
      )}

      <tr className="hover:bg-primary/10 pt-8 border-b-1 border-primary/25 text-base-content/80 hover:text-primary animate-color">
        <td>{request.request.request_id}</td>
        <td>{request.request.username}</td>
        <td>{request.request.name}</td>
        <td>
          {request.event_count}{" "}
          <span className="opacity-40 px-1">
            {request.event_count > 1 ? "days" : "day"}
          </span>
        </td>
        <td>
          {reqStartDate.getDate()}.{reqStartDate.getMonth() + 1}.
          {reqStartDate.getFullYear()}
        </td>
        <td>
          {reqEndDate.getDate()}.{reqEndDate.getMonth() + 1}.
          {reqEndDate.getFullYear()}
        </td>
        <td className="align-middle">
          {request.conflicts && (
            <div className="flex flex-wrap space-x-2 space-y-2">
              {request.conflicts.map((c) => (
                <span className="badge badge-soft badge-info badge-sm cursor-default hover:text-white/60 hover:bg-info/20 animate-color">
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
  onDecline: ({ state, reason }: { state: State; reason: string }) => void;
}) {
  const [text, setText] = useState("");

  return (
    <div className="fixed z-10 inset-0 bg-opacity-50 backdrop-blur-md flex items-center justify-center">
      <dialog className="modal modal-open">
        <div className="modal-box">
          <button
            onClick={onClose}
            className="absolute cursor-pointer right-3 top-3 icon-outlined items-end text-xl justify-end"
          >
            close
          </button>
          <h1 className="text-xl font-bold">Reject Requeust</h1>
          {message && <h4>{message}</h4>}
          <div className="modal-backdrop">
            <h5>Reject Reason</h5>
            <textarea
              onChange={(e) => setText(e.target.value)}
              className="textarea w-full text-white textarea-bordered"
            ></textarea>
            <br />
            <button
              id="reject-btn"
              className="btn btn-error text-xl animate-color"
              onClick={() => onDecline({ state: "declined", reason: text })}
            >
              Decline
            </button>
          </div>
        </div>
      </dialog>
    </div>
  );
}
