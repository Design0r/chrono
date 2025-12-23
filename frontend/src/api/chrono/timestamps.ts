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
