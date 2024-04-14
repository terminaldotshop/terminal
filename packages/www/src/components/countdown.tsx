import { createCountdownFromNow } from '@solid-primitives/date'
import { type Component, type JSX } from 'solid-js'

type CountdownProps = {
  date: string | Date
} & JSX.HTMLAttributes<HTMLDivElement>

const CountdownComponent: Component<CountdownProps> = (props) => {
  const [countdown] = createCountdownFromNow(props.date)

  return (
    <div {...props} classList={{ [props.class ?? '']: !!props.class }}>
      {`${countdown.days}d ${countdown.hours}h ${countdown.minutes}m ${countdown.seconds}s`}
    </div>
  )
}

export default CountdownComponent
