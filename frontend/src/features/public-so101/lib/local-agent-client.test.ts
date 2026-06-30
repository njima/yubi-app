import {
  DEFAULT_LOCAL_AGENT_BASE_URL,
  buildLocalAgentUrl,
  normalizeLocalAgentHealth,
  type LocalAgentHealthResponse,
} from "./local-agent-client";

const successfulHealth: LocalAgentHealthResponse = {
  ok: true,
  name: "lelab",
  version: "0.1.0",
  robotType: "so101",
  capabilities: ["health", "motor_check"],
};

const unsupportedHealth: LocalAgentHealthResponse = {
  ok: true,
  name: "lelab",
  version: "0.1.0",
  robotType: "unknown",
  capabilities: ["health"],
};

const unavailableHealth: LocalAgentHealthResponse = {
  ok: false,
  name: "lelab",
  version: "0.1.0",
  robotType: "so101",
  capabilities: ["health"],
};

if (DEFAULT_LOCAL_AGENT_BASE_URL !== "http://127.0.0.1:32101") {
  throw new Error("unexpected default local agent URL");
}

if (
  buildLocalAgentUrl("http://127.0.0.1:32101/", "/health").toString() !==
  "http://127.0.0.1:32101/health"
) {
  throw new Error("local agent URL builder should normalize slashes");
}

if (normalizeLocalAgentHealth(successfulHealth).status !== "available") {
  throw new Error("SO101 health response should be available");
}

if (normalizeLocalAgentHealth(unsupportedHealth).status !== "unsupported") {
  throw new Error("non-SO101 health response should be unsupported");
}

if (normalizeLocalAgentHealth(unavailableHealth).status !== "unavailable") {
  throw new Error("ok=false health response should be unavailable");
}
