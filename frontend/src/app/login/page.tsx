import { Chrome } from "lucide-react";
import { redirect } from "next/navigation";

import { auth, signIn } from "@/auth";
import { Button } from "@/components/ui/button";

export default async function LoginPage() {
  const session = await auth();
  if (session?.user?.id) {
    redirect("/dashboard");
  }

  return (
    <main className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <div className="mx-auto flex min-h-screen w-full max-w-sm flex-col justify-center px-6">
        <div className="mb-8">
          <h1 className="text-2xl font-semibold text-gray-950 dark:text-gray-50">
            Yubi App
          </h1>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            Sign in to continue.
          </p>
        </div>
        <form
          action={async () => {
            "use server";
            await signIn("google", { redirectTo: "/dashboard" });
          }}
        >
          <Button type="submit" className="w-full">
            <Chrome className="h-4 w-4" />
            Continue with Google
          </Button>
        </form>
      </div>
    </main>
  );
}
