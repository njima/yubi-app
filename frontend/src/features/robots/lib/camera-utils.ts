import type { LayoutCamera } from "@/shared/lib/layout-registry";
import type { RobotLayoutConfig } from "@/shared/lib/layout-types";

export type Camera = LayoutCamera;

export interface StreamConfig {
  quality?: number;
  width?: number;
  height?: number;
  type?: string;
}

/**
 * Extract layout config from robot_config.layout
 * Returns undefined if missing or malformed (caller uses default layout).
 */
export function extractRobotLayout(
  robotConfig?: Record<string, unknown> | null
): RobotLayoutConfig | undefined {
  if (!robotConfig) return undefined;
  try {
    const layout = robotConfig.layout;
    if (typeof layout !== "object" || layout === null) return undefined;
    const l = layout as Record<string, unknown>;
    // Validate: each top-level key (e.g., "teleoperation") maps to a views
    // object where each view has optional main_area/sidebar/sections arrays.
    for (const pageKey of Object.keys(l)) {
      const views = l[pageKey];
      if (typeof views !== "object" || views === null) return undefined;
      for (const view of Object.values(views as Record<string, unknown>)) {
        if (typeof view !== "object" || view === null) return undefined;
        const v = view as Record<string, unknown>;
        if (v.main_area !== undefined && !Array.isArray(v.main_area))
          return undefined;
        if (v.sidebar !== undefined && !Array.isArray(v.sidebar))
          return undefined;
        if (v.sections !== undefined && !Array.isArray(v.sections))
          return undefined;
      }
    }
    return l as RobotLayoutConfig;
  } catch {
    console.warn("[extractRobotLayout] malformed layout config, using default");
    return undefined;
  }
}

/**
 * Extract MJPEG stream config (quality, resolution) from robot_config.stream
 */
export function extractStreamConfig(
  robotConfig?: Record<string, unknown> | null
): StreamConfig | undefined {
  if (!robotConfig) return undefined;
  const stream = robotConfig.stream;
  if (typeof stream !== "object" || stream === null) return undefined;
  const s = stream as Record<string, unknown>;
  return {
    quality: typeof s.quality === "number" ? s.quality : undefined,
    width: typeof s.width === "number" ? s.width : undefined,
    height: typeof s.height === "number" ? s.height : undefined,
    type: typeof s.type === "string" ? s.type : undefined,
  };
}

/**
 * Extract host and port from robot_config
 */
export function extractHostPort(robotConfig?: Record<string, unknown> | null): {
  host: string | undefined;
  port: number | undefined;
  rosbridgePort: number | undefined;
} {
  if (!robotConfig) {
    return { host: undefined, port: undefined, rosbridgePort: undefined };
  }

  const host =
    typeof robotConfig.host === "string" ? robotConfig.host : undefined;
  const port =
    typeof robotConfig.port === "number" ? robotConfig.port : undefined;
  const rosbridgePort =
    typeof robotConfig.rosbridge_port === "number"
      ? robotConfig.rosbridge_port
      : undefined;

  return { host, port, rosbridgePort };
}

/**
 * Extract camera list from robot_config
 * Supports multiple config formats for flexibility
 */
export function extractCameras(
  robotConfig?: Record<string, unknown> | null
): Camera[] {
  if (!robotConfig) {
    return [];
  }

  // Format 1: { cameras: [{ namespace, name }] }
  if (Array.isArray(robotConfig.cameras)) {
    return robotConfig.cameras
      .filter(
        (cam): cam is Camera =>
          typeof cam === "object" &&
          cam !== null &&
          "namespace" in cam &&
          typeof cam.namespace === "string"
      )
      .map((cam) => ({
        namespace: cam.namespace,
        name: cam.name as string | undefined,
      }));
  }

  // Format 2: { namespace: string, name?: string }
  if (typeof robotConfig.namespace === "string") {
    return [
      {
        namespace: robotConfig.namespace,
        name: robotConfig.name as string | undefined,
      },
    ];
  }

  // Format 3: { [key]: { namespace: string } }
  const cameras: Camera[] = [];
  for (const [key, value] of Object.entries(robotConfig)) {
    if (
      typeof value === "object" &&
      value !== null &&
      "namespace" in value &&
      typeof value.namespace === "string"
    ) {
      const cameraValue = value as { namespace: string; name?: string };
      cameras.push({
        namespace: cameraValue.namespace,
        name: cameraValue.name || key,
      });
    }
  }

  return cameras;
}
