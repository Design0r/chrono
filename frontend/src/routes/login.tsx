import { useMutation } from "@tanstack/react-query";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useAuth } from "../auth";
import { useForm } from "react-hook-form";
import type { LoginRequest } from "../types/auth";

export const Route = createFileRoute("/login")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const auth = useAuth();
  const { register, handleSubmit } = useForm<LoginRequest>();
  const mutation = useMutation({
    mutationFn: auth.login,
    onSuccess: async () => {
      console.log("hello");
      await router.invalidate();
      router.navigate({ to: "." });
    },
  });

  return (
    <div className="flex my-10">
      <div className="align-middle flex m-auto">
        <div>
          <h1 className="font-bold text-xl">Log in</h1>
          <br />
          <form
            className="w-max"
            onSubmit={handleSubmit((data: LoginRequest) =>
              mutation.mutate(data),
            )}
          >
            <div className="w-lg">
              <label htmlFor="email">Email</label>
              <br />
              <input
                className="input w-full input-bordered"
                type="email"
                required
                {...register("email")}
              />
              <br />
              <br />
            </div>
            <div>
              <label htmlFor="password">Password</label>
              <br />
              <input
                className="input w-full input-bordered"
                type="password"
                {...register("password")}
                required
              />
              <br />
              <br />
            </div>
            <button
              className="btn text-white btn-primary bg-primary/80 hover:bg-primary animate-color"
              type="submit"
              disabled={mutation.isPending}
            >
              {mutation.isPending ? "Logging in..." : "Log in"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
