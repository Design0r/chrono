import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { ErrorPage } from "../components/ErrorPage";
import { useRef, useState } from "react";

export const Route = createFileRoute("/_auth/_admin/tokens")({
  component: RouteComponent,
});

function RouteComponent() {
  const { chrono } = Route.useRouteContext();
  const selectRef = useRef<HTMLSelectElement>(null!);
  const queryClient = useQueryClient();

  const usersQ = useQuery({
    queryKey: ["users", "vacation"],
    queryFn: () => chrono.users.getUsers(),
    staleTime: 1000 * 60 * 30, // 30min
    gcTime: 1000 * 60 * 60 * 1, // 1h
  });

  const [token, setToken] = useState(0);

  const mutation = useMutation({
    mutationKey: ["tokens"],
    mutationFn: (token: number) =>
      chrono.tokens.createTokens(
        Number.parseInt(selectRef.current.value),
        token,
      ),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["users"] }),
  });

  if (usersQ.isPending) return <LoadingSpinnerPage />;
  if (usersQ.isError) return <ErrorPage error={usersQ.error} />;

  return (
    <div className="flex my-10">
      <div className="align-middle flex m-auto">
        <div>
          <div className="flex flex-col space-y-2">
            <select
              ref={selectRef}
              className="col-span-1 select min-w-56 min-h-14 h-full focus:border-0 border-0 bg-base-300 hover:bg-base-300 transition-color text-lg rounded-xl"
              name="filter"
            >
              {usersQ.data!.map((u) => (
                <option key={u.id} value={u.id}>
                  {u.username}
                </option>
              ))}
            </select>
            <label>
              <input
                className="input input-bordered"
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
        </div>
      </div>
    </div>
  );
}
