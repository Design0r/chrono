import type { LoginRequest } from "../../types/auth";
import type { ChronoResponse } from "../../types/response";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export async function login(data: LoginRequest): Promise<ChronoResponse> {
  const form = new FormData();
  form.append("username", data.username);
  form.append("password", data.password);

  const response = await fetch(CHRONO_URL + "/login", {
    method: "POST",
    body: form,
    credentials: "include",
  });

  return await returnOrError(response);
}

export async function logout(): Promise<ChronoResponse> {
  const response = await fetch(CHRONO_URL + "/logout", {
    method: "POST",
    credentials: "include",
  });

  return await returnOrError(response);
}
