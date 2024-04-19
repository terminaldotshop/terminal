import {
  type Component,
  type JSX,
} from 'solid-js'

import PendingComponent from '@components/pending'

type ShopProps = {} & JSX.HTMLAttributes<HTMLDivElement>

const ShopComponent: Component<ShopProps> = () => {
  return <PendingComponent />
}

export default ShopComponent
