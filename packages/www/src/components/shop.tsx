import {
  Match,
  Show,
  Switch,
  createSignal,
  onCleanup,
  onMount,
  type Component,
  type JSX,
} from 'solid-js'
import TeaserComponent, { release } from '@components/teaser'
import PendingComponent from '@components/pending'

type ShopProps = {} & JSX.HTMLAttributes<HTMLDivElement>

const ShopComponent: Component<ShopProps> = () => {
  let sshTimeout: ReturnType<typeof setTimeout> | undefined

  const [ready, setReady] = createSignal(false)
  const [done, setDone] = createSignal<boolean>(false)

  onMount(() => {
    setDone(new Date() > release)
    setReady(true)
  })

  sshTimeout = setInterval(() => {
    const now = new Date()
    setDone(now > release)

    if (now > release) clearInterval(sshTimeout)
  }, 1000)

  onCleanup(() => clearInterval(sshTimeout))

  return (
    <Show when={ready()}>
      <Switch>
        <Match when={done()}>
          <PendingComponent />
        </Match>
        <Match when={!done()}>
          <TeaserComponent />
        </Match>
      </Switch>
    </Show>
  )
}

export default ShopComponent
