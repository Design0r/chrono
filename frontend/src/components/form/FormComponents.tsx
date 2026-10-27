import type { FormHTMLAttributes, ReactNode } from "react";
import { useFormContext } from "../../lib/form-context";
import { LoadingSpinner } from "../LoadingSpinner";

export function Form(props: Omit<FormHTMLAttributes<HTMLFormElement>, "onSubmit">) {
  const form = useFormContext();
  return (
    <form
      {...props}
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        void form.handleSubmit();
      }}
    />
  );
}

export function SubmitButton({
  className,
  children,
  pending,
}: {
  className?: string;
  children: ReactNode;
  pending?: ReactNode;
}) {
  const form = useFormContext();
  return (
    <form.Subscribe selector={(s) => s.isSubmitting}>
      {(isSubmitting) => (
        <button type="submit" className={className} disabled={isSubmitting}>
          {isSubmitting ? (pending ?? <LoadingSpinner size="sm" />) : children}
        </button>
      )}
    </form.Subscribe>
  );
}

export function FormErrors({ className }: { className?: string }) {
  const form = useFormContext();
  return (
    <form.Subscribe
      selector={(s) => ({
        errors: s.errors,
        visible: s.isTouched || s.submissionAttempts > 0,
      })}
    >
      {({ errors, visible }) =>
        visible && errors.length > 0 ? (
          <p role="alert" className={className}>
            {errors.map(String).join(" ")}
          </p>
        ) : null
      }
    </form.Subscribe>
  );
}
