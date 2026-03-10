import {useMutation, useQueryClient} from '@tanstack/react-query'
import {Link, useNavigate, useParams} from '@tanstack/react-router'
import {ChronoClient} from '../api/chrono/client'
import type {User} from '../types/auth'
import type {EventUser, Month} from '../types/response'
import {hexToHSL, hsla} from '../utils/colors'
import {LoadingSpinner} from './LoadingSpinner'
import {useToast} from './Toast'

export function CalendarNavigation({
	year: currYear,
	month: currMonth,
	monthName,
	compact,
}: {
	year: number
	month: number
	monthName: string
	compact?: boolean
}) {
	let year = currYear
	let prevYear = currYear
	let nextYear = currYear
	let nextMonth = currMonth + 1
	let prevMonth = currMonth - 1
	if (prevMonth <= 0) {
		prevMonth = 12
		prevYear--
	}

	if (nextMonth > 12) {
		nextMonth = nextMonth % 12
		nextYear++
	}

	return (
		<div
			className={
				compact
					? 'flex justify-center items-center w-fit p-2 rounded-lg gap-2 pr-6 max-w-xs'
					: 'col-span-2 flex justify-start space-x-2 bg-base-300 p-2 rounded-lg gap-4'
			}
		>
			<div className='flex items-center gap-2 w-full'>
				<div className='flex justify-center items-center'>
					<Link
						to='/calendar/$year/$month'
						params={{year: prevYear.toString(), month: prevMonth.toString()}}
						className='btn btn-sm btn-soft btn-primary bg-primary/10 hover:bg-primary/25 hover:text-primary w-9 h-9 rounded-lg border justify-center items-center icon-outlined animate-color duration-300'
						search={(prev) => prev}
					>
						arrow_back
					</Link>
				</div>
				<div className='flex justify-center items-center'>
					<Link
						to='/calendar/$year/$month'
						params={{year: nextYear.toString(), month: nextMonth.toString()}}
						className='btn btn-sm btn-soft btn-primary bg-primary/10 hover:bg-primary/25 hover:text-primary w-9 h-9 rounded-lg border icon-outlined animate-all duration-300'
						search={(prev) => prev}
					>
						arrow_forward
					</Link>
				</div>
				<div className='pl-4'>
					<p className='text-sm lg:text-base whitespace-nowrap'>
						{monthName.substring(0, 3)} {year}
					</p>
				</div>
			</div>
		</div>
	)
}

export function UserFilter({
	users,
	userFilter,
	setUserFilter,
}: {
	users: User[]
	userFilter?: string
	setUserFilter: (value?: string) => void
}) {
	const navigate = useNavigate()
	const params = useParams({from: '/_auth/calendar/$year/$month'})

	return (
		<select
			defaultValue={userFilter}
			onChange={(e) => {
				const filtered = e.target.value === 'allUsers' ? undefined : e.target.value
				setUserFilter(filtered)
				navigate({
					to: '/calendar/$year/$month',
					search: (prev) => ({...prev, user: filtered}),
					params: params,
				})
			}}
			className='w-full col-span-1 cursor-pointer bg-base-300 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-12 lg:h-full pl-4 text-base border-0 rounded-lg animate-all'
		>
			<option value='allUsers'>All Users</option>
			{users.map((u, i) => (
				<option key={i} value={u.username}>
					{u.username}
				</option>
			))}
		</select>
	)
}

export function EventFilter({
	events,
	eventFilter,
	setEventFilter,
}: {
	events: string[]
	eventFilter?: string
	setEventFilter: (value: string | undefined) => void
}) {
	const navigate = useNavigate()
	const params = useParams({from: '/_auth/calendar/$year/$month'})

	return (
		<select
			className='w-full col-span-1 cursor-pointer bg-base-300 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-full pl-4 text-base border-0 rounded-lg animate-all'
			defaultValue={eventFilter}
			onChange={(e) => {
				const filtered = e.target.value === 'allEvents' ? undefined : e.target.value
				setEventFilter(filtered)
				navigate({
					to: '/calendar/$year/$month',
					search: (prev) => ({...prev, event: filtered}),
					params: params,
				})
			}}
		>
			<option value={'allEvents'}>All Events</option>
			{events.map((e, i) => (
				<option key={i} value={e.toLowerCase()}>
					{e}
				</option>
			))}
		</select>
	)
}

