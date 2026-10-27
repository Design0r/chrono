import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { currentUserId, userQueries } from "../api/queries";

export const Route = createFileRoute("/_auth/_admin")({
  beforeLoad: async ({ context: { auth, queryClient } }) => {
    const user = await queryClient.ensureQueryData(userQueries.detail(currentUserId(auth.userId)));
    if (!user.is_superuser) {
      throw redirect({
        to: "/",
      });
    }
  },
  component: AdminLayout,
});

function AdminLayout() {
  return <Outlet />;
}
