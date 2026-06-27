import { getUserSession } from "@/lib/auth/session";
import { SessionProvider } from "@/lib/auth/session-context";

import { TopNav } from "@/components/layout/top-nav";
import { Toaster } from "@/components/ui/sonner";

export default async function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const userSession = await getUserSession();

  return (
    <SessionProvider userSession={userSession}>
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <TopNav />
        <main className="container mx-auto px-4 py-6">{children}</main>
        <Toaster />
      </div>
    </SessionProvider>
  );
}
