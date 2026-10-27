import type { InputHTMLAttributes, SelectHTMLAttributes, TextareaHTMLAttributes } from "react";
import { useFieldContext } from "../../lib/form-context";

type ControlProps<T> = Omit<T, "name" | "value" | "defaultValue" | "onChange" | "onBlur">;

function useControl<T>() {
  const field = useFieldContext<T>();
  const invalid = field.state.meta.isTouched && field.state.meta.errors.length > 0;
  return {
    field,
    common: {
      id: field.name,
      name: field.name,
      onBlur: field.handleBlur,
      "aria-invalid": invalid || undefined,
    },
  };
}

export function TextInput(props: ControlProps<InputHTMLAttributes<HTMLInputElement>>) {
  const { field, common } = useControl<string>();
  return (
    <input
      type="text"
      {...common}
      {...props}
      value={field.state.value}
      onChange={(e) => field.handleChange(e.target.value)}
    />
  );
}

export function NumberInput(props: ControlProps<InputHTMLAttributes<HTMLInputElement>>) {
  const { field, common } = useControl<number>();
  return (
    <input
      {...common}
      {...props}
      type="number"
      value={Number.isNaN(field.state.value) ? "" : field.state.value}
      onChange={(e) => field.handleChange(e.target.valueAsNumber)}
    />
  );
}

export function Select(props: ControlProps<SelectHTMLAttributes<HTMLSelectElement>>) {
  const { field, common } = useControl<string>();
  return (
    <select
      {...common}
      {...props}
      value={field.state.value}
      onChange={(e) => field.handleChange(e.target.value)}
    />
  );
}

export function Checkbox(props: ControlProps<InputHTMLAttributes<HTMLInputElement>>) {
  const { field, common } = useControl<boolean>();
  return (
    <input
      {...common}
      {...props}
      type="checkbox"
      checked={field.state.value}
      onChange={(e) => field.handleChange(e.target.checked)}
    />
  );
}

export function TextArea(props: ControlProps<TextareaHTMLAttributes<HTMLTextAreaElement>>) {
  const { field, common } = useControl<string>();
  return (
    <textarea
      {...common}
      {...props}
      value={field.state.value}
      onChange={(e) => field.handleChange(e.target.value)}
    />
  );
}
