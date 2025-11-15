import type { BatchRequest, PatchRequestForm } from "../../types/response";
import { returnOrError } from "../error";
import { CHRONO_URL } from "./chrono";

export class ApiRequests {
  async getRequests(): Promise<BatchRequest[]> {
    const response = await fetch(CHRONO_URL + `/requests`, {
      method: "GET",
      credentials: "include",
    });

    const r = await returnOrError(response);
    return r.data as BatchRequest[];
  }

  async patchRequest(data: PatchRequestForm): Promise<void> {
    const form = new FormData();
    form.append("user_id", data.userId.toString());
    form.append("state", data.state);
    form.append("reason", data.reason);
    form.append("start_date", data.start_date);
    form.append("end_date", data.end_date);

    const response = await fetch(CHRONO_URL + `/requests`, {
      method: "PATCH",
      credentials: "include",
      body: form,
    });

    await returnOrError(response);
  }
}
