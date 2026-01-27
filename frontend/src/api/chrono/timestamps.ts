import type { Timestamp, WorkTime } from "../../types/response";
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

  async getAllForUser(start?: string, end?: string): Promise<Timestamp[]> {
    let url = "";
    if (start && end) url = `?startDate=${start}&endDate=${end}`;
    else if (start) url = `?startDate=${start}`;
    else if (end) url = `?endDate=${end}`;

    const response = await fetch(CHRONO_URL + `/timestamps${url}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Timestamp[];
  }

  async getAll(start?: string, end?: string): Promise<Timestamp[]> {
    let url = "";
    if (start && end) url = `?startDate=${start}&endDate=${end}`;
    else if (start) url = `?startDate=${start}`;
    else if (end) url = `?endDate=${end}`;

    const response = await fetch(CHRONO_URL + `/timestamps/all${url}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Timestamp[];
  }

  async getForYear(year: number): Promise<Timestamp[]> {
    const response = await fetch(CHRONO_URL + `/timestamps/year/${year}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Timestamp[];
  }

  async getForMonth(year: number, month: number): Promise<Timestamp[]> {
    const response = await fetch(
      CHRONO_URL + `/timestamps/year/${year}/month/${month}`,
      {
        method: "GET",
        credentials: "include",
      },
    );

    const r = await returnOrError(response);
    return r.data as Timestamp[];
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

  async getWorkHours(year: number): Promise<WorkTime> {
    const response = await fetch(CHRONO_URL + `/timestamps/worked/${year}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as WorkTime;
  }

  async getWorkHoursForAllUsers(
    year: number,
  ): Promise<Record<number, WorkTime>> {
    const response = await fetch(
      CHRONO_URL + `/timestamps/worked/${year}/all`,
      {
        method: "GET",
        credentials: "include",
      },
    );

    const r = await returnOrError(response);
    return r.data as Record<number, WorkTime>;
  }
}
