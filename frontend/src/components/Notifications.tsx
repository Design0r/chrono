import { Notification } from "../types/response";

type NotificationsProps = {
  notifications: Notification[];
};

export function Notifications({
  notifications,
}: {
  notifications: Notification;
}) {
  <div className="indicator">
    @NotificationIndicator(len(notifications))
    <details className="dropdown dropdown-end">
      <summary
        hx-get="/notifications"
        hx-swap="outerHTML"
        hx-target="#notification-container"
        role="button"
        className="btn btn-ghost px-5 border-1.5 border-white/2 py-1 hover:bg-info/20 rounded-full text-xl icon-outlined bg-base-100 animate-color"
      >
        notifications
      </summary>
      @NotificationContainer([]domain.Notification{})
    </details>
  </div>;
}

export function NotificationIndicator({value}: {value: number}) {
	if (value > 0) {
		<span
			className="indicator-item font-bold rounded-full align-items-start text-white/85 border border-error/50 badge badge-error backdrop-blur-md bg-error/40 p-0 h-6 px-2 pointer-events-none"
		>{ value }</span>
	} else {
		<span></span>
	}
}

templ NotificationContainer(notifications []domain.Notification) {
	<ul id="notification-container" tabindex="0" class="mt-1.5 min-w-64 pt-4 pb-3 px-3 dropdown-content menu bg-info/20 backdrop-blur-xl rounded-box z-10 drop-shadow-xl">
		<p class="px-3 pb-2 text-lg font-bold">Notifications</p>
		<hr class="border-base-200/80 pb-2"/>
		for _, n := range notifications {
			<li id={ fmt.Sprintf("notification-%v", n.ID) } class="py-1">
				<div class="hover:text-white">
					<p>{ n.Message }</p>
					<button
						hx-patch={ fmt.Sprintf("/notifications/%v", n.ID) }
						hx-swap="delete"
						hx-target={ fmt.Sprintf("#notification-%v", n.ID) }
						class="btn btn-soft border-0 hover:border-0 bg-base-200/50 text-xl font-semibold icon-outlined hover:bg-primary hover:text-neutral animate-color"
					>close</button>
				</div>
			</li>
		}
		<button
			class="mt-4 btn btn-soft rounded-xl text-neutral border-0 hover:border-0 bg-primary/90 font-semibold hover:bg-primary hover:text-neutral animate-color"
			hx-patch="/notifications"
			hx-swap="delete"
			hx-target="#notification-container"
		>Clear All</button>
	</ul>
}