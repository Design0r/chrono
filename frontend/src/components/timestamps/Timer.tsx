import {useEffect, useState} from 'react'
import {formatCounter, secondsToCounter, type TimeCounter} from '../../lib/timestamp-utils'

export function Timer({
	startUnix,
	paused,
	onUpdate,
}: {
	startUnix: number
	paused: boolean
	onUpdate: (seconds: number) => void
}) {
	const [timer, setTimer] = useState<TimeCounter>(() => secondsToCounter(0))

	useEffect(() => {
		function tick() {
			const elapsedSeconds = (Date.now() - startUnix) / 1000
			onUpdate(elapsedSeconds)
			setTimer(secondsToCounter(elapsedSeconds))
		}

		tick()

		if (paused) return

		const interval = setInterval(tick, 1000)
		return () => clearInterval(interval)
	}, [startUnix, paused])

	const f = formatCounter(timer)

	return (
		<div
			className={`font-mono text-3xl tabular-nums tracking-tight text-center transition-opacity duration-300 ${
				paused
					? 'text-base-content/60'
					: 'text-error *:even:text-error/70 *:even:animate-pulse'
			}`}
		>
			<span>{f.hours}</span>
			<span>:</span>
			<span>{f.minutes}</span>
			<span>:</span>
			<span>{f.seconds}</span>
		</div>
	)
}
