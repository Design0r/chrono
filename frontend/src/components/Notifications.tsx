import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { ChronoClient } from "../api/chrono/client";
import type { Notification } from "../types/response";
import { useToast } from "./Toast";

export function Notifications() {
  const chrono = new ChronoClient();
  const { addErrorToast } = useToast();
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const notifs = useQuery({
    queryKey: ["notifcations"],
    queryFn: () => chrono.notifications.get(),
    staleTime: 1000 * 60 * 10, // 10min
    gcTime: 1000 * 60 * 20, // 20min
  });

  const mutation = useMutation({
    mutationKey: ["notifications", "clear"],
    mutationFn: () => chrono.notifications.clearAll(),
    onError: (e) => addErrorToast(e),
    onSuccess: () => setNotifications([]),
  });

  useEffect(() => {
    if (notifs.isError) return;
    setNotifications(notifs.data || []);
  }, [notifs.data, notifs.isError]);

  useEffect(() => {
    if (notifs.isError) addErrorToast(notifs.error);
  }, [notifs.isError]);

  return (
    <div className="indicator">
      {notifications.length > 0 && (
        <span className="indicator-item font-bold rounded-full align-items-start text-white/85 border border-error/50 badge badge-error backdrop-blur-md bg-error/40 p-0 h-6 px-2 pointer-events-none">
          {notifications.length}
        </span>
      )}
      <div className="dropdown dropdown-end">
        <button className="btn btn-ghost px-5 border-1.5 border-white/2 py-1 hover:bg-info/20 rounded-full text-xl icon-outlined bg-base-100 animate-color">
          notifications
        </button>

        <ul
          tabIndex={0}
          className="mt-1.5 min-w-64 pt-4 pb-3 px-3 dropdown-content menu bg-info/20 backdrop-blur-xl rounded-box z-10 drop-shadow-xl"
        >
          <p className="px-3 pb-2 text-lg font-bold">Notifications</p>
          <hr className="border-base-200/80 pb-2" />
          {notifications.map((n, i) => (
            <NotificationElement
              key={i}
              onClear={(id: number) =>
                setNotifications((prev) => prev.filter((n) => n.id !== id))
              }
              notification={n}
            />
          ))}
          <button
            onClick={() => mutation.mutate()}
            className="mt-4 btn btn-soft rounded-xl text-neutral border-0 hover:border-0 bg-primary/90 font-semibold hover:bg-primary hover:text-neutral animate-color"
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
  const chrono = new ChronoClient();
  const mutation = useMutation({
    mutationKey: ["notification", notification.id],
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
