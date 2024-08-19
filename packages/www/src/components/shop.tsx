import { type Component, type JSX } from 'solid-js'

import SshComponent from '@components/ssh'

type ShopProps = { apiUrl: string } & JSX.HTMLAttributes<HTMLDivElement>

const ShopComponent: Component<ShopProps> = (props) => {
  return <SshComponent apiUrl={props.apiUrl} />
}

export default ShopComponent
