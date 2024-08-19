import { splitProps, type Component, type JSX } from 'solid-js'

type CaretProps = { blink?: boolean } & JSX.HTMLAttributes<HTMLDivElement>

const CaretComponent: Component<CaretProps> = (props) => {
  const [local, others] = splitProps(props, ['class', 'classList', 'blink'])
  return (
    <div
      {...others}
      classList={{
        ...local.classList,
        'bg-orange w-[10px] h-[21px] shrink-0': true,
        'animate-blink': local.blink ?? true,
        [local.class ?? '']: true,
      }}
    ></div>
  )
}

export default CaretComponent
