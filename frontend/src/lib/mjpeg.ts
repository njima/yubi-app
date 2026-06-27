interface MjpegUrlOptions {
  host: string;
  port: number;
  namespace: string;
  topic?: string;
  quality?: number;
  width?: number;
  height?: number;
  type?: string;
}

export function buildMjpegStreamUrl(options: MjpegUrlOptions): string {
  const {
    host,
    port,
    namespace,
    topic = "image_raw",
    quality,
    width,
    height,
    type,
  } = options;
  const topicPath = type ? `/${namespace}` : `/${namespace}/${topic}`;
  let url = `http://${host}:${port}/stream?topic=${topicPath}`;
  if (quality !== undefined) url += `&quality=${quality}`;
  if (width !== undefined) url += `&width=${width}`;
  if (height !== undefined) url += `&height=${height}`;
  if (type !== undefined) url += `&type=${type}`;
  return url;
}

export function buildMjpegViewerUrl(options: MjpegUrlOptions): string {
  const { host, port, namespace, topic = "image_raw" } = options;
  const topicPath = `/${namespace}/${topic}`;
  return `http://${host}:${port}/stream_viewer?topic=${topicPath}`;
}