export function VacationCounter({
	pending,
	used,
	remaining,
}: {
	pending: number
	used: number
	remaining: number
}) {
	return (
		<div className='flex px-3 w-full justify-center text-base bg-base-300 items-center rounded-lg align-middle h-12 lg:h-full text-center'>
			<div className='w-1/3 truncate! tooltip text-accent/90' data-tip='remaining'>
				{remaining} remaining
			</div>
			<div
				className='w-1/3 lg:w-fit border-x border-info/15 px-2 mx-2 truncate! tooltip text-warning'
				data-tip='used'
			>
				{used} used
			</div>{' '}
			<div className='w-1/3 truncate! tooltip text-primary' data-tip='pending'>
				{pending} pending
			</div>
		</div>
	)
}

// export function WeekdayHeader({label, highlighted = true}: {label: string; highlighted?: boolean}) {
// 	return (
// 		<div
// 			className={`lg:block hidden truncate text-sm ${highlighted ? 'text-primary' : ''} rounded-xl p-1 text-center lg:text-lg`}
// 		>
// 			{label}
// 		</div>
// 	)
// }

export function Event({
	event,
	currUser,
	index,
	totalCount,
}: {
	event: EventUser
	currUser: User
	index: number
	totalCount: number
}) {
	const chrono = new ChronoClient()
	const {addToast, addErrorToast} = useToast()
	const queryClient = useQueryClient()

	const hsl = hexToHSL(event.user.color)
	const bgColor = hsla(...hsl, 0.2)
	const borderColor = hsla(...hsl, 0.3)

	const currDate = new Date()
	const eventDate = new Date(event.event.scheduled_at)
	const isInFuture = eventDate >= currDate
	const isFromCurrUser = event.event.user_id === currUser.id
	const isAdmin = currUser.is_superuser
	const isHoliday = event.user.id === 1
	const positionClass =
		totalCount === 1
			? 'rounded-xl border'
			: index === 0
				? 'rounded-t-xl border-l border-r border-t'
				: index === totalCount - 1
					? 'rounded-b-xl border-l border-r border-b '
					: 'rounded-none border-l border-r'
	const positionDeleteClass =
		totalCount === 1
			? 'hover:rounded-xl rounded-xl'
			: index === 0
				? 'hover:rounded-t-xl hover:rounded-b-none rounded-xl'
				: index === totalCount - 1
					? 'hover:rounded-b-xl hover:rounded-t-none rounded-xl'
					: 'hover:rounded-none rounded-xl'
	const shortName = isHoliday
		? event.event.name
		: chrono.events.getShortEventName(event.event.name)

	const isDeletable =
		isAdmin || (isFromCurrUser && (isInFuture || event.event.state !== 'accepted'))

	const mutation = useMutation({
		mutationKey: ['deleteEvent', event.event.id],
		mutationFn: () => chrono.events.deleteEvent(event.event.id),
		onSuccess: () => {
			addToast(`Successfully deleted event ${event.event.name}`, 'success')
			queryClient.invalidateQueries({queryKey: ['month']})
		},
		onError: (error) => addErrorToast(error),
		retry: false,
	})

	return (
		<div className='indicator w-full'>
			<div
				style={{backgroundColor: bgColor, borderColor: borderColor}}
				className={`group gap-2 relative ${!isHoliday && 'flex'} text-center py-1 w-full ${positionClass}`}
			>
				{isDeletable && (
					<>
						<span className='flex items-center justify-center text-transparent rounded-lg group-hover:mix-blend-revert group-hover:text-base-content absolute top-0 left-0 w-full h-full icon-outlined animate-all'>
							<button
								onClick={() => mutation.mutate()}
								className={`${positionDeleteClass} flex items-center justify-center h-7 w-16 border border-transparent hover:border-error text-center cursor-pointer icon-outlined group-hover:bg-neutral/30 group-hover:text-base-content hover:text-error-content hover:duration-1000 hover:bg-error/90 hover:w-full hover:h-full animate-all z-999`}
							>
								delete
							</button>
						</span>
					</>
				)}
				<div
					className={`${!isHoliday ? 'bg-black/30 min-w-10 rounded-lg px-2' : 'truncate! whitespace-nowrap px-2 '} ml-1 text-base-content ${isDeletable ? 'group-hover:opacity-0' : ''} animate-all`}
				>
					{shortName}
				</div>
				{!isHoliday && (
					<div className='flex justify-between items-center w-full'>
						<div
							className={`text-accent/60 text-[15px] truncate ${isDeletable && 'group-hover:text-white/0'} animate-all`}
						>
							{event.user.username}
						</div>

						<div
							className={`flex mr-1 w-3 aspect-square rounded-full items-center ${isDeletable && 'group-hover:opacity-0'}`}
						>
							<span
								className={
									event.event.state === 'pending'
										? 'bg-accent status status-sm status-accent animate-ping'
										: event.event.state === 'declined'
											? 'status status-md status-error'
											: 'status status-md status-success'
								}
							></span>
						</div>
					</div>
				)}
			</div>
		</div>
	)
}

