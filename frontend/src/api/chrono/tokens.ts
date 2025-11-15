import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiTokens {
  async createTokens(userId: number, tokens: number): Promise<void> {
    const form = new FormData();
    form.append("filter", userId.toString());
    form.append("token", tokens.toString());

    const response = await fetch(CHRONO_URL + `/tokens`, {
      method: "POST",
      credentials: "include",
      body: form,
    });

    await returnOrError(response);
  }
}
