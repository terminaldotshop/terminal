import { type Component, type JSX } from 'solid-js'
import Line from '@components/line'
import Editor from '@components/editor'

type FaqProps = {
  faqs: { question: string; answer: string }[]
} & JSX.HTMLAttributes<HTMLDivElement>

const FaqComponent: Component<FaqProps> = (props) => {
  return (
    <Editor class="max-w-xl">
      {props.faqs.map((faq, i) => (
        <Line number={i + 1}>
          <p>
            <span class="text-white mr-2">{faq.question}</span># {faq.answer}
          </p>
        </Line>
      ))}
    </Editor>
  )
}

export default FaqComponent
