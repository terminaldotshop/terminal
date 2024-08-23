import {
  createSignal,
  type Component,
  type JSX,
  type ParentProps,
  Show,
  type ParentComponent,
} from 'solid-js'
import Line, { type State } from '@components/line'
import Caret from './caret'

type EasterEggArtisanProps = {} & JSX.HTMLAttributes<HTMLDivElement>
type SpanProps = ParentProps & JSX.HTMLAttributes<HTMLSpanElement>

const B: Component<SpanProps> = ({ class: className, ...props }) => {
  return (
    <span class={`text-blue-11 ${className}`} {...props}>
      {props.children}
    </span>
  )
}

const W: Component<SpanProps> = (props) => {
  return (
    <span class="text-white" {...props}>
      {props.children}
    </span>
  )
}

const G: Component<SpanProps> = (props) => {
  return (
    <span class="text-lime" {...props}>
      {props.children}
    </span>
  )
}

const Gray: Component<SpanProps> = (props) => {
  return (
    <span class="text-gray-11" {...props}>
      {props.children}
    </span>
  )
}

const R: Component<SpanProps> = (props) => {
  return (
    <span class="text-orange" {...props}>
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

const Expand = () => {
  return (
    <svg
      width="13"
      height="14"
      viewBox="0 0 13 14"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      class="ml-3 inline-block"
    >
      <path d="M6.5 14L13 0H0L6.5 14Z" fill="#959595" />
    </svg>
  )
}

const Attributes: ParentComponent<{ title: string }> = (props) => {
  return (
    <p class="text-white px-16">
      <R>#</R>
      <W>{props.title}</W>
      <R>: </R>
      {props.children}
    </p>
  )
}

const StringAttributes: ParentComponent<{ title: string }> = (props) => {
  return (
    <p class="text-white px-16">
      <R>#</R>
      <W>{props.title}</W>
      <R>:</R>
      <R> "</R>
      <G>{props.children}</G>
      <R>"</R>
    </p>
  )
}

const ConstantAttributes: ParentComponent<{ title: string; blue?: boolean }> = (
  props,
) => {
  return (
    <p class="text-white px-16">
      <R>#</R>
      <W>{props.title}</W>
      <R>: </R>
      <Show when={props.blue} fallback={<R>{props.children}</R>}>
        <B>{props.children}</B>
      </Show>
    </p>
  )
}

const EasterEggArtisanComponent: Component<EasterEggArtisanProps> = () => {
  const [state, setState] = createSignal<State>('normal')
  const report = () => {
    setState('error')
  }

  return (
    <div class="font-light px-4">
      <p class="text-white">
        <B>array:3</B>
        <R> [</R>
      </p>
      <p class="text-white px-4">
        <R>"</R>
        <G>Coffee</G>
        <R>" </R>
        <R>=&gt;</R>
        <B> Illumina...\Collection</B>
        <R> {`{`}</R>
        <Gray>#2086</Gray>
        <Expand />
      </p>
      <p class="text-white px-8">
        <R>#</R>
        <W>items</W>
        <R>:</R>
        <B> array:1</B>
        <R> [</R>
        <Expand />
      </p>
      <p class="text-white px-12">
        <B>0</B>
        <R> =&gt; </R>
        <B>App\Coffee</B>
        <R> {`{`}</R>
        <Gray>#2088</Gray>
        <Expand />
      </p>
      <StringAttributes title="connection">mysql</StringAttributes>
      <StringAttributes title="table">coffees</StringAttributes>
      <StringAttributes title="primaryKey">id</StringAttributes>
      <StringAttributes title="incrementing">int</StringAttributes>
      <ConstantAttributes title="with">true</ConstantAttributes>
      <ConstantAttributes title="withCount">[]</ConstantAttributes>
      <ConstantAttributes title="perPage">[]</ConstantAttributes>
      <ConstantAttributes title="exists" blue>
        15
      </ConstantAttributes>
      <ConstantAttributes title="wasRecentlyCreated">true</ConstantAttributes>
      <ConstantAttributes title="attributes">false</ConstantAttributes>
      <Attributes title="original">
        <B> array:2</B>
        <R> [ ...2]</R>
      </Attributes>
      <Attributes title="changes">
        <B> array:2</B>
        <R> [ ...2]</R>
      </Attributes>
      <ConstantAttributes title="casts">[]</ConstantAttributes>
      <ConstantAttributes title="classCastCache">[]</ConstantAttributes>
      <ConstantAttributes title="dates">[]</ConstantAttributes>
      <ConstantAttributes title="dateformat">null</ConstantAttributes>
      <ConstantAttributes title="appends">[]</ConstantAttributes>
      <ConstantAttributes title="dispatchedEvents">[]</ConstantAttributes>
      <ConstantAttributes title="observables">[]</ConstantAttributes>
      <ConstantAttributes title="relations">[]</ConstantAttributes>
      <ConstantAttributes title="touches">[]</ConstantAttributes>
      <ConstantAttributes title="timestamps">true</ConstantAttributes>
      <ConstantAttributes title="hidden">[]</ConstantAttributes>
      <ConstantAttributes title="visible">[]</ConstantAttributes>
      <ConstantAttributes title="fillable">[]</ConstantAttributes>
      <Attributes title="guarded">
        <B> array:1</B>
        <R> [ ...1]</R>
      </Attributes>
      <p class="text-white px-12">
        <R>{`}`}</R>
      </p>
      <p class="text-white px-8">
        <R>]</R>
      </p>
      <p class="text-white px-4">
        <R>{`}`}</R>
      </p>
      <p class="text-white">
        <R>]</R>
      </p>
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
    </div>
  )
}

export default EasterEggArtisanComponent
