"use client";

import { Plus, Trash2, X } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export type ParameterDefinition = z.infer<typeof schemas.TaskVersionParameter>;

interface TaskVersionParametersEditorProps {
  value: ParameterDefinition[];
  onChange: (parameters: ParameterDefinition[]) => void;
}

export function TaskVersionParametersEditor({
  value,
  onChange,
}: TaskVersionParametersEditorProps) {
  const { t } = useTranslation();
  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <label className="text-sm font-medium">
          {t("parametersEditor.parameters")}
        </label>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="gap-1"
          onClick={() => onChange([...value, { key: "", values: [] }])}
        >
          <Plus className="h-3 w-3" />
          {t("parametersEditor.addParameter")}
        </Button>
      </div>

      {value.length === 0 && (
        <p className="text-sm text-gray-500 italic">
          {t("parametersEditor.noParameters")}
        </p>
      )}

      {value.map((param, index) => (
        <ParameterRow
          key={index}
          param={param}
          onChange={(updated) => {
            const next = [...value];
            next[index] = updated;
            onChange(next);
          }}
          onRemove={() => {
            onChange(value.filter((_, i) => i !== index));
          }}
        />
      ))}
    </div>
  );
}

interface ParameterRowProps {
  param: ParameterDefinition;
  onChange: (param: ParameterDefinition) => void;
  onRemove: () => void;
}

function ParameterRow({ param, onChange, onRemove }: ParameterRowProps) {
  const { t } = useTranslation();
  const [newValue, setNewValue] = useState("");

  const addValue = () => {
    const trimmed = newValue.trim();
    if (trimmed && !param.values.includes(trimmed)) {
      onChange({ ...param, values: [...param.values, trimmed] });
      setNewValue("");
    }
  };

  const removeValue = (valueToRemove: string) => {
    onChange({
      ...param,
      values: param.values.filter((v) => v !== valueToRemove),
    });
  };

  return (
    <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3 space-y-2">
      <div className="flex items-center gap-2">
        <Input
          placeholder={t("parametersEditor.keyPlaceholder")}
          value={param.key}
          onChange={(e) => onChange({ ...param, key: e.target.value })}
          className="flex-1"
        />
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={onRemove}
          className="text-red-500 hover:text-red-700"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>

      <div className="flex flex-wrap gap-1">
        {param.values.map((v) => (
          <span
            key={v}
            className="inline-flex items-center gap-1 rounded-full bg-blue-100 dark:bg-blue-900 px-2.5 py-0.5 text-xs font-medium text-blue-800 dark:text-blue-200"
          >
            {v}
            <button
              type="button"
              onClick={() => removeValue(v)}
              className="hover:text-red-600"
            >
              <X className="h-3 w-3" />
            </button>
          </span>
        ))}
      </div>

      <div className="flex gap-2">
        <Input
          placeholder={t("parametersEditor.addValuePlaceholder")}
          value={newValue}
          onChange={(e) => setNewValue(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              addValue();
            }
          }}
          className="flex-1"
        />
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={addValue}
          disabled={!newValue.trim()}
        >
          {t("parametersEditor.add")}
        </Button>
      </div>
    </div>
  );
}
