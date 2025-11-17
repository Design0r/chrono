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
        year: new Date().getFullYear(),
      }),
    staleTime: 1000 * 60 * 30, // 30min
    gcTime: 1000 * 60 * 60 * 1, // 1h
  });

  return (
    <div className="container  max-w-xs space-y-2 mt-6 justify-center items m-auto flex flex-col">
      <button className="btn btn-warning">Clear token table</button>
      <button className="btn btn-warning">
        Create tokens for accepted events
      </button>
      <button className="btn btn-warning">Generate default user color</button>
      <button className="btn btn-warning">Clear sessions table</button>

      {!usersQ.isPending && (
        <div className="flex flex-col border-warning border-2 rounded-2xl p-4 gap-2">
          <p>Change Password</p>
          <select className="select h-full bg-base-300 hover:bg-base-100 transition-color text-lg rounded-xl">
            {usersQ.data?.map((u) => (
              <option value={u.id}>{u.username}</option>
            ))}
          </select>
          <input className="input" type="password" name="password" />
          <button className="btn btn-warning" type="submit">
            Submit
          </button>
        </div>
      )}
      <div className="flex flex-col border-warning border-2 rounded-2xl p-4 gap-2">
        <p>Delete Chrono Events By Name</p>
        <div>
          <input className="input" type="text" name="eventName" />
          <button className="btn btn-warning">Delete</button>
        </div>
      </div>
    </div>
  );
}
