import { ApiAuth } from "./auth";
import { ApiEvents } from "./events";
import { ApiRequests } from "./requests";
import { ApiSettings } from "./settings";
import { ApiTokens } from "./tokens";
import { ApiUsers } from "./users";

export class ChronoClient {
  auth: ApiAuth;
  users: ApiUsers;
  events: ApiEvents;
  requests: ApiRequests;
  tokens: ApiTokens;
  settings: ApiSettings;

  constructor() {
    this.auth = new ApiAuth();
    this.users = new ApiUsers();
    this.events = new ApiEvents();
    this.requests = new ApiRequests();
    this.tokens = new ApiTokens();
    this.settings = new ApiSettings();
  }
}
