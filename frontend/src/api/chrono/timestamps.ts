import type { Timestamp } from "../../types/response";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiTimestamps {
  async start(): Promise<Timestamp> {
    const response = await fetch(CHRONO_URL + `/timestamps`, {
      method: "POST",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Timestamp;
  }

  async update(ts: Timestamp): Promise<Timestamp> {
    const form = new FormData();
    form.append("id", ts.id.toString());
    form.append("start_time", ts.start_time);
    form.append("end_time", ts.end_time || "");
    form.append("user_id", ts.user_id.toString());

    const response = await fetch(CHRONO_URL + `/timestamps/${ts.id}`, {
      method: "PUT",
      credentials: "include",
      body: form,
    });

    const r = await returnOrError(response);
    return r.data as Timestamp;
  }

  async stop(id: number): Promise<Timestamp> {
    const response = await fetch(CHRONO_URL + `/timestamps/${id}`, {
      method: "PATCH",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Timestamp;
  }

  async getForToday(): Promise<Timestamp[]> {
    const response = await fetch(CHRONO_URL + `/timestamps/day`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Timestamp[];
  }

  async getLatest(): Promise<Timestamp> {
    const response = await fetch(CHRONO_URL + `/timestamps/latest`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Timestamp;
  }
}
