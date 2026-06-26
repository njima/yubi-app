"use client";

import {
  parseAsBoolean,
  parseAsInteger,
  parseAsString,
  useQueryState,
} from "nuqs";
import { Suspense, useMemo } from "react";
import { useTranslation } from "react-i18next";

import { useFormatAbsoluteTime } from "@/shared/hooks/use-date-formatters";
import { DEFAULT_PAGE_SIZE } from "@/shared/lib/pagination";
import { Checkbox } from "@/shared/ui/checkbox";
import { Label } from "@/shared/ui/label";
import { SearchableSelect } from "@/shared/ui/searchable-select";

import {
  ApiKeyDataTable,
  CreateApiKeyDialog,
  getApiKeyColumns,
  useApiKeysQuery,
} from "@/features/api-keys";
import { useRobotQuery } from "@/features/robots";
import { useRobotSearchOptions } from "@/features/robots/hooks/use-robot-search-options";
import { usePermission } from "@/features/users";

function ApiKeysContent() {
  const { t } = useTranslation();

  const canList = usePermission("api_key:list");
  const canRevoke = usePermission("api_key:revoke");
  const canCreate = usePermission("api_key:create");

  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(1));
  const [limit, setLimit] = useQueryState(
    "limit",
    parseAsInteger.withDefault(DEFAULT_PAGE_SIZE)
  );
  const [robotId, setRobotId] = useQueryState("robot_id", parseAsString);
  const [includeRevoked, setIncludeRevoked] = useQueryState(
    "include_revoked",
    parseAsBoolean.withDefault(false)
  );

  const robotSearch = useRobotSearchOptions({ enabled: canList });
  const formatDateTime = useFormatAbsoluteTime();

  const { data: prefilledRobot, isLoading: isPrefilledRobotLoading } =
    useRobotQuery(robotId ?? "", { enabled: !!robotId && canList });

  const { data, isLoading, error } = useApiKeysQuery(
    {
      page,
      limit,
      robotId: robotId || undefined,
      includeRevoked: includeRevoked || undefined,
    },
    { enabled: canList }
  );

  const keys = data?.api_keys ?? [];
  const pagination = data?.pagination;
  const totalPages = pagination
    ? Math.max(1, Math.ceil(pagination.count / pagination.limit))
    : 1;

  const columns = useMemo(
    () => getApiKeyColumns({ canRevoke, t, formatDateTime }),
    [canRevoke, t, formatDateTime]
  );

  if (!canList) {
    return (
      <div className="p-6">
        <p className="text-gray-600 dark:text-gray-400">403 Forbidden</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-start justify-between gap-4 flex-wrap">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            {t("apiKeysPage.title")}
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            {t("apiKeysPage.subtitle")}
          </p>
        </div>
        {canCreate && (
          <CreateApiKeyDialog
            prefilledRobotId={robotId ?? undefined}
            prefilledRobotName={prefilledRobot?.name}
            triggerDisabled={!!robotId && isPrefilledRobotLoading}
          />
        )}
      </div>

      <div className="flex flex-wrap items-center gap-4">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">
            {t("apiKeysPage.filterRobot")}:
          </span>
          <SearchableSelect
            value={robotId ?? ""}
            onValueChange={(v) => {
              setRobotId(v === "" ? null : v);
              setPage(1);
              robotSearch.onValueChange(v);
            }}
            options={[
              { value: "", label: t("apiKeysPage.filterRobotAll") },
              ...robotSearch.options,
            ]}
            onSearch={robotSearch.onSearch}
            isLoading={robotSearch.isLoading}
            selectedLabel={robotId ? robotSearch.selectedLabel : undefined}
            placeholder={t("apiKeysPage.filterRobotAll")}
            disabled={isLoading}
            className="min-w-[200px]"
          />
        </div>

        <div className="flex items-center gap-2">
          <Checkbox
            id="include-revoked"
            checked={includeRevoked}
            onCheckedChange={(value) => {
              setIncludeRevoked(value === true);
              setPage(1);
            }}
          />
          <Label htmlFor="include-revoked" className="text-sm cursor-pointer">
            {t("apiKeysPage.includeRevoked")}
          </Label>
        </div>
      </div>

      {error ? (
        <div className="rounded-md border border-red-200 bg-red-50 p-4 text-sm text-red-800 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
          {t("apiKeysPage.errorLoading", { message: error.message })}
        </div>
      ) : (
        <ApiKeyDataTable
          columns={columns}
          data={keys}
          isLoading={isLoading}
          totalCount={pagination?.count}
          page={page}
          totalPages={totalPages}
          onPageChange={(p) => setPage(p)}
          limit={limit}
          onLimitChange={(l) => {
            setLimit(l);
            setPage(1);
          }}
        />
      )}
    </div>
  );
}

export default function ApiKeysPage() {
  const { t } = useTranslation();

  return (
    <Suspense
      fallback={
        <div className="p-8 text-center text-gray-600 dark:text-gray-400">
          {t("common.loading")}
        </div>
      }
    >
      <ApiKeysContent />
    </Suspense>
  );
}
