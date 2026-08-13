import { useQuery } from "@tanstack/react-query";
import {
  Link,
  useLocation,
  useNavigate,
  type LinkProps,
  type RegisteredRouter,
} from "@tanstack/react-router";
import type { CSSProperties, ReactNode } from "react";
import type { ChronoClient } from "../api/chrono/client";
import { useAuth } from "../auth";
import { isoToDateLocal } from "../lib/timestamp-utils";
import { Avatar } from "./Avatar";
import { Notifications } from "./Notifications";

const ADMIN_MENU_ID = "admin-menu";

type AdminPath = "/requests" | "/tokens" | "/export" | "/settings";

function AdminMenu() {
  const navigate = useNavigate();

  const go = (to: AdminPath) => {
    document.getElementById(ADMIN_MENU_ID)?.hidePopover();
    void navigate({ to });
  };

  return (
    <div className="inline-block!  flex-none!">
      <button
        type="button"
        className="btn btn-soft py-6 border-0 bg-primary/10 rounded-xl hover:bg-secondary/10 max-lg:min-w-24"
        popoverTarget={ADMIN_MENU_ID}
        style={{ anchorName: "--admin-anchor" } as CSSProperties}
      >
        <span className="font-normal text-base tracking-wide opacity-90">
          Admin
        </span>
        <span className="icon-outlined pr-0.5">arrow_drop_down</span>
      </button>
      {/* Unter lg sitzt der Anchor im Dock am unteren Rand: dort muss das Menü
          nach oben und am rechten Rand ausgerichtet aufgehen, sonst liegt es
          außerhalb des Viewports. flip-block fängt den Rest ab. */}
      <ul
        id={ADMIN_MENU_ID}
        popover="auto"
        className="dropdown max-lg:dropdown-top max-lg:dropdown-end [position-try-fallbacks:flip-block] menu w-52 rounded-box bg-base-100 p-2 shadow-sm"
        style={{ positionAnchor: "--admin-anchor" } as CSSProperties}
      >
        <li>
          <button type="button" onClick={() => go("/requests")}>
            <span className="icon-outlined pr-0.25">mark_chat_unread</span>
            Requests
          </button>
        </li>
        <li>
          <button type="button" onClick={() => go("/tokens")}>
            <span className="icon-outlined pr-0.25">local_activity</span>
            Tokens
          </button>
        </li>
        <li>
          <button type="button" onClick={() => go("/export")}>
            <span className="icon-outlined pr-0.25">file_export</span>
            Export
          </button>
        </li>
        <li>
          <button type="button" onClick={() => go("/settings")}>
            <span className="icon-outlined pr-0.25">settings</span>
            Settings
          </button>
        </li>
      </ul>
    </div>
  );
}

export function Header({ chrono }: { chrono: ChronoClient }) {
  const auth = useAuth();

  const userQ = useQuery({
    queryKey: ["user", auth.userId],
    queryFn: () => {
      if (!auth.userId) throw new Error("no user id found");
      return chrono.users.getUserById(auth.userId);
    },
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
    retry: false,
  });

  const date = new Date();
  const startDate = new Date(Date.UTC(date.getFullYear(), 0, 1, 0, 0, 0, 0));
  const endDate = new Date(Date.UTC(date.getFullYear() + 1, 0, 1, 0, 0, 0, 0));

  return (
    <div className="mb-4 mx-auto p-4 lg:px-4">
      <div className="navbar flex justify-between">
        <div className="flex items-center">
          <div className="pr-14">
            {/* width/height entsprechen der viewBox von chrono.svg (639x103).
                Ohne die Attribute kennt der Browser das Seitenverhältnis erst
                nach dem Laden und die Navbar springt beim Rendern. */}
            <img
              className="w-40"
              alt="chrono logo"
              src="/chrono.svg"
              width={639}
              height={103}
            />
          </div>

          {auth.isAuthenticated && (
            <ul
              className="z-20! max-lg:dock w-full flex justify-between max-lg:border-t max-lg:border-info/10 max-lg:bg-base-200/70! backdrop-blur-lg overflow-x-auto overflow-y-hidden gap-1 md:gap-2 lg:w-fit
						*:flex *:flex-col! *:lg:flex-row! *:lg:gap-2 *:lg:items-center"
            >
              <MenuButton to="/">
                <span className="icon-outlined pr-0.25">home</span>
                <span className="font-normal text-base tracking-wide opacity-90">
                  Home
                </span>
              </MenuButton>
              <MenuButton
                to="/calendar/$year/$month"
                params={{
                  year: date.getFullYear().toString(),
                  month: (date.getMonth() + 1).toString(),
                }}
              >
                <span className="icon-outlined pr-0.25">calendar_today</span>
                <span className="font-normal text-base tracking-wide opacity-90">
                  Calendar
                </span>
              </MenuButton>
              <MenuButton
                to="/timestamps"
                search={{
                  startDate: isoToDateLocal(startDate.toISOString()),
                  endDate: isoToDateLocal(endDate.toISOString()),
                }}
              >
                <span className="icon-outlined pr-0.25">timer</span>
                <span className="font-normal text-base tracking-wide opacity-90">
                  Timestamps
                </span>
              </MenuButton>
              <MenuButton to="/team">
                <span className="icon-outlined pr-0.25">group</span>
                <span className="font-normal text-base tracking-wide opacity-90">
                  Team
                </span>
              </MenuButton>
              <span className="text-sm px-4 tracking-wide opacity-60">|</span>
              {userQ.data?.is_superuser && <AdminMenu />}
            </ul>
          )}
        </div>
        <div className="flex items-center justify-end gap-6">
          {!auth.isAuthenticated ? (
            <>
              <a href="/login" className="btn btn-ghost">
                Login
              </a>
              {
                <a href="/signup" className="btn btn-ghost">
                  Signup
                </a>
              }
            </>
          ) : (
            <>
              <Notifications />
              <Avatar user={userQ.data} />
            </>
          )}
        </div>
      </div>
    </div>
  );
}

interface MenuButtonProps extends LinkProps<RegisteredRouter> {
  children?: ReactNode | ReactNode[];
  /** Überschreibt den Pfad, ab dem der Eintrag als aktiv gilt. */
  activePath?: string;
}

/**
 * Statischer Teil eines Route-Patterns, also alles vor dem ersten Parameter:
 * '/calendar/$year/$month' -> '/calendar'
 */
function staticPrefix(to: string): string {
  const segments = to.split("/").filter(Boolean);
  const params = segments.findIndex((s) => s.startsWith("$"));
  const staticSegments = params === -1 ? segments : segments.slice(0, params);
  return staticSegments.length === 0 ? "/" : `/${staticSegments.join("/")}`;
}

export function MenuButton({
  children,
  to,
  activePath,
  ...props
}: MenuButtonProps) {
  const pathname = useLocation({
    select: (location) => location.pathname,
  });

  // Gegen den statischen Teil des Patterns prüfen, nicht gegen das Pattern selbst.
  // So bleibt der Eintrag auch bei abweichenden Params aktiv (anderer Monat im
  // Kalender) und unabhängig von den Search-Params des Links.
  const base = activePath ?? staticPrefix(String(to ?? "/"));
  const isActive =
    base === "/"
      ? pathname === "/"
      : pathname === base || pathname.startsWith(`${base}/`);

  return (
    <Link
      to={to}
      {...props}
      className={`btn btn-ghost py-6 hover:bg-accent/5 border-0 max-lg:min-w-24 ${isActive ? "text-primary" : ""}`}
    >
      {children}
    </Link>
  );
}
