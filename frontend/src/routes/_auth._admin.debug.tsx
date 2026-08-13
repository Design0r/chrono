import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/_admin/debug")({
  component: RouteComponent,
});

export function RouteComponent() {
  const { chrono } = Route.useRouteContext();

  const usersQ = useQuery({
    queryKey: ["users", "vacation"],
    queryFn: () =>
      chrono.users.getUsers({
        vacation: { year: new Date().getFullYear() },
      }),
    staleTime: 1000 * 60 * 30, // 30min
    gcTime: 1000 * 60 * 60 * 1, // 1h
    retry: false,
  });

  return (
    <div className="max-w-xl space-y-4 justify-center items m-auto flex flex-col">
      <div className="flex flex-col gap-2">
        <button className="btn btn-error">Clear token table</button>
        <button className="btn btn-error">
          Create tokens for accepted events
        </button>
        <button className="btn btn-error">Generate default user color</button>
        <button className="btn btn-error">Clear sessions table</button>
      </div>

      {!usersQ.isPending && (
        <div className="flex flex-col border-error border-2 rounded-2xl p-4 gap-2">
          <p>Change Password</p>
          <select className="w-full col-span-1 cursor-pointer bg-base-300 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-10 pl-4 text-base border-0 rounded-lg animate-all">
            {usersQ.data?.map((u) => (
              <option value={u.id}>{u.username}</option>
            ))}
          </select>
          <input
            className="input max-w-xl w-full"
            type="password"
            name="password"
          />
          <button className="btn btn-error" type="submit">
            Submit
          </button>
        </div>
      )}
      <div className="flex flex-col border-error border-2 rounded-2xl p-4 gap-2">
        <p>Delete Chrono Events By Name</p>
        <div className="max-w-xl flex w-full space-x-2">
          <input className="input w-full" type="text" name="eventName" />
          <button className="btn btn-error">Delete</button>
        </div>
      </div>
    </div>
  );
}
