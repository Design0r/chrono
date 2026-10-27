import { useMutation } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { chrono } from "../api/chrono/client";
import { useToast } from "../components/Toast";
import { useAppForm } from "../lib/form";

export const Route = createFileRoute("/_auth/_admin/export")({
  component: RouteComponent,
});

function RouteComponent() {
  const [currentYear] = useState(() => new Date().getFullYear());
  const toast = useToast();

  const mutation = useMutation({
    mutationFn: (y: number) => chrono.export.download(y),
    onSuccess: () => toast.addToast("Created report", "success"),
    onError: (error) => toast.addErrorToast(error),
  });

  const form = useAppForm({
    defaultValues: { year: currentYear },
    onSubmit: ({ value }) => mutation.mutateAsync(value.year).catch(() => {}),
  });

  return (
    <div className="container justify-center">
      <div className="justify-center flex flex-col mx-auto max-w-sm">
        <form.AppForm>
          <form.Form className="flex flex-col space-y-2">
            <label>
              <form.AppField name="year">
                {(field) => (
                  <field.NumberInput
                    aria-label="Year"
                    className="input input-bordered w-full"
                    step="1"
                    required
                  />
                )}
              </form.AppField>
            </label>
            <form.SubmitButton className="btn text-white btn-primary bg-primary/80 hover:bg-primary animate-color">
              Krankheitstage Jahr
            </form.SubmitButton>
          </form.Form>
        </form.AppForm>
      </div>
    </div>
  );
}
