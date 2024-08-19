import { type Component, type JSX } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { getLineNumber } from './editor'

export type State = 'normal' | 'success' | 'error' | 'busy'

type LineProps = {
  href?: string
  state?: State
} & JSX.HTMLAttributes<HTMLDivElement>

const LineComponent: Component<LineProps> = (props) => {
  const number = getLineNumber()

  return (
    <Dynamic
      component={props.href ? 'a' : 'div'}
      tabindex="0"
      href={props.href}
      target={props.href ? '_blank' : undefined}
      {...props}
      classList={{
        ...props.classList,
        'group flex items-center text-gray-10 hover:bg-gray-5 px-6': true,
        'active:border-orange active:border-l-2 active:pl-[22px] active:text-gray-11 active:bg-gray-6':
          true,
        '[&>svg]:hover:block [&>svg]:active:text-gray-11': !!props.href,
        '!border-green-11 !bg-green-5 border-l-2 pl-[22px]':
          props.state === 'success',
        '!border-red-11 !bg-red-5 border-l-2 pl-[22px]':
          props.state === 'error',
        '!border-blue-11 !bg-blue-5 border-l-2 pl-[22px]':
          props.state === 'busy',
        'focus:border-orange focus:border-l-2 focus:pl-[22px] focus:text-gray-11 focus:bg-gray-6 focus:outline-none':
          true,
        'has-[:focus]:border-orange has-[:focus]:border-l-2 has-[:focus]:pl-[22px] has-[:focus]:text-gray-11 has-[:focus]:bg-gray-6 has-[:focus]:outline-none':
          true,
        'pointer-events-none': props.state && props.state !== 'normal',
        [props.class ?? '']: !!props.class,
      }}
    >
      <div class="text-center pr-10 self-start">{number}</div>
      {props.children}
      <svg
        class="hidden w-5 h-5 ml-6 text-gray-10 shrink-0"
        viewBox="0 0 20 20"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <g id="square-arrow-top-right, open, new, link, open link, box, arrow">
          <path
            class="stroke-current"
            d="M15.2083 11.6667V16.875H3.125V4.79167H7.70833M11.4583 3.125H16.875V8.54167M9.16667 10.8333L16.25 3.75"
            stroke-width="1.5"
            stroke-linecap="square"
          ></path>
        </g>
      </svg>
    </Dynamic>
  )
}

export default LineComponent
