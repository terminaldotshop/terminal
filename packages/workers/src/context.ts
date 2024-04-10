import { AsyncLocalStorage } from "node:async_hooks";

export function createContext<T>() {
  const storage = new AsyncLocalStorage<T>();
  return {
    use() {
      return storage.getStore();
    },
    with(value: T, fn: () => void) {
      return storage.run(value, fn);
    },
  };
}
