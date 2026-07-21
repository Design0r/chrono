import {useEffect, useRef, useState} from 'react'
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

	// Ref, damit ein neu erzeugter Callback das Interval nicht neu aufsetzt.
	const onUpdateRef = useRef(onUpdate)
	useEffect(() => {
		onUpdateRef.current = onUpdate
	})

	useEffect(() => {
		function tick() {
			const elapsedSeconds = paused ? 0 : (Date.now() - startUnix) / 1000
			onUpdateRef.current(elapsedSeconds)
			setTimer(secondsToCounter(elapsedSeconds))
		}

		tick()

		if (paused) return

		const interval = setInterval(tick, 1000)

		// Hintergrund-Tabs drosseln Intervalle; beim Zurückkehren sofort neu rechnen,
		// damit die Anzeige nicht kurzzeitig veraltet wirkt.
		function onVisible() {
			if (document.visibilityState === 'visible') tick()
		}
		document.addEventListener('visibilitychange', onVisible)
		window.addEventListener('focus', onVisible)

		return () => {
			clearInterval(interval)
			document.removeEventListener('visibilitychange', onVisible)
			window.removeEventListener('focus', onVisible)
		}
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
