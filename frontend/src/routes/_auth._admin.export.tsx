import { useMutation } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useToast } from "../components/Toast";

export const Route = createFileRoute("/_auth/_admin/export")({
  component: RouteComponent,
});

function RouteComponent() {
  const { chrono } = Route.useRouteContext();
  const currYear = new Date().getFullYear();
  const [year, setYear] = useState(currYear);
  const toast = useToast();

  const mutation = useMutation({
    mutationKey: ["export", year],
    mutationFn: (y: number) => chrono.export.download(y),
    onSuccess: () => toast.addToast("Created report", "success"),
    onError: (error) => toast.addErrorToast(error),
    retry: false,
  });

  return (
    <div className="container justify-center">
      <div className="justify-center flex flex-col mx-auto max-w-sm">
        <div className="flex flex-col space-y-2">
          <label>
            <input
              className="input input-bordered w-full"
              name="year"
              type="number"
              step="1"
              defaultValue={currYear}
              onChange={(e) => setYear(Number(e.target.value))}
            />
          </label>
          <button
            onClick={() => mutation.mutate(year)}
            className="btn text-white btn-primary bg-primary/80 hover:bg-primary animate-color"
          >
            Krankheitstage Jahr
          </button>
        </div>
      </div>
    </div>
  );
}
