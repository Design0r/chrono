import { Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { Header } from "../components/Header";
import { ToastProvider } from "../components/Toast";
import type { RouterContext } from "../main";

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
});

function RootComponent() {
  const { chrono } = Route.useRouteContext();
  return (
    <div>
      <ToastProvider>
        <Header chrono={chrono} />
        <div className="container h-full justify-center mx-auto ">
          <Outlet />
        </div>
        <TanStackRouterDevtools position="bottom-right" />
      </ToastProvider>
    </div>
  );
}
