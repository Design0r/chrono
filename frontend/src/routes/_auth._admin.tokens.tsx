import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useRef, useState } from "react";
import { chrono } from "../api/chrono/client";
import { userQueries } from "../api/queries";
import { useToast } from "../components/Toast";

export const Route = createFileRoute("/_auth/_admin/tokens")({
  loader: ({ context: { queryClient } }) => queryClient.ensureQueryData(userQueries.list()),
  component: RouteComponent,
});

function RouteComponent() {
  const selectRef = useRef<HTMLSelectElement>(null!);
  const queryClient = useQueryClient();
  const { addToast, addErrorToast } = useToast();

  const { data: users } = useSuspenseQuery(userQueries.list());

  const [token, setToken] = useState(0);

  const mutation = useMutation({
    mutationFn: (token: number) =>
      chrono.tokens.createTokens(Number.parseInt(selectRef.current.value), token),
    onSuccess: (_, variables) => {
      addToast(`Added ${variables} tokens`, "success");
      return queryClient.invalidateQueries({ queryKey: userQueries.all() });
    },
    onError: (error) => addErrorToast(error),
  });

  return (
    <div className="mx-auto flex w-full max-w-xl flex-col space-y-2">
      <select
        ref={selectRef}
        className="w-full cursor-pointer bg-base-300 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-10 text-base border-0 rounded-md animate-all"
        name="filter"
      >
        {users.map((u) => (
          <option key={u.id} label={u.username} value={u.id}>
            {u.username}
          </option>
        ))}
      </select>
      <label>
        <input
          className="input w-full input-bordered"
          name="token"
          type="number"
          step="0.5"
          defaultValue={0}
          onChange={(e) => setToken(Number.parseFloat(e.target.value))}
        />
      </label>
      <button
        onClick={() => mutation.mutate(token)}
        className="btn text-white btn-primary bg-primary/80 hover:bg-primary animate-color"
      >
        Add Token
      </button>
    </div>
  );
}
