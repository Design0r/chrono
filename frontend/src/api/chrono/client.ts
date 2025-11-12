import { ApiAuth } from "./auth";
import { ApiUsers } from "./users";

export class ChronoClient {
  auth: ApiAuth;
  users: ApiUsers;
  constructor() {
    this.auth = new ApiAuth();
    this.users = new ApiUsers();
  }
}
