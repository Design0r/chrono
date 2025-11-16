import { useMutation } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";

export const Route = createFileRoute("/_auth/_admin/export")({
  component: RouteComponent,
});

function RouteComponent() {
  const { chrono } = Route.useRouteContext();
  const currYear = new Date().getFullYear();
  const [year, setYear] = useState(currYear);

  const mutation = useMutation({
    mutationKey: ["export", year],
    mutationFn: (y: number) => chrono.export.download(y),
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
              onChange={(e) => setYear(Number(e))}
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
