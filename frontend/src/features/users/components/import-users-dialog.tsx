"use client";

import {
  Upload,
  Download,
  CheckCircle,
  AlertCircle,
  SkipForward,
} from "lucide-react";
import { useCallback, useRef, useState } from "react";
import { toast } from "sonner";

import type {
  UserImportValidationResponse,
  UserImportResponse,
} from "@/lib/api/backend-client";

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
  useValidateUserImportMutation,
  useImportUsersMutation,
} from "../hooks/use-import-users-mutation";

type Step = "upload" | "preview" | "complete";

export function ImportUsersDialog() {
  const [open, setOpen] = useState(false);
  const [step, setStep] = useState<Step>("upload");
  const [csvContent, setCsvContent] = useState<string>("");
  const [fileName, setFileName] = useState<string>("");
  const [validation, setValidation] =
    useState<UserImportValidationResponse | null>(null);
  const [importResult, setImportResult] = useState<UserImportResponse | null>(
    null
  );
  const [showErrors, setShowErrors] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const validateMutation = useValidateUserImportMutation();
  const importMutation = useImportUsersMutation();

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
    if (!isOpen) reset();
  };

  const processFile = useCallback((file: File) => {
    if (file.size > 5 * 1024 * 1024) {
      toast.error("File too large", {
        description: "Maximum file size is 5MB",
      });
      return;
    }

    setFileName(file.name);
    const reader = new FileReader();
    reader.onload = (event) => {
      setCsvContent(event.target?.result as string);
    };
    reader.readAsText(file, "UTF-8");
  }, []);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    processFile(file);
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
        toast.error("Validation failed", { description: error.message });
      },
    });
  };

  const handleImport = () => {
    importMutation.mutate(csvContent, {
      onSuccess: (result) => {
        setImportResult(result);
        setStep("complete");
        toast.success(`${result.imported_count} users imported successfully`);
      },
      onError: (error) => {
        toast.error("Import failed", { description: error.message });
      },
    });
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant="outline">
          <Upload className="mr-2 h-4 w-4" />
          Import
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg max-h-[90vh] flex flex-col">
        {step === "upload" && (
          <>
            <DialogHeader>
              <DialogTitle>Import Users from CSV</DialogTitle>
              <DialogDescription>
                Upload a CSV file to invite users in bulk. Each valid row
                triggers an invitation email.
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
                  if (file) processFile(file);
                }}
              >
                <Upload className="mx-auto h-8 w-8 text-gray-400 mb-2" />
                {fileName ? (
                  <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
                    {fileName}
                  </p>
                ) : (
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    Click to select or drag and drop a CSV file
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
                href="/web/templates/user-import-template.csv"
                download
                className="inline-flex items-center text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
              >
                <Download className="mr-1 h-4 w-4" />
                Download sample template
              </a>
            </div>
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={() => handleOpenChange(false)}>
                Cancel
              </Button>
              <Button
                onClick={handleValidate}
                disabled={!csvContent || validateMutation.isPending}
              >
                {validateMutation.isPending ? "Validating..." : "Validate"}
              </Button>
            </div>
          </>
        )}

        {step === "preview" && validation && (
          <>
            <DialogHeader>
              <DialogTitle>Import Preview</DialogTitle>
              <DialogDescription>
                Review the validation results before importing.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-3 py-4 flex-1 overflow-y-auto">
              <div className="flex items-center gap-2 text-sm">
                <CheckCircle className="h-4 w-4 text-green-500" />
                <span>
                  {validation.summary.valid_count} users ready to import
                </span>
              </div>
              {validation.summary.duplicate_count > 0 && (
                <div className="flex items-center gap-2 text-sm">
                  <SkipForward className="h-4 w-4 text-yellow-500" />
                  <span>
                    {validation.summary.duplicate_count} duplicates (will be
                    skipped)
                  </span>
                </div>
              )}
              {validation.summary.error_count > 0 && (
                <div className="flex items-center gap-2 text-sm">
                  <AlertCircle className="h-4 w-4 text-red-500" />
                  <span>{validation.summary.error_count} errors</span>
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
                    {showErrors ? "Hide details" : "Show details"}
                  </button>
                  {showErrors && (
                    <div className="mt-2 max-h-48 overflow-y-auto space-y-2">
                      {validation.error_rows.map((row) => (
                        <div
                          key={row.row_number}
                          className="text-xs bg-red-50 dark:bg-red-950 p-2 rounded"
                        >
                          <span className="font-medium">
                            Row {row.row_number}
                            {row.email ? ` (${row.email})` : ""}:
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
                            Row {row.row_number}
                            {row.email ? ` (${row.email})` : ""}:
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
                Back
              </Button>
              <Button
                onClick={handleImport}
                disabled={
                  validation.summary.valid_count === 0 ||
                  importMutation.isPending
                }
              >
                {importMutation.isPending
                  ? "Importing..."
                  : `Import ${validation.summary.valid_count} users`}
              </Button>
            </div>
          </>
        )}

        {step === "complete" && importResult && (
          <>
            <DialogHeader>
              <DialogTitle>Import Complete</DialogTitle>
            </DialogHeader>
            <div className="space-y-3 py-4">
              <div className="flex items-center gap-2 text-sm">
                <CheckCircle className="h-4 w-4 text-green-500" />
                <span>{importResult.imported_count} users imported</span>
              </div>
              {importResult.skipped_count > 0 && (
                <div className="flex items-center gap-2 text-sm">
                  <SkipForward className="h-4 w-4 text-yellow-500" />
                  <span>{importResult.skipped_count} skipped</span>
                </div>
              )}
              {importResult.error_count > 0 && (
                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-sm">
                    <AlertCircle className="h-4 w-4 text-red-500" />
                    <span>{importResult.error_count} errors</span>
                  </div>
                  <div className="max-h-48 overflow-y-auto space-y-1">
                    {importResult.errors.map((e) => (
                      <div
                        key={e.row_number}
                        className="text-xs bg-red-50 dark:bg-red-950 p-2 rounded"
                      >
                        <span className="font-medium">
                          Row {e.row_number}
                          {e.email ? ` (${e.email})` : ""}:
                        </span>{" "}
                        {e.errors.join("; ")}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
            <div className="flex justify-end">
              <Button onClick={() => handleOpenChange(false)}>Close</Button>
            </div>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
