import { useSearchParams } from "@solidjs/router";

export default function DefaultCallbackPage() {
  const [searchParams] = useSearchParams();

  return (
    <div class="h-full flex justify-center items-center mx-5">
      <div class="card w-screen max-w-[480px] shadow-xl">
        <div class="card-body">
          <div id="header" class="text-center mb-5">
            {/* Just Kidding */}
            <h1 class="text-xl font-bold">Default Callback</h1>
            <p>
              If you see this page, it means some genius developer forgot to set the redirect address, so you visited
              this default callback address.
              General Douglas MacArthur, a five-star general in the United States, commented on this: "If I let my
              soldiers use default callbacks, they'd rather die."
              The large documentary film "Callback Legend" is currently in theaters.
            </p>
          </div>

          <div class="text-center">
            <p>Authorization Code</p>
            <code>{searchParams["code"]}</code>
          </div>
        </div>
      </div>
    </div>
  );
}