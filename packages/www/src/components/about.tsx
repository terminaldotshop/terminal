import { type Component, type JSX } from 'solid-js'
import Line from '@components/line'
import Editor from '@components/editor'
import Caret from '@components/caret'

type AboutProps = {} & JSX.HTMLAttributes<HTMLDivElement>

const AboutComponent: Component<AboutProps> = () => {
  const heroes = [
    'thdxr',
    'adamdotdev',
    'theprimeagen',
    'teej_dv',
    'iamdavidhill',
  ]

  return (
    <Editor>
      <Line number={1}>
        <h1>
          # Amazingly awesome products for developers brought to you by a group
          of talented, good looking, and humble heroes...
        </h1>
      </Line>
      {heroes.map((hero, index) => (
        <Line number={index + 2} href={`https://x.com/${hero}`}>
          <h2>{`# @${hero}`}</h2>
        </Line>
      ))}
      <Line number={heroes.length + 2}>
        <div class="flex items-center gap-1.5">
          <h3 class="text-white whitespace-nowrap">Terminal Products, Inc.</h3>
          <Caret />
        </div>
      </Line>
    </Editor>
  )
}

export default AboutComponent
