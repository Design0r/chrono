import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth")({
  beforeLoad: ({ context }) => {
    if (!context.auth.isAuthenticated) {
      throw redirect({
        to: "/login",
      });
    }
  },
  component: AuthLayout,
});

function AuthLayout() {
  return (
    <div className="container h-full justify-center mx-auto ">
      <Outlet />
    </div>
  );
}
