import { useEffect, useState } from "react"

/**
 * Support any types of state.
 * For string state, please use `useLocalState`.
 */
export function useLocalJsonState<T>(
  defaultState: T,
  name: string
): [T, (state: T) => void] {
  let initialState = defaultState
  const localItem = localStorage.getItem(name)
  if (localItem !== null) {
    initialState = JSON.parse(localItem)
  }
  const [state, setState] = useState(initialState)
  useEffect(() => {
    localStorage.setItem(name, JSON.stringify(state))
  }, [state])
  return [state, setState]
}
