export function Header() {
  return (
    <div className="mb-4 mx-auto p-4 lg:px-4">
      <div className="navbar flex justify-between">
        <div className="flex items-center">
          <div className="pr-14">
            <img className="w-40" alt="" src="chrono.svg" />
          </div>
          <div
            className="!z-20 max-lg:dock max-lg:border-t max-lg:border-accent/15 max-lg:!bg-base-100/50 backdrop-blur-xl overflow-x-auto flex gap-4 lg:w-fit 
						*:flex *:!flex-col *:lg:!flex-row *:lg:gap-2 *:lg:items-center"
          >
            <a href="/">
              <span className="icon-outlined">home</span>
              <span className="font-medium text-base">Home</span>
            </a>
            <a href={"/"}>
              <span className="icon-outlined">calendar_today</span>
              <span className="font-medium text-base">Calendar</span>
            </a>
            <a href="/team">
              <span className="icon-outlined">group</span>
              <span className="font-medium text-base">Team</span>
            </a>
          </div>
        </div>
        <div className="flex items-center justify-end gap-6">
          <a href="/login" className="btn btn-ghost">
            Login
          </a>
          <a href="/signup" className="btn btn-ghost">
            Signup
          </a>
        </div>
      </div>
    </div>
  );
}
