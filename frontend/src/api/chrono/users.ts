import type { User, UserWithVacation } from "../../types/auth";
import type { ProfileEditForm, TeamEditForm } from "../../types/forms";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiUsers {
  async getUserById(
    id: number,
    vacation: { year: number } | null = null,
  ): Promise<User> {
    const params = new URLSearchParams(
      vacation ? { vacation: "true", year: String(vacation.year) } : {},
    ).toString();
    const response = await fetch(
      CHRONO_URL + `/users/${id}` + (params ? "?" + params : ""),
      {
        method: "GET",
        credentials: "include",
      },
    );

    const r = await returnOrError(response);
    return r.data as User | UserWithVacation;
  }

  async getUsers(
    vacation: null | { year: number } = null,
  ): Promise<User[] | UserWithVacation[]> {
    const params = new URLSearchParams(
      vacation ? { vacation: "true", year: String(vacation.year) } : {},
    ).toString();

    const response = await fetch(
      CHRONO_URL + "/users" + (params ? "?" + params : ""),
      {
        method: "GET",
        credentials: "include",
      },
    );

    const r = await returnOrError(response);
    return r.data;
  }

  async updateUser(
    userId: number,
    data: ProfileEditForm | TeamEditForm,
  ): Promise<User> {
    const form = new FormData();
    Object.entries(data).map(([k, v]) => form.append(k, v.toString()));

    console.log("hello");

    const response = await fetch(CHRONO_URL + `/users/${userId}`, {
      method: "PATCH",
      credentials: "include",
      body: form,
    });

    const r = await returnOrError(response);
    return r.data;
  }

  getRoles() {
    return ["admin", "user", "guest"] as const;
  }
}
