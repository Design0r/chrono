import {
  Link,
  useLocation,
  type LinkProps,
  type RegisteredRouter,
} from "@tanstack/react-router";
import { useAuth } from "../auth";
import { Avatar } from "./Avatar";

export function Header() {
  const auth = useAuth();

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
              <MenuButton to="/calendar">
                <span className="icon-outlined">calendar_today</span>
                <span className="font-medium text-base">Calendar</span>
              </MenuButton>
              <MenuButton to="/team">
                <span className="icon-outlined">group</span>
                <span className="font-medium text-base">Team</span>
              </MenuButton>
            </div>
          )}
        </div>
        <div className="flex items-center justify-end gap-6">
          {!auth.isAuthenticated ? (
            <>
              <a href="/login" className="btn btn-ghost">
                Login
              </a>
              <a href="/signup" className="btn btn-ghost">
                Signup
              </a>
            </>
          ) : (
            <Avatar user={auth.user} />
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
      className={`btn btn-ghost py-6 hover:bg-accent/5 border-0 max-lg:min-w-24 ${pathname.includes(to!) && "text-primary"}`}
    >
      {children && Array.isArray(children) ? (
        <>{...children}</>
      ) : (
        <>{children}</>
      )}
    </Link>
  );
}
