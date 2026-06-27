import { Fragment } from "react";

interface ParameterizedNameProps {
  name: string;
  parameterValues?: Record<string, string> | null;
}

/**
 * Renders a subtask name with resolved parameter values in bold.
 * Searches the name for any value present in parameterValues and wraps matches in <strong>.
 */
export function ParameterizedName({
  name,
  parameterValues,
}: ParameterizedNameProps) {
  if (!parameterValues || Object.keys(parameterValues).length === 0) {
    return <>{name}</>;
  }

  const values = Object.values(parameterValues).filter(Boolean);
  if (values.length === 0) {
    return <>{name}</>;
  }

  // Build regex matching any parameter value, longest first to avoid partial matches
  const sorted = [...values].sort((a, b) => b.length - a.length);
  const escaped = sorted.map((v) => v.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"));
  const regex = new RegExp(`(${escaped.join("|")})`, "g");

  const parts = name.split(regex);

  const valueSet = new Set(values);

  return (
    <>
      {parts.map((part, i) =>
        valueSet.has(part) ? (
          <strong key={i}>{part}</strong>
        ) : (
          <Fragment key={i}>{part}</Fragment>
        )
      )}
    </>
  );
}
