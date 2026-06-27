"use client";

import { ChevronDown, Download } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { ExportOperatorYieldDialog } from "@/features/reporting";

import {
  ExportEpisodesDialog,
  type ExportEpisodesInitialFilters,
} from "./export-episodes-dialog";

type ExportMenuProps = {
  initialFilters?: ExportEpisodesInitialFilters;
};

export function ExportMenu({ initialFilters }: ExportMenuProps = {}) {
  const { t } = useTranslation();
  const [openEpisodes, setOpenEpisodes] = useState(false);
  const [openYield, setOpenYield] = useState(false);
  const [episodesOpenKey, setEpisodesOpenKey] = useState(0);

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" />
            {t("exportMenu.label")}
            <ChevronDown className="ml-2 h-4 w-4 opacity-50" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem
            onSelect={() => {
              setEpisodesOpenKey((k) => k + 1);
              setOpenEpisodes(true);
            }}
          >
            {t("exportMenu.episodes")}
          </DropdownMenuItem>
          <DropdownMenuItem onSelect={() => setOpenYield(true)}>
            {t("exportMenu.operatorYield")}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* key bumps on each open so the dialog remounts and re-reads initialFilters from the current list state. */}
      <ExportEpisodesDialog
        key={episodesOpenKey}
        open={openEpisodes}
        onOpenChange={setOpenEpisodes}
        showTrigger={false}
        initialFilters={initialFilters}
      />
      <ExportOperatorYieldDialog
        open={openYield}
        onOpenChange={setOpenYield}
        showTrigger={false}
      />
    </>
  );
}
