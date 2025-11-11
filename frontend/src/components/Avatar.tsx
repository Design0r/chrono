import { useAuth } from "../auth";
import type { User } from "../types/auth";
import { hexToHSL } from "../utils/colors";

export function Avatar({ user }: { user: User }) {
  const auth = useAuth();

  let initial = "?";
  if (user.username.length > 0) {
    initial = user.username.slice(0, 1);
  }
  const [h, s, l] = hexToHSL(user.color);
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
        {user.is_superuser && (
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
