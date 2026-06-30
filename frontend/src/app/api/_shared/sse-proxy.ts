import { NextResponse } from "next/server";

import { getActiveOrganizationId, getUserId } from "@/lib/auth/session";

const BACKEND_API_URL = process.env.BACKEND_API_URL || "http://backend:8000";

/**
 * Creates an SSE proxy response that forwards a backend SSE stream
 * to the client with X-User-ID authentication.
 *
 * Cancellation: returning backendResponse.body directly to Next.js does
 * not propagate downstream client disconnects back to the upstream
 * fetch(), which leaks one TCP socket per closed EventSource. Wrapping
 * the body in a manual ReadableStream whose cancel() aborts an owned
 * AbortController fixes that — both a mid-stream client disconnect and
 * an abort-before-streaming-starts tear down the upstream socket.
 */
export async function proxySSEStream(
  request: Request,
  backendPath: string
): Promise<Response> {
  const userId = await getUserId();
  const activeOrganizationId = await getActiveOrganizationId();
  const backendUrl = `${BACKEND_API_URL}${backendPath}`;

  const authHeaders: Record<string, string> = {
    "X-User-ID": userId,
  };
  if (activeOrganizationId) {
    authHeaders["X-Organization-ID"] = activeOrganizationId;
  }

  const upstream = new AbortController();
  const onClientAbort = () => upstream.abort();
  if (request.signal.aborted) {
    upstream.abort();
  } else {
    request.signal.addEventListener("abort", onClientAbort, { once: true });
  }

  let backendResponse: Response;
  try {
    backendResponse = await fetch(backendUrl, {
      headers: authHeaders,
      cache: "no-store",
      signal: upstream.signal,
    });
  } catch {
    request.signal.removeEventListener("abort", onClientAbort);
    return NextResponse.json(
      { error: "Failed to connect to backend SSE stream" },
      { status: 502 }
    );
  }

  if (!backendResponse.ok) {
    request.signal.removeEventListener("abort", onClientAbort);
    upstream.abort();
    return NextResponse.json(
      { error: `Backend error: ${backendResponse.statusText}` },
      { status: backendResponse.status }
    );
  }

  if (!backendResponse.body) {
    request.signal.removeEventListener("abort", onClientAbort);
    upstream.abort();
    return NextResponse.json(
      { error: "No stream body from backend" },
      { status: 502 }
    );
  }

  const reader = backendResponse.body.getReader();

  const stream = new ReadableStream<Uint8Array>({
    async pull(controller) {
      try {
        const { done, value } = await reader.read();
        if (done) {
          controller.close();
          request.signal.removeEventListener("abort", onClientAbort);
          return;
        }
        controller.enqueue(value);
      } catch (err) {
        controller.error(err);
        upstream.abort();
        request.signal.removeEventListener("abort", onClientAbort);
      }
    },
    async cancel() {
      upstream.abort();
      try {
        await reader.cancel();
      } catch {
        // upstream already torn down — ignore
      }
      request.signal.removeEventListener("abort", onClientAbort);
    },
  });

  return new Response(stream, {
    headers: {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache",
      Connection: "keep-alive",
      "X-Accel-Buffering": "no",
    },
  });
}
