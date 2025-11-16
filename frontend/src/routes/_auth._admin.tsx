import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/_admin")({
  beforeLoad: async ({ context }) => {
    const user = await context.auth.getUser();
    console.log(user);
    if (!user?.is_superuser) {
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
