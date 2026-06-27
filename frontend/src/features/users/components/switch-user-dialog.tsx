"use client";

import { ArrowLeftRight } from "lucide-react";
import { useState, useTransition } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { useSession } from "@/lib/auth/session-context";
import { switchUser } from "@/lib/auth/switch-user";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";

import { useUsersQuery } from "../hooks/use-users-query";

export function SwitchUserDialog() {
  const { t } = useTranslation();
  const { userSession } = useSession();
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [isPending, startTransition] = useTransition();

  const { data, isLoading } = useUsersQuery(
    { search: search || undefined, limit: 50 },
    { enabled: open }
  );

  const users = data?.users ?? [];

  function handleSelect(userId: string) {
    startTransition(async () => {
      try {
        await switchUser(userId);
        setOpen(false);
        window.location.reload();
      } catch {
        toast.error(t("userMenu.switchError"));
      }
    });
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <DropdownMenuItem
          onSelect={(e) => {
            e.preventDefault();
            setOpen(true);
          }}
        >
          <ArrowLeftRight className="mr-2 h-4 w-4" />
          {t("userMenu.switchAccount")}
        </DropdownMenuItem>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t("userMenu.switchAccount")}</DialogTitle>
        </DialogHeader>

        <Input
          placeholder={t("userMenu.searchUsers")}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="mb-4"
        />

        <div className="max-h-64 overflow-y-auto space-y-1">
          {isLoading ? (
            <div className="space-y-2">
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"
                />
              ))}
            </div>
          ) : users.length === 0 ? (
            <p className="text-sm text-gray-500 text-center py-4">
              {t("userMenu.noUsersFound")}
            </p>
          ) : (
            users.map((user) => {
              const isActive = user.user_id === userSession?.userId;
              const initials = user.display_name
                .split(" ")
                .map((n: string) => n[0])
                .join("")
                .toUpperCase()
                .slice(0, 2);

              return (
                <button
                  key={user.user_id}
                  onClick={() => handleSelect(user.user_id)}
                  disabled={isPending || isActive}
                  className={`w-full flex items-center gap-3 px-3 py-2 rounded-md text-left transition-colors ${
                    isActive
                      ? "bg-blue-50 dark:bg-blue-900/20 cursor-default"
                      : "hover:bg-gray-100 dark:hover:bg-gray-800 cursor-pointer"
                  } ${isPending ? "opacity-50" : ""}`}
                >
                  <Avatar className="h-8 w-8">
                    <AvatarFallback
                      className={`text-xs ${isActive ? "bg-blue-600 text-white" : "bg-gray-400 text-white"}`}
                    >
                      {initials}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">
                      {user.display_name}
                    </p>
                    <p className="text-xs text-gray-500 truncate">
                      {user.email}
                    </p>
                  </div>
                  {isActive && (
                    <span className="text-xs text-blue-600 dark:text-blue-400 font-medium">
                      {t("userMenu.current")}
                    </span>
                  )}
                </button>
              );
            })
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
