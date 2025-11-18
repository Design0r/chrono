import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { ErrorPage } from "../components/ErrorPage";
import {
  LoadingSpinner,
  LoadingSpinnerPage,
} from "../components/LoadingSpinner";
import { StatCard, StatCardElement } from "../components/StatCard";
import { TitleSection } from "../components/TitleSection";
import { useToast } from "../components/Toast";
import { VacationGraph } from "../components/VacationGraph";
import type { UserWithVacation } from "../types/auth";
import type { WorkTime } from "../types/response";
import { dayOfYear, daysInYear } from "../utils/calendar";

export const Route = createFileRoute("/_auth/")({
  component: Home,
});

function Home() {
  const { chrono, auth } = Route.useRouteContext();
  const year = new Date().getFullYear();
  const { addErrorToast } = useToast();

  const userQ = useQuery({
    queryKey: ["user", auth.userId, "vacation", year],
    queryFn: () =>
      chrono.users.getUserById(auth.userId!, {
        year: year,
      }),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
  });

  const vacationQ = useQuery({
    queryKey: ["vacationGraph"],
    queryFn: () => chrono.events.getVacationGraph(year),
    staleTime: 1000 * 60 * 1, // 1min
    gcTime: 1000 * 60 * 30, // 30min
  });

  const aworkQ = useQuery({
    queryKey: ["awork", year],
    queryFn: () => chrono.awork.getWorkTimesforYear(year),
    staleTime: 1000 * 60 * 60, // 1h
    gcTime: 1000 * 60 * 60 * 2, // 2h
  });

  const [awork, setAwork] = useState<WorkTime | undefined>();

  const queries = [userQ, vacationQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  useEffect(() => {
    if (aworkQ.isError) {
      addErrorToast(aworkQ.error);
      return;
    }
    if (aworkQ.isPending) return;
    setAwork(aworkQ.data);
  }, [aworkQ.isError, aworkQ.isPending, aworkQ.data]);

  if (anyPending) return <LoadingSpinnerPage />;
  if (firstError) return <ErrorPage error={firstError} />;

  const user = userQ.data! as UserWithVacation;
  const vacation = vacationQ.data!;

  const daysYear = daysInYear(year);
  const currDay = dayOfYear();
  const vacRemainingPercent =
    (user.vacation_remaining / user.vacation_days) * 100;
  const yearRemainingPercent = (currDay / daysYear) * 100;
  const vacTakenPercent =
    user.vacation_remaining > user.vacation_days
      ? (user.vacation_used / (user.vacation_used + user.vacation_remaining)) *
        100
      : 100 - vacRemainingPercent;

  return (
    <div className="flex flex-col container mx-auto justify-center align-middle gap-6 p-4">
      <div className="text-[48px] pl-2 text-primary font-light mb-2">
        <span className="animate-pulse font-medium text-white pr-1"> Hej </span>
        {user?.username}
      </div>
      <TitleSection title="Your vacation">
        <StatCard>
          <StatCardElement
            title="Vacation remaining"
            subtitle={`${vacRemainingPercent.toFixed(2)}% remaining`}
          >
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-primary">
              {user.vacation_remaining} d
            </span>
          </StatCardElement>
          <StatCardElement
            title="Vacation taken"
            subtitle={`${vacTakenPercent.toFixed(2)}% taken`}
          >
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-secondary opacity-40">
              {user.vacation_used} d
            </span>
          </StatCardElement>
          <StatCardElement
            title="Vacation total"
            subtitle={`${user.vacation_days} days total`}
          >
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-secondary opacity-40">
              {user.vacation_days} d
            </span>
          </StatCardElement>
          <StatCardElement
            title="Vacation pending"
            subtitle={`${user.pending_events} events pending`}
          >
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-accent">
              <span className="animate-pulse text-primary">
                {user.pending_events}
              </span>{" "}
              d
            </span>
          </StatCardElement>
        </StatCard>
      </TitleSection>

      <TitleSection title="Your worktimes">
        {aworkQ.isPending ? (
          <div className="justify-center">
            <LoadingSpinner />
          </div>
        ) : (
          <>
            {awork && (
              <>
                <StatCard>
                  <StatCardElement
                    title="Hours horked this year"
                    subtitle={`${awork.worked.toFixed(2)}h`}
                  >
                    <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-primary">
                      {awork.worked.toFixed(2)}h
                    </span>
                  </StatCardElement>
                  <StatCardElement title="Expected work" subtitle="">
                    <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-secondary opacity-40">
                      {awork.expected * 8} h
                    </span>
                  </StatCardElement>
                  <StatCardElement
                    title="Vacation hours"
                    subtitle={`${awork.vacation * 8} h`}
                  >
                    <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-secondary opacity-40">
                      {awork.vacation * 8} h
                    </span>
                  </StatCardElement>
                  <StatCardElement title="Expected (With Vacation)" subtitle="">
                    <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-accent">
                      <span className="animate-pulse text-primary">
                        {awork.expected * 8 - awork.vacation * 8}
                      </span>{" "}
                      h
                    </span>
                  </StatCardElement>
                </StatCard>
              </>
            )}
          </>
        )}
      </TitleSection>
      <TitleSection title="Year progession">
        <StatCard>
          <StatCardElement
            title="Days this year"
            subtitle={`${daysYear} days total`}
          >
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-secondary opacity-40">
              {daysYear} d
            </span>
          </StatCardElement>
          <StatCardElement
            title="Days passed"
            subtitle={`${currDay} days passed`}
          >
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-accent">
              {currDay} d
            </span>
          </StatCardElement>
          <StatCardElement
            title="Days completed"
            subtitle={`${(100 - yearRemainingPercent).toFixed(2)}% remaining`}
          >
            <span className="-mb-1 pt-1.5 stat-value max-sm:text-2xl text-accent">
              {yearRemainingPercent.toFixed(2)} %
            </span>
          </StatCardElement>
          <StatCardElement
            title="Days progress"
            subtitle={`${(100 - yearRemainingPercent).toFixed(2)}% remaining`}
          >
            <progress
              className="progress progress-primary mb-3 mt-2 h-3.5"
              value={yearRemainingPercent}
              max="100"
              role="progressbar"
            />
          </StatCardElement>
        </StatCard>
      </TitleSection>
      <TitleSection title="Team Vacation">
        <VacationGraph
          yearOffset={vacation.year_offset}
          gaps={vacation.month_gaps}
          data={vacation.vacation_data}
        />
      </TitleSection>
    </div>
  );
}
