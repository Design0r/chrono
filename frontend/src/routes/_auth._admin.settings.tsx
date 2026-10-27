import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { chrono } from "../api/chrono/client";
import { settingsQueries } from "../api/queries";
import { useToast } from "../components/Toast";
import type { Settings } from "../types/response";

export const Route = createFileRoute("/_auth/_admin/settings")({
  loader: ({ context: { queryClient } }) => queryClient.ensureQueryData(settingsQueries.detail()),
  component: SettingsComponent,
});

function SettingsComponent() {
  const { addToast, addErrorToast } = useToast();
  const queryClient = useQueryClient();

  const { data: settings } = useSuspenseQuery(settingsQueries.detail());

  const mutation = useMutation({
    mutationFn: (s: Settings) => chrono.settings.updateSettings(s),
    onSuccess: (updated) => {
      queryClient.setQueryData(settingsQueries.detail().queryKey, updated);
      addToast("Updated settings", "success");
    },
    onError: (error) => addErrorToast(error),
  });

  return (
    <div>
      <div className="space-y-2 bg-base-300 p-4 max-w-sm w-full rounded-xl grid grid-cols-2 justify-center mx-auto">
        <div className="col-start-1 mb-0 ">
          <p className="mb-0 text-lg">Signup enabled</p>
        </div>
        <div className="flex justify-end col-start-2">
          <input
            defaultChecked={settings.signup_enabled}
            onChange={(e) => {
              mutation.mutate({
                ...settings,
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
