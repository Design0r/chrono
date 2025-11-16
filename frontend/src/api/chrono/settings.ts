import type { Settings } from "../../types/response";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiSettings {
  async getSettings(): Promise<Settings> {
    const response = await fetch(CHRONO_URL + `/settings`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Settings;
  }

  async updateSettings(settings: Settings): Promise<Settings> {
    const form = new FormData();
    Object.entries(settings).map(([k, v]) => form.append(k, v.toString()));

    const response = await fetch(CHRONO_URL + `/settings`, {
      method: "PATCH",
      credentials: "include",
      body: form,
    });

    const r = await returnOrError(response);
    return r.data as Settings;
  }
}
