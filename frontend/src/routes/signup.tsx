import { useMutation } from "@tanstack/react-query";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useAuth } from "../auth";
import { useToast } from "../components/Toast";
import { useAppForm } from "../lib/form";
import type { SignupRequest } from "../types/auth";

export const Route = createFileRoute("/signup")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const auth = useAuth();
  const { addToast, addErrorToast } = useToast();

  const mutation = useMutation({
    mutationFn: (data: SignupRequest) => auth.signup(data),
    onSuccess: async () => {
      addToast("Successfully signed up", "success");
      await router.invalidate();
      await router.navigate({ to: "/" });
    },
    onError: (error) => addErrorToast(error),
  });

  const form = useAppForm({
    defaultValues: { username: "", email: "", password: "" } satisfies SignupRequest,
    onSubmit: ({ value }) => mutation.mutateAsync(value).catch(() => {}),
  });

  return (
    <div className="flex my-10">
      <div className="align-middle flex m-auto">
        <div>
          <h1 className="font-bold text-xl">Sign up</h1>
          <br />
          <form.AppForm>
            <form.Form className="w-xs md:w-max">
              <div className="w-xs md:w-lg">
                <label htmlFor="username">Username</label>
                <br />
                <form.AppField name="username">
                  {(field) => (
                    <field.TextInput
                      className="input w-full input-bordered"
                      autoComplete="username"
                      required
                    />
                  )}
                </form.AppField>
                <br />
                <br />
              </div>
              <div>
                <label htmlFor="email">Email</label>
                <br />
                <form.AppField name="email">
                  {(field) => (
                    <field.TextInput
                      className="input w-full input-bordered"
                      type="email"
                      autoComplete="email"
                      required
                    />
                  )}
                </form.AppField>
                <br />
                <br />
              </div>
              <div>
                <label htmlFor="password">Password</label>
                <br />
                <form.AppField name="password">
                  {(field) => (
                    <field.TextInput
                      className="input w-full input-bordered"
                      type="password"
                      autoComplete="new-password"
                      required
                    />
                  )}
                </form.AppField>
                <br />
                <br />
              </div>
              <form.SubmitButton className="btn text-white btn-primary bg-primary/80 hover:bg-primary animate-color">
                Sign up
              </form.SubmitButton>
            </form.Form>
          </form.AppForm>
        </div>
      </div>
    </div>
  );
}
