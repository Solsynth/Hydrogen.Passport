import { createContext, useContext, useState } from "react";
import { request } from "../scripts/request.ts";

const WellKnownContext = createContext<any>(null);

export function WellKnownProvider(props: any) {
  const [wellKnown, setWellKnown] = useState<any>(null);

  async function readWellKnown() {
    const res = await request("/.well-known");
    setWellKnown(await res.json());
  }

  return (
    <WellKnownContext.Provider value={{ wellKnown, readWellKnown }}>
      {props.children}
    </WellKnownContext.Provider>
  );
}

export function useWellKnown() {
  return useContext(WellKnownContext);
}