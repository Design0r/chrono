import type { User } from "../../types/auth";
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
}
