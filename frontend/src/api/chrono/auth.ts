import type { LoginRequest, SignupRequest } from "../../types/auth";
import type { ChronoResponse } from "../../types/response";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiAuth {
  async login(data: LoginRequest): Promise<ChronoResponse> {
    const form = new FormData();
    form.append("email", data.email);
    form.append("password", data.password);

    const response = await fetch(CHRONO_URL + "/login", {
      method: "POST",
      body: form,
      credentials: "include",
    });

    return await returnOrError(response);
  }

  async signup(data: SignupRequest): Promise<ChronoResponse> {
    const form = new FormData();
    form.append("email", data.email);
    form.append("password", data.password);
    form.append("username", data.username);

    const response = await fetch(CHRONO_URL + "/signup", {
      method: "POST",
      body: form,
      credentials: "include",
    });

    return await returnOrError(response);
  }

  async logout(): Promise<ChronoResponse> {
    const response = await fetch(CHRONO_URL + "/logout", {
      method: "POST",
      credentials: "include",
    });

    return await returnOrError(response);
  }
}