export function Day({
	date,
	day,
	month,
	year,
	events,
	selectedEvent,
	currUser,
}: {
	date: number
	day: string
	month: number
	year: number
	events: EventUser[]
	selectedEvent: string
	currUser: User
}) {
	const now = new Date()
	const isToday =
		now.getDate() === date && now.getMonth() + 1 === month && now.getFullYear() === year

	const chrono = new ChronoClient()
	const queryClient = useQueryClient()
	const {addToast, addErrorToast} = useToast()

	const mutation = useMutation({
		mutationKey: ['createEvent', year, month, date, selectedEvent],
		retry: false,
		mutationFn: () =>
			chrono.events.createEvent({
				year: year,
				month: month,
				day: date,
				event: selectedEvent,
			}),
		onSuccess: () => {
			addToast(`Successfully created event ${selectedEvent}`, 'success')
			queryClient.invalidateQueries({queryKey: ['month']})
		},
		onError: (error) => addErrorToast(error),
	})

	return (
		<div
			className={
				isToday
					? 'border bg-primary/35 border-primary/1 max-lg:rounded-xl flex flex-col overflow-hidden'
					: 'bg-base-300 max-lg:rounded-xl flex flex-col overflow-hidden'
			}
		>
			<div className='p-2 pr-4 '>
				<div className='flex justify-between items-center'>
					<div className='text-base-content/85 text-md w-9.5 h-7 flex items-center justify-center rounded-xl bg-base-200/50 font-medium'>
						{date}
					</div>
					<span className='lg:hidden pl-1.75 text-info/70 tracking-wide'>{day}</span>
					<span className='hidden lg:inline pl-1.75 text-info/70 text-sm'>
						{day.substring(0, 3)}
					</span>
				</div>

				{/* <div className='hidden lg:block'>
					{date}{' '}
					<span className='pl-1 pt-1.25 text-accent/40 text-[14px]'>
						
					</span>
				</div> */}
			</div>
			<div className='flex items-end flex-col px-2 h-full bg-base-200/70 rounded-t-none'>
				<div className='flex flex-col gap-0.5 flex-1 mt-4 w-full'>
					{events.map((e, i) => (
						<Event
							key={i}
							event={e}
							currUser={currUser}
							index={i}
							totalCount={events.length}
						/>
					))}
				</div>
				<button
					onClick={() => mutation.mutate()}
					className='mb-1.75 mt-2.75 p-1.25 cursor-pointer transition-all duration-300 border-primary/0 hover:border-primary/50 border border-dashed hover:border-dashed hover:bg-base-200/20 flex items-center justify-center focus-within:outline-primary/50 bg-base-200 lg:bg-base-200/80 rounded-xl text-primary/15 hover:text-primary w-full hover:icon-filled'
				>
					{mutation.isPending ? (
						<div className='w-5 flex'>
							<LoadingSpinner />
						</div>
					) : (
						<span className='icon-outlined hover:icon-outlined leading-5'>add</span>
					)}
				</button>
			</div>
		</div>
	)
}

export function Calendar({
	month,
	eventFilter,
	userFilter,
	selectedEvent,
	currUser,
}: {
	month: Month
	eventFilter: string | undefined
	userFilter: string | undefined
	selectedEvent: string
	currUser: User
}) {
	return (
		<div className='my-12 rounded-xl *:border *:border-base-200 *:backdrop-brightness-115 lg:my-8 lg:mt-0 mx-auto grid grid-cols-1 max-lg:gap-y-4 lg:gap-px lg:grid-cols-7 overflow-x-scroll'>
			{Array.from({length: month.offset}).map((_, i) => (
				<div key={i} className='hidden lg:block'></div>
			))}
			{month.days.map((d, i) => {
				const date = new Date(d.date)

				return (
					<Day
						selectedEvent={selectedEvent}
						key={i}
						date={d.number}
						day={d.name}
						month={date.getMonth() + 1}
						year={date.getFullYear()}
						currUser={currUser}
						events={
							d.events?.filter((e) => {
								let event = true
								let user = true
								if (eventFilter && e.user.id !== 1) {
									event = e.event.name === eventFilter
								}
								if (userFilter && e.user.id !== 1) {
									user = e.user.username === userFilter
								}

								return event && user
							}) || []
						}
					/>
				)
			})}
		</div>
	)
}
