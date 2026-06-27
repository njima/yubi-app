"use client";

import {
  Upload,
  Download,
  CheckCircle,
  AlertCircle,
  SkipForward,
} from "lucide-react";
import { useCallback, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import type { TaskImportValidationResponse } from "@/lib/api/backend-client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import {
  useValidateTaskImportMutation,
  useImportTasksMutation,
} from "../hooks/use-import-tasks-mutation";

type Step = "upload" | "preview" | "complete";

export function ImportTasksDialog() {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [step, setStep] = useState<Step>("upload");
  const [csvContent, setCsvContent] = useState<string>("");
  const [fileName, setFileName] = useState<string>("");
  const [validation, setValidation] =
    useState<TaskImportValidationResponse | null>(null);
  const [importResult, setImportResult] = useState<{
    imported_count: number;
    skipped_count: number;
    error_count: number;
  } | null>(null);
  const [showErrors, setShowErrors] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const validateMutation = useValidateTaskImportMutation();
  const importMutation = useImportTasksMutation();

  const reset = useCallback(() => {
    setStep("upload");
    setCsvContent("");
    setFileName("");
    setValidation(null);
    setImportResult(null);
    setShowErrors(false);
    validateMutation.reset();
    importMutation.reset();
  }, [validateMutation, importMutation]);

  const handleOpenChange = (isOpen: boolean) => {
    setOpen(isOpen);
    if (!isOpen) {
      reset();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 5 * 1024 * 1024) {
      toast.error(t("importTasksDialog.fileTooLarge"), {
        description: t("importTasksDialog.maxFileSize"),
      });
      return;
    }

    setFileName(file.name);
    const reader = new FileReader();
    reader.onload = (event) => {
      const text = event.target?.result as string;
      setCsvContent(text);
    };
    reader.readAsText(file, "UTF-8");

    // Reset file input so same file can be selected again
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleValidate = () => {
    validateMutation.mutate(csvContent, {
      onSuccess: (result) => {
        setValidation(result);
        setStep("preview");
      },
      onError: (error) => {
        toast.error(t("importTasksDialog.validationFailed"), {
          description: error.message,
        });
      },
    });
  };

  const handleImport = () => {
    importMutation.mutate(csvContent, {
      onSuccess: (result) => {
        setImportResult({
          imported_count: result.imported_count,
          skipped_count: result.skipped_count,
          error_count: result.error_count,
        });
        setStep("complete");
        toast.success(
          t("importTasksDialog.importedSuccessfully", {
            count: result.imported_count,
          })
        );
      },
      onError: (error) => {
        toast.error(t("importTasksDialog.importFailed"), {
          description: error.message,
        });
      },
    });
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant="outline">
          <Upload className="mr-2 h-4 w-4" />
          {t("dialog.import")}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg max-h-[90vh] flex flex-col">
        {step === "upload" && (
          <>
            <DialogHeader>
              <DialogTitle>{t("importTasksDialog.title")}</DialogTitle>
              <DialogDescription>
                {t("importTasksDialog.description")}
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div
                className="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-8 text-center cursor-pointer hover:border-gray-400 dark:hover:border-gray-500 transition-colors"
                onClick={() => fileInputRef.current?.click()}
                onDragOver={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                }}
                onDrop={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  const file = e.dataTransfer.files[0];
                  if (file) {
                    const fakeEvent = {
                      target: { files: [file] },
                    } as unknown as React.ChangeEvent<HTMLInputElement>;
                    handleFileChange(fakeEvent);
                  }
                }}
              >
                <Upload className="mx-auto h-8 w-8 text-gray-400 mb-2" />
                {fileName ? (
                  <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
                    {fileName}
                  </p>
                ) : (
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {t("importTasksDialog.clickToSelect")}
                  </p>
                )}
                <input
                  ref={fileInputRef}
                  type="file"
                  accept=".csv"
                  className="hidden"
                  onChange={handleFileChange}
                />
              </div>

              <a
                href="/web/templates/task-import-template.csv"
                download
                className="inline-flex items-center text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
              >
                <Download className="mr-1 h-4 w-4" />
                {t("importTasksDialog.downloadTemplate")}
              </a>
            </div>
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={() => handleOpenChange(false)}>
                {t("dialog.cancel")}
              </Button>
              <Button
                onClick={handleValidate}
                disabled={!csvContent || validateMutation.isPending}
              >
                {validateMutation.isPending
                  ? t("dialog.validating")
                  : t("dialog.validate")}
              </Button>
            </div>
          </>
        )}

        {step === "preview" && validation && (
          <>
            <DialogHeader>
              <DialogTitle>{t("importTasksDialog.previewTitle")}</DialogTitle>
              <DialogDescription>
                {t("importTasksDialog.previewDescription")}
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-3 py-4 flex-1 overflow-y-auto">
              <div className="flex items-center gap-2 text-sm">
                <CheckCircle className="h-4 w-4 text-green-500" />
                <span>
                  {t("importTasksDialog.readyToImport", {
                    count: validation.summary.valid_count,
                  })}
                </span>
              </div>
              {validation.summary.duplicate_count > 0 && (
                <div className="flex items-center gap-2 text-sm">
                  <SkipForward className="h-4 w-4 text-yellow-500" />
                  <span>
                    {t("importTasksDialog.duplicates", {
                      count: validation.summary.duplicate_count,
                    })}
                  </span>
                </div>
              )}
              {validation.summary.error_count > 0 && (
                <div className="flex items-center gap-2 text-sm">
                  <AlertCircle className="h-4 w-4 text-red-500" />
                  <span>
                    {t("importTasksDialog.errorsCount", {
                      count: validation.summary.error_count,
                    })}
                  </span>
                </div>
              )}

              {(validation.error_rows.length > 0 ||
                validation.duplicate_rows.length > 0) && (
                <div>
                  <button
                    type="button"
                    className="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                    onClick={() => setShowErrors(!showErrors)}
                  >
                    {showErrors
                      ? t("importTasksDialog.hideDetails")
                      : t("importTasksDialog.showDetails")}
                  </button>
                  {showErrors && (
                    <div className="mt-2 max-h-48 overflow-y-auto space-y-2">
                      {validation.error_rows.map((row) => (
                        <div
                          key={row.row_number}
                          className="text-xs bg-red-50 dark:bg-red-950 p-2 rounded"
                        >
                          <span className="font-medium">
                            {t("importTasksDialog.row")} {row.row_number}
                            {row.name ? ` (${row.name})` : ""}:
                          </span>{" "}
                          {row.errors.join("; ")}
                        </div>
                      ))}
                      {validation.duplicate_rows.map((row) => (
                        <div
                          key={row.row_number}
                          className="text-xs bg-yellow-50 dark:bg-yellow-950 p-2 rounded"
                        >
                          <span className="font-medium">
                            {t("importTasksDialog.row")} {row.row_number}
                            {row.name ? ` (${row.name})` : ""}:
                          </span>{" "}
                          {row.errors.join("; ")}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              )}
            </div>
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={reset}>
                {t("dialog.back")}
              </Button>
              <Button
                onClick={handleImport}
                disabled={
                  validation.summary.valid_count === 0 ||
                  importMutation.isPending
                }
              >
                {importMutation.isPending
                  ? t("dialog.importing")
                  : t("importTasksDialog.importCount", {
                      count: validation.summary.valid_count,
                    })}
              </Button>
            </div>
          </>
        )}

        {step === "complete" && importResult && (
          <>
            <DialogHeader>
              <DialogTitle>{t("importTasksDialog.completeTitle")}</DialogTitle>
            </DialogHeader>
            <div className="space-y-3 py-4">
              <div className="flex items-center gap-2 text-sm">
                <CheckCircle className="h-4 w-4 text-green-500" />
                <span>
                  {t("importTasksDialog.imported", {
                    count: importResult.imported_count,
                  })}
                </span>
              </div>
              {importResult.skipped_count > 0 && (
                <div className="flex items-center gap-2 text-sm">
                  <SkipForward className="h-4 w-4 text-yellow-500" />
                  <span>
                    {t("importTasksDialog.skipped", {
                      count: importResult.skipped_count,
                    })}
                  </span>
                </div>
              )}
              {importResult.error_count > 0 && (
                <div className="flex items-center gap-2 text-sm">
                  <AlertCircle className="h-4 w-4 text-red-500" />
                  <span>
                    {t("importTasksDialog.errorsCount", {
                      count: importResult.error_count,
                    })}
                  </span>
                </div>
              )}
            </div>
            <div className="flex justify-end">
              <Button onClick={() => handleOpenChange(false)}>
                {t("dialog.close")}
              </Button>
            </div>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
