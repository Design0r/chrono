import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { chrono } from "../api/chrono/client";
import { notificationQueries } from "../api/queries";
import type { Notification } from "../types/response";
import { useToast } from "./Toast";

export function Notifications() {
  const { addErrorToast } = useToast();
  const queryClient = useQueryClient();

  const notifs = useQuery(notificationQueries.list());

  const mutation = useMutation({
    mutationFn: () => chrono.notifications.clearAll(),
    onError: (e) => addErrorToast(e),
    onSuccess: () => queryClient.setQueryData(notificationQueries.list().queryKey, []),
  });

  const notifications: Notification[] = notifs.isError ? [] : (notifs.data ?? []);

  useEffect(() => {
    if (notifs.isError) addErrorToast(notifs.error);
  }, [notifs.isError, notifs.error, addErrorToast]);

  return (
    <div className="indicator">
      {notifications.length > 0 && (
        <span className="indicator-item font-bold rounded-full align-items-start text-white/85 border border-error/50 badge badge-error backdrop-blur-md bg-error/40 p-0 h-6 px-2 pointer-events-none">
          {notifications.length}
        </span>
      )}
      <div className="dropdown dropdown-end">
        <button
          tabIndex={0}
          className="btn btn-ghost rounded-full text-xl icon-outlined bg-base-300 animate-color"
        >
          notifications
        </button>

        <ul
          tabIndex={-1}
          className="mt-1.5 min-w-64 pt-4 pb-3 px-3 dropdown-content menu bg-base-300/50 backdrop-blur-xl rounded-box z-10 drop-shadow-xl"
        >
          <p className="px-3 pb-2 text-lg font-bold">Notifications</p>
          <hr className="border-base-300/80 pb-2" />
          {notifications.map((n, i) => (
            <NotificationElement
              key={i}
              onClear={(id: number) =>
                queryClient.setQueryData(notificationQueries.list().queryKey, (prev) =>
                  (prev ?? []).filter((n) => n.id !== id),
                )
              }
              notification={n}
            />
          ))}
          <hr className="border-base-300/80 pb-2" />
          <button
            onClick={() => mutation.mutate()}
            className="btn btn-soft bg-warning/70 text-warning-content shadow-none border-0 hover:bg-warning animate-all hover:text-warning-content rounded-xl"
          >
            Clear All
          </button>
        </ul>
      </div>
    </div>
  );
}

export function NotificationElement({
  notification,
  onClear,
}: {
  notification: Notification;
  onClear: (id: number) => void;
}) {
  const { addErrorToast } = useToast();
  const mutation = useMutation({
    mutationFn: () => chrono.notifications.clear(notification.id),
    onError: (e) => addErrorToast(e),
  });

  return (
    <li className="py-1">
      <div className="hover:text-white">
        <p>{notification.message}</p>
        <button
          onClick={() => {
            mutation.mutate();
            onClear(notification.id);
          }}
          className="btn btn-soft border-0 hover:border-0 bg-base-200/50 text-xl font-semibold icon-outlined hover:bg-primary hover:text-neutral animate-color"
        >
          close
        </button>
      </div>
    </li>
  );
}
