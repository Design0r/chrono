import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { ChronoClient } from "./api/chrono/client";
import { AuthProvider, useAuth, type AuthContext } from "./auth";
import { ErrorPage } from "./components/ErrorPage";
import "./css/index.css";
import { routeTree } from "./routeTree.gen";

const queryClient = new QueryClient();
const chrono = new ChronoClient();

export type RouterContext = {
  auth: ReturnType<typeof useAuth>;
  queryClient: QueryClient;
  chrono: ChronoClient;
};

// Set up a Router instance
const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  scrollRestoration: true,
  context: {
    auth: undefined! as unknown as AuthContext,
    chrono: chrono,
    queryClient: queryClient,
  },

  defaultNotFoundComponent: () => (
    <ErrorPage error={{ name: "404", message: "Not Found" }} />
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
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <ReactQueryDevtools />
        <InnerApp />
      </AuthProvider>
    </QueryClientProvider>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
