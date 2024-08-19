import { type Component, type JSX } from 'solid-js'
type EasterEgg404Props = {} & JSX.HTMLAttributes<HTMLDivElement>

const EasterEgg404Component: Component<EasterEgg404Props> = () => {
  return (
    <div class="mx-auto text-center flex flex-col gap-2">
      <h1 class="text-5xl text-white">404</h1>
      <p class="text-gray-11">Page not found</p>
      <div class="text-red-11 text-sm mx-auto">
        <a href="/report" target="_blank" class="flex gap-2">
          Report bug
          <svg
            class="size-5 text-inherit shrink-0"
            viewBox="0 0 20 20"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <g id="square-arrow-top-right, open, new, link, open link, box, arrow">
              <path
                class="stroke-current"
                d="M15.2083 11.6667V16.875H3.125V4.79167H7.70833M11.4583 3.125H16.875V8.54167M9.16667 10.8333L16.25 3.75"
                stroke-width="1.5"
                stroke-linecap="square"
              ></path>
            </g>
          </svg>
        </a>
      </div>
    </div>
  )
}

export default EasterEgg404Component
