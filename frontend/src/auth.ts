import NextAuth from "next-auth";
import Google from "next-auth/providers/google";

import { provisionGoogleUserSession } from "@/lib/auth/google-provisioning";

export const { handlers, auth, signIn, signOut } = NextAuth({
  providers: [
    Google({
      clientId: process.env.AUTH_GOOGLE_ID ?? process.env.GOOGLE_CLIENT_ID,
      clientSecret:
        process.env.AUTH_GOOGLE_SECRET ?? process.env.GOOGLE_CLIENT_SECRET,
    }),
  ],
  pages: {
    signIn: "/login",
  },
  callbacks: {
    async jwt({ token, account, profile }) {
      if (account?.provider !== "google") {
        return token;
      }

      const googleSub =
        typeof profile?.sub === "string" ? profile.sub : token.sub;
      const email =
        typeof profile?.email === "string" ? profile.email : token.email;
      const name =
        typeof profile?.name === "string"
          ? profile.name
          : (token.name ?? email ?? "Google User");
      const avatarUrl =
        typeof profile?.picture === "string"
          ? profile.picture
          : typeof token.picture === "string"
            ? token.picture
            : undefined;

      if (!googleSub || !email) {
        throw new Error("Google profile is missing required identity fields");
      }

      const yubiSession = await provisionGoogleUserSession({
        googleSub,
        email,
        name,
        avatarUrl,
      });

      token.yubiUserId = yubiSession.user_id;
      token.activeOrganizationId = yubiSession.active_organization_id;
      token.activeOrganizationName = yubiSession.active_organization_name;
      token.activeRole = yubiSession.active_role;

      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        session.user.id =
          typeof token.yubiUserId === "string" ? token.yubiUserId : "";
        session.user.activeOrganizationId =
          typeof token.activeOrganizationId === "string"
            ? token.activeOrganizationId
            : undefined;
        session.user.activeOrganizationName =
          typeof token.activeOrganizationName === "string"
            ? token.activeOrganizationName
            : undefined;
        session.user.activeRole =
          typeof token.activeRole === "number" ? token.activeRole : undefined;
      }
      return session;
    },
  },
});
