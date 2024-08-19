import {
  createContext,
  type Component,
  type JSX,
  type ParentProps,
  useContext,
  createSignal,
} from 'solid-js'

type EditorProps = ParentProps & JSX.HTMLAttributes<HTMLDivElement>

const LineNumberContext = createContext<() => number>()
export function getLineNumber() {
  return useContext(LineNumberContext)?.()
}

const EditorComponent: Component<EditorProps> = (props) => {
  const [getLineNumber, setLatestLineNumber] = createSignal(1)

  function getNextLineNumber() {
    const latest = getLineNumber()
    setLatestLineNumber(latest + 1)
    return latest
  }

  return (
    <LineNumberContext.Provider value={getNextLineNumber}>
      <div
        {...props}
        classList={{
          'leading-10': true,
          [props.class ?? '']: !!props.class,
        }}
      >
        {props.children}
      </div>
    </LineNumberContext.Provider>
  )
}

export default EditorComponent
