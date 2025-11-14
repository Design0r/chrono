import type { Month, VacationGraph } from "../../types/response";
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

  async getVacationGraph(year: number): Promise<VacationGraph> {
    const response = await fetch(CHRONO_URL + `/events/${year}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as VacationGraph;
  }
}
