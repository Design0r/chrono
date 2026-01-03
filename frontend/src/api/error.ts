import { logoutOutsideReact } from "../auth";
import type { ChronoResponse } from "../types/response";

export async function returnOrError(
  response: Response,
): Promise<ChronoResponse> {
  const data: ChronoResponse = await response.json();

  if (!response.ok) {
    if (response.status === 401) {
      await logoutOutsideReact();
    }
    throw new Error(data.message);
  }

  return data;
}
