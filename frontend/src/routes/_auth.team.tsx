import { useMutation, useQuery, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import {
  columnVisibilityFeature,
  createColumnHelper,
  createSortedRowModel,
  metaHelper,
  rowSortingFeature,
  sortFn_alphanumeric,
  sortFn_text,
  tableFeatures,
  useTable,
} from "@tanstack/react-table";
import { useMemo, useState } from "react";
import { chrono } from "../api/chrono/client";
import { aworkQueries, currentUserId, userQueries } from "../api/queries";
import { useToast } from "../components/Toast";
import type { User, UserWithVacation } from "../types/auth";
import type { TeamEditForm } from "../types/forms";
import type { AworkUser } from "../types/response";
import { hexToHSL, hsla } from "../utils/colors";
import { capitalize } from "../utils/string";

export const Route = createFileRoute("/_auth/team")({
  loader: async ({ context: { auth, queryClient } }) => {
    const userId = currentUserId(auth.userId);
    const year = new Date().getFullYear();
    const [currUser] = await Promise.all([
      queryClient.ensureQueryData(userQueries.detail(userId)),
      queryClient.ensureQueryData(userQueries.listWithVacation(year, { includeInactive: true })),
    ]);
    if (currUser.is_superuser) void queryClient.prefetchQuery(aworkQueries.users());
    return { userId, year };
  },
  component: TeamComponent,
});

const noAworkUsers: AworkUser[] = [];

function TeamComponent() {
  const { userId, year } = Route.useLoaderData();

  const { data: users } = useSuspenseQuery(
    userQueries.listWithVacation(year, { includeInactive: true }),
  );
  const { data: currUser } = useSuspenseQuery(userQueries.detail(userId));
  const aworkUsersQ = useQuery({ ...aworkQueries.users(), enabled: currUser.is_superuser });
  const aworkUsers = aworkUsersQ.data ?? noAworkUsers;

  return (
    <div className="pb-8">
      <TeamTable
        users={users.filter((u) => u.enabled)}
        currUser={currUser}
        aworkUsers={aworkUsers}
      />
      {currUser.is_superuser && (
        <div className="flex flex-col gap-8 pt-20">
          <hr />
          <h1 className="text-xl">Inactive Users</h1>
          <TeamTable
            users={users.filter((u) => !u.enabled)}
            currUser={currUser}
            aworkUsers={aworkUsers}
          />
        </div>
      )}
    </div>
  );
}

type TeamRow = UserWithVacation & {
  awork_name: string;
  draft: TeamEditForm | null;
};

type TeamTableMeta = {
  roles: readonly string[];
  aworkUsers: AworkUser[];
  savingUserId: number | null;
  startEdit: (user: TeamRow) => void;
  updateDraft: (userId: number, patch: Partial<TeamEditForm>) => void;
  cancelEdit: (userId: number) => void;
  save: (userId: number, draft: TeamEditForm) => void;
};

type TeamColumnMeta = {
  className?: string;
  cellClassName?: string;
  rowHeader?: boolean;
};

const features = tableFeatures({
  rowSortingFeature,
  columnVisibilityFeature,
  sortedRowModel: createSortedRowModel(),
  sortFns: { alphanumeric: sortFn_alphanumeric, text: sortFn_text },
  tableMeta: metaHelper<TeamTableMeta>(),
  columnMeta: metaHelper<TeamColumnMeta>(),
});

const column = createColumnHelper<typeof features, TeamRow>();

const selectClassName =
  "w-full col-span-1 cursor-pointer bg-base-300 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-full pl-4 text-base border-0 rounded-md animate-all";

const columns = column.columns([
  column.accessor("id", {
    header: "ID",
    meta: { rowHeader: true },
    cell: ({ row }) => <UserIdAvatar user={row.original} />,
  }),
  column.accessor("username", {
    header: "Name",
    meta: { rowHeader: true },
  }),
  column.accessor("vacation_days", {
    header: "Vacation Days",
    meta: { className: "text-center" },
    cell: ({ row, getValue, table }) => {
      const { draft, id } = row.original;
      if (!draft) return getValue();
      return (
        <input
          className="input input-bordered"
          type="number"
          step={0.5}
          required
          defaultValue={draft.vacation_days}
          onChange={(e) =>
            table.options.meta!.updateDraft(id, { vacation_days: Number(e.target.value) })
          }
        />
      );
    },
  }),
  column.accessor("vacation_used", {
    header: "Used",
    meta: { className: "text-center" },
  }),
  column.accessor("vacation_remaining", {
    header: "Remaining",
    meta: { className: "text-center" },
  }),
  column.accessor("email", {
    header: "Email",
  }),
  column.accessor("role", {
    header: "Role",
    cell: ({ row, getValue, table }) => {
      const { draft, id } = row.original;
      if (!draft) return getValue();
      return (
        <select
          className={selectClassName}
          required
          value={draft.role}
          onChange={(e) => table.options.meta!.updateDraft(id, { role: e.target.value })}
        >
          {table.options.meta!.roles.map((r) => (
            <option key={r} value={r}>
              {capitalize(r)}
            </option>
          ))}
        </select>
      );
    },
  }),
  column.accessor("enabled", {
    header: "Enabled",
    meta: { cellClassName: "icon-outlined" },
    cell: ({ row, getValue, table }) => {
      const { draft, id } = row.original;
      if (!draft) {
        return (
          <input
            type="checkbox"
            className="checkbox cursor-not-allowed"
            checked={getValue()}
            readOnly
          />
        );
      }
      return (
        <input
          type="checkbox"
          className="checkbox"
          checked={draft.enabled}
          onChange={(e) => table.options.meta!.updateDraft(id, { enabled: e.target.checked })}
        />
      );
    },
  }),
  column.accessor("awork_name", {
    id: "awork",
    header: "Awork",
    cell: ({ row, getValue, table }) => {
      const { draft, id } = row.original;
      if (!draft) return getValue();
      return (
        <select
          className={selectClassName}
          value={draft.awork_id}
          onChange={(e) => table.options.meta!.updateDraft(id, { awork_id: e.target.value })}
        >
          <option value="">-</option>
          {table.options.meta!.aworkUsers.map((a) => (
            <option key={a.id} value={a.id}>
              {a.firstName} {a.lastName}
            </option>
          ))}
        </select>
      );
    },
  }),
  column.display({
    id: "actions",
    header: "",
    cell: ({ row, table }) => <RowActions user={row.original} meta={table.options.meta!} />,
  }),
]);

function TeamTable({
  users,
  currUser,
  aworkUsers,
}: {
  users: UserWithVacation[];
  currUser: User;
  aworkUsers: AworkUser[];
}) {
  const queryClient = useQueryClient();
  const { addToast, addErrorToast } = useToast();
  const [drafts, setDrafts] = useState<Record<number, TeamEditForm>>({});

  const cancelEdit = (userId: number) =>
    setDrafts((d) => Object.fromEntries(Object.entries(d).filter(([id]) => Number(id) !== userId)));

  const mutation = useMutation({
    mutationFn: ({ userId, data }: { userId: number; data: TeamEditForm }) =>
      chrono.users.updateUser(userId, data),
    onSuccess: async (_, { userId }) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: userQueries.all() }),
        queryClient.invalidateQueries({ queryKey: aworkQueries.workTimes() }),
      ]);
      cancelEdit(userId);
      addToast("Successfully updated profile settings", "success");
    },
    onError: (error) => addErrorToast(error),
  });

  const data = useMemo(
    () =>
      users.map((u) => ({
        ...u,
        awork_name: aworkName(aworkUsers, u.awork_id),
        draft: drafts[u.id] ?? null,
      })),
    [users, aworkUsers, drafts],
  );

  const table = useTable({
    features,
    columns,
    data,
    getRowId: (u) => String(u.id),
    state: {
      columnVisibility: {
        awork: currUser.is_superuser,
        actions: currUser.is_superuser,
      },
    },
    meta: {
      roles: chrono.users.getRoles(),
      aworkUsers,
      savingUserId: mutation.isPending ? (mutation.variables?.userId ?? null) : null,
      startEdit: (u) =>
        setDrafts((d) => ({
          ...d,
          [u.id]: {
            vacation_days: u.vacation_days,
            role: u.role,
            enabled: u.enabled,
            awork_id: u.awork_id || "",
          },
        })),
      updateDraft: (userId, patch) =>
        setDrafts((d) => (d[userId] ? { ...d, [userId]: { ...d[userId], ...patch } } : d)),
      cancelEdit,
      save: (userId, draft) => mutation.mutate({ userId, data: draft }),
    },
  });

  return (
    <div className="overflow-x-auto bg-info/3">
      <table className="table table-zebra table-md">
        <thead>
          {table.getHeaderGroups().map((group) => (
            <tr key={group.id} className="text-base-content border-b-1.5 border-primary/65">
              {group.headers.map((header) => {
                const sorted = header.column.getIsSorted();
                return (
                  <th
                    key={header.id}
                    className={header.column.columnDef.meta?.className}
                    aria-sort={
                      sorted === "asc" ? "ascending" : sorted === "desc" ? "descending" : undefined
                    }
                  >
                    {header.isPlaceholder ? null : header.column.getCanSort() ? (
                      <button
                        type="button"
                        className="inline-flex items-center gap-1 cursor-pointer hover:text-primary animate-color"
                        onClick={header.column.getToggleSortingHandler()}
                      >
                        <table.FlexRender header={header} />
                        <span
                          className={`icon-outlined text-base! ${sorted ? "" : "invisible"}`}
                          aria-hidden
                        >
                          {sorted === "desc" ? "arrow_downward" : "arrow_upward"}
                        </span>
                      </button>
                    ) : (
                      <table.FlexRender header={header} />
                    )}
                  </th>
                );
              })}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr
              key={row.id}
              className={
                row.original.vacation_days > 0
                  ? "hover:bg-primary/10 pt-8 border-b border-primary/25 text-base-content/80 hover:text-primary animate-color"
                  : "bg-error/35 hover:bg-error/50 pt-8 border-2 border-error/70 text-base-content/80 animate-color"
              }
            >
              {row.getVisibleCells().map((cell) => {
                const { className, cellClassName, rowHeader } = cell.column.columnDef.meta ?? {};
                const Cell = rowHeader ? "th" : "td";
                return (
                  <Cell
                    key={cell.id}
                    className={[className, cellClassName].filter(Boolean).join(" ") || undefined}
                  >
                    <table.FlexRender cell={cell} />
                  </Cell>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function UserIdAvatar({ user }: { user: User }) {
  const hsl = hexToHSL(user.color);
  return (
    <div tabIndex={0} role="button" className="avatar placeholder p-0 ">
      <div
        className="w-8 flex items-center justify-center border rounded-full text-neutral-content"
        style={{ backgroundColor: hsla(...hsl, 0.2), borderColor: hsla(...hsl, 0.6) }}
      >
        <div className="text-sm text-center" style={{ color: hsla(...hsl, 1) }}>
          {user.id}
        </div>
      </div>
    </div>
  );
}

function RowActions({ user, meta }: { user: TeamRow; meta: TeamTableMeta }) {
  if (!user.draft) {
    return (
      <button
        onClick={() => meta.startEdit(user)}
        className="btn btn-soft btn-ghost animate-color icon-outlined"
      >
        edit
      </button>
    );
  }
  const draft = user.draft;
  return (
    <>
      <button
        onClick={() => meta.save(user.id, draft)}
        disabled={meta.savingUserId === user.id}
        className="btn btn-soft btn-success animate-color icon-outlined"
      >
        check
      </button>
      <button
        onClick={() => meta.cancelEdit(user.id)}
        className="btn btn-soft btn-error animate-color icon-outlined"
      >
        close
      </button>
    </>
  );
}

function aworkName(aworkUsers: AworkUser[], aworkId: string | null) {
  if (!aworkId) return "-";
  const a = aworkUsers.find((a) => a.id === aworkId);
  return a ? `${a.firstName} ${a.lastName}` : aworkId;
}
