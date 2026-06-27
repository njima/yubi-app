"use client";

import { Component, type ErrorInfo, type ReactNode } from "react";

import { cn } from "@/shared/lib/utils";

import { getLayoutComponent } from "@/features/robots/lib/teleop-layout-registry";
import type { LayoutContext } from "@/features/robots/lib/teleop-layout-registry";
import type {
  LayoutItem,
  LayoutRow,
  PageLayoutConfig,
} from "@/features/robots/lib/teleop-layout-types";

interface LayoutRendererProps {
  layout: PageLayoutConfig;
  context: LayoutContext;
  /** When true, main_area fills available height (fullscreen mode) */
  fillHeight?: boolean;
}

export function LayoutRenderer({
  layout,
  context,
  fillHeight = false,
}: LayoutRendererProps) {
  // Mode 1: main_area + sidebar
  if (layout.main_area) {
    return (
      <div
        className={cn(
          "grid grid-cols-1 lg:grid-cols-3",
          fillHeight ? "h-full gap-2" : "gap-4"
        )}
      >
        <div
          className={cn(
            "lg:col-span-2",
            fillHeight ? "overflow-hidden" : "space-y-2"
          )}
          style={fillHeight ? { maxHeight: "calc(100vh - 6rem)" } : undefined}
        >
          {layout.main_area.map((row, i) => {
            // In fullscreen, single-item rows get 50% height, multi-item rows get 45%
            const rowMaxHeight = fillHeight
              ? row.items.length <= 1
                ? "50vh"
                : "45vh"
              : undefined;
            return (
              <RowRenderer
                key={i}
                row={row}
                context={context}
                maxHeight={rowMaxHeight}
              />
            );
          })}
        </div>
        {layout.sidebar && (
          <div className="space-y-4">
            {layout.sidebar.map((item, i) => (
              <ItemRenderer key={i} item={item} context={context} />
            ))}
          </div>
        )}
      </div>
    );
  }

  // Mode 2: sections (single-column stack)
  if (layout.sections) {
    return (
      <div className="space-y-6">
        {layout.sections.map((row, i) => (
          <RowRenderer key={i} row={row} context={context} />
        ))}
      </div>
    );
  }

  return null;
}

// --- Internal components ---

function RowRenderer({
  row,
  context,
  maxHeight,
}: {
  row: LayoutRow;
  context: LayoutContext;
  maxHeight?: string;
}) {
  const totalSpan = row.items.reduce((sum, item) => sum + (item.span ?? 1), 0);
  const gapRem = (row.gap ?? 2) * 0.25;

  return (
    <div
      className={maxHeight ? "overflow-hidden" : undefined}
      style={{
        display: "grid",
        gridTemplateColumns: `repeat(${totalSpan}, minmax(0, 1fr))`,
        gap: `${gapRem}rem`,
        maxHeight,
      }}
    >
      {row.items.map((item, i) => {
        const span = item.span ?? 1;
        const key = `${item.type}-${i}`;
        return (
          <div
            key={key}
            style={span > 1 ? { gridColumn: `span ${span}` } : undefined}
          >
            <ItemRenderer item={item} context={context} />
          </div>
        );
      })}
    </div>
  );
}

function ItemRenderer({
  item,
  context,
}: {
  item: LayoutItem;
  context: LayoutContext;
}) {
  const render = getLayoutComponent(item.type);
  if (!render) {
    console.warn(`[LayoutRenderer] component "${item.type}" not registered`);
    return null;
  }
  return (
    <ComponentErrorBoundary componentId={item.type}>
      {render(context, item)}
    </ComponentErrorBoundary>
  );
}

// --- Error boundary for layout components ---

class ComponentErrorBoundary extends Component<
  { componentId: string; children: ReactNode },
  { hasError: boolean }
> {
  constructor(props: { componentId: string; children: ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error(
      `[LayoutRenderer] component "${this.props.componentId}" threw:`,
      error,
      info
    );
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600 dark:border-red-800 dark:bg-red-950 dark:text-red-400">
          Component error: {this.props.componentId}
        </div>
      );
    }
    return this.props.children;
  }
}
