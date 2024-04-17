import {
  Match,
  Switch,
  createSignal,
  onCleanup,
  type Component,
  type JSX,
} from 'solid-js'
import SshComponent from '@components/ssh'
import TeaserComponent, { release } from '@components/teaser'

type ShopProps = {} & JSX.HTMLAttributes<HTMLDivElement>

const ShopComponent: Component<ShopProps> = () => {
  let sshTimeout: ReturnType<typeof setTimeout> | undefined

  const [done, setDone] = createSignal<boolean>(new Date() > release)

  sshTimeout = setInterval(() => {
    const now = new Date()
    setDone(now > release)

    if (now > release) clearInterval(sshTimeout)
  }, 1000)

  onCleanup(() => clearInterval(sshTimeout))

  return (
    <Switch fallback={<TeaserComponent />}>
      <Match when={done()}>
        <SshComponent />
      </Match>
      <Match when={!done()}>
        <TeaserComponent />
      </Match>
    </Switch>
  )
}

export default ShopComponent
