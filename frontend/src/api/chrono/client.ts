import { ApiAuth } from "./auth";
import { ApiAwork } from "./awork";
import { ApiEvents } from "./events";
import { ApiExport } from "./export";
import { ApiRequests } from "./requests";
import { ApiSettings } from "./settings";
import { ApiTokens } from "./tokens";
import { ApiUsers } from "./users";

export class ChronoClient {
  auth = new ApiAuth();
  users = new ApiUsers();
  events = new ApiEvents();
  requests = new ApiRequests();
  tokens = new ApiTokens();
  settings = new ApiSettings();
  export = new ApiExport();
  awork = new ApiAwork();
}
