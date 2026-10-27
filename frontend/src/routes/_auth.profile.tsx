import { useMutation, useQuery, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import dayjs from "dayjs";
import { useState } from "react";
import { chrono } from "../api/chrono/client";
import { aworkQueries, currentUserId, userQueries } from "../api/queries";
import { useToast } from "../components/Toast";
import { useAppForm } from "../lib/form";
import { dateLocalToIso, isoToDateLocal } from "../lib/timestamp-utils";
import type { User } from "../types/auth";
import type { ProfileEditForm } from "../types/forms";
import type { AworkUser } from "../types/response";

export const Route = createFileRoute("/_auth/profile")({
  loader: async ({ context: { auth, queryClient } }) => {
    const userId = currentUserId(auth.userId);
    const user = await queryClient.ensureQueryData(userQueries.detail(userId));
    if (user.is_superuser) void queryClient.prefetchQuery(aworkQueries.users());
    return { userId };
  },
  component: ProfileComponent,
});

function ProfileComponent() {
  const { userId } = Route.useLoaderData();
  const [edit, setEdit] = useState(false);
  const { addToast, addErrorToast } = useToast();
  const queryClient = useQueryClient();

  const { data: user } = useSuspenseQuery(userQueries.detail(userId));

  // Das awork-Verzeichnis gibt es nur für Admins.
  const aworkUserQ = useQuery({ ...aworkQueries.users(), enabled: user.is_superuser });

  const mutation = useMutation({
    mutationFn: ({ userId, data }: { userId: number; data: ProfileEditForm }) =>
      chrono.users.updateUser(userId, data),
    onSuccess: async () => {
      addToast("Successfully updated profile settings", "success");
      setEdit(false);
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: userQueries.all() }),
        queryClient.invalidateQueries({ queryKey: aworkQueries.workTimes() }),
      ]);
    },
    onError: (error) => addErrorToast(error),
  });

  const createdAt = dayjs(user.created_at).toString();
  const editedAt = dayjs(user.edited_at).toString();
  const startDate = dayjs(user.start_date).format("DD.MM.YYYY");

  return (
    <>
      {edit ? (
        <ProfileEditor
          user={user}
          aworkUsers={aworkUserQ.data ?? []}
          createdAt={createdAt}
          editedAt={editedAt}
          onSave={(data) => mutation.mutateAsync({ userId: user.id, data })}
        />
      ) : (
        <div>
          <div className="container mx-auto max-w-lg">
            <div className="space-y-2 bg-base-100 rounded-xl px-8 py-10">
              <h1 className="font-bold text-xl mb-0">Profile</h1>
              <div className="divider mb-4 mt-0"></div>
              <div className="grid gap-4 grid-cols-2">
                <p>Username</p>
                <p>{user.username}</p>
                <p>Email</p>
                <p>{user.email}</p>
                <p>Awork ID</p>
                <p>{user.awork_id}</p>
                <p>Workday Hours</p>
                <p>{user.workday_hours}</p>
                <p>Workdays / Week</p>
                <p>{user.workdays_week}</p>
                <p>Admin</p>
                <p>{String(user.is_superuser)}</p>
                <p>Yearly Vacation</p>
                <p>{user.vacation_days}</p>
                <p>Color</p>
                <p className="max-w-40 rounded-full " style={{ backgroundColor: user.color }}></p>
                <p>Start Date</p>
                <p>{startDate}</p>
                <p>Joined</p>
                <p>{createdAt}</p>
                <p>Last Edit</p>
                <p>{editedAt}</p>
              </div>
            </div>
          </div>

          <div className="container mx-auto max-w-lg px-6 mt-4">
            <button
              onClick={() => {
                setEdit((prev) => !prev);
              }}
              type="submit"
              className="btn btn-soft btn-primary animate-color"
            >
              <span className="icon-outlined pr-2">{edit ? "save" : "edit"}</span>
              {edit ? "Save" : "Edit"}
            </button>
          </div>
        </div>
      )}
    </>
  );
}

