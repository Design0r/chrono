import { createFormHook } from "@tanstack/react-form";
import { Form, FormErrors, SubmitButton } from "../components/form/FormComponents";
import { Checkbox, NumberInput, Select, TextArea, TextInput } from "../components/form/FormFields";
import { fieldContext, formContext } from "./form-context";

export const { useAppForm } = createFormHook({
  fieldContext,
  formContext,
  fieldComponents: { TextInput, NumberInput, Select, Checkbox, TextArea },
  formComponents: { Form, SubmitButton, FormErrors },
});
