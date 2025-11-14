import { createFileRoute } from "@tanstack/react-router";
import { StatCard, StatCardElement } from "../components/StatCard";
import { TitleSection } from "../components/TitleSection";
import { VacationGraph } from "../components/VacationGraph";
import type { User } from "../types/auth";

export const Route = createFileRoute("/_auth/")({
  component: Home,
  loader: async ({ context }) => {
    const { user, refreshUser } = context.auth;
    const { queryClient, chrono } = context;
    let u: User | null = null;
    if (!user) {
      u = await refreshUser();
    }

    const data = await queryClient.ensureQueryData({
      queryKey: ["vacation_graph"],
      queryFn: async () =>
        await chrono.events.getVacationGraph(new Date().getFullYear()),
    });

    return { user: user ? user : u, data: data };
  },
});

function Home() {
  const { user, data } = Route.useLoaderData();

  return (
    <div className="flex flex-col container mx-auto justify-center align-middle gap-6 p-4">
      <div className="text-[48px] pl-2 text-primary font-light mb-2">
        <span className="animate-pulse font-medium text-white pr-1"> Hej </span>
        {user?.username}
      </div>
      <TitleSection title="Your vacation stats">
        <StatCard>
          <StatCardElement title="Vacation remaining" subtitle="43% remaining">
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-primary">
              13 d
            </span>
          </StatCardElement>
          <StatCardElement title="Vacation taken" subtitle="57% taken">
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-secondary opacity-40">
              17 d
            </span>
          </StatCardElement>
          <StatCardElement title="Vacation total" subtitle="30 days total">
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-secondary opacity-40">
              30 d
            </span>
          </StatCardElement>
          <StatCardElement title="Vacation pending" subtitle="5 events pending">
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-accent">
              <span className="animate-pulse text-primary">5</span> d
            </span>
          </StatCardElement>
        </StatCard>
      </TitleSection>
      <TitleSection title="Your year stats">
        <StatCard>
          <StatCardElement title="Days this year" subtitle="365 days total">
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-secondary opacity-40">
              365 d
            </span>
          </StatCardElement>
          <StatCardElement title="Days passed" subtitle="309 days passed">
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-accent">
              309 d
            </span>
          </StatCardElement>
          <StatCardElement title="Days completed" subtitle="15.34% remaining">
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-accent">
              84.66 %
            </span>
          </StatCardElement>
          <StatCardElement title="Vacation pending" subtitle="5 events pending">
            <progress
              className="progress progress-primary mb-3 mt-2 h-3.5"
              value={84.66}
              max="100"
              role="progressbar"
            />
          </StatCardElement>
        </StatCard>
      </TitleSection>
      <TitleSection title="Team Vacation">
        <VacationGraph
          yearOffset={data.year_offset}
          gaps={data.month_gaps}
          data={data.vacation_data}
        />
      </TitleSection>
    </div>
  );
}
