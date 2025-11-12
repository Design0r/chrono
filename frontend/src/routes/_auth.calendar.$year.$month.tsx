import { createFileRoute } from "@tanstack/react-router";
import {
  Calendar,
  CalendarNavigation,
  EventFilter,
  UserFilter,
  VacationCounter,
} from "../components/Calendar";

export const Route = createFileRoute("/_auth/calendar/$year/$month")({
  component: RouteComponent,
  loader: async ({ context: { chrono, queryClient }, params }) => {
    const users = await queryClient.ensureQueryData({
      queryKey: ["users"],
      queryFn: async () =>
        await chrono.users.getUsers({ year: Number.parseInt(params.year) }),
    });

    const month = await queryClient.ensureQueryData({
      queryKey: ["month", params.year, params.month],
      queryFn: async () =>
        await chrono.events.getEventsForMonth(
          Number.parseInt(params.year),
          Number.parseInt(params.month),
        ),
    });

    return { users: users, month: month };
  },
});

function RouteComponent() {
  const { users, month } = Route.useLoaderData();
  const params = Route.useParams();

  return (
    <div>
      <div className="grid grid-cols-7">
        <div className="px-6 col-span-7 grid grid-cols-1 grid-rows-4 lg:grid-rows-1 items-center gap-y-2 lg:gap-y-0 lg:gap-x-2 mt-2 lg:mb-16 lg:grid-cols-7 lg:px-4 ">
          <select
            className=" col-span-1 lg:col-span-1 cursor-pointer bg-base-100 select hover:text-white border-0 hover:bg-[#6F78EA] text-center focus:outline-0 h-full w-full text-lg rounded-xl animate-color"
            name="eventName"
            id="eventName"
          >
            <option value="urlaub">Urlaub</option>
            <option value="urlaub halbtags">Urlaub Halbtags</option>
            <option value="workation">Workation</option>
            <option value="krank">Krank</option>
            <option value="home office">Home Office</option>
          </select>
          <CalendarNavigation
            year={Number.parseInt(params.year)}
            month={Number.parseInt(params.month)}
          />
          <div className="row-span-2 col-span-2 lg:row-span-1 lg:col-span-4 h-full text-lg">
            <div className="h-full items-center rounded-xl bg-base-200">
              <div className="grid grid-cols-2 lg:grid-cols-4 w-full h-full gap-x-2 gap-y-2 lg:gap-y-0">
                <div className="col-span-1 w-full justify-center ">
                  <UserFilter users={users} />
                </div>
                <div className="col-span-1 w-full justify-center ">
                  <EventFilter
                    events={[
                      "Krank",
                      "Home Office",
                      "Urlaub",
                      "Urlaub Halbtags",
                      "Workation",
                    ]}
                  />
                </div>
                <div className="flex col-span-2 w-full items-center align-middle h-full ">
                  <VacationCounter pending={3} used={5} remaining={6} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <Calendar month={month} />
    </div>
  );
}
