import { useCallback, useState } from "react"

export function useLocalState<T extends string = string>(defaultState: T, name: string): [T, (state: T) => void] {
  const initialState = localStorage.getItem(name) as T || defaultState
  const [state, setState] = useState(initialState)
  const setStateWithLocal = useCallback(
    (newState: T) => {
      console.log(`${name} changed to ${newState}`)
      setState(newState)
      localStorage.setItem(name, newState)
    },
    [setState]
  )
  return [state, setStateWithLocal]
}
