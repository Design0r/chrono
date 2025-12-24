import { useMutation } from "@tanstack/react-query";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { useAuth } from "../auth";
import { LoadingSpinner } from "../components/LoadingSpinner";
import { useToast } from "../components/Toast";
import type { SignupRequest } from "../types/auth";

export const Route = createFileRoute("/signup")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const auth = useAuth();
  const { register, handleSubmit } = useForm<SignupRequest>();
  const { addToast, addErrorToast } = useToast();

  const mutation = useMutation({
    mutationFn: (data: SignupRequest) => auth.signup(data),
    onSuccess: async () => {
      addToast("Successfully signed up", "success");
      await router.invalidate();
      await router.navigate({ to: "/" });
    },
    onError: (error) => addErrorToast(error),
    retry: false,
  });

  return (
    <div className="flex my-10">
      <div className="align-middle flex m-auto">
        <div>
          <h1 className="font-bold text-xl">Sign up</h1>
          <br />
          <form
            className="w-xs md:w-max"
            onSubmit={handleSubmit((data: SignupRequest) =>
              mutation.mutate(data),
            )}
          >
            <div className="w-xs md:w-lg">
              <label htmlFor="username">Username</label>
              <br />
              <input
                className="input w-full input-bordered"
                type="text"
                required
                {...register("username")}
              />
              <br />
              <br />
            </div>
            <div>
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
              {mutation.isPending ? <LoadingSpinner /> : "Sign up"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
