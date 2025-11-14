import { useRouter } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import type { User } from "../types/auth";
import { hexToHSL, hsla } from "../utils/colors";

export function Avatar({ user }: { user: User | null }) {
  const router = useRouter();
  const [initial, setInitial] = useState("?");

  useEffect(() => {
    if (!user) return;
    if (user.username.length > 0) {
      setInitial(user.username.slice(0, 1));
    }
  }, [user]);

  const hsl = hexToHSL(user ? user.color : "#000");
  const borderColor = hsla(...hsl, 0.6);
  const bgColor = hsla(...hsl, 0.2);
  const textColor = hsla(...hsl, 1);

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
