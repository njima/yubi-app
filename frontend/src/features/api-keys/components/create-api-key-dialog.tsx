"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { KeyRound } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { SearchableSelect } from "@/components/ui/searchable-select";

import { useRobotSearchOptions } from "@/features/robots/hooks/use-robot-search-options";
import { useMeQuery } from "@/features/users/hooks/use-me-query";

import { RawKeyDisplayDialog } from "./raw-key-display-dialog";
import { useCreateApiKeyMutation } from "../hooks/use-create-api-key-mutation";

interface CreateApiKeyDialogProps {
  /**
   * When provided, the robot field is prefilled with this id and the input is
   * disabled so the operator cannot change it. Used by the robot detail page
   * "Issue API key for this robot" entry point.
   *
   * Always pass `prefilledRobotName` together so the dropdown can render a
   * human-readable label without an extra fetch.
   */
  prefilledRobotId?: string;
  /** Display label for the prefilled robot. Required iff `prefilledRobotId` is set. */
  prefilledRobotName?: string;
  /**
   * Disable the trigger button. Use this while the caller is still resolving
   * `prefilledRobotName` for a `prefilledRobotId` taken from the URL so the
   * dialog never opens with an empty/disabled robot field.
   */
  triggerDisabled?: boolean;
}

export function CreateApiKeyDialog({
  prefilledRobotId,
  prefilledRobotName,
  triggerDisabled,
}: CreateApiKeyDialogProps = {}) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [rawKey, setRawKey] = useState<string | null>(null);

  const { data: me } = useMeQuery();
  const { mutateAsync, isPending } = useCreateApiKeyMutation();
  const robotSearch = useRobotSearchOptions();

  const formSchema = z.object({
    name: z
      .string()
      .min(1, t("createApiKeyDialog.fieldNamePlaceholder"))
      .max(255),
    robot_id: z.string().min(1, t("createApiKeyDialog.fieldRobotPlaceholder")),
  });
  type FormValues = z.infer<typeof formSchema>;

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      robot_id: prefilledRobotId ?? "",
    },
  });

  const handleSubmit = async (values: FormValues) => {
    try {
      const result = await mutateAsync({
        name: values.name,
        robot_id: values.robot_id,
        expires_at: null,
      });
      setRawKey(result.key);
      setOpen(false);
      form.reset({
        name: "",
        robot_id: prefilledRobotId ?? "",
      });
    } catch {
      // toast already raised
    }
  };

  return (
    <>
      <Dialog
        open={open}
        onOpenChange={(next) => {
          setOpen(next);
          form.reset({
            name: "",
            robot_id: prefilledRobotId ?? "",
          });
        }}
      >
        <DialogTrigger asChild>
          <Button className="gap-1" disabled={triggerDisabled}>
            <KeyRound className="h-4 w-4" />
            {t("createApiKeyDialog.trigger")}
          </Button>
        </DialogTrigger>

        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{t("createApiKeyDialog.title")}</DialogTitle>
            <DialogDescription>
              {t("createApiKeyDialog.description")}
            </DialogDescription>
          </DialogHeader>

          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(handleSubmit)}
              className="space-y-4"
            >
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("createApiKeyDialog.fieldName")}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t(
                          "createApiKeyDialog.fieldNamePlaceholder"
                        )}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="robot_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("createApiKeyDialog.fieldRobot")}</FormLabel>
                    <FormControl>
                      <SearchableSelect
                        value={field.value}
                        onValueChange={(v) => {
                          field.onChange(v);
                          robotSearch.onValueChange(v);
                        }}
                        options={robotSearch.options}
                        onSearch={robotSearch.onSearch}
                        isLoading={robotSearch.isLoading}
                        selectedLabel={
                          prefilledRobotId
                            ? prefilledRobotName
                            : field.value
                              ? robotSearch.selectedLabel
                              : undefined
                        }
                        placeholder={t(
                          "createApiKeyDialog.fieldRobotPlaceholder"
                        )}
                        disabled={!!prefilledRobotId}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {me?.display_name && (
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {t("createApiKeyDialog.issuer", { name: me.display_name })}
                </p>
              )}

              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setOpen(false)}
                  disabled={isPending}
                >
                  {t("createApiKeyDialog.cancel")}
                </Button>
                <Button type="submit" disabled={isPending}>
                  {isPending
                    ? t("createApiKeyDialog.submitting")
                    : t("createApiKeyDialog.submit")}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      <RawKeyDisplayDialog rawKey={rawKey} onClose={() => setRawKey(null)} />
    </>
  );
}
