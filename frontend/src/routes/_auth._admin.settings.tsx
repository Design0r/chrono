import { useMutation, useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { ErrorPage } from "../components/ErrorPage";
import type { Settings } from "../types/response";
import { useToast } from "../components/Toast";

export const Route = createFileRoute("/_auth/_admin/settings")({
  component: SettingsComponent,
});

function SettingsComponent() {
  const { chrono } = Route.useRouteContext();
  const { addToast, addErrorToast } = useToast();

  const settingQ = useQuery({
    queryKey: ["settings"],
    queryFn: () => chrono.settings.getSettings(),
    staleTime: 1000 * 30, // 30s
    gcTime: 1000 * 60 * 5, // 5m
  });

  const mutation = useMutation({
    mutationKey: ["settings", "signup_enabled"],
    mutationFn: (s: Settings) => chrono.settings.updateSettings(s),
    onSuccess: () => addToast("Updated settings", "success"),
    onError: (error) => addErrorToast(error),
  });

  if (settingQ.isPending) return <LoadingSpinnerPage />;
  if (settingQ.isError) return <ErrorPage error={settingQ.error} />;

  return (
    <div className="pt-2 m-2">
      <div className="space-y-2 bg-base-300 p-4 max-w-sm w-full rounded-xl grid grid-cols-2 justify-center mx-auto">
        <div className="col-start-1 mb-0 ">
          <p className="mb-0 text-lg">Signup enabled</p>
        </div>
        <div className="flex justify-end col-start-2">
          <input
            defaultChecked={settingQ.data!.signup_enabled}
            onChange={(e) => {
              mutation.mutate({
                ...settingQ.data,
                signup_enabled: e.target.checked,
              });
            }}
            type="checkbox"
            className="toggle border-error text-error checked:border-success checked:text-success"
          />
        </div>
      </div>
    </div>
  );
}
