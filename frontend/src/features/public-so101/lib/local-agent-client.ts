export const DEFAULT_LOCAL_AGENT_BASE_URL = "http://127.0.0.1:32101";

export type LocalAgentStatus =
  | "available"
  | "unavailable"
  | "unsupported"
  | "timeout";

export interface LocalAgentHealthResponse {
  ok: boolean;
  name?: string;
  version?: string;
  robotType?: string;
  capabilities?: string[];
}

export interface LocalAgentHealth {
  status: LocalAgentStatus;
  name: string;
  version: string;
  robotType: string;
  capabilities: string[];
  message: string;
}

export interface CheckLocalAgentHealthOptions {
  baseUrl?: string;
  timeoutMs?: number;
  fetcher?: typeof fetch;
}

const DEFAULT_TIMEOUT_MS = 2500;

export function buildLocalAgentUrl(baseUrl: string, path: string) {
  const normalizedBaseUrl = baseUrl.endsWith("/") ? baseUrl : `${baseUrl}/`;
  const normalizedPath = path.startsWith("/") ? path.slice(1) : path;

  return new URL(normalizedPath, normalizedBaseUrl);
}

export function getLocalAgentBaseUrl() {
  return (
    process.env.NEXT_PUBLIC_SO101_LOCAL_AGENT_URL ??
    DEFAULT_LOCAL_AGENT_BASE_URL
  );
}

export function normalizeLocalAgentHealth(
  response: LocalAgentHealthResponse
): LocalAgentHealth {
  const name = response.name ?? "local-agent";
  const version = response.version ?? "unknown";
  const robotType = response.robotType ?? "unknown";
  const capabilities = response.capabilities ?? [];

  if (!response.ok) {
    return {
      status: "unavailable",
      name,
      version,
      robotType,
      capabilities,
      message: "Local agent is not ready.",
    };
  }

  if (robotType.toLowerCase() !== "so101") {
    return {
      status: "unsupported",
      name,
      version,
      robotType,
      capabilities,
      message: `Connected agent reports ${robotType}, not SO101.`,
    };
  }

  return {
    status: "available",
    name,
    version,
    robotType,
    capabilities,
    message: "SO101 local agent is ready.",
  };
}

export async function checkLocalAgentHealth({
  baseUrl = getLocalAgentBaseUrl(),
  timeoutMs = DEFAULT_TIMEOUT_MS,
  fetcher = fetch,
}: CheckLocalAgentHealthOptions = {}): Promise<LocalAgentHealth> {
  const controller = new AbortController();
  const timeout = globalThis.setTimeout(() => controller.abort(), timeoutMs);

  try {
    const response = await fetcher(buildLocalAgentUrl(baseUrl, "/health"), {
      method: "GET",
      signal: controller.signal,
      headers: {
        Accept: "application/json",
      },
    });

    if (!response.ok) {
      return {
        status: "unavailable",
        name: "local-agent",
        version: "unknown",
        robotType: "unknown",
        capabilities: [],
        message: `Local agent returned HTTP ${response.status}.`,
      };
    }

    const health = (await response.json()) as LocalAgentHealthResponse;
    return normalizeLocalAgentHealth(health);
  } catch (error) {
    if (error instanceof DOMException && error.name === "AbortError") {
      return {
        status: "timeout",
        name: "local-agent",
        version: "unknown",
        robotType: "unknown",
        capabilities: [],
        message: "Timed out while connecting to local agent.",
      };
    }

    return {
      status: "unavailable",
      name: "local-agent",
      version: "unknown",
      robotType: "unknown",
      capabilities: [],
      message: "Local agent is not reachable.",
    };
  } finally {
    globalThis.clearTimeout(timeout);
  }
}
