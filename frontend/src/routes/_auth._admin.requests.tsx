import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { requestQueries } from "../api/queries";
import { RequestRow } from "../components/Requests";

export const Route = createFileRoute("/_auth/_admin/requests")({
  loader: ({ context: { queryClient } }) => queryClient.ensureQueryData(requestQueries.list()),
  component: RouteComponent,
});

function RouteComponent() {
  const { data } = useSuspenseQuery(requestQueries.list());

  return (
    <div>
      <div className="overflow-x-auto bg-info/3">
        <table className="table table-zebra table-md ">
          <thead>
            <tr className="bg-base-200 text-base-content border-b-1.5 border-primary/65">
              <th>ID</th>
              <th>Name</th>
              <th>Event Type</th>
              <th>Number of Days</th>
              <th>Event Start Date</th>
              <th>Event End Date</th>
              <th>Conflicts</th>
              <th className="text-center">Actions</th>
            </tr>
          </thead>
          <tbody>
            {data.map((r) => (
              <RequestRow key={r.request.request_id} request={r} />
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
