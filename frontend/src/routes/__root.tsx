import { Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { Header } from "../components/Header";
import type { RouterContext } from "../main";

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
});

function RootComponent() {
  const { chrono } = Route.useRouteContext();
  return (
    <div>
      <Header chrono={chrono} />
      <Outlet />
      <TanStackRouterDevtools position="bottom-right" />
    </div>
  );
}
