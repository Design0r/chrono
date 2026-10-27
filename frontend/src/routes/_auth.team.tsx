import { useMutation, useQuery, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import {
  FlexRender,
  columnVisibilityFeature,
  createColumnHelper,
  createSortedRowModel,
  metaHelper,
  rowSortingFeature,
  sortFn_alphanumeric,
  sortFn_text,
  tableFeatures,
  useTable,
  type Row,
} from "@tanstack/react-table";
import { createContext, useContext, useMemo, useState } from "react";
import { chrono } from "../api/chrono/client";
import { aworkQueries, currentUserId, userQueries } from "../api/queries";
import { useToast } from "../components/Toast";
import { useAppForm } from "../lib/form";
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
};

type TeamTableMeta = {
  roles: readonly string[];
  aworkUsers: AworkUser[];
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
    cell: ({ getValue }) => <VacationDaysCell value={getValue()} />,
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
    cell: ({ getValue, table }) => (
      <RoleCell value={getValue()} roles={table.options.meta!.roles} />
    ),
  }),
  column.accessor("enabled", {
    header: "Enabled",
    meta: { cellClassName: "icon-outlined" },
    cell: ({ getValue }) => <EnabledCell value={getValue()} />,
  }),
  column.accessor("awork_name", {
    id: "awork",
    header: "Awork",
    cell: ({ getValue, table }) => (
      <AworkCell value={getValue()} aworkUsers={table.options.meta!.aworkUsers} />
    ),
  }),
  column.display({
    id: "actions",
    header: "",
    cell: () => <RowActions />,
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
  const data = useMemo(
    () => users.map((u) => ({ ...u, awork_name: aworkName(aworkUsers, u.awork_id) })),
    [users, aworkUsers],
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
            <TeamTableRow key={row.id} row={row} />
          ))}
        </tbody>
      </table>
    </div>
  );
}

function rowFormValues(user: User): TeamEditForm {
  return {
    vacation_days: user.vacation_days,
    role: user.role,
    enabled: user.enabled,
    awork_id: user.awork_id || "",
  };
}

function useTeamRowForm(user: User, onSave: (data: TeamEditForm) => Promise<unknown>) {
  return useAppForm({
    defaultValues: rowFormValues(user),
    onSubmit: ({ value }) => onSave(value),
  });
}

type RowEdit = {
  user: User;
  editing: boolean;
  formId: string;
  form: ReturnType<typeof useTeamRowForm>;
  startEdit: () => void;
  cancelEdit: () => void;
};

const RowEditContext = createContext<RowEdit | null>(null);

function useRowEdit() {
  const rowEdit = useContext(RowEditContext);
  if (!rowEdit) throw new Error("useRowEdit must be used inside a team table row");
  return rowEdit;
}

function TeamTableRow({ row }: { row: Row<typeof features, TeamRow> }) {
  const user = row.original;
  const queryClient = useQueryClient();
  const { addToast, addErrorToast } = useToast();
  const [editing, setEditing] = useState(false);

  const mutation = useMutation({
    mutationFn: (data: TeamEditForm) => chrono.users.updateUser(user.id, data),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: userQueries.all() }),
        queryClient.invalidateQueries({ queryKey: aworkQueries.workTimes() }),
      ]);
      addToast("Successfully updated profile settings", "success");
    },
    onError: (error) => addErrorToast(error),
  });

  const form = useTeamRowForm(user, (data) =>
    mutation.mutateAsync(data).then(
      () => setEditing(false),
      () => {},
    ),
  );

  const rowEdit: RowEdit = {
    user,
    editing,
    formId: `team-row-${user.id}`,
    form,
    startEdit: () => {
      form.reset(rowFormValues(user));
      setEditing(true);
    },
    cancelEdit: () => setEditing(false),
  };

  return (
    <RowEditContext value={rowEdit}>
      <tr
        className={
          user.vacation_days > 0
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
              <FlexRender cell={cell} />
            </Cell>
          );
        })}
      </tr>
    </RowEditContext>
  );
}

function VacationDaysCell({ value }: { value: number }) {
  const { user, editing, form, formId } = useRowEdit();
  if (!editing) return value;
  return (
    <form.AppField name="vacation_days">
      {(field) => (
        <field.NumberInput
          form={formId}
          id={`vacation_days-${user.id}`}
          aria-label={`Vacation days of ${user.username}`}
          className="input input-bordered"
          step={0.5}
          min={0}
          required
        />
      )}
    </form.AppField>
  );
}

function RoleCell({ value, roles }: { value: string; roles: readonly string[] }) {
  const { user, editing, form, formId } = useRowEdit();
  if (!editing) return value;
  return (
    <form.AppField name="role">
      {(field) => (
        <field.Select
          form={formId}
          id={`role-${user.id}`}
          aria-label={`Role of ${user.username}`}
          className={selectClassName}
          required
        >
          {roles.map((r) => (
            <option key={r} value={r}>
              {capitalize(r)}
            </option>
          ))}
        </field.Select>
      )}
    </form.AppField>
  );
}

function EnabledCell({ value }: { value: boolean }) {
  const { user, editing, form, formId } = useRowEdit();
  if (!editing) {
    return (
      <input type="checkbox" className="checkbox cursor-not-allowed" checked={value} readOnly />
    );
  }
  return (
    <form.AppField name="enabled">
      {(field) => (
        <field.Checkbox
          form={formId}
          id={`enabled-${user.id}`}
          aria-label={`${user.username} enabled`}
          className="checkbox"
        />
      )}
    </form.AppField>
  );
}

function AworkCell({ value, aworkUsers }: { value: string; aworkUsers: AworkUser[] }) {
  const { user, editing, form, formId } = useRowEdit();
  if (!editing) return value;
  return (
    <form.AppField name="awork_id">
      {(field) => (
        <field.Select
          form={formId}
          id={`awork_id-${user.id}`}
          aria-label={`Awork user of ${user.username}`}
          className={selectClassName}
        >
          <option value="">-</option>
          {aworkUsers.map((a) => (
            <option key={a.id} value={a.id}>
              {a.firstName} {a.lastName}
            </option>
          ))}
        </field.Select>
      )}
    </form.AppField>
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

function RowActions() {
  const { editing, form, formId, startEdit, cancelEdit } = useRowEdit();
  if (!editing) {
    return (
      <button
        type="button"
        onClick={startEdit}
        className="btn btn-soft btn-ghost animate-color icon-outlined"
      >
        edit
      </button>
    );
  }
  return (
    <form.AppForm>
      <form.Form id={formId} className="flex gap-1">
        <form.SubmitButton className="btn btn-soft btn-success animate-color icon-outlined">
          check
        </form.SubmitButton>
        <button
          type="button"
          onClick={cancelEdit}
          className="btn btn-soft btn-error animate-color icon-outlined"
        >
          close
        </button>
      </form.Form>
    </form.AppForm>
  );
}

function aworkName(aworkUsers: AworkUser[], aworkId: string | null) {
  if (!aworkId) return "-";
  const a = aworkUsers.find((a) => a.id === aworkId);
  return a ? `${a.firstName} ${a.lastName}` : aworkId;
}
