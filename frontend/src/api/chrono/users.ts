import type { User, UserWithVacation } from "../../types/auth";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiUsers {
  async getUserById(id: number): Promise<User> {
    const response = await fetch(CHRONO_URL + `/users/${id}`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as User;
  }

  async getUsers(
    vacation: null | { year: number } = null
  ): Promise<User[] | UserWithVacation[]> {
    const params = new URLSearchParams(
      vacation ? { vacation: "true", year: String(vacation.year) } : {}
    ).toString();

    const response = await fetch(
      CHRONO_URL + "/users" + (params ? "?" + params : ""),
      {
        method: "GET",
        credentials: "include",
      }
    );

    const r = await returnOrError(response);
    return r.data;
  }
}
