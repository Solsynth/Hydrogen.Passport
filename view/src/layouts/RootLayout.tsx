import Navbar from "./shared/Navbar.tsx";

export default function RootLayout(props: any) {
  return (
    <div>
      <Navbar />

      <main class="h-[calc(100vh-68px)]">{props.children}</main>
    </div>
  );
}