import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import type { User, UserWithVacation } from "../types/auth";
import { hexToHSL, hsla } from "../utils/colors";

export const Route = createFileRoute("/_auth/team")({
  component: RouteComponent,
  loader: async ({ context: { chrono, queryClient } }) => {
    const data = queryClient.ensureQueryData({
      queryKey: ["users"],
      queryFn: () => chrono.users.getUsers({ year: new Date().getFullYear() }),
    });

    return data;
  },
});

function RouteComponent() {
  const users = Route.useLoaderData() as UserWithVacation[];
  const {
    auth: { user },
  } = Route.useRouteContext();

  console.log(users);

  return (
    <div className="p-2 my-2">
      <div className="overflow-x-auto bg-info/3">
        <table className="table table-zebra table-md">
          <thead>
            <tr className="bg-base-200 text-base-content border-b-1.5 border-primary/65">
              <th>ID</th>
              <th>Name</th>
              <th className="text-center">Vacation Days</th>
              <th className="text-center">Used</th>
              <th className="text-center">Remaining</th>
              <th>Email</th>
              <th>Role</th>
              <th>Enabled</th>
              {user?.is_superuser && (
                <>
                  <th></th>
                </>
              )}
            </tr>
          </thead>
          <tbody>
            {user &&
              users.map((u, i) => (
                <TableRow key={i} user={u} currUser={user} />
              ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function TableRow({
  user,
  currUser,
}: {
  user: UserWithVacation;
  currUser: User;
}) {
  const [edit, setEdit] = useState(false);

  const hsl = hexToHSL(user ? user.color : "#000");
  const borderColor = hsla(...hsl, 0.6);
  const bgColor = hsla(...hsl, 0.2);
  const textColor = hsla(...hsl, 1);

  return (
    <tr
      className={
        user.vacation_days > 0
          ? "hover:bg-primary/10 pt-8 border-b border-primary/25 text-base-content/80 hover:text-primary animate-color"
          : "bg-error/35 hover:bg-error/50 pt-8 border-2 border-error/70 text-base-content/80 animate-color"
      }
    >
      <th>
        <div tabIndex={0} role="button" className="avatar placeholder p-0 ">
          <div
            className="w-8 flex items-center justify-center border rounded-full text-neutral-content"
            style={{ backgroundColor: bgColor, borderColor: borderColor }}
          >
            <div className="text-sm text-center" style={{ color: textColor }}>
              {user.id}
            </div>
          </div>
        </div>
      </th>
      <th>{user.username}</th>
      <td className="text-center">{user.vacation_days}</td>
      <td className="text-center">{user.vacation_used}</td>
      <td className="text-center">{user.vacation_remaining}</td>
      <td>{user.email}</td>
      <td>
        <span>{user.role}</span>
      </td>

      <td className="icon-outlined">
        <span className="icon-outlined cursor-not-allowed ">
          {user.enabled ? "check_box" : "check_box_outline_blank"}
        </span>
      </td>
      {currUser.is_superuser && (
        <td>
          {edit ? (
            <>
              <button
                onClick={() => setEdit(false)}
                className="btn btn-soft btn-success animate-color icon-outlined"
              >
                check
              </button>
              <button
                onClick={() => setEdit(false)}
                className="btn btn-soft btn-error animate-color icon-outlined"
              >
                close
              </button>
            </>
          ) : (
            <button
              onClick={() => setEdit(true)}
              className="btn btn-soft btn-ghost animate-color icon-outlined"
            >
              edit
            </button>
          )}
        </td>
      )}
    </tr>
  );
}
