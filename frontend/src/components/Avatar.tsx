import { useQuery } from "@tanstack/react-query";
import { useRouter } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { ChronoClient } from "../api/chrono/client";
import { useAuth } from "../auth";
import type { User } from "../types/auth";
import { hexToHSL } from "../utils/colors";

export function Avatar() {
  const auth = useAuth();
  const router = useRouter();
  const chrono = new ChronoClient();
  const [initial, setInitial] = useState("");
  const [user, setUser] = useState<User | null>(null);

  const { data } = useQuery({
    queryKey: ["user", auth.userId],
    queryFn: () => chrono.users.getUserById(auth.userId!),
    staleTime: 60_000,
  });

  useEffect(() => {
    if (!data) return;
    if (data.username.length > 0) {
      setInitial(data.username.slice(0, 1));
    }
    setUser(data);
  }, [data]);

  const [h, s, l] = hexToHSL(user ? user.color : "#000");
  const bgColor = `hsla(${h.toFixed(1)}, ${(s * 100).toFixed(1)}%, ${l * 100}%, 0.25)`;
  const borderColor = `hsla(${h.toFixed(1)}, ${(s * 100).toFixed(1)}%, ${l * 100}%, 0.6)`;
  const textColor = `hsla(${h.toFixed(1)}, ${(s * 100).toFixed(1)}%, ${l * 100}%, 0.6)`;

  return (
    <div className="dropdown dropdown-end pr-2">
      <div
        tabIndex={0}
        role="button"
        className="avatar avatar-placeholder cursor-pointer"
      >
        <div
          className="w-10 border rounded-full text-neutral-content"
          style={{ backgroundColor: bgColor, borderColor: borderColor }}
        >
          <span className="text-xl pt-px" style={{ color: textColor }}>
            {initial}
          </span>
        </div>
      </div>
      <ul
        tabIndex={0}
        className="dropdown-content mt-1.5 min-w-40 pt-4 pb-3 px-3 menu bg-info/20 backdrop-blur-xl rounded-box z-10 drop-shadow-xl animate-color"
      >
        <li>
          <a className="py-2.5" href="/profile">
            Profile
          </a>
        </li>
        {user?.is_superuser && (
          <li>
            <a className="py-2.5" href="/settings">
              Settings
            </a>
          </li>
        )}
        <li>
          <button
            onClick={async () => {
              await auth.logout();
              await router.invalidate();
              await router.navigate({ to: "/login" });
            }}
            className="w-full text-left py-1"
          >
            Logout
          </button>
        </li>
      </ul>
    </div>
  );
}
