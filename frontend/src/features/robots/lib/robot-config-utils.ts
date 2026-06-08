/**
 * Merge the dedicated Host field into the parsed Advanced Settings JSON.
 * Host wins on key collision; returns null on invalid JSON.
 */
export function mergeRobotConfigWithHost(
  robotConfigRaw: unknown,
  host: string
): Record<string, unknown> | null {
  let parsedConfig: Record<string, unknown> = {};
  if (typeof robotConfigRaw === "string") {
    if (robotConfigRaw.trim() !== "") {
      try {
        parsedConfig = JSON.parse(robotConfigRaw);
      } catch {
        return null;
      }
    }
  } else if (robotConfigRaw && typeof robotConfigRaw === "object") {
    parsedConfig = robotConfigRaw as Record<string, unknown>;
  }
  return { ...parsedConfig, host };
}

/**
 * Split the stored robot_config into the dedicated Host field and the
 * remaining Advanced Settings.
 */
export function splitHostFromRobotConfig(robotConfig: unknown): {
  host: string;
  advancedSettings: Record<string, unknown> | undefined;
} {
  if (!robotConfig || typeof robotConfig !== "object") {
    return { host: "", advancedSettings: undefined };
  }
  const { host, ...advancedSettings } = robotConfig as Record<string, unknown>;
  return {
    host: typeof host === "string" ? host : "",
    advancedSettings: Object.keys(advancedSettings).length
      ? advancedSettings
      : undefined,
  };
}
