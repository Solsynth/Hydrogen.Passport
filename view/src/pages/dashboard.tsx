import { useUserinfo } from "../stores/userinfo.tsx";

export default function DashboardPage() {
  const userinfo = useUserinfo();

  return (
    <div class="container mx-auto pt-12">
      <h1 class="text-2xl font-bold">Welcome, {userinfo?.displayName}</h1>
      <p>What's a nice day!</p>
    </div>
  );
}