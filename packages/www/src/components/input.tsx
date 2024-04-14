import {
  type Component,
  type JSX,
  createSignal,
  Show,
  Switch,
  Match,
  createEffect,
} from 'solid-js'
import Caret from '@components/caret'
import { autofocus } from '@solid-primitives/autofocus'
import type { State } from '@components/line'

// ensures it doesn't get tree shaken
autofocus

type InputProps = {
  ref?: HTMLElement
  autofocus?: boolean
  labelledby?: string
  state?: State
  value?: string
  message?: string
  result?: string
  readonly?: boolean
  onReturn?: (value: string) => void
} & JSX.HTMLAttributes<HTMLSpanElement>

const InputComponent: Component<InputProps> = (props) => {
  let ref: HTMLSpanElement | undefined

  const [visible, setVisible] = createSignal<boolean>(true)
  const [before, setBefore] = createSignal<string>()
  const [after, setAfter] = createSignal<string>()
  const [blink, setBlink] = createSignal<boolean>(true)

  let blinkTimeout: ReturnType<typeof setTimeout> | undefined = undefined

  const update = () => {
    if (props.readonly) return

    const selection = document.getSelection()
    const visible = selection?.isCollapsed ?? false
    const position = visible ? selection?.anchorOffset : undefined

    if (blinkTimeout) clearTimeout(blinkTimeout)
    setBlink(false)

    if (position !== undefined) {
      const beforeText = ref?.innerText
        .substring(0, position)
        .replace(/ /g, '&nbsp')
      const afterText = ref?.innerText
        .substring(position)
        .replace(/ /g, '&nbsp')

      setBefore(beforeText)
      setAfter(afterText)
    }

    setVisible(visible)

    blinkTimeout = setTimeout(() => {
      setBlink(true)
    }, 200)
  }

  createEffect(() => {
    if (props.state === 'normal') ref?.focus()
  }, [props.state])

  const submit = (ev: KeyboardEvent) => {
    if (props.state === 'busy') return

    const span = ev.target as HTMLSpanElement
    if (ev.key === 'Enter') {
      ev.preventDefault()

      if (props.onReturn && span.innerText) props.onReturn(span.innerText)
    }
  }

  return (
    <div
      classList={{
        ...props.classList,
        'relative group/input': true,
        'flex items-center overflow-x-scroll no-scrollbar pr-[10px]': true,
        '!overflow-visible': props.state !== 'normal' || props.readonly,
        [props.class ?? '']: !!props.class,
      }}
    >
      <span
        use:autofocus
        contenteditable={props.state === 'normal' && !props.readonly}
        {...props}
        ref={ref}
        role="textbox"
        aria-labelledby={props.labelledby}
        classList={{
          'text-white leading-10 flex gap-2 flex-wrap': true,
          'focus:outline-none whitespace-nowrap caret-transparent': true,
          hidden: props.state !== 'normal',
        }}
        innerText={props.readonly ? props.value : ''}
        onInput={update}
        onKeyUp={update}
        onSelect={update}
        onMouseMove={update}
        onMouseUp={update}
        onTouchStart={update}
        onPaste={update}
        onCut={update}
        onKeyPress={submit}
        // onChange={update}
        // onKeyDown={update}
      ></span>
      <Show when={props.state !== 'normal'}>
        <span class="text-white leading-10 flex gap-2 flex-wrap">
          <div classList={{ 'animate-shake': props.state === 'error' }}>
            {props.value}
          </div>
          <div
            classList={{
              'w-4 h-4 self-center': true,
              'text-blue-11': props.state === 'busy',
              'text-green-11': props.state === 'success',
              'text-red-11': props.state === 'error',
            }}
          >
            <Switch>
              <Match when={props.state === 'busy'}>
                <svg
                  class="animate-spin h-4 w-4 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  ></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
              </Match>
              <Match when={props.state === 'success'}>
                <svg
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
              <Match when={props.state === 'error'}>
                <svg
                  viewBox="0 0 15 15"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    class="fill-current"
                    d="M11.7816 4.03157C12.0062 3.80702 12.0062 3.44295 11.7816 3.2184C11.5571 2.99385 11.193 2.99385 10.9685 3.2184L7.50005 6.68682L4.03164 3.2184C3.80708 2.99385 3.44301 2.99385 3.21846 3.2184C2.99391 3.44295 2.99391 3.80702 3.21846 4.03157L6.68688 7.49999L3.21846 10.9684C2.99391 11.193 2.99391 11.557 3.21846 11.7816C3.44301 12.0061 3.80708 12.0061 4.03164 11.7816L7.50005 8.31316L10.9685 11.7816C11.193 12.0061 11.5571 12.0061 11.7816 11.7816C12.0062 11.557 12.0062 11.193 11.7816 10.9684L8.31322 7.49999L11.7816 4.03157Z"
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                  ></path>
                </svg>
              </Match>
            </Switch>
          </div>
          <span class="text-gray-10">{props.message}</span>
        </span>
      </Show>
      <div
        class="absolute inset-0 flex items-center pointer-events-none"
        aria-hidden="true"
      >
        <div class="flex items-center leading-10 whitespace-nowrap focus:outline-none">
          <span
            class="text-transparent"
            innerHTML={before() ?? (props.readonly ? props.value : '')}
          ></span>
          <Caret
            blink={blink()}
            classList={{
              'hidden group-has-[:focus]/input:block': true,
              '!block ml-1.5': props.readonly,
              '!hidden': !visible() || props.state !== 'normal',
            }}
          />
          <span class="text-transparent" innerHTML={after()}></span>
        </div>
      </div>
    </div>
  )
}

export default InputComponent
