import {
  type Component,
  type JSX,
  type ComponentProps,
  splitProps,
  Show,
} from 'solid-js'
import Line from '@components/line'
import Input from '@components/input'

type InputLineProps = {} & ComponentProps<typeof Line> &
  ComponentProps<typeof Input> &
  JSX.HTMLAttributes<HTMLSpanElement>

const InputLineComponent: Component<InputLineProps> = (props) => {
  let inputRef: HTMLElement | undefined
  const [input, line] = splitProps(props, ['autofocus'])

  const handleFocus = () => {
    inputRef?.focus()
    if (props.onClick) (props as any).onClick()
  }

  return (
    <Line {...line} tabindex="-1" onClick={handleFocus}>
      <Input ref={inputRef} {...input} {...line} />
      <Show when={props.state === 'normal'}>{props.children}</Show>
    </Line>
  )
}

export default InputLineComponent
