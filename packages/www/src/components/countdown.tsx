import { createCountdownFromNow } from '@solid-primitives/date'
import { createSignal, onMount, type Component, type JSX } from 'solid-js'

type CountdownProps = {
  date: string | Date
} & JSX.HTMLAttributes<HTMLDivElement>

const CountdownComponent: Component<CountdownProps> = (props) => {
  const [ready, setReady] = createSignal(false)
  const [countdown] = createCountdownFromNow(props.date)

  onMount(() => setReady(true))

  return (
    <div
      {...props}
      classList={{ [props.class ?? '']: !!props.class, hidden: !ready() }}
    >
      {`${Math.max(countdown.days ?? 0, 0)}d ${Math.max(countdown.hours ?? 0, 0)}h ${Math.max(countdown.minutes ?? 0, 0)}m ${Math.max(countdown.seconds ?? 0, 0)}s`}
    </div>
  )
}

export default CountdownComponent
