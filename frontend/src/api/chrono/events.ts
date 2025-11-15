import type { EventUser, Month, VacationGraph } from "../../types/response";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiEvents {
  async getEventsForMonth(year: number, month: number): Promise<Month> {
    const response = await fetch(CHRONO_URL + `/events/${year}/${month}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as Month;
  }

  async createEvent({
    year,
    month,
    day,
    event,
  }: {
    year: number;
    month: number;
    day: number;
    event: string;
  }): Promise<EventUser> {
    const form = new FormData();
    form.append("year", year.toString());
    form.append("month", month.toString());
    form.append("day", day.toString());
    form.append("eventName", event);

    const response = await fetch(CHRONO_URL + `/events`, {
      method: "POST",
      credentials: "include",
      body: form,
    });

    const r = await returnOrError(response);
    return r.data as EventUser;
  }

  async deleteEvent(id: number): Promise<void> {
    const response = await fetch(CHRONO_URL + `/events/${id}`, {
      method: "DELETE",
      credentials: "include",
    });

    await returnOrError(response);
  }

  async getVacationGraph(year: number): Promise<VacationGraph> {
    const response = await fetch(CHRONO_URL + `/events/${year}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as VacationGraph;
  }

  getEventTypes() {
    return [
      "Krank",
      "Home Office",
      "Urlaub",
      "Urlaub Halbtags",
      "Workation",
    ] as const;
  }

  getShortEventName(name: string) {
    return name
      .split(" ")
      .reduce((acc, curr) => acc + curr[0].toUpperCase(), "");
  }
}
