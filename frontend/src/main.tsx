import { QueryClient } from "@tanstack/react-query";
import * as TanStackQueryProvider from "./integrations/tanstack-query/root-provider.tsx";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { ChronoClient } from "./api/chrono/client";
import { AuthProvider, useAuth, type AuthContext } from "./auth";
import { ErrorPage } from "./components/ErrorPage";
import "./css/index.css";
import { routeTree } from "./routeTree.gen";

const chrono = new ChronoClient();

export type RouterContext = {
  auth: ReturnType<typeof useAuth>;
  queryClient: QueryClient;
  chrono: ChronoClient;
};

const TanStackQueryProviderContext = TanStackQueryProvider.getContext();
// Set up a Router instance
const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  scrollRestoration: true,
  context: {
    auth: undefined! as unknown as AuthContext,
    chrono: chrono,
    ...TanStackQueryProviderContext,
  },
  defaultStructuralSharing: true,
  defaultPreloadStaleTime: 0,

  defaultNotFoundComponent: () => (
    <ErrorPage error={{ name: "404", message: "Not Found" }} />
  ),

  defaultErrorComponent: () => (
    <ErrorPage
      error={{
        name: "",
        message: "I dont know what you did, but it wasn't good",
      }}
    />
  ),
});

// Register things for typesafety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function InnerApp() {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

function App() {
  return (
    <TanStackQueryProvider.Provider {...TanStackQueryProviderContext}>
      <AuthProvider>
        <InnerApp />
      </AuthProvider>
    </TanStackQueryProvider.Provider>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
