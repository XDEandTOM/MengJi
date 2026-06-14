export function isImage(val?: string): boolean {
  return !!val?.startsWith("/uploads/") || !!val?.startsWith("http")
}
