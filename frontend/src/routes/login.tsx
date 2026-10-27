import { useMutation } from "@tanstack/react-query";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useAuth } from "../auth";
import { useToast } from "../components/Toast";
import { useAppForm } from "../lib/form";
import type { LoginRequest } from "../types/auth";

export const Route = createFileRoute("/login")({
  component: LoginComponent,
});

function LoginComponent() {
  const router = useRouter();
  const auth = useAuth();
  const { addToast, addErrorToast } = useToast();
  const mutation = useMutation({
    mutationFn: async (data: LoginRequest) => await auth.login(data),
    onSuccess: async () => {
      addToast("Successfully logged in", "success");
      await router.invalidate();
      await router.navigate({ to: "/" });
    },
    onError: (error) => addErrorToast(error),
  });

  const form = useAppForm({
    defaultValues: { email: "", password: "" } satisfies LoginRequest,
    onSubmit: ({ value }) => mutation.mutateAsync(value).catch(() => {}),
  });

  return (
    <div className="flex my-10">
      <div className="align-middle flex m-auto">
        <div>
          <h1 className="font-bold text-xl">Log in</h1>
          <br />
          <form.AppForm>
            <form.Form className="w-xs md:w-max">
              <div className="w-xs md:w-lg">
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
                      autoComplete="current-password"
                      required
                    />
                  )}
                </form.AppField>
                <br />
                <br />
              </div>
              <form.SubmitButton className="btn text-white btn-primary bg-primary/80 hover:bg-primary animate-color">
                Log in
              </form.SubmitButton>
            </form.Form>
          </form.AppForm>
        </div>
      </div>
    </div>
  );
}
