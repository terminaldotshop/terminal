import {
  createSignal,
  type Component,
  type JSX,
  type ComponentProps,
  Index,
} from 'solid-js'
import Line from '@components/line'
import Editor from '@components/editor'
import InputLine from '@components/input-line'
import { hc } from 'hono/client'
import type { AppType } from '../../../workers/src/api.ts'

type PendingProps = {} & JSX.HTMLAttributes<HTMLDivElement>
type InputLineProps = ComponentProps<typeof InputLine>
type State = InputLineProps['state']

const PendingComponent: Component<PendingProps> = () => {
  const client = hc<AppType>(import.meta.env.PUBLIC_API_URL)
  let stateTimeout: ReturnType<typeof setTimeout> | undefined

  const [state, setState] = createSignal<State>('normal')
  const [message, setMessage] = createSignal<string>()
  const [lines, setLines] = createSignal<InputLineProps[]>([])

  const restore = (timeout = 2500) => {
    if (stateTimeout) clearTimeout(stateTimeout)

    stateTimeout = setTimeout(() => {
      setState('normal')
      setMessage(undefined)
    }, timeout)
  }

  const submitEmail = async (email: string) => {
    if (
      !email ||
      !email.match(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/)
    ) {
      setState('error')
      setMessage('# invalid email address')
      restore()
      return
    }

    setState('busy')

    try {
      await client.api.subscription.$post({
        json: { email },
      })
      setState('success')
      setMessage('# signed up')
      setLines([{ state: 'normal' }])
    } catch (err) {
      console.error(err)
      setState('error')
      setMessage((err as Error).message)
      restore()
      return
    }
  }

  return (
    <Editor>
      <Line number={1}>
        <h1># nil blend has almost finished roasting</h1>
      </Line>
      <Line number={2}>
        <h2># order your coffee via the command line later this week</h2>
      </Line>
      <Line number={3}>
        <label id="email-label" for="email-input">
          # be the first to know, enter your email address below...
        </label>
      </Line>
      <InputLine
        autofocus
        state={state()}
        message={message()}
        number={4}
        labelledby="email-label"
        onReturn={submitEmail}
      />
      <Index each={lines()}>
        {(line, i) => <InputLine {...line()} number={i + 5} />}
      </Index>
    </Editor>
  )
}

export default PendingComponent
