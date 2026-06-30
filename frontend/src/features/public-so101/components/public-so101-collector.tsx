"use client";

import {
  Activity,
  Cable,
  CheckCircle2,
  Circle,
  Download,
  Gauge,
  Hand,
  Play,
  RotateCcw,
  ShieldCheck,
  SlidersHorizontal,
  Square,
  UserRound,
} from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";

import { cn } from "@/lib/utils";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

type AgentState = "offline" | "connected";
type StepState = "idle" | "running" | "done";
type RecordingState = "idle" | "recording" | "ready";

interface WorkflowState {
  agent: AgentState;
  motorCheck: StepState;
  calibration: StepState;
  recording: RecordingState;
  taskName: string;
}

const defaultTask = "Pick and place household object";

const genericTasks = [
  defaultTask,
  "Open drawer and retrieve item",
  "Wipe tabletop area",
];

const defaultState: WorkflowState = {
  agent: "offline",
  motorCheck: "idle",
  calibration: "idle",
  recording: "idle",
  taskName: defaultTask,
};

export function PublicSo101Collector() {
  const [state, setState] = useState<WorkflowState>(defaultState);

  const canRunChecks = state.agent === "connected";
  const canTeleoperate =
    state.agent === "connected" &&
    state.motorCheck === "done" &&
    state.calibration === "done";
  const canDownload = state.recording === "ready";

  const completedSteps = useMemo(() => {
    return [
      state.agent === "connected",
      state.motorCheck === "done",
      state.calibration === "done",
      state.recording === "ready",
    ].filter(Boolean).length;
  }, [state]);

  function connectAgent() {
    setState((current) => ({ ...current, agent: "connected" }));
  }

  function runMotorCheck() {
    if (!canRunChecks) return;
    setState((current) => ({ ...current, motorCheck: "running" }));
    window.setTimeout(() => {
      setState((current) => ({ ...current, motorCheck: "done" }));
    }, 450);
  }

  function runCalibration() {
    if (!canRunChecks) return;
    setState((current) => ({ ...current, calibration: "running" }));
    window.setTimeout(() => {
      setState((current) => ({ ...current, calibration: "done" }));
    }, 450);
  }

  function startRecording() {
    if (!canTeleoperate) return;
    setState((current) => ({ ...current, recording: "recording" }));
  }

  function stopRecording() {
    if (state.recording !== "recording") return;
    setState((current) => ({ ...current, recording: "ready" }));
  }

  function resetWorkflow() {
    setState(defaultState);
  }

  function downloadManifest() {
    const manifest = {
      robot_type: "so101",
      mode: "public_guest",
      storage: "local_only",
      task_name: state.taskName,
      motor_check: state.motorCheck,
      calibration: state.calibration,
      recording: state.recording,
      generated_at: new Date().toISOString(),
    };
    const blob = new Blob([JSON.stringify(manifest, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = "so101-public-collection-manifest.json";
    anchor.click();
    URL.revokeObjectURL(url);
  }

  return (
    <main className="min-h-screen bg-gray-50 text-gray-950 dark:bg-gray-950 dark:text-gray-50">
      <div className="mx-auto flex min-h-screen w-full max-w-7xl flex-col gap-5 px-4 py-5 sm:px-6 lg:px-8">
        <header className="flex flex-col gap-4 border-b border-gray-200 pb-5 dark:border-gray-800 md:flex-row md:items-center md:justify-between">
          <div>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-md bg-gray-900 text-white dark:bg-gray-100 dark:text-gray-950">
                <Hand className="h-5 w-5" />
              </div>
              <div>
                <h1 className="text-xl font-semibold sm:text-2xl">
                  SO101 Public Collector
                </h1>
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  Guest workflow for local-only checks, calibration,
                  teleoperation, and download.
                </p>
              </div>
            </div>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <Badge
              variant={state.agent === "connected" ? "success" : "outline"}
            >
              {state.agent === "connected"
                ? "Agent connected"
                : "Agent offline"}
            </Badge>
            <Button variant="outline" asChild>
              <Link href="/login">
                <UserRound className="h-4 w-4" />
                Sign in
              </Link>
            </Button>
          </div>
        </header>

        <section className="grid gap-4 lg:grid-cols-[280px_minmax(0,1fr)]">
          <aside className="space-y-4">
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-base">Session</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <ProgressRow
                  icon={<Cable className="h-4 w-4" />}
                  label="Local agent"
                  done={state.agent === "connected"}
                />
                <ProgressRow
                  icon={<Gauge className="h-4 w-4" />}
                  label="Motor check"
                  done={state.motorCheck === "done"}
                />
                <ProgressRow
                  icon={<SlidersHorizontal className="h-4 w-4" />}
                  label="Calibration"
                  done={state.calibration === "done"}
                />
                <ProgressRow
                  icon={<Download className="h-4 w-4" />}
                  label="Local artifact"
                  done={state.recording === "ready"}
                />
                <Separator />
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">
                    Completed
                  </span>
                  <span className="font-medium">{completedSteps}/4</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="text-base">Privacy</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 text-sm text-gray-600 dark:text-gray-400">
                <div className="flex items-start gap-2">
                  <ShieldCheck className="mt-0.5 h-4 w-4 text-emerald-600" />
                  <p>No database user is created for this guest session.</p>
                </div>
                <div className="flex items-start gap-2">
                  <ShieldCheck className="mt-0.5 h-4 w-4 text-emerald-600" />
                  <p>Collection state stays in the browser and local agent.</p>
                </div>
              </CardContent>
            </Card>
          </aside>

          <section className="space-y-4">
            <Tabs defaultValue="setup" className="gap-4">
              <TabsList className="w-full justify-start overflow-x-auto rounded-md bg-white dark:bg-gray-900">
                <TabsTrigger value="setup">Setup</TabsTrigger>
                <TabsTrigger value="checks">Checks</TabsTrigger>
                <TabsTrigger value="teleop">Teleop</TabsTrigger>
                <TabsTrigger value="download">Download</TabsTrigger>
              </TabsList>

              <TabsContent value="setup">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Cable className="h-5 w-5" />
                      Local Agent
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="grid gap-4 lg:grid-cols-[1fr_260px]">
                    <div className="space-y-3">
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        The browser connects to a local SO101 bridge. This shell
                        uses a mocked connection until the LeRobot bridge
                        contract is wired in.
                      </p>
                      <div className="grid gap-2 text-sm sm:grid-cols-3">
                        <Metric label="Bridge" value="mock-local" />
                        <Metric label="Robot" value="SO101" />
                        <Metric label="Storage" value="Local" />
                      </div>
                    </div>
                    <div className="flex flex-col gap-2">
                      <Button
                        onClick={connectAgent}
                        disabled={state.agent === "connected"}
                      >
                        <Cable className="h-4 w-4" />
                        Connect agent
                      </Button>
                      <Button variant="outline" onClick={resetWorkflow}>
                        <RotateCcw className="h-4 w-4" />
                        Reset session
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="checks">
                <div className="grid gap-4 xl:grid-cols-2">
                  <ActionCard
                    icon={<Gauge className="h-5 w-5" />}
                    title="Motor check"
                    description="Confirm joints respond before calibration."
                    state={state.motorCheck}
                    disabled={!canRunChecks}
                    actionLabel="Run motor check"
                    onAction={runMotorCheck}
                  />
                  <ActionCard
                    icon={<SlidersHorizontal className="h-5 w-5" />}
                    title="Calibration"
                    description="Capture neutral and range references locally."
                    state={state.calibration}
                    disabled={!canRunChecks}
                    actionLabel="Run calibration"
                    onAction={runCalibration}
                  />
                </div>
              </TabsContent>

              <TabsContent value="teleop">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Activity className="h-5 w-5" />
                      Teleoperation
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid gap-3 sm:grid-cols-3">
                      {genericTasks.map((task) => (
                        <button
                          key={task}
                          type="button"
                          onClick={() =>
                            setState((current) => ({
                              ...current,
                              taskName: task,
                            }))
                          }
                          className={cn(
                            "min-h-20 rounded-md border px-3 py-3 text-left text-sm transition-colors",
                            state.taskName === task
                              ? "border-gray-950 bg-gray-100 dark:border-gray-50 dark:bg-gray-800"
                              : "border-gray-200 bg-white hover:bg-gray-50 dark:border-gray-800 dark:bg-gray-900 dark:hover:bg-gray-800"
                          )}
                        >
                          <span className="font-medium">{task}</span>
                        </button>
                      ))}
                    </div>
                    <div className="flex flex-wrap items-center gap-2">
                      <Button
                        onClick={startRecording}
                        disabled={!canTeleoperate || state.recording !== "idle"}
                      >
                        <Play className="h-4 w-4" />
                        Start recording
                      </Button>
                      <Button
                        variant="outline"
                        onClick={stopRecording}
                        disabled={state.recording !== "recording"}
                      >
                        <Square className="h-4 w-4" />
                        Stop recording
                      </Button>
                      <StatusBadge state={state.recording} />
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="download">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Download className="h-5 w-5" />
                      Local Download
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                    <div className="space-y-1">
                      <p className="text-sm font-medium">
                        {canDownload
                          ? "A local manifest is ready."
                          : "Finish a recording to enable download."}
                      </p>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        This downloads a local JSON manifest now. Dataset zip
                        export will be backed by the local agent contract.
                      </p>
                    </div>
                    <Button onClick={downloadManifest} disabled={!canDownload}>
                      <Download className="h-4 w-4" />
                      Download manifest
                    </Button>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </section>
        </section>
      </div>
    </main>
  );
}

function ProgressRow({
  icon,
  label,
  done,
}: {
  icon: React.ReactNode;
  label: string;
  done: boolean;
}) {
  return (
    <div className="flex items-center justify-between gap-3 text-sm">
      <div className="flex items-center gap-2">
        <span className="text-gray-500">{icon}</span>
        <span>{label}</span>
      </div>
      {done ? (
        <CheckCircle2 className="h-4 w-4 text-emerald-600" />
      ) : (
        <Circle className="h-4 w-4 text-gray-300" />
      )}
    </div>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-md border border-gray-200 bg-gray-50 px-3 py-2 dark:border-gray-800 dark:bg-gray-950">
      <div className="text-xs text-gray-500 dark:text-gray-400">{label}</div>
      <div className="mt-1 truncate text-sm font-medium">{value}</div>
    </div>
  );
}

function ActionCard({
  icon,
  title,
  description,
  state,
  disabled,
  actionLabel,
  onAction,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
  state: StepState;
  disabled: boolean;
  actionLabel: string;
  onAction: () => void;
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          {icon}
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-gray-600 dark:text-gray-400">
          {description}
        </p>
        <div className="flex flex-wrap items-center gap-2">
          <Button onClick={onAction} disabled={disabled || state === "running"}>
            <Play className="h-4 w-4" />
            {actionLabel}
          </Button>
          <StatusBadge state={state} />
        </div>
      </CardContent>
    </Card>
  );
}

function StatusBadge({ state }: { state: StepState | RecordingState }) {
  if (state === "done" || state === "ready") {
    return <Badge variant="success">Ready</Badge>;
  }
  if (state === "running" || state === "recording") {
    return <Badge variant="info">Running</Badge>;
  }
  return <Badge variant="outline">Idle</Badge>;
}
