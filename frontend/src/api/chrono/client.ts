import { ApiAuth } from "./auth";
import { ApiEvents } from "./events";
import { ApiUsers } from "./users";

export class ChronoClient {
  auth: ApiAuth;
  users: ApiUsers;
  events: ApiEvents;
  constructor() {
    this.auth = new ApiAuth();
    this.users = new ApiUsers();
    this.events = new ApiEvents();
  }
}
