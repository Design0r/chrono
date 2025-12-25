import { Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { Header } from "../components/Header";
import { ToastProvider } from "../components/Toast";
import type { RouterContext } from "../main";
import TanStackQueryDevtools from "../integrations/tanstack-query/devtools";
import { TanStackDevtools } from "@tanstack/react-devtools";

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
        <TanStackDevtools
          config={{
            position: "bottom-right",
          }}
          plugins={[
            {
              name: "Tanstack Router",
              render: <TanStackRouterDevtoolsPanel />,
            },
            TanStackQueryDevtools,
          ]}
        />
      </ToastProvider>
    </div>
  );
}
