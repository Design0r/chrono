import {isWeekend} from '../../lib/timestamp-utils'

const NORMAL_HOURS = 8
const OVERTIME_CAP_HOURS = 4
const BAR_TOTAL_HOURS = NORMAL_HOURS + OVERTIME_CAP_HOURS
const NORMAL_FRACTION = NORMAL_HOURS / BAR_TOTAL_HOURS
const OVERTIME_FRACTION = OVERTIME_CAP_HOURS / BAR_TOTAL_HOURS

/** Horizontal bar: 8h normal (solid left, dashed right). &lt;8h: turquoise gap. &gt;8h: overtime bar. Weekend: overtime bar only. */
export function DayBar({
	totalSeconds,
	dayDate,
}: {
	totalSeconds: number
	dayDate: Date
}) {
	const weekend = isWeekend(dayDate)
	const normalSeconds = NORMAL_HOURS * 3600
	const overtimeCapSeconds = OVERTIME_CAP_HOURS * 3600
	const barTotalSeconds = BAR_TOTAL_HOURS * 3600

	const workedNormal = Math.min(totalSeconds, normalSeconds)
	const workedOvertime = Math.min(Math.max(0, totalSeconds - normalSeconds), overtimeCapSeconds)

	const normalZonePercent = NORMAL_FRACTION * 100
	const redWidthPercent = (workedOvertime / overtimeCapSeconds) * OVERTIME_FRACTION * 100

	const underEight = workedNormal < normalSeconds
	const turquoiseStartPercent = (workedNormal / normalSeconds) * NORMAL_FRACTION * 100
	const turquoiseWidthPercent = underEight
		? ((normalSeconds - workedNormal) / normalSeconds) * NORMAL_FRACTION * 100
		: 0

	const weekendBarWidthPercent = weekend
		? Math.min(100, (totalSeconds / barTotalSeconds) * 100)
		: 0

	if (weekend) {
		return (
			<div className='w-[80%] min-h-6 flex items-center' aria-hidden>
				<div className='relative h-5 w-full'>
					{weekendBarWidthPercent > 0 && (
						<div
							className='absolute top-0 left-0 bottom-0 rounded-r-sm border-r-2 border-[hsl(344,84%,60%)]'
							style={{
								width: `${weekendBarWidthPercent}%`,
								background:
									'linear-gradient(90deg, hsl(222 30% 28% / 0.0) 0%, hsl(344 63% 36% / 1) 100%)',
							}}
						/>
					)}
				</div>
			</div>
		)
	}

	return (
		<div className='w-[80%] min-h-6 flex items-center' aria-hidden>
			<div
				className='relative h-5 w-full border-l border-info/50'
				style={{
					background: `linear-gradient(90deg, #1B213344 0%, #242C41 ${normalZonePercent}%, transparent ${normalZonePercent}%)`,
				}}
			>
				<div
					className='absolute top-0 bottom-0 w-0 border-l border-dashed border-white'
					style={{left: `${normalZonePercent}%`}}
				/>
				{turquoiseWidthPercent > 0 && (
					<div
						className='absolute top-0 bottom-0'
						style={{
							left: `${turquoiseStartPercent}%`,
							width: `${turquoiseWidthPercent}%`,
							background:
								'linear-gradient(90deg, #44E1EF 0%, #44E1EF 2px, #227495 2px, #12161f 100%)',
						}}
					/>
				)}
				{redWidthPercent > 0 && (
					<div
						className='absolute top-0 bottom-0 border-r-2 border-[hsl(344,84%,60%)]'
						style={{
							left: `${normalZonePercent}%`,
							width: `${redWidthPercent}%`,
							background:
								'linear-gradient(90deg, hsl(222 30% 28% / 0.0) 0%, hsl(344 63% 36% / 1) 100%)',
						}}
					/>
				)}
			</div>
		</div>
	)
}
