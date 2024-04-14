import { type Component, type JSX } from 'solid-js'

type EditorProps = {
  children: JSX.Element
} & JSX.HTMLAttributes<HTMLDivElement>

const EditorComponent: Component<EditorProps> = (props) => {
  return (
    <div
      {...props}
      classList={{
        'leading-10': true,
        [props.class ?? '']: !!props.class,
      }}
    >
      {props.children}
    </div>
  )
}

export default EditorComponent
