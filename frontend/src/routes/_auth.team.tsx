import { useMutation, useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { ErrorPage } from "../components/ErrorPage";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { useToast } from "../components/Toast";
import type { User, UserWithVacation } from "../types/auth";
import type { TeamEditForm } from "../types/forms";
import { hexToHSL, hsla } from "../utils/colors";
import { capitalize } from "../utils/string";

export const Route = createFileRoute("/_auth/team")({
  component: TeamComponent,
});

function TeamComponent() {
  const { chrono, auth } = Route.useRouteContext();

  const usersQ = useQuery({
    queryKey: ["users", "vacation"],
    queryFn: () =>
      chrono.users.getUsers({
        year: new Date().getFullYear(),
      }),
    staleTime: 1000 * 60 * 30, // 30min
    gcTime: 1000 * 60 * 60 * 1, // 1h
    retry: false,
  });

  const currUserQ = useQuery({
    queryKey: ["user", auth.userId],
    queryFn: () => chrono.users.getUserById(auth.userId!),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
    retry: false,
  });

  const queries = [usersQ, currUserQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending) return <LoadingSpinnerPage />;
  if (firstError) return <ErrorPage error={firstError} />;

  const users = usersQ.data! as UserWithVacation[];
  const user = currUserQ.data! as UserWithVacation;

  return (
    <div className="p-2 my-2">
      <div className="overflow-x-auto bg-info/3">
        <table className="table table-zebra table-md">
          <thead>
            <tr className="text-base-content border-b-1.5 border-primary/65">
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
  const { chrono, queryClient } = Route.useRouteContext();
  const { addToast, addErrorToast } = useToast();
  const [vacDays, setVacDays] = useState(user.vacation_days);
  const [role, setRole] = useState(user.role);
  const [enabled, setEnabled] = useState(user.enabled);

  const mutation = useMutation({
    mutationKey: ["user", "update", user.id],
    mutationFn: ({ userId, data }: { userId: number; data: TeamEditForm }) =>
      chrono.users.updateUser(userId, data),
    onSuccess: () => {
      addToast("Successfully updated profile settings", "success");
      setEdit(false);
      return queryClient.invalidateQueries({ queryKey: ["users"] });
    },
    onError: (error) => addErrorToast(error),
    retry: false,
  });

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
      <td className="text-center">
        {edit ? (
          <input
            className="input input-bordered"
            type="number"
            step={0.5}
            required
            defaultValue={user.vacation_days}
            onChange={(e) => setVacDays(Number(e.target.value))}
          />
        ) : (
          <>{user.vacation_days}</>
        )}
      </td>
      <td className="text-center">{user.vacation_used}</td>
      <td className="text-center">{user.vacation_remaining}</td>
      <td>{user.email}</td>
      <td>
        {edit ? (
          <select
            required
            defaultValue={user.role}
            onChange={(e) => setRole(e.target.value)}
          >
            {chrono.users.getRoles().map((u, i) => (
              <option key={i} value={u}>
                {capitalize(u)}
              </option>
            ))}
          </select>
        ) : (
          <>{user.role}</>
        )}
      </td>

      <td className="icon-outlined">
        {edit ? (
          <input
            type="checkbox"
            className="checkbox"
            onChange={(e) => setEnabled(Boolean(e.target.value))}
            defaultChecked={user.enabled}
          />
        ) : (
          <input
            type="checkbox"
            className="checkbox cursor-not-allowed"
            defaultChecked={user.enabled}
            contentEditable={false}
          />
        )}
      </td>
      {currUser.is_superuser && (
        <td>
          {edit ? (
            <>
              <button
                onClick={() => {
                  mutation.mutate({
                    userId: user.id,
                    data: { vacation_days: vacDays, enabled, role },
                  });
                  setEdit(false);
                }}
                className="btn btn-soft btn-success animate-color icon-outlined"
              >
                check
              </button>
              <button
                onClick={() => {
                  setEdit(false);
                  setRole(user.role);
                  setVacDays(user.vacation_days);
                  setEnabled(user.enabled);
                }}
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
