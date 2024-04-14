import { createSignal, type Component, type JSX } from 'solid-js'
import Line, { type State } from '@components/line'
import Editor from '@components/editor'
import Caret from '@components/caret'
import InputLine from '@components/input-line'

type ShopProps = {} & JSX.HTMLAttributes<HTMLDivElement>

const ShopComponent: Component<ShopProps> = () => {
  let sshTimeout: ReturnType<typeof setTimeout> | undefined

  const [state, setState] = createSignal<State>('normal')
  const [message, setMessage] = createSignal<string>()

  const restore = (timeout = 2500) => {
    if (sshTimeout) clearTimeout(sshTimeout)

    sshTimeout = setTimeout(() => {
      setState('normal')
      setMessage(undefined)
    }, timeout)
  }

  const copy = () => {
    console.log('copy')
    if (sshTimeout) clearTimeout(sshTimeout)

    navigator.clipboard.writeText('ssh terminal.shop')
    setState('success')
    setMessage('# copied to clipboard')
    restore()

    // sshLine.dataset.state = 'success'
    // sshMessage.innerText = '# copied to clipboard'
    // sshSupplement.classList.add('!block')
    // sshMessage.classList.add('!block')
    // copyIcon.classList.add('hidden')
    // checkIcon.classList.remove('hidden')

    // sshTimeout = setTimeout(() => {
    // sshLine.removeAttribute('data-state')
    // sshMessage.innerHTML = originalCopyMessage
    // sshMessage.classList.remove('!block')
    // sshSupplement.classList.remove('!block')
    // copyIcon.classList.remove('hidden')
    // checkIcon.classList.add('hidden')
    // }, 2500)
  }

  return (
    <Editor>
      <Line number={1}>
        <h1>
          # use the command below to order your delicious 12oz bag of Nil Blend
          coffee
        </h1>
      </Line>
      <InputLine
        readonly
        state={state()}
        message={message()}
        value="ssh terminal.shop"
        number={2}
        class="group/ssh"
        onReturn={copy}
        onClick={copy}
      >
        <div
          id="ssh-supplement"
          classList={{
            'hidden w-5 h-5 mx-6 text-gray-10 shrink-0': true,
            'group-hover/ssh:block group-active/ssh:text-gray-11 group-focus/ssh:block group-focus/ssh:text-gray-11':
              true,
          }}
        >
          <svg
            id="copy-icon"
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
          <svg
            id="check-icon"
            class="hidden text-green-11 w-5 h-5"
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
        </div>
        <span
          id="ssh-message"
          class="hidden group-hover/ssh:block group-focus/ssh:block whitespace-nowrap"
        >
          # copy to clipboard
        </span>
      </InputLine>
    </Editor>
  )
}

export default ShopComponent
