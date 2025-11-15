import {
  Link,
  useLocation,
  type LinkProps,
  type RegisteredRouter,
} from "@tanstack/react-router";
import { useAuth } from "../auth";
import { Avatar } from "./Avatar";
import { useQuery } from "@tanstack/react-query";
import type { ChronoClient } from "../api/chrono/client";
import { LoadingSpinner } from "./LoadingSpinner";
import { ErrorPage } from "./ErrorPage";

export function Header({ chrono }: { chrono: ChronoClient }) {
  const auth = useAuth();

  const userQ = useQuery({
    queryKey: ["user", auth.userId],
    queryFn: () => chrono.users.getUserById(auth.userId!),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
  });

  const settingQ = useQuery({
    queryKey: ["settings"],
    queryFn: () => chrono.settings.getSettings(),
    staleTime: 1000 * 30, // 30s
    gcTime: 1000 * 60 * 5, // 5m
  });

  const queries = [userQ, settingQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending) return <LoadingSpinner />;
  if (firstError) return <ErrorPage error={firstError} />;

  const date = new Date();

  return (
    <div className="mb-4 mx-auto p-4 lg:px-4">
      <div className="navbar flex justify-between">
        <div className="flex items-center">
          <div className="pr-14">
            <img className="w-40" alt="chrono logo" src="chrono.svg" />
          </div>

          {auth.isAuthenticated && (
            <div
              className="z-20! max-lg:dock max-lg:border-t max-lg:border-accent/15 max-lg:bg-base-100/50! backdrop-blur-xl overflow-x-auto flex gap-4 lg:w-fit 
						*:flex *:flex-col! *:lg:flex-row! *:lg:gap-2 *:lg:items-center"
            >
              <MenuButton to="/">
                <span className="icon-outlined">home</span>
                <span className="font-medium text-base">Home</span>
              </MenuButton>
              <MenuButton
                to="/calendar/$year/$month"
                params={{
                  year: date.getFullYear().toString(),
                  month: (date.getMonth() + 1).toString(),
                }}
              >
                <span className="icon-outlined">calendar_today</span>
                <span className="font-medium text-base">Calendar</span>
              </MenuButton>
              <MenuButton to="/team">
                <span className="icon-outlined">group</span>
                <span className="font-medium text-base">Team</span>
              </MenuButton>
              {userQ.data!.is_superuser && (
                <>
                  <MenuButton to="/requests">
                    <span className="icon-outlined">mark_chat_unread</span>
                    <span className="font-medium text-base">Requests</span>
                  </MenuButton>
                  <MenuButton to="/tokens">
                    <span className="icon-outlined">local_activity</span>
                    <span className="font-medium text-base">Tokens</span>
                  </MenuButton>
                  <MenuButton to="/settings">
                    <span className="icon-outlined">settings</span>
                    <span className="font-medium text-base">Settings</span>
                  </MenuButton>
                  <MenuButton to="/export">
                    <span className="icon-outlined">file_export</span>
                    <span className="font-medium text-base">Export</span>
                  </MenuButton>
                </>
              )}
            </div>
          )}
        </div>
        <div className="flex items-center justify-end gap-6">
          {!auth.isAuthenticated ? (
            <>
              <a href="/login" className="btn btn-ghost">
                Login
              </a>
              {settingQ.data!.signup_enabled && (
                <a href="/signup" className="btn btn-ghost">
                  Signup
                </a>
              )}
            </>
          ) : (
            <Avatar user={userQ.data!} />
          )}
        </div>
      </div>
    </div>
  );
}

interface MenuButtonProps extends LinkProps<RegisteredRouter> {
  children?: React.ReactNode | React.ReactNode[];
}

export function MenuButton({ children, to, ...props }: MenuButtonProps) {
  const pathname = useLocation({
    select: (location) => location.pathname,
  });

  return (
    <Link
      to={to}
      {...props}
      className={`btn btn-ghost py-6 hover:bg-accent/5 border-0 max-lg:min-w-24 ${pathname === to && "text-primary"}`}
    >
      {children && Array.isArray(children) ? (
        <>{...children}</>
      ) : (
        <>{children}</>
      )}
    </Link>
  );
}
