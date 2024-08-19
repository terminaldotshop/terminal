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

type EasterEggDarkModeProps = {} & JSX.HTMLAttributes<HTMLDivElement>
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

const EasterEggDarkModeComponent: Component<EasterEggDarkModeProps> = () => {
  const [state, setState] = createSignal<State>('normal')
  const report = () => {
    setState('error')
  }

  return (
    <Editor class="font-light">
      <Line class="!text-black">
        <p>
          <R>Error:</R> Failed <B>to</B> load resource: the server responded{' '}
          <B>with</B> a status <B>of</B> <P>404</P> (<R>Not</R> Found)
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          at HTMLScriptElement.onError (index.js:<P>12</P>)
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          <R>Error:</R> Uncaught ReferenceError: myFunction <R>is not</R>{' '}
          defined
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          at main.js:<P>5</P>
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          at onLoad (index.js:<P>15</P>)
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          <R>Warning:</R> Unexpected token &lt; <B>in</B> JSON at position{' '}
          <P>0</P>
        </p>
      </Line>
      <Line class="!text-black">
        <p>{`at JSON.parse (<anonymous>)`}</p>
      </Line>
      <Line class="!text-black">
        <p>
          at fetchData (api.js:<P>22</P>)
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          at onLoad (index.js:<P>15</P>)
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          <R>Error:</R> NetworkError <B>when</B> attempting <B>to</B> fetch
          resource.
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          at fetch (fetch.js:<P>8</P>)
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          at loadResource (loader.js:<P>14</P>)
        </p>
      </Line>
      <Line class="!text-black">
        <p>
          at onLoad (index.js:<P>15</P>)
        </p>
      </Line>
      <Line state={state()} href="/report" onClick={report} class="!text-black">
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

export default EasterEggDarkModeComponent
