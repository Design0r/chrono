import type { AworkUser, WorkTime } from "../../types/response";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiAwork {
  async getWorkTimesforYear(year: number): Promise<WorkTime> {
    const response = await fetch(CHRONO_URL + `/awork/${year}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as WorkTime;
  }

  async getUsers(): Promise<AworkUser[]> {
    const response = await fetch(CHRONO_URL + `/awork/users`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as AworkUser[];
  }
}
