import { DefaultSession } from "next-auth";

declare module "next-auth" {
  interface Session {
    user: {
      id: string;
      activeOrganizationId?: string;
      activeOrganizationName?: string;
      activeRole?: number;
    } & DefaultSession["user"];
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    yubiUserId: string;
    activeOrganizationId?: string;
    activeOrganizationName?: string;
    activeRole?: number;
  }
}
