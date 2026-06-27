// --- Layout structure ---

export interface PageLayoutConfig {
  title?: string;
  // Mode 1: two-column with sidebar (teleoperation style)
  main_area?: LayoutRow[];
  sidebar?: LayoutItem[];
  // Mode 2: single-column stack (task detail style)
  sections?: LayoutRow[];
}

export interface LayoutRow {
  type: "row";
  gap?: number;
  items: LayoutItem[];
}

/**
 * Layout item. `type` is the identifier:
 * - `"camera"` — renders a camera feed (requires `ref`, optional `overlay`)
 * - Any other string — looks up a registered component by type name
 */
export type LayoutItem = CameraLayoutItem | ComponentLayoutItem;

export interface CameraLayoutItem {
  type: "camera";
  ref: string;
  span?: number;
  overlay?: boolean;
}

export interface ComponentLayoutItem {
  type: string; // component identifier (e.g., "task-information", "status-card")
  span?: number;
}

export function isCameraItem(item: LayoutItem): item is CameraLayoutItem {
  return item.type === "camera" && "ref" in item;
}

// --- Multi-view layout config ---

/** Multiple named views for a page section (e.g., "default", "subtasks") */
export interface ViewsConfig {
  [viewName: string]: PageLayoutConfig;
}

/** Full layout config per robot. Keys match URL path segments. */
export interface RobotLayoutConfig {
  teleoperation?: ViewsConfig;
  [key: string]: ViewsConfig | undefined;
}

// --- Default layouts ---

export const DEFAULT_TELEOP_VIEWS: ViewsConfig = {
  default: {
    title: "Teleoperation Console",
    main_area: [
      {
        type: "row",
        items: [{ type: "camera", ref: "*main*", span: 2, overlay: true }],
      },
      {
        type: "row",
        items: [
          { type: "camera", ref: "*left*", span: 1 },
          { type: "camera", ref: "*right*", span: 1 },
        ],
      },
    ],
    sidebar: [
      { type: "task-information" },
      { type: "status-card" },
      { type: "3d-model" },
    ],
  },
  subtasks: {
    title: "Subtasks Overview",
    main_area: [
      {
        type: "row",
        items: [{ type: "subtask-detail-list", span: 2 }],
      },
    ],
    sidebar: [{ type: "teach-me-card" }, { type: "status-card" }],
  },
};
