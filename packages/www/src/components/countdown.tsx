import { createCountdownFromNow } from '@solid-primitives/date'
import { type Component, type JSX } from 'solid-js'

type CountdownProps = {
  date: string | Date
} & JSX.HTMLAttributes<HTMLDivElement>

const CountdownComponent: Component<CountdownProps> = (props) => {
  const [countdown] = createCountdownFromNow(props.date)

  return (
    <div {...props} classList={{ [props.class ?? '']: !!props.class }}>
      {`${Math.max(countdown.days ?? 0, 0)}d ${Math.max(countdown.hours ?? 0, 0)}h ${Math.max(countdown.minutes ?? 0, 0)}m ${Math.max(countdown.seconds ?? 0, 0)}s`}
    </div>
  )
}

export default CountdownComponent