function ProfileEditor({
  user,
  aworkUsers,
  createdAt,
  editedAt,
  onSave,
}: {
  user: User;
  aworkUsers: AworkUser[];
  createdAt: string;
  editedAt: string;
  onSave: (data: ProfileEditForm) => Promise<unknown>;
}) {
  const form = useAppForm({
    defaultValues: {
      username: user.username,
      email: user.email,
      awork_id: user.awork_id || "",
      workday_hours: user.workday_hours,
      workdays_week: user.workdays_week,
      password: "",
      color: user.color,
      start_date: isoToDateLocal(user.start_date),
    },
    onSubmit: ({ value }) => {
      const { awork_id, ...rest } = value;
      return onSave({
        ...rest,
        ...(user.is_superuser ? { awork_id } : {}),
        start_date: dateLocalToIso(value.start_date),
      }).catch(() => {});
    },
  });

  return (
    <form.AppForm>
      <form.Form>
        <div className="container justify-center flex">
          <div className="space-y-2 bg-base-100 rounded-xl px-8 py-10">
            <h1 className="font-bold text-xl mb-0">Profile</h1>
            <div className="divider mb-4 mt-0"></div>
            <div className="grid gap-4 grid-cols-2">
              <label htmlFor="username">Username</label>
              <form.AppField name="username">
                {(field) => <field.TextInput className="input input-bordered" required />}
              </form.AppField>
              <label htmlFor="email">Email</label>
              <form.AppField name="email">
                {(field) => (
                  <field.TextInput className="input input-bordered" type="email" required />
                )}
              </form.AppField>
              <label htmlFor="awork_id">Awork ID</label>
              {user.is_superuser ? (
                <form.AppField name="awork_id">
                  {(field) => (
                    <field.Select className="w-full col-span-1 cursor-pointer bg-base-300 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-full pl-4 text-base border-0 rounded-md animate-all">
                      <option value="">-</option>
                      {aworkUsers.map((a) => (
                        <option key={a.id} value={a.id}>
                          {a.firstName} {a.lastName}
                        </option>
                      ))}
                    </field.Select>
                  )}
                </form.AppField>
              ) : (
                <p>{user.awork_id || "-"}</p>
              )}

              <label htmlFor="workday_hours">Workday Hours</label>
              <form.AppField name="workday_hours">
                {(field) => <field.NumberInput className="input input-bordered" required />}
              </form.AppField>

              <label htmlFor="workdays_week">Workdays / Week </label>
              <form.AppField name="workdays_week">
                {(field) => <field.NumberInput className="input input-bordered" required />}
              </form.AppField>

              <label htmlFor="password">New Password</label>
              <form.AppField name="password">
                {(field) => (
                  <field.TextInput
                    className="input input-bordered"
                    type="password"
                    autoComplete="new-password"
                  />
                )}
              </form.AppField>
              <p>Admin</p>
              <p>{String(user.is_superuser)}</p>
              <p>Yearly Vacation Days</p>
              <p>{user.vacation_days}</p>
              <label htmlFor="color">Color</label>
              <form.AppField name="color">
                {(field) => (
                  <field.TextInput className="input input-bordered" type="color" required />
                )}
              </form.AppField>
              <label htmlFor="start_date">Start Date</label>
              <form.AppField name="start_date">
                {(field) => (
                  <field.TextInput className="input input-bordered" type="date" required />
                )}
              </form.AppField>
              <p>Joined</p>
              <p>{createdAt}</p>
              <p>Last Edit</p>
              <p>{editedAt}</p>
            </div>
          </div>
        </div>

        <div className="mx-auto max-w-lg px-6 mt-4">
          <form.SubmitButton className="btn btn-soft btn-primary animate-color">
            <span className="icon-outlined pr-2">save</span>
            Save
          </form.SubmitButton>
        </div>
      </form.Form>
    </form.AppForm>
  );
}
