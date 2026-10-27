import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { chrono } from "../api/chrono/client";
import { userQueries } from "../api/queries";
import { useToast } from "../components/Toast";
import { useAppForm } from "../lib/form";

export const Route = createFileRoute("/_auth/_admin/tokens")({
  loader: ({ context: { queryClient } }) => queryClient.ensureQueryData(userQueries.list()),
  component: RouteComponent,
});

function RouteComponent() {
  const queryClient = useQueryClient();
  const { addToast, addErrorToast } = useToast();

  const { data: users } = useSuspenseQuery(userQueries.list());

  const mutation = useMutation({
    mutationFn: ({ userId, tokens }: { userId: number; tokens: number }) =>
      chrono.tokens.createTokens(userId, tokens),
    onSuccess: (_, { tokens }) => {
      addToast(`Added ${tokens} tokens`, "success");
      return queryClient.invalidateQueries({ queryKey: userQueries.all() });
    },
    onError: (error) => addErrorToast(error),
  });

  const form = useAppForm({
    defaultValues: { userId: String(users[0]?.id ?? ""), tokens: 0 },
    onSubmit: ({ value }) =>
      mutation.mutateAsync({ userId: Number(value.userId), tokens: value.tokens }).catch(() => {}),
  });

  return (
    <form.AppForm>
      <form.Form className="mx-auto flex w-full max-w-xl flex-col space-y-2">
        <form.AppField name="userId">
          {(field) => (
            <field.Select
              required
              aria-label="User"
              className="w-full cursor-pointer bg-base-300 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-10 text-base border-0 rounded-md animate-all"
            >
              {users.map((u) => (
                <option key={u.id} label={u.username} value={u.id}>
                  {u.username}
                </option>
              ))}
            </field.Select>
          )}
        </form.AppField>
        <form.AppField name="tokens">
          {(field) => (
            <field.NumberInput
              aria-label="Tokens"
              className="input w-full input-bordered"
              step="0.5"
              required
            />
          )}
        </form.AppField>
        <form.SubmitButton className="btn text-white btn-primary bg-primary/80 hover:bg-primary animate-color">
          Add Token
        </form.SubmitButton>
      </form.Form>
    </form.AppForm>
  );
}
