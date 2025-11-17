import type { Notification } from "../../types/response";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiNotifications {
  async get(): Promise<Notification[]> {
    const response = await fetch(CHRONO_URL + `/notifications`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Notification[];
  }

  async clear(id: number): Promise<void> {
    const response = await fetch(CHRONO_URL + `/notifications/${id}`, {
      method: "PATCH",
      credentials: "include",
    });

    await returnOrError(response);
  }

  async clearAll(): Promise<void> {
    const response = await fetch(CHRONO_URL + `/notifications`, {
      method: "PATCH",
      credentials: "include",
    });

    await returnOrError(response);
  }
}
