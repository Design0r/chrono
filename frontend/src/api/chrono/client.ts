import { ApiAuth } from "./auth";
import { ApiEvents } from "./events";
import { ApiRequests } from "./requests";
import { ApiUsers } from "./users";

export class ChronoClient {
  auth: ApiAuth;
  users: ApiUsers;
  events: ApiEvents;
  requests: ApiRequests;

  constructor() {
    this.auth = new ApiAuth();
    this.users = new ApiUsers();
    this.events = new ApiEvents();
    this.requests = new ApiRequests();
  }
}
