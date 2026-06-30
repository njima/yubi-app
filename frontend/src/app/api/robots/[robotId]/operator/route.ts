import { NextResponse } from "next/server";

import { clearRobotOperator } from "@/lib/api/backend-client";
import { handleApiError } from "@/lib/api/response-helpers";
import { getActiveOrganizationId, getUserId } from "@/lib/auth/session";

const BACKEND_API_URL = process.env.BACKEND_API_URL || "http://backend:8000";

type Params = Promise<{ robotId: string }>;

/**
 * GET /web/api/robots/:robotId/operator
 * Read-only check for current operator (no side effects)
 */
export async function GET(request: Request, { params }: { params: Params }) {
  try {
    const { robotId } = await params;
    const userId = await getUserId();
    const activeOrganizationId = await getActiveOrganizationId();
    const authHeaders: Record<string, string> = { "X-User-ID": userId };
    if (activeOrganizationId) {
      authHeaders["X-Organization-ID"] = activeOrganizationId;
    }

    const response = await fetch(
      `${BACKEND_API_URL}/api/robots/${robotId}/operator`,
      {
        headers: authHeaders,
        cache: "no-store",
      }
    );

    if (response.status === 204) {
      return new NextResponse(null, { status: 204 });
    }

    if (response.ok) {
      const operator = await response.json();
      return NextResponse.json(operator);
    }

    return new NextResponse(null, { status: response.status });
  } catch (error) {
    return handleApiError(error);
  }
}

/**
 * PUT /web/api/robots/:robotId/operator
 * Set or refresh the active operator (heartbeat).
 * Must pass 409 body through (contains current operator info).
 */
export async function PUT(request: Request, { params }: { params: Params }) {
  try {
    const { robotId } = await params;
    const body = await request.json();
    const userId = await getUserId();
    const activeOrganizationId = await getActiveOrganizationId();
    const authHeaders: Record<string, string> = {
      "Content-Type": "application/json",
      "X-User-ID": userId,
    };
    if (activeOrganizationId) {
      authHeaders["X-Organization-ID"] = activeOrganizationId;
    }

    const response = await fetch(
      `${BACKEND_API_URL}/api/robots/${robotId}/operator`,
      {
        method: "PUT",
        headers: authHeaders,
        body: JSON.stringify(body),
      }
    );

    if (response.status === 409) {
      const operator = await response.json();
      return NextResponse.json(operator, { status: 409 });
    }

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      return NextResponse.json(
        { error: error.message || response.statusText },
        { status: response.status }
      );
    }

    return new NextResponse(null, { status: 200 });
  } catch (error) {
    return handleApiError(error);
  }
}

/**
 * DELETE /web/api/robots/:robotId/operator
 * Clear the active operator
 */
export async function DELETE(
  _request: Request,
  { params }: { params: Params }
) {
  try {
    const { robotId } = await params;
    await clearRobotOperator(robotId);
    return new NextResponse(null, { status: 204 });
  } catch (error) {
    return handleApiError(error);
  }
}
