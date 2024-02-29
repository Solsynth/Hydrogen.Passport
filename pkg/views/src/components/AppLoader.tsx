import { ReactNode, useEffect } from "react";
import { useWellKnown } from "@/stores/wellKnown.tsx";
import { useUserinfo } from "@/stores/userinfo.tsx";

export default function AppLoader({ children }: { children: ReactNode }) {
  const { readWellKnown } = useWellKnown();
  const { readProfiles } = useUserinfo();

  useEffect(() => {
    Promise.all([readWellKnown(), readProfiles()]);
  }, []);

  return children;
}