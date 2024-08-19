import {
  createSignal,
  type Component,
  type JSX,
  type ParentProps,
  Show,
} from 'solid-js'
import Line, { type State } from '@components/line'
import Editor from '@components/editor'
import Caret from './caret'

type EasterEggSegfaultProps = {} & JSX.HTMLAttributes<HTMLDivElement>
type SpanProps = ParentProps & JSX.HTMLAttributes<HTMLSpanElement>

const B: Component<SpanProps> = (props) => {
  return (
    <span class="text-blue-11" {...props}>
      {props.children}
    </span>
  )
}

const R: Component<SpanProps> = (props) => {
  return (
    <span class="text-red-11" {...props}>
      {props.children}
    </span>
  )
}

const P: Component<SpanProps> = (props) => {
  return (
    <span class="text-purple" {...props}>
      {props.children}
    </span>
  )
}

const EasterEggSegfaultComponent: Component<EasterEggSegfaultProps> = () => {
  const [state, setState] = createSignal<State>('normal')
  const report = () => {
    setState('error')
  }

  return (
    <Editor class="font-light">
      <div class="fixed top-1/3 pointer-events-none inset-x-0 h-screen bg-white mix-blend-difference"></div>
      <Line>
        <p class="text-white">
          <R>panic:</R> runtime <R>error:</R> invalid memory <B>address</B> or
          nil pointer <B>dereference</B>
        </p>
      </Line>
      <Line>
        <p class="text-white">
          [<R>signal SIGSEGV:</R> segmentation violation <B>code=0x1</B>{' '}
          addr=0x0 pc=
          <P>0x4808f6</P>]
        </p>
      </Line>
      <Line />
      <Line>
        <p class="text-white">
          <B>goroutine</B> 1 [<P>running</P>]:
        </p>
      </Line>
      <Line>
        <p class="text-white">
          coffee.<B>main()</B>
        </p>
      </Line>
      <Line>
        <p class="text-white">
          <span class="whitespace-pre">{`\t`}</span>
          <span>
            <B>$HOME</B>
            <span class="break-all">/code/terminal.shop/cmd/coffee.go:</span>
            <P>1337</P> <B>+0x16</B>
          </span>
        </p>
      </Line>
      <Line>
        <p class="text-white">
          <B>exit</B> status <P>2</P>
        </p>
      </Line>
      <Line state={state()} href="/report" onClick={report}>
        <div class="flex gap-2 items-center">
          <R>Report Bug?</R> <Caret class="bg-red-11" />
          <Show when={state() === 'error'}>
            <div class="flex gap-6 items-center">
              <svg
                viewBox="0 0 18 18"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                class="stroke-inherit size-4"
              >
                <path
                  d="M3.47538 7.17173L1.50008 6.49984M3.47538 10.4582H1.29175M3.47538 13.5363L1.50008 14.2082M14.5251 7.17173L16.5001 6.49984M14.5251 10.4582H16.7087M14.5251 13.5363L16.5001 14.2082M9.00008 10.4582V16.2915M5.45841 5.24984V4.83317C5.45841 2.87716 7.04407 1.2915 9.00008 1.2915C10.9561 1.2915 12.5417 2.87716 12.5417 4.83317V5.24984M14.2084 5.45817H3.79175V11.4998C3.79175 14.3763 6.1236 16.7082 9.00008 16.7082C11.8766 16.7082 14.2084 14.3763 14.2084 11.4998V5.45817Z"
                  stroke="currentColor"
                  stroke-opacity="0.62"
                  stroke-width="1.5"
                  stroke-linecap="square"
                />
              </svg>
              # thank you for reporting
            </div>
          </Show>
        </div>
      </Line>
    </Editor>
  )
}

export default EasterEggSegfaultComponent
