import { createContext, useContext } from "solid-js";
import { createStore } from "solid-js/store";
import { request } from "../scripts/request.ts";

const WellKnownContext = createContext<any>();

const [wellKnown, setWellKnown] = createStore<any>(null);

export async function readWellKnown() {
  const res = await request("/.well-known")
  setWellKnown(await res.json())
}

export function WellKnownProvider(props: any) {
  return (
    <WellKnownContext.Provider value={wellKnown}>
      {props.children}
    </WellKnownContext.Provider>
  );
}

export function useWellKnown() {
  return useContext(WellKnownContext);
}