import {
  createSignal,
  type Component,
  type JSX,
  Index,
  type ComponentProps,
  Switch,
  Match,
} from 'solid-js'
import Editor from '@components/editor'
import Line from '@components/line'
import InputLine from '@components/input-line'
import { hc } from 'hono/client'
import type { Routes } from '@terminal/functions/api'

type SshProps = { apiUrl: string } & JSX.HTMLAttributes<HTMLDivElement>
type InputLineProps = ComponentProps<typeof InputLine>
type State = InputLineProps['state']

const SshComponent: Component<SshProps> = (props) => {
  let sshTimeout: ReturnType<typeof setTimeout> | undefined

  const [state, setState] = createSignal<State>('normal')

  const restore = (timeout = 2500) => {
    if (sshTimeout) clearTimeout(sshTimeout)

    sshTimeout = setTimeout(() => {
      setState('normal')
    }, timeout)
  }

  const copy = () => {
    if (sshTimeout) clearTimeout(sshTimeout)

    navigator.clipboard.writeText('ssh terminal.shop')
    setState('success')
    restore()
  }

  const client = hc<Routes>(props.apiUrl)
  let stateTimeout: ReturnType<typeof setTimeout> | undefined

  const [emailState, setEmailState] = createSignal<State>('normal')
  const [emailMessage, setEmailMessage] = createSignal<string>()
  const [lines, setLines] = createSignal<InputLineProps[]>([])

  const restoreEmail = (timeout = 2500) => {
    if (stateTimeout) clearTimeout(stateTimeout)

    stateTimeout = setTimeout(() => {
      setEmailState('normal')
      setEmailMessage(undefined)
    }, timeout)
  }

  const submitEmail = async (email: string) => {
    if (
      !email ||
      !email.match(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/)
    ) {
      setEmailState('error')
      setEmailMessage('# invalid email address')
      restoreEmail()
      return
    }

    setEmailState('busy')

    try {
      await client.email.subscription.$post({ json: { email } })
      setEmailState('success')
      setEmailMessage('# signed up')
      setLines([{ state: 'normal' }])
    } catch (err) {
      console.error(err)
      setEmailState('error')
      setEmailMessage((err as Error).message)
      restoreEmail()
      return
    }
  }

  return (
    <Editor>
      <Line>
        <h1>
          # use the command below to order your delicious whole bean coffee
        </h1>
      </Line>
      <Line state={state()} class="group/ssh" onClick={copy}>
        <span class="text-white">ssh terminal.shop</span>
        <div
          id="ssh-supplement"
          classList={{
            'hidden w-5 h-5 mx-6 text-gray-10 shrink-0': true,
            'group-hover/ssh:block group-active/ssh:text-gray-11 group-focus/ssh:block group-focus/ssh:text-gray-11':
              true,
          }}
        >
          <Switch>
            <Match when={state() === 'normal'}>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="w-5 h-5"
                viewBox="0 0 20 20"
                fill="none"
              >
                <path
                  d="M6.45833 6.45833V3.125H16.875V13.55H13.5417M13.5417 6.45833V16.875H3.125V6.45833H13.5417Z"
                  class="stroke-current"
                  stroke-width="1.5"
                  stroke-linecap="round"
                />
              </svg>
            </Match>
            <Match when={state() === 'success'}>
              <svg
                class="text-green-11 w-5 h-5"
                viewBox="0 0 20 20"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <g id="check-small, checkmark-small">
                  <path
                    id="vector"
                    class="stroke-current"
                    d="M5 11.9651L8.37838 14.7522L15 5.83331"
                    stroke-width="1.5"
                    stroke-linecap="square"
                  />
                </g>
              </svg>
            </Match>
          </Switch>
        </div>
        <span
          id="ssh-message"
          class="hidden group-hover/ssh:block group-focus/ssh:block whitespace-nowrap"
        >
          <Switch>
            <Match when={state() === 'normal'}># copy to clipboard</Match>
            <Match when={state() === 'success'}># copied to clipboard</Match>
          </Switch>
        </span>
      </Line>
      <Line>cat ~/.ssh/known_hosts</Line>
      <Line>
        <span
          class="break-all"
          innerText="terminal.shop ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEzsOgEiuiTQEUZnMORRmhMHDSAo8VBUl/g55Ec6ZaKM"
        ></span>
      </Line>
      <Line />
      <Line>
        <label id="email-label" for="email-input">
          # sign up for updates, enter your email address below...
        </label>
      </Line>
      <InputLine
        autofocus
        state={emailState()}
        message={emailMessage()}
        labelledby="email-label"
        onReturn={submitEmail}
      />
      <Index each={lines()}>{(line) => <InputLine {...line()} />}</Index>
    </Editor>
  )
}

export default SshComponent
