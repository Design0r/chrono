import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import dayjs from "dayjs";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { ErrorPage } from "../components/ErrorPage";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { useToast } from "../components/Toast";
import type { ProfileEditForm } from "../types/forms";

export const Route = createFileRoute("/_auth/profile")({
  component: ProfileComponent,
});

function ProfileComponent() {
  const { auth, chrono } = Route.useRouteContext();
  const [edit, setEdit] = useState(false);
  const { register, handleSubmit } = useForm<ProfileEditForm>();
  const { addToast, addErrorToast } = useToast();
  const queryClient = useQueryClient();

  const userQ = useQuery({
    queryKey: ["user", auth.userId],
    queryFn: () => chrono.users.getUserById(auth.userId!),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
  });

  const aworkUserQ = useQuery({
    queryKey: ["aworkUsers"],
    queryFn: () => chrono.awork.getUsers(),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
  });

  const mutation = useMutation({
    mutationKey: ["profile", userQ.data?.id],
    mutationFn: ({ userId, data }: { userId: number; data: ProfileEditForm }) =>
      chrono.users.updateUser(userId, data),
    onSuccess: () => {
      addToast("Successfully updated profile settings", "success");
      setEdit(false);
      return queryClient.invalidateQueries({ queryKey: ["user"] });
    },
    onError: (error) => addErrorToast(error),
  });

  if (userQ.isPending) return <LoadingSpinnerPage />;
  if (userQ.isError) return <ErrorPage error={userQ.error} />;

  const user = userQ.data!;
  const createdAt = dayjs(user.created_at).toString();
  const editedAt = dayjs(user.edited_at).toString();

  return (
    <>
      {edit ? (
        <div>
          <form
            onSubmit={handleSubmit((data: ProfileEditForm) => {
              return mutation.mutate({
                userId: user.id,
                data,
              });
            })}
          >
            <div className="container justify-center flex">
              <div className="space-y-2 bg-base-100 rounded-xl px-8 py-10">
                <h1 className="font-bold text-xl mb-0">Profile</h1>
                <div className="divider mb-4 mt-0"></div>
                <div className="grid gap-4 grid-cols-2">
                  <label htmlFor="username">Username</label>
                  <input
                    className="input input-bordered"
                    type="text"
                    required
                    defaultValue={user.username}
                    {...register("username")}
                  />
                  <label htmlFor="email">Email</label>
                  <input
                    className="input input-bordered"
                    type="email"
                    required
                    defaultValue={user.email}
                    {...register("email")}
                  />
                  <label htmlFor="awork_id">Awork ID</label>
                  <select
                    defaultValue={user.awork_id || ""}
                    className="w-full col-span-1 cursor-pointer p-2 select pl-3 text-center h-full  rounded-md "
                    {...register("awork_id")}
                  >
                    <option value="">-</option>
                    {aworkUserQ.data?.map((a) => (
                      <option key={a.id} value={a.id}>
                        {a.firstName} {a.lastName}
                      </option>
                    ))}
                  </select>

                  <label htmlFor="password">New Password</label>
                  <input
                    className="input input-bordered"
                    type="password"
                    {...register("password")}
                  />
                  <p>Admin</p>
                  <p>{String(user.is_superuser)}</p>
                  <p>Yearly Vacation</p>
                  <p>{user.vacation_days}</p>
                  <label htmlFor="color">Color</label>
                  <input
                    className="input input-bordered"
                    type="color"
                    required
                    defaultValue={user.color}
                    {...register("color")}
                  />
                  <p>Joined</p>
                  <p>{createdAt}</p>
                  <p>Last Edit</p>
                  <p>{editedAt}</p>
                </div>
              </div>
            </div>

            <div className="container mx-auto max-w-lg px-6 mt-4">
              <button
                type="submit"
                className="btn btn-soft btn-primary animate-color"
              >
                <span className="icon-outlined pr-2">save</span>
                Save
              </button>
            </div>
          </form>
        </div>
      ) : (
        <div>
          <div className="container mx-auto max-w-lg px-6">
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
                <p>Admin</p>
                <p>{String(user.is_superuser)}</p>
                <p>Yearly Vacation</p>
                <p>{user.vacation_days}</p>
                <p>Color</p>
                <p
                  className="max-w-40 rounded-full "
                  style={{ backgroundColor: user.color }}
                ></p>
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
              <span className="icon-outlined pr-2">
                {edit ? "save" : "edit"}
              </span>
              {edit ? "Save" : "Edit"}
            </button>
          </div>
        </div>
      )}
    </>
  );
}
