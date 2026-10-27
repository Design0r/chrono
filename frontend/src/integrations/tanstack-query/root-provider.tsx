import { QueryClientProvider, type QueryClient } from "@tanstack/react-query";
import { createQueryClient } from "../../api/queries";

export function getContext() {
  const queryClient = createQueryClient();
  return {
    queryClient,
  };
}

export function Provider({
  children,
  queryClient,
}: {
  children: React.ReactNode;
  queryClient: QueryClient;
}) {
  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
}
