import { userinfo } from "../stores/userinfo.ts";

export default function Dashboard() {
  return (
    <div class="container mx-auto pt-12">
      <h1 class="text-2xl font-bold">Welcome, {userinfo.displayName}</h1>
      <p>What's a nice day!</p>
    </div>
  )
}