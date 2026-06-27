"use client";

import { Pencil } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { useUserRoleLabel } from "@/shared/hooks/use-status-labels";
import { USER_ROLE, type UserRoleValue } from "@/shared/lib/status-constants";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useLocationSearchOptions } from "@/features/locations";
import { useSiteSearchOptions } from "@/features/sites";

import { useUpdateUserMutation } from "../hooks/use-update-user-mutation";
import { useUpdateUserRoleMutation } from "../hooks/use-update-user-role-mutation";

type UserResponse = z.infer<typeof schemas.UserResponse>;

const roleOptions: UserRoleValue[] = [
  USER_ROLE.ADMIN,
  USER_ROLE.DATA_ENGINEER,
  USER_ROLE.MANAGER,
  USER_ROLE.OPERATOR,
  USER_ROLE.VIEWER,
];

interface EditUserDialogProps {
  user: UserResponse;
  currentUser?: UserResponse;
  children?: React.ReactNode;
}

export function EditUserDialog({
  user,
  currentUser,
  children,
}: EditUserDialogProps) {
  const { t } = useTranslation();
  const getRoleLabel = useUserRoleLabel();
  const [open, setOpen] = useState(false);
  const [selectedRole, setSelectedRole] = useState<string>(
    String(user.role ?? USER_ROLE.VIEWER)
  );
  const [selectedLocationIds, setSelectedLocationIds] = useState<string[]>(
    user.locations?.map((l) => l.location_id) ?? []
  );

  const { mutateAsync: mutateRole, isPending: isRolePending } =
    useUpdateUserRoleMutation();
  const { mutateAsync: mutateUser, isPending: isUserPending } =
    useUpdateUserMutation();
  const [selectedSiteId, setSelectedSiteId] = useState<string>("");
  const {
    options: siteOptions,
    isLoading: sitesLoading,
    onSearch: onSiteSearch,
    selectedLabel: siteSelectedLabel,
    onValueChange: onSiteSelectChange,
  } = useSiteSearchOptions();
  const { options: locationOptions, isLoading: locationsLoading } =
    useLocationSearchOptions({
      site_id: selectedSiteId || undefined,
    });

  const isPending = isRolePending || isUserPending;

  const isSelfAdminEditing =
    currentUser?.role === USER_ROLE.ADMIN &&
    currentUser.user_id === user.user_id;

  const initialRole = String(user.role ?? USER_ROLE.VIEWER);
  const initialLocationIds = user.locations?.map((l) => l.location_id) ?? [];

  const roleChanged = selectedRole !== initialRole;
  const locationIdsChanged =
    selectedLocationIds.length !== initialLocationIds.length ||
    selectedLocationIds.some((id) => !initialLocationIds.includes(id));
  const hasChanges = roleChanged || locationIdsChanged;

  const handleSave = async () => {
    if (!hasChanges) return;

    const promises: Promise<unknown>[] = [];

    if (roleChanged) {
      promises.push(
        mutateRole({
          userId: user.user_id,
          data: { role: Number(selectedRole) as UserRoleValue },
        })
      );
    }

    if (locationIdsChanged) {
      promises.push(
        mutateUser({
          userId: user.user_id,
          data: { location_ids: selectedLocationIds },
        })
      );
    }

    await Promise.all(promises);
    setOpen(false);
  };

  const toggleLocation = (locationId: string) => {
    setSelectedLocationIds((prev) =>
      prev.includes(locationId)
        ? prev.filter((id) => id !== locationId)
        : [...prev, locationId]
    );
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        if (nextOpen) {
          setSelectedRole(String(user.role ?? USER_ROLE.VIEWER));
          setSelectedLocationIds(
            user.locations?.map((l) => l.location_id) ?? []
          );
        }
        setOpen(nextOpen);
      }}
    >
      <DialogTrigger asChild>
        {children || (
          <Button size="sm" variant="ghost">
            <Pencil className="h-4 w-4" />
            <span className="sr-only">Edit user</span>
          </Button>
        )}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("editUserDialog.title")}</DialogTitle>
          <DialogDescription>
            {t("editUserDialog.description", { name: user.display_name })}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Role */}
          <div className="space-y-2">
            <p className="text-sm font-medium text-gray-700 dark:text-gray-300">
              {t("editUserDialog.role")}
            </p>
            <SearchableSelect
              value={selectedRole}
              onValueChange={setSelectedRole}
              options={roleOptions.map((role) => ({
                value: String(role),
                label: getRoleLabel(role),
              }))}
              placeholder={t("editUserDialog.selectRole")}
              disabled={isSelfAdminEditing}
            />
            {isSelfAdminEditing && (
              <p className="text-xs text-gray-500 dark:text-gray-400">
                {t("editUserDialog.adminCannotChange")}
              </p>
            )}
          </div>

          {/* Site Filter */}
          <div className="space-y-2">
            <p className="text-sm font-medium text-gray-700 dark:text-gray-300">
              {t("editUserDialog.filterBySite")}
            </p>
            <SearchableSelect
              value={selectedSiteId}
              onValueChange={(value) => {
                setSelectedSiteId(value);
                onSiteSelectChange(value);
              }}
              options={[
                { value: "", label: t("editUserDialog.allSites") },
                ...siteOptions,
              ]}
              onSearch={onSiteSearch}
              isLoading={sitesLoading}
              selectedLabel={selectedSiteId ? siteSelectedLabel : undefined}
              placeholder={t("editUserDialog.allSites")}
            />
          </div>

          {/* Locations */}
          <div className="space-y-2">
            <p className="text-sm font-medium text-gray-700 dark:text-gray-300">
              {t("editUserDialog.locations")}
            </p>
            {locationOptions.length > 0 ? (
              <div className="max-h-48 overflow-y-auto space-y-2 rounded-md border border-gray-200 dark:border-gray-700 p-3">
                {locationOptions.map((loc) => (
                  <label
                    key={loc.value}
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <Checkbox
                      checked={selectedLocationIds.includes(loc.value)}
                      onCheckedChange={() => toggleLocation(loc.value)}
                    />
                    <span className="text-sm text-gray-700 dark:text-gray-300">
                      {loc.label}
                    </span>
                  </label>
                ))}
              </div>
            ) : (
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {locationsLoading
                  ? t("editUserDialog.loadingLocations")
                  : t("editUserDialog.noLocations")}
              </p>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => setOpen(false)}
            disabled={isPending}
          >
            {t("dialog.cancel")}
          </Button>
          <Button
            type="button"
            onClick={handleSave}
            disabled={isPending || !hasChanges}
          >
            {isPending ? t("dialog.saving") : t("dialog.save")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
